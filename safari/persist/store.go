package persist

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/masteroz/vcmp-go-server/safari/apidef"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS players (
  uid TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  last_seen INTEGER NOT NULL,
  password_hash TEXT NOT NULL DEFAULT '',
  registered INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS player_stats (
  uid TEXT PRIMARY KEY REFERENCES players(uid),
  escort_pts INTEGER DEFAULT 0,
  defend_pts INTEGER DEFAULT 0,
  marks INTEGER DEFAULT 0,
  rounds_played INTEGER DEFAULT 0,
  rounds_won INTEGER DEFAULT 0,
  preferred_pack INTEGER NOT NULL DEFAULT 1,
  preferred_team INTEGER NOT NULL DEFAULT 0,
  preferred_skin INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS match_history (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  ended_at INTEGER NOT NULL,
  winner_team INTEGER NOT NULL,
  escort_score INTEGER NOT NULL,
  defend_score INTEGER NOT NULL,
  hydra_survived_secs INTEGER NOT NULL,
  mode TEXT NOT NULL DEFAULT 'patrol'
);
`

type Store struct {
	db *sql.DB
}

func OpenStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close()
		return nil, err
	}
	_, _ = db.Exec(`ALTER TABLE player_stats ADD COLUMN preferred_pack INTEGER NOT NULL DEFAULT 1`)
	_, _ = db.Exec(`ALTER TABLE player_stats ADD COLUMN preferred_team INTEGER NOT NULL DEFAULT 0`)
	_, _ = db.Exec(`ALTER TABLE player_stats ADD COLUMN preferred_skin INTEGER NOT NULL DEFAULT 0`)
	_, _ = db.Exec(`ALTER TABLE players ADD COLUMN password_hash TEXT NOT NULL DEFAULT ''`)
	_, _ = db.Exec(`ALTER TABLE players ADD COLUMN registered INTEGER NOT NULL DEFAULT 0`)
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) IsRegistered(uid string) (bool, error) {
	var registered int
	err := s.db.QueryRow(`SELECT registered FROM players WHERE uid = ?`, uid).Scan(&registered)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return registered != 0, nil
}

func (s *Store) RegisterAccount(uid, name, passwordHash string) error {
	now := time.Now().Unix()
	_, err := s.db.Exec(
		`INSERT INTO players (uid, name, last_seen, password_hash, registered)
		 VALUES (?, ?, ?, ?, 1)
		 ON CONFLICT(uid) DO UPDATE SET
		   name = excluded.name,
		   last_seen = excluded.last_seen,
		   password_hash = excluded.password_hash,
		   registered = 1`,
		uid, name, now, passwordHash,
	)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO player_stats (uid) VALUES (?) ON CONFLICT(uid) DO NOTHING`, uid)
	return err
}

func (s *Store) UpsertPlayer(uid, name string) error {
	now := time.Now().Unix()
	_, err := s.db.Exec(
		`INSERT INTO players (uid, name, last_seen) VALUES (?, ?, ?)
		 ON CONFLICT(uid) DO UPDATE SET name=excluded.name, last_seen=excluded.last_seen`,
		uid, name, now,
	)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		`INSERT INTO player_stats (uid) VALUES (?) ON CONFLICT(uid) DO NOTHING`,
		uid,
	)
	return err
}

func (s *Store) TopEscortLeaderboard(limit int) ([]LeaderboardEntry, error) {
	return s.queryLeaderboard(limit, "escort_pts")
}

func (s *Store) TopDefendLeaderboard(limit int) ([]LeaderboardEntry, error) {
	return s.queryLeaderboard(limit, "defend_pts")
}

