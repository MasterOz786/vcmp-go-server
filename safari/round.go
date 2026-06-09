package safari

import (
	"fmt"
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
	resetAt       time.Time
}

func NewRound() *Round {
	return &Round{
		State: RoundIdle,
		Hydra: NewHydra(),
	}
}

func (r *Round) Start(api API, mapCfg MapConfig, roundMinutes int) {
	r.State = RoundActive
	r.Score = Scoring{}
	r.StartedAt = time.Now()
	r.EndsAt = r.StartedAt.Add(time.Duration(roundMinutes) * time.Minute)
	r.WinnerTeam = 0
	r.EndReason = ""
	r.resetAt = time.Time{}
	r.Hydra.lastMinute = 0

	vid := r.Hydra.Spawn(api, mapCfg)
	if vid < 0 {
		api.Log("[safari] failed to spawn Hydra")
		r.State = RoundIdle
		return
	}
	api.Broadcast(ColourCyan, fmt.Sprintf("SAFARI ROUND START — Protect the Hydra! (%d min)", roundMinutes))
	api.Log(fmt.Sprintf("[safari] round started, hydra vehicle=%d", vid))
}

func (r *Round) End(api API, winnerTeam int, reason string) {
	if r.State != RoundActive {
		return
	}
	r.State = RoundEnded
	r.WinnerTeam = winnerTeam
	r.EndReason = reason
	r.resetAt = time.Now().Add(15 * time.Second)
	r.Hydra.Destroy(api)

	teamName := "Escort"
	colour := ColourGreen
	if winnerTeam == TeamDefend {
		teamName = "Defenders"
		colour = ColourRed
	}
	api.Broadcast(colour, fmt.Sprintf("ROUND OVER — %s win! %s (Escort %d - Defend %d)",
		teamName, reason, r.Score.EscortScore, r.Score.DefendScore))
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
		return true
	}
	return false
}

func (r *Round) CheckTimer(api API) (ended bool, winner int, reason string) {
	if r.State != RoundActive {
		return false, 0, ""
	}
	if time.Now().After(r.EndsAt) {
		winner = r.Score.WinnerByScore()
		reason = fmt.Sprintf("Time expired (Escort %d - Defend %d)", r.Score.EscortScore, r.Score.DefendScore)
		return true, winner, reason
	}
	return false, 0, ""
}
