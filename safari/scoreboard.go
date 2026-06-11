package safari

import "time"

const (
	scoreboardStateAuto    int32 = -1
	scoreboardStateHidden  int32 = 0
	scoreboardStatePreview int32 = 4
)

func (e *Engine) scoreboardState(forceState int32) int32 {
	if forceState >= 0 {
		return forceState
	}
	switch e.round.State {
	case RoundActive:
		if e.round.Paused {
			return 3
		}
		return 1
	case RoundEnded:
		return 2
	default:
		return 0
	}
}

func (e *Engine) buildScoreboardPacket(forceState int32) []byte {
	state := e.scoreboardState(forceState)

	left := e.round.TimeLeft()
	mins := int32(left / time.Minute)
	secs := int32(left.Seconds()) % 60

	hydraHP := float32(0)
	cpIdx, cpTotal := int32(0), int32(0)
	if e.round.State == RoundActive && e.round.Hydra.VehicleID >= 0 {
		hydraHP = e.round.Hydra.Health(e.api)
		cpIdx = int32(e.round.Hydra.Index)
		cpTotal = int32(len(e.round.Hydra.Waypoints))
	}

	s := NewStreamWriter()
	s.WriteInt(PacketScoreboard)
	s.WriteInt(int32(e.round.Score.EscortScore))
	s.WriteInt(int32(e.round.Score.DefendScore))
	s.WriteInt(state)
	s.WriteInt(mins)
	s.WriteInt(secs)
	s.WriteFloat(hydraHP)
	s.WriteInt(cpIdx)
	s.WriteInt(cpTotal)
	return s.Bytes()
}

func (e *Engine) SendScoreboardTo(playerID int, forceState int32) {
	if !e.api.IsConnected(playerID) {
		return
	}
	_ = e.api.SendScriptData(playerID, e.buildScoreboardPacket(forceState))
}

func (e *Engine) SendScoreboardHide(playerID int) {
	e.SendScoreboardTo(playerID, scoreboardStateHidden)
}

func (e *Engine) BroadcastScoreboard() {
	payload := e.buildScoreboardPacket(scoreboardStateAuto)
	for _, playerID := range e.teams.ConnectedIDs() {
		if e.api.IsConnected(playerID) {
			_ = e.api.SendScriptData(playerID, payload)
		}
	}
}
