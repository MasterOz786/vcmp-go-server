package safari

import (
	"log"
	"sync"
)

type DBJob interface {
	Run(store *Store) error
}

type upsertPlayerJob struct {
	uid  string
	name string
}

func (j upsertPlayerJob) Run(store *Store) error {
	return store.UpsertPlayer(j.uid, j.name)
}

type addMarkJob struct {
	uid string
}

func (j addMarkJob) Run(store *Store) error {
	return store.AddMark(j.uid)
}

type recordRoundJob struct {
	players      []RoundPlayerRecord
	winnerTeam   int
	escortScore  int
	defendScore  int
	survivedSecs int
}

func (j recordRoundJob) Run(store *Store) error {
	return store.RecordRoundEndSimple(j.players, j.winnerTeam, j.escortScore, j.defendScore, j.survivedSecs)
}

type savePackJob struct {
	uid  string
	pack int
	w    *DBWorker
}

func (j savePackJob) Run(store *Store) error {
	if err := store.SetPreferredPack(j.uid, j.pack); err != nil {
		return err
	}
	j.w.setCachedPack(j.uid, j.pack)
	return nil
}

type prefetchPackJob struct {
	uid string
	w   *DBWorker
}

func (j prefetchPackJob) Run(store *Store) error {
	pack, err := store.GetPreferredPack(j.uid)
	if err != nil {
		return err
	}
	j.w.setCachedPackIfAbsent(j.uid, pack)
	return nil
}

type prefetchStatsJob struct {
	uid string
	w   *DBWorker
}

func (j prefetchStatsJob) Run(store *Store) error {
	st, err := store.GetStats(j.uid)
	if err != nil {
		return err
	}
	j.w.setCachedStats(j.uid, st)
	return nil
}

type DBWorker struct {
	store  *Store
	jobs   chan DBJob
	done   chan struct{}
	closed chan struct{}

	mu         sync.RWMutex
	packCache  map[string]int
	statsCache map[string]PlayerStats
}

func NewDBWorker(store *Store, queueSize int) *DBWorker {
	return &DBWorker{
		store:      store,
		jobs:       make(chan DBJob, queueSize),
		done:       make(chan struct{}),
		closed:     make(chan struct{}),
		packCache:  make(map[string]int),
		statsCache: make(map[string]PlayerStats),
	}
}

func (w *DBWorker) Start() {
	go func() {
		defer close(w.closed)
		for {
			select {
			case <-w.done:
				return
			case job, ok := <-w.jobs:
				if !ok {
					return
				}
				if err := job.Run(w.store); err != nil {
					log.Printf("[safari-db] job error: %v", err)
				}
			}
		}
	}()
}

func (w *DBWorker) Enqueue(job DBJob) {
	select {
	case w.jobs <- job:
	default:
		log.Printf("[safari-db] job queue full, dropping job")
	}
}

func (w *DBWorker) setCachedPackIfAbsent(uid string, pack int) {
	if pack < 1 || pack > 2 {
		pack = 1
	}
	w.mu.Lock()
	if _, exists := w.packCache[uid]; !exists {
		w.packCache[uid] = pack
	}
	w.mu.Unlock()
}

func (w *DBWorker) setCachedPack(uid string, pack int) {
	if pack < 1 || pack > 2 {
		pack = 1
	}
	w.mu.Lock()
	w.packCache[uid] = pack
	w.mu.Unlock()
}

func (w *DBWorker) CachedPreferredPack(uid string) (int, bool) {
	w.mu.RLock()
	pack, ok := w.packCache[uid]
	w.mu.RUnlock()
	if !ok || pack < 1 || pack > 2 {
		return 1, false
	}
	return pack, true
}

func (w *DBWorker) PrefetchPreferredPack(uid string) {
	if uid == "" {
		return
	}
	w.Enqueue(prefetchPackJob{uid: uid, w: w})
}

func (w *DBWorker) setCachedStats(uid string, st PlayerStats) {
	w.mu.Lock()
	w.statsCache[uid] = st
	w.mu.Unlock()
}

func (w *DBWorker) CachedStats(uid string) (PlayerStats, bool) {
	w.mu.RLock()
	st, ok := w.statsCache[uid]
	w.mu.RUnlock()
	return st, ok
}

func (w *DBWorker) PrefetchStats(uid string) {
	if uid == "" {
		return
	}
	w.Enqueue(prefetchStatsJob{uid: uid, w: w})
}

func (w *DBWorker) InvalidateStats(uid string) {
	w.mu.Lock()
	delete(w.statsCache, uid)
	w.mu.Unlock()
}

func (w *DBWorker) SavePreferredPack(uid string, pack int) {
	w.setCachedPack(uid, pack)
	w.Enqueue(savePackJob{uid: uid, pack: pack, w: w})
}

func (w *DBWorker) Stop() {
	close(w.done)
	<-w.closed
}
