package safari

import "fmt"

func (e *Engine) SendShowPacks(playerID int) {
	if !e.api.IsConnected(playerID) {
		return
	}
	team := e.teams.Team(playerID)
	if team == 0 {
		e.ensurePlayerSession(playerID)
		team = e.teams.Team(playerID)
	}
	s := NewStreamWriter()
	s.WriteInt(PacketShowPacks)
	s.WriteInt(int32(team))
	s.WriteInt(int32(e.teams.Pack(playerID)))
	if err := e.api.SendScriptData(playerID, s.Bytes()); err != nil {
		e.api.Log(fmt.Sprintf("[safari-stream] SHOW_PACKS to %d failed: %v", playerID, err))
	}
}

func (e *Engine) SendHidePacks(playerID int) {
	if !e.api.IsConnected(playerID) {
		return
	}
	s := NewStreamWriter()
	s.WriteInt(PacketHidePacks)
	_ = e.api.SendScriptData(playerID, s.Bytes())
}

func (e *Engine) handleSelectPack(playerID, pack int) {
	if pack < 1 || pack > MaxPack {
		e.api.Send(playerID, ColourYellow, fmt.Sprintf("Pack must be 1 to %d.", MaxPack))
		return
	}
	if e.round.State == RoundActive && e.teams.HasSpawnedThisRound(playerID) {
		e.api.Send(playerID, ColourYellow, "Cannot change pack after spawning this round.")
		return
	}
	e.ensurePlayerSession(playerID)
	if e.teams.Team(playerID) == 0 {
		e.api.Send(playerID, ColourRed, "You are not assigned to a team yet.")
		return
	}
	e.ApplyPack(playerID, pack)
	e.SendHidePacks(playerID)
	team := e.teams.Team(playerID)
	var name string
	if team == TeamEscort {
		name = EscortPacks()[pack].Name
	} else {
		name = DefendPacks()[pack].Name
	}
	e.api.Send(playerID, ColourGreen, fmt.Sprintf("Loadout equipped: %s", name))
}

func (e *Engine) handleRequestShowPacks(playerID int) {
	if e.round.State == RoundActive && e.teams.HasSpawnedThisRound(playerID) {
		e.api.Send(playerID, ColourYellow, "Cannot change pack after spawning this round.")
		return
	}
	e.SendShowPacks(playerID)
}

func (e *Engine) BroadcastRoundEndStats(winnerTeam int, reason string) {
	escortN, defendN := 0, 0
	for _, id := range e.teams.ConnectedIDs() {
		if !e.api.IsConnected(id) {
			continue
		}
		switch e.teams.Team(id) {
		case TeamEscort:
			escortN++
		case TeamDefend:
			defendN++
		}
	}
	escortShare, defendShare := 0, 0
	if escortN > 0 {
		escortShare = e.round.Score.EscortScore / escortN
	}
	if defendN > 0 {
		defendShare = e.round.Score.DefendScore / defendN
	}

	var players []struct {
		name   string
		team   int
		points int
		kills  int
		deaths int
	}
	for _, id := range e.teams.ConnectedIDs() {
		if !e.api.IsConnected(id) {
			continue
		}
		s := e.teams.session(id)
		if s == nil {
			continue
		}
		pts := escortShare
		if s.Team == TeamDefend {
			pts = defendShare
		}
		players = append(players, struct {
			name   string
			team   int
			points int
			kills  int
			deaths int
		}{
			name:   e.api.PlayerName(id),
			team:   s.Team,
			points: pts,
			kills:  s.RoundKills,
			deaths: s.RoundDeaths,
		})
	}

	s := NewStreamWriter()
	s.WriteInt(PacketRoundEndStats)
	s.WriteInt(int32(winnerTeam))
	s.WriteInt(int32(e.round.Score.EscortScore))
	s.WriteInt(int32(e.round.Score.DefendScore))
	s.WriteString(reason)
	s.WriteInt(int32(len(players)))
	for _, p := range players {
		s.WriteString(p.name)
		s.WriteInt(int32(p.team))
		s.WriteInt(int32(p.points))
		s.WriteInt(int32(p.kills))
		s.WriteInt(int32(p.deaths))
	}

	payload := s.Bytes()
	for _, playerID := range e.teams.ConnectedIDs() {
		if e.api.IsConnected(playerID) {
			_ = e.api.SendScriptData(playerID, payload)
		}
	}
}

func (e *Engine) BroadcastHideRoundStats() {
	s := NewStreamWriter()
	s.WriteInt(PacketRoundEndStats)
	s.WriteInt(-1) // hide signal
	payload := s.Bytes()
	for _, playerID := range e.teams.ConnectedIDs() {
		if e.api.IsConnected(playerID) {
			_ = e.api.SendScriptData(playerID, payload)
		}
	}
}
