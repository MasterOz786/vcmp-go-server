package safari

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CommandResult struct {
	Handled bool
	Deny    bool
}

func normalizeCommand(raw string) (string, bool) {
	cmd := strings.TrimSpace(raw)
	if cmd == "" {
		return "", false
	}
	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}
	return cmd, true
}

func (e *Engine) HandleCommand(playerID int, raw string) CommandResult {
	cmd, ok := normalizeCommand(raw)
	if !ok {
		return CommandResult{}
	}
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return CommandResult{}
	}
	name := strings.ToLower(parts[0])
	args := parts[1:]

	switch name {
	case "/help":
		e.sendHelp(playerID)
		return CommandResult{Handled: true, Deny: true}
	case "/startsafari":
		if !e.api.IsAdmin(playerID) {
			e.api.Send(playerID, ColourRed, "Admin only.")
			return CommandResult{Handled: true, Deny: true}
		}
		if e.round.State == RoundActive {
			e.api.Send(playerID, ColourYellow, "Round already active.")
			return CommandResult{Handled: true, Deny: true}
		}
		e.startRound()
		return CommandResult{Handled: true, Deny: true}
	case "/stopsafari":
		if !e.api.IsAdmin(playerID) {
			e.api.Send(playerID, ColourRed, "Admin only.")
			return CommandResult{Handled: true, Deny: true}
		}
		if e.round.State == RoundActive {
			e.endRound(TeamDefend, "Round stopped by admin.")
		}
		return CommandResult{Handled: true, Deny: true}
	case "/pausesafari":
		if !e.api.IsAdmin(playerID) {
			e.api.Send(playerID, ColourRed, "Admin only.")
			return CommandResult{Handled: true, Deny: true}
		}
		if e.round.State != RoundActive {
			e.api.Send(playerID, ColourYellow, "No active round.")
			return CommandResult{Handled: true, Deny: true}
		}
		e.TogglePause()
		return CommandResult{Handled: true, Deny: true}
	case "/autostart":
		if !e.api.IsAdmin(playerID) {
			e.api.Send(playerID, ColourRed, "Admin only.")
			return CommandResult{Handled: true, Deny: true}
		}
		if len(args) == 0 {
			state := "off"
			if e.AutostartEnabled() {
				state = "on"
			}
			e.api.Send(playerID, ColourWhite, "Autostart is "+state+".")
			return CommandResult{Handled: true, Deny: true}
		}
		switch strings.ToLower(args[0]) {
		case "on", "1", "true":
			e.SetAutostart(true)
		case "off", "0", "false":
			e.SetAutostart(false)
		default:
			e.api.Send(playerID, ColourYellow, "Usage: /autostart on|off")
		}
		return CommandResult{Handled: true, Deny: true}
	case "/pack":
		return e.cmdPack(playerID, args)
	case "/mark":
		return e.cmdMark(playerID, args)
	case "/status":
		e.cmdStatus(playerID)
		return CommandResult{Handled: true, Deny: true}
	case "/stats":
		e.cmdStats(playerID)
		return CommandResult{Handled: true, Deny: true}
	default:
		return CommandResult{}
	}
}

func (e *Engine) sendHelp(playerID int) {
	lines := []string{
		"Project Safari: Hydra Warfare",
		"/pack 1|2 — choose loadout",
		"/mark [name] — Escort: designate target",
		"/status — round info",
		"/stats — your persistent stats",
		"/startsafari — admin: start round",
		"/stopsafari — admin: stop round",
		"/pausesafari — admin: pause/resume round",
		"/autostart on|off — admin: toggle autostart",
	}
	for _, l := range lines {
		e.api.Send(playerID, ColourWhite, l)
	}
}

func (e *Engine) cmdPack(playerID int, args []string) CommandResult {
	if len(args) != 1 {
		e.api.Send(playerID, ColourYellow, "Usage: /pack 1|2")
		return CommandResult{Handled: true, Deny: true}
	}
	pack, err := strconv.Atoi(args[0])
	if err != nil || pack < 1 || pack > 2 {
		e.api.Send(playerID, ColourYellow, "Pack must be 1 or 2.")
		return CommandResult{Handled: true, Deny: true}
	}
	if e.round.State == RoundActive && e.teams.HasSpawnedThisRound(playerID) {
		e.api.Send(playerID, ColourYellow, "Cannot change pack after spawning this round.")
		return CommandResult{Handled: true, Deny: true}
	}
	e.ensurePlayerSession(playerID)
	if e.teams.Team(playerID) == 0 {
		e.api.Send(playerID, ColourRed, "You are not assigned to a team yet.")
		return CommandResult{Handled: true, Deny: true}
	}
	e.ApplyPack(playerID, pack)
	team := e.teams.Team(playerID)
	var name string
	if team == TeamEscort {
		name = EscortPacks()[pack].Name
	} else {
		name = DefendPacks()[pack].Name
	}
	e.api.Send(playerID, ColourGreen, fmt.Sprintf("Loadout equipped: %s", name))
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) cmdMark(playerID int, args []string) CommandResult {
	if e.round.State != RoundActive {
		e.announceTo(playerID, MsgMarkFail, "No active round.")
		return CommandResult{Handled: true, Deny: true}
	}
	target := strings.Join(args, " ")
	ok, msg := e.marking.TryMark(e.api, e.db, playerID, target, &e.round.Score)
	if !ok {
		e.announceTo(playerID, MsgMarkFail, msg)
	} else {
		e.announce(MsgMarkSuccess, msg)
		e.teams.SyncScores(e.api, e.round.Score)
	}
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) cmdStatus(playerID int) {
	switch e.round.State {
	case RoundIdle:
		e.api.Send(playerID, ColourWhite, "Safari idle. Waiting for round start.")
	case RoundActive:
		e.api.Send(playerID, ColourCyan, e.formatStatusLine())
	case RoundEnded:
		e.api.Send(playerID, ColourWhite, fmt.Sprintf("Round ended: %s", e.round.EndReason))
	}
}

func (e *Engine) cmdStats(playerID int) {
	uid := e.api.PlayerUID(playerID)
	if uid == "" {
		e.api.Send(playerID, ColourRed, "Could not load your UID.")
		return
	}
	st, ok := e.db.CachedStats(uid)
	if !ok {
		e.db.PrefetchStats(uid)
		e.api.Send(playerID, ColourYellow, "Stats loading — try again in a moment.")
		return
	}
	e.api.Send(playerID, ColourWhite, fmt.Sprintf(
		"Stats — Escort pts: %d | Defend pts: %d | Marks: %d | Rounds: %d | Wins: %d",
		st.EscortPts, st.DefendPts, st.Marks, st.RoundsPlayed, st.RoundsWon,
	))
}

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "0:00"
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}
