package gameplay

import (
	"fmt"
	"strconv"
	"time"

	"github.com/masteroz/vcmp-go-server/safari/apidef"
)

type Round struct {
	State        apidef.RoundState
	Score        apidef.Scoring
	Hydra        *Hydra
	EndsAt       time.Time
	StartedAt    time.Time
	WinnerTeam   int
	EndReason    string
	Paused       bool
	pauseStarted time.Time
	resetAt      time.Time
	roundMinutes int
}

func NewRound() *Round {
	return &Round{
		State: apidef.RoundIdle,
		Hydra: NewHydra(),
	}
}

func (r *Round) Start(api apidef.API, mapCfg apidef.MapConfig, roundMinutes int, teams *Teams, marking *Marking, hydraModel int) bool {
	r.State = apidef.RoundActive
	r.Score = apidef.Scoring{}
	r.StartedAt = time.Now()
	r.EndsAt = r.StartedAt.Add(time.Duration(roundMinutes) * time.Minute)
	r.WinnerTeam = 0
	r.EndReason = ""
	r.resetAt = time.Time{}
	r.Paused = false
	r.pauseStarted = time.Time{}
	r.roundMinutes = roundMinutes
	r.Hydra.lastMinute = 0

	vid := r.Hydra.Spawn(api, mapCfg, hydraModel)
	if vid < 0 {
		api.Log("[safari] failed to spawn Hydra")
		r.State = apidef.RoundIdle
		return false
	}
	r.Bootstrap(api, mapCfg, teams, marking)
	api.Log(fmt.Sprintf("[safari] round started, hydra vehicle=%d", vid))
	return true
}

func (r *Round) Bootstrap(api apidef.API, mapCfg apidef.MapConfig, teams *Teams, marking *Marking) {
	teams.ResetRoundState()
	marking.ResetAll()
	teams.TeleportToSpawns(api, mapCfg)
	teams.ApplyLoadouts(api)
	teams.SyncScores(api, r.Score)
}

func (r *Round) End(api apidef.API, winnerTeam int, reason string) {
	if r.State != apidef.RoundActive {
		return
	}
	r.State = apidef.RoundEnded
	r.WinnerTeam = winnerTeam
	r.EndReason = reason
	r.resetAt = time.Now().Add(15 * time.Second)
	r.Paused = false
	r.Hydra.Destroy(api)
}

func (r *Round) TogglePause() bool {
	if r.State != apidef.RoundActive {
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
	if r.State != apidef.RoundActive {
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
	if r.State == apidef.RoundEnded && !r.resetAt.IsZero() {
		end = r.resetAt.Add(-15 * time.Second)
	}
	return int(end.Sub(r.StartedAt).Seconds())
}

func (r *Round) MaybeReset() bool {
	if r.State == apidef.RoundEnded && !r.resetAt.IsZero() && time.Now().After(r.resetAt) {
		r.State = apidef.RoundIdle
		r.WinnerTeam = 0
		r.EndReason = ""
		r.Score = apidef.Scoring{}
		r.Paused = false
		return true
	}
	return false
}

func (r *Round) CheckTimer() (ended bool, winner int, reason string) {
	if r.State != apidef.RoundActive || r.Paused {
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
