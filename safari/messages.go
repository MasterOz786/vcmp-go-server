package safari

import "fmt"

// EventKind identifies a player-facing announcement (Flag Raids ShowStateMessages).
type EventKind int

const (
	MsgServerReady EventKind = iota
	MsgRoundStart
	MsgRoundEnd
	MsgRoundReset
	MsgCheckpoint
	MsgHydraDestroyed
	MsgHydraObjective
	MsgHydraMinute
	MsgMarkSuccess
	MsgMarkFail
	MsgKillFeed
	MsgStatusBroadcast
	MsgPause
	MsgUnpause
	MsgAutostart
)

func formatMessage(kind EventKind, args []string) (uint32, string) {
	switch kind {
	case MsgServerReady:
		return ColourCyan, "Project Safari loaded. /help for commands. Admin: /startsafari"
	case MsgRoundStart:
		if len(args) >= 1 {
			return ColourCyan, fmt.Sprintf("SAFARI ROUND START — Protect the Hydra! (%s min)", args[0])
		}
		return ColourCyan, "SAFARI ROUND START — Protect the Hydra!"
	case MsgRoundEnd:
		if len(args) >= 4 {
			return parseColour(args[0]), fmt.Sprintf("ROUND OVER — %s win! %s (Escort %s - Defend %s)",
				args[1], args[2], args[3], args[4])
		}
		return ColourWhite, "ROUND OVER"
	case MsgRoundReset:
		return ColourWhite, "Ready for next round. /startsafari"
	case MsgCheckpoint:
		if len(args) >= 1 {
			return ColourYellow, args[0]
		}
		return ColourYellow, "Checkpoint reached!"
	case MsgHydraDestroyed:
		return ColourRed, "Hydra destroyed! Defenders win."
	case MsgHydraObjective:
		return ColourGreen, "Hydra reached its objective! Escort wins."
	case MsgHydraMinute:
		if len(args) >= 1 {
			return ColourYellow, args[0]
		}
		return ColourYellow, "Hydra holding route (+escort)"
	case MsgMarkSuccess:
		if len(args) >= 1 {
			return ColourGreen, fmt.Sprintf("[ESCORT] Target marked: %s (+15)", args[0])
		}
		return ColourGreen, "[ESCORT] Target marked (+15)"
	case MsgMarkFail:
		if len(args) >= 1 {
			return ColourYellow, args[0]
		}
		return ColourYellow, "Mark failed."
	case MsgKillFeed:
		if len(args) >= 1 {
			return ColourWhite, args[0]
		}
		return ColourWhite, ""
	case MsgStatusBroadcast:
		if len(args) >= 1 {
			return ColourCyan, args[0]
		}
		return ColourCyan, ""
	case MsgPause:
		return ColourYellow, "Safari round paused."
	case MsgUnpause:
		return ColourGreen, "Safari round resumed."
	case MsgAutostart:
		if len(args) >= 1 {
			return ColourWhite, fmt.Sprintf("Autostart %s.", args[0])
		}
		return ColourWhite, "Autostart toggled."
	default:
		return ColourWhite, ""
	}
}

func parseColour(name string) uint32 {
	switch name {
	case "green":
		return ColourGreen
	case "red":
		return ColourRed
	default:
		return ColourWhite
	}
}

func (e *Engine) announce(kind EventKind, args ...string) {
	colour, msg := formatMessage(kind, args)
	if msg == "" {
		return
	}
	e.api.Broadcast(colour, msg)
	e.api.Log("[safari] " + msg)
}

func (e *Engine) announceTo(playerID int, kind EventKind, args ...string) {
	colour, msg := formatMessage(kind, args)
	if msg == "" {
		return
	}
	e.api.Send(playerID, colour, msg)
}

func (e *Engine) formatStatusLine() string {
	hp := e.round.Hydra.Health(e.api)
	idx := e.round.Hydra.Index
	total := len(e.round.Hydra.Waypoints)
	left := e.round.TimeLeft()
	return fmt.Sprintf(
		"[Safari] Hydra HP: %.0f/%.0f | Checkpoint: %d/%d | Escort %d - Defend %d | Time: %s",
		hp, HydraMaxHP, idx, total,
		e.round.Score.EscortScore, e.round.Score.DefendScore,
		formatDuration(left),
	)
}
