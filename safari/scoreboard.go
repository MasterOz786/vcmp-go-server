package safari

import "time"

func (e *Engine) BroadcastScoreboard() {
	var state int32
	switch e.round.State {
	case RoundActive:
		if e.round.Paused {
			state = 3
		} else {
			state = 1
		}
	case RoundEnded:
		state = 2
	default:
		state = 0
	}

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

	payload := s.Bytes()
	for _, playerID := range e.teams.ConnectedIDs() {
		if e.api.IsConnected(playerID) {
			_ = e.api.SendScriptData(playerID, payload)
		}
	}
}
