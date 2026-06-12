package persist

import (
	"log"
	"sync"

	"github.com/masteroz/vcmp-go-server/safari/apidef"
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

type prefetchRegisteredJob struct {
	uid string
	w   *DBWorker
}

func (j prefetchRegisteredJob) Run(store *Store) error {
	ok, err := store.IsRegistered(j.uid)
	if err != nil {
		return err
	}
	j.w.setCachedRegistered(j.uid, ok)
	return nil
}

type refreshLeaderboardJob struct {
	w    *DBWorker
	done chan struct{}
}

func (j refreshLeaderboardJob) Run(store *Store) error {
	escort, err := store.TopEscortLeaderboard(10)
	if err != nil {
		log.Printf("[safari-db] escort leaderboard query failed: %v", err)
		escort = nil
	}
	defend, err := store.TopDefendLeaderboard(10)
	if err != nil {
		log.Printf("[safari-db] defend leaderboard query failed: %v", err)
		defend = nil
	}
	j.w.setLeaderboards(escort, defend)
	if j.done != nil {
		close(j.done)
	}
	return nil
}

type DBWorker struct {
	store  *Store
	jobs   chan DBJob
	done   chan struct{}
	closed chan struct{}

	mu              sync.RWMutex
	packCache       map[string]int
	statsCache      map[string]PlayerStats
	registeredCache map[string]bool
	escortLB        []LeaderboardEntry
	defendLB        []LeaderboardEntry
	lbReady         bool
}

func NewDBWorker(store *Store, queueSize int) *DBWorker {
	return &DBWorker{
		store:           store,
		jobs:            make(chan DBJob, queueSize),
		done:            make(chan struct{}),
		closed:          make(chan struct{}),
		packCache:       make(map[string]int),
		statsCache:      make(map[string]PlayerStats),
		registeredCache: make(map[string]bool),
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

func (w *DBWorker) EnqueueMark(uid string) {
	w.Enqueue(addMarkJob{uid: uid})
}

func (w *DBWorker) EnqueueUpsertPlayer(uid, name string) {
	w.Enqueue(upsertPlayerJob{uid: uid, name: name})
}

func (w *DBWorker) EnqueueRecordRound(players []RoundPlayerRecord, winnerTeam, escortScore, defendScore, survivedSecs int) {
	w.Enqueue(recordRoundJob{
		players:      players,
		winnerTeam:   winnerTeam,
		escortScore:  escortScore,
		defendScore:  defendScore,
		survivedSecs: survivedSecs,
	})
}

func clampPack(pack int) int {
	if pack < 1 || pack > apidef.MaxPack {
		return 1
	}
	return pack
}

func (w *DBWorker) setCachedPackIfAbsent(uid string, pack int) {
	pack = clampPack(pack)
	w.mu.Lock()
	if _, exists := w.packCache[uid]; !exists {
		w.packCache[uid] = pack
	}
	w.mu.Unlock()
}

func (w *DBWorker) setCachedPack(uid string, pack int) {
	pack = clampPack(pack)
	w.mu.Lock()
	w.packCache[uid] = pack
	w.mu.Unlock()
}

func (w *DBWorker) CachedPreferredPack(uid string) (int, bool) {
	w.mu.RLock()
	pack, ok := w.packCache[uid]
	w.mu.RUnlock()
	if !ok {
		return 1, false
	}
	return clampPack(pack), true
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

func (w *DBWorker) setCachedRegistered(uid string, registered bool) {
	w.mu.Lock()
	w.registeredCache[uid] = registered
	w.mu.Unlock()
}

func (w *DBWorker) CachedRegistered(uid string) (bool, bool) {
	w.mu.RLock()
	v, ok := w.registeredCache[uid]
	w.mu.RUnlock()
	return v, ok
}

func (w *DBWorker) PrefetchRegistered(uid string) {
	if uid == "" {
		return
	}
	w.Enqueue(prefetchRegisteredJob{uid: uid, w: w})
}

func (w *DBWorker) IsRegistered(uid string) (bool, error) {
	if uid == "" {
		return false, nil
	}
	if v, ok := w.CachedRegistered(uid); ok {
		return v, nil
	}
	return w.store.IsRegistered(uid)
}

func (w *DBWorker) RegisterAccount(uid, name, passwordHash string) error {
	if err := w.store.RegisterAccount(uid, name, passwordHash); err != nil {
		return err
	}
	w.setCachedRegistered(uid, true)
	return nil
}

func (w *DBWorker) setLeaderboards(escort, defend []LeaderboardEntry) {
	w.mu.Lock()
	w.escortLB = escort
	w.defendLB = defend
	w.lbReady = true
	w.mu.Unlock()
}

func (w *DBWorker) Leaderboards() (escort, defend []LeaderboardEntry, ok bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if !w.lbReady {
		return nil, nil, false
	}
	return w.escortLB, w.defendLB, true
}

func (w *DBWorker) RefreshLeaderboardsAsync() <-chan struct{} {
	done := make(chan struct{})
	w.Enqueue(refreshLeaderboardJob{w: w, done: done})
	return done
}

func (w *DBWorker) Stop() {
	close(w.done)
	<-w.closed
}
