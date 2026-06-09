package safari

import "log"

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

type getStatsJob struct {
	uid   string
	reply chan<- PlayerStats
	err   chan<- error
}

func (j getStatsJob) Run(store *Store) error {
	st, err := store.GetStats(j.uid)
	if err != nil {
		j.err <- err
		return err
	}
	j.reply <- st
	return nil
}

type DBWorker struct {
	store  *Store
	jobs   chan DBJob
	done   chan struct{}
	closed chan struct{}
}

func NewDBWorker(store *Store, queueSize int) *DBWorker {
	return &DBWorker{
		store:  store,
		jobs:   make(chan DBJob, queueSize),
		done:   make(chan struct{}),
		closed: make(chan struct{}),
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
					if _, isGet := job.(getStatsJob); !isGet {
						log.Printf("[safari-db] job error: %v", err)
					}
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

func (w *DBWorker) GetStats(uid string) (PlayerStats, error) {
	reply := make(chan PlayerStats, 1)
	errCh := make(chan error, 1)
	w.Enqueue(getStatsJob{uid: uid, reply: reply, err: errCh})
	select {
	case st := <-reply:
		return st, nil
	case err := <-errCh:
		return PlayerStats{UID: uid}, err
	}
}

func (w *DBWorker) Stop() {
	close(w.done)
	<-w.closed
}
