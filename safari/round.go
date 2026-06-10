package safari

import (
	"fmt"
	"strconv"
	"time"
)

type Round struct {
	State         RoundState
	Score         Scoring
	Hydra         *Hydra
	EndsAt        time.Time
	StartedAt     time.Time
	WinnerTeam    int
	EndReason     string
	Paused        bool
	pauseStarted  time.Time
	resetAt       time.Time
	roundMinutes  int
}

func NewRound() *Round {
	return &Round{
		State: RoundIdle,
		Hydra: NewHydra(),
	}
}

func (r *Round) Start(api API, mapCfg MapConfig, roundMinutes int, teams *Teams, marking *Marking) bool {
	r.State = RoundActive
	r.Score = Scoring{}
	r.StartedAt = time.Now()
	r.EndsAt = r.StartedAt.Add(time.Duration(roundMinutes) * time.Minute)
	r.WinnerTeam = 0
	r.EndReason = ""
	r.resetAt = time.Time{}
	r.Paused = false
	r.pauseStarted = time.Time{}
	r.roundMinutes = roundMinutes
	r.Hydra.lastMinute = 0

	vid := r.Hydra.Spawn(api, mapCfg)
	if vid < 0 {
		api.Log("[safari] failed to spawn Hydra")
		r.State = RoundIdle
		return false
	}
	r.Bootstrap(api, mapCfg, teams, marking)
	api.Log(fmt.Sprintf("[safari] round started, hydra vehicle=%d", vid))
	return true
}

func (r *Round) Bootstrap(api API, mapCfg MapConfig, teams *Teams, marking *Marking) {
	teams.ResetRoundState()
	marking.ResetAll()
	teams.TeleportToSpawns(api, mapCfg)
	teams.ApplyLoadouts(api)
	teams.SyncScores(api, r.Score)
}

func (r *Round) End(api API, winnerTeam int, reason string) {
	if r.State != RoundActive {
		return
	}
	r.State = RoundEnded
	r.WinnerTeam = winnerTeam
	r.EndReason = reason
	r.resetAt = time.Now().Add(15 * time.Second)
	r.Paused = false
	r.Hydra.Destroy(api)
}

func (r *Round) TogglePause() bool {
	if r.State != RoundActive {
		return r.Paused
	}
	if !r.Paused {
		r.Paused = true
		r.pauseStarted = time.Now()
	} else {
		if !r.pauseStarted.IsZero() {
			r.EndsAt = r.EndsAt.Add(time.Since(r.pauseStarted))
		}
		r.Paused = false
		r.pauseStarted = time.Time{}
	}
	return r.Paused
}

func (r *Round) TimeLeft() time.Duration {
	if r.State != RoundActive {
		return 0
	}
	left := time.Until(r.EndsAt)
	if left < 0 {
		return 0
	}
	return left
}

func (r *Round) SurvivedSecs() int {
	if r.StartedAt.IsZero() {
		return 0
	}
	end := time.Now()
	if r.State == RoundEnded && !r.resetAt.IsZero() {
		end = r.resetAt.Add(-15 * time.Second)
	}
	return int(end.Sub(r.StartedAt).Seconds())
}

func (r *Round) MaybeReset() bool {
	if r.State == RoundEnded && !r.resetAt.IsZero() && time.Now().After(r.resetAt) {
		r.State = RoundIdle
		r.WinnerTeam = 0
		r.EndReason = ""
		r.Score = Scoring{}
		r.Paused = false
		return true
	}
	return false
}

func (r *Round) CheckTimer() (ended bool, winner int, reason string) {
	if r.State != RoundActive || r.Paused {
		return false, 0, ""
	}
	if time.Now().After(r.EndsAt) {
		winner = r.Score.WinnerByScore()
		reason = fmt.Sprintf("Time expired (Escort %d - Defend %d)", r.Score.EscortScore, r.Score.DefendScore)
		return true, winner, reason
	}
	return false, 0, ""
}

func (r *Round) RoundMinutesStr() string {
	return strconv.Itoa(r.roundMinutes)
}