func (s *Store) queryLeaderboard(limit int, ptsColumn string) ([]LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	query := fmt.Sprintf(`
		SELECT p.name, ps.%s, ps.marks, ps.rounds_won
		FROM player_stats ps
		JOIN players p ON p.uid = ps.uid
		ORDER BY ps.%s DESC, ps.marks DESC, ps.rounds_won DESC, p.name ASC
		LIMIT ?`, ptsColumn, ptsColumn)
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []LeaderboardEntry
	for rows.Next() {
		var row LeaderboardEntry
		if err := rows.Scan(&row.Name, &row.Points, &row.Marks, &row.Wins); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (s *Store) GetStats(uid string) (PlayerStats, error) {
	var st PlayerStats
	st.UID = uid
	err := s.db.QueryRow(
		`SELECT escort_pts, defend_pts, marks, rounds_played, rounds_won
		 FROM player_stats WHERE uid = ?`, uid,
	).Scan(&st.EscortPts, &st.DefendPts, &st.Marks, &st.RoundsPlayed, &st.RoundsWon)
	if err == sql.ErrNoRows {
		return st, nil
	}
	return st, err
}

func (s *Store) GetPreferredPack(uid string) (int, error) {
	var pack int
	err := s.db.QueryRow(
		`SELECT preferred_pack FROM player_stats WHERE uid = ?`, uid,
	).Scan(&pack)
	if err == sql.ErrNoRows {
		return 1, nil
	}
	if err != nil {
		return 1, err
	}
	if pack < 1 || pack > apidef.MaxPack {
		return 1, nil
	}
	return pack, nil
}

func (s *Store) SetPreferredPack(uid string, pack int) error {
	if pack < 1 || pack > apidef.MaxPack {
		pack = 1
	}
	_, err := s.db.Exec(
		`INSERT INTO player_stats (uid, preferred_pack) VALUES (?, ?)
		 ON CONFLICT(uid) DO UPDATE SET preferred_pack = excluded.preferred_pack`,
		uid, pack,
	)
	return err
}

func clampPreferredTeam(team int) int {
	if team == apidef.TeamEscort || team == apidef.TeamDefend {
		return team
	}
	return 0
}

func clampPreferredSkin(skin int) int {
	if skin < 0 {
		return 0
	}
	if skin >= apidef.MaxSkin {
		return apidef.MaxSkin - 1
	}
	return skin
}

func (s *Store) GetPreferredTeam(uid string) (int, error) {
	var team int
	err := s.db.QueryRow(
		`SELECT preferred_team FROM player_stats WHERE uid = ?`, uid,
	).Scan(&team)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return clampPreferredTeam(team), nil
}

func (s *Store) SetPreferredTeam(uid string, team int) error {
	team = clampPreferredTeam(team)
	_, err := s.db.Exec(
		`INSERT INTO player_stats (uid, preferred_team) VALUES (?, ?)
		 ON CONFLICT(uid) DO UPDATE SET preferred_team = excluded.preferred_team`,
		uid, team,
	)
	return err
}

func (s *Store) GetPreferredSkin(uid string) (int, error) {
	var skin int
	err := s.db.QueryRow(
		`SELECT preferred_skin FROM player_stats WHERE uid = ?`, uid,
	).Scan(&skin)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return clampPreferredSkin(skin), nil
}

func (s *Store) SetPreferredSkin(uid string, skin int) error {
	skin = clampPreferredSkin(skin)
	_, err := s.db.Exec(
		`INSERT INTO player_stats (uid, preferred_skin) VALUES (?, ?)
		 ON CONFLICT(uid) DO UPDATE SET preferred_skin = excluded.preferred_skin`,
		uid, skin,
	)
	return err
}

func (s *Store) AddMark(uid string) error {
	_, err := s.db.Exec(`UPDATE player_stats SET marks = marks + 1 WHERE uid = ?`, uid)
	return err
}

func (s *Store) RecordRoundEndSimple(players []RoundPlayerRecord, winnerTeam, escortScore, defendScore, survivedSecs int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now().Unix()
	if _, err := tx.Exec(
		`INSERT INTO match_history (ended_at, winner_team, escort_score, defend_score, hydra_survived_secs, mode)
		 VALUES (?, ?, ?, ?, ?, 'patrol')`,
		now, winnerTeam, escortScore, defendScore, survivedSecs,
	); err != nil {
		return err
	}

	for _, p := range players {
		won := 0
		if p.Team == winnerTeam {
			won = 1
		}
		if _, err := tx.Exec(
			`UPDATE player_stats SET
			   escort_pts = escort_pts + ?,
			   defend_pts = defend_pts + ?,
			   marks = marks + ?,
			   rounds_played = rounds_played + 1,
			   rounds_won = rounds_won + ?
			 WHERE uid = ?`,
			p.EscortPts, p.DefendPts, p.MarksAdded, won, p.UID,
		); err != nil {
			return fmt.Errorf("update stats for %s: %w", p.UID, err)
		}
	}
	return tx.Commit()
}
