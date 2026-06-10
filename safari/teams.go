package safari

import "fmt"

type Teams struct {
	countEscort int
	countDefend int
	sessions    map[int]*PlayerSession
}

func NewTeams() *Teams {
	return &Teams{sessions: make(map[int]*PlayerSession)}
}

func (t *Teams) session(playerID int) *PlayerSession {
	return t.sessions[playerID]
}

func (t *Teams) Assign(api API, playerID, pack int) int {
	if s := t.sessions[playerID]; s != nil {
		return s.Team
	}
	if pack < 1 || pack > MaxPack {
		pack = 1
	}
	team := TeamEscort
	if t.countEscort > t.countDefend {
		team = TeamDefend
	}
	t.sessions[playerID] = newPlayerSession(team, pack)
	api.SetPlayerTeam(playerID, team)
	if team == TeamEscort {
		t.countEscort++
	} else {
		t.countDefend++
	}
	return team
}

func (t *Teams) Remove(playerID int) {
	s, ok := t.sessions[playerID]
	if !ok {
		return
	}
	if s.Team == TeamEscort {
		t.countEscort--
	} else if s.Team == TeamDefend {
		t.countDefend--
	}
	delete(t.sessions, playerID)
}

func (t *Teams) Team(playerID int) int {
	if s := t.sessions[playerID]; s != nil {
		return s.Team
	}
	return 0
}

func (t *Teams) SetPack(playerID, pack int) bool {
	if pack < 1 || pack > MaxPack {
		return false
	}
	s := t.sessions[playerID]
	if s == nil {
		return false
	}
	s.Pack = pack
	return true
}

func (t *Teams) Pack(playerID int) int {
	if s := t.sessions[playerID]; s != nil {
		return s.Pack
	}
	return 1
}

func (t *Teams) HasSpawnedThisRound(playerID int) bool {
	if s := t.sessions[playerID]; s != nil {
		return s.HasSpawnedThisRound
	}
	return false
}

func (t *Teams) MarkSpawned(playerID int) {
	if s := t.sessions[playerID]; s != nil {
		s.HasSpawnedThisRound = true
	}
}

func (t *Teams) ResetRoundState() {
	for _, s := range t.sessions {
		s.SpawnIndex = 0
		s.HasSpawnedThisRound = false
	}
}

func (t *Teams) AdvanceSpawn(playerID int) {
	if s := t.sessions[playerID]; s != nil {
		s.SpawnIndex++
	}
}

func (t *Teams) CountEscort() int { return t.countEscort }
func (t *Teams) CountDefend() int { return t.countDefend }

func (t *Teams) AllowClassRequest(playerID, classIndex int, roundActive bool) bool {
	if classIndex < 0 || classIndex > 3 {
		return false
	}
	if !roundActive {
		return true
	}
	team := t.Team(playerID)
	if team == 0 {
		return true
	}
	switch team {
	case TeamEscort:
		return t.countEscort <= t.countDefend
	case TeamDefend:
		return t.countDefend <= t.countEscort
	}
	return true
}

func (t *Teams) RoleName(team int) string {
	if team == TeamEscort {
		return "Escort"
	}
	return "Defender"
}

func (t *Teams) Welcome(api API, playerID int) {
	team := t.Team(playerID)
	if team == 0 {
		return
	}
	colour := ColourGreen
	if team == TeamDefend {
		colour = ColourRed
	}
	api.Send(playerID, colour, fmt.Sprintf("Project Safari: you are on team %s. Use /pack 1|2 and /help.", t.RoleName(team)))
}

func (t *Teams) TeleportToSpawns(api API, mapCfg MapConfig) {
	spIdx := 0
	for playerID, s := range t.sessions {
		if !api.IsConnected(playerID) {
			continue
		}
		var spawns []Vec3
		switch s.Team {
		case TeamEscort:
			spawns = mapCfg.EscortSpawns
		case TeamDefend:
			spawns = mapCfg.DefendSpawns
		default:
			continue
		}
		if len(spawns) == 0 {
			continue
		}
		pos := spawns[spIdx%len(spawns)]
		_ = api.SetPlayerPosition(playerID, pos)
		s.SpawnIndex = spIdx % len(spawns)
		spIdx++
	}
}

func (t *Teams) ApplyLoadouts(api API) {
	for playerID, s := range t.sessions {
		if !api.IsConnected(playerID) || s.Team == 0 {
			continue
		}
		if api.IsSpawned(playerID) {
			ApplyLoadout(api, playerID, s.Team, s.Pack)
			EnforceAllowed(api, playerID, s.Team, s.Pack)
		}
	}
}

func (t *Teams) SyncScores(api API, score Scoring) {
	for playerID, s := range t.sessions {
		if !api.IsConnected(playerID) {
			continue
		}
		pts := score.EscortScore
		if s.Team == TeamDefend {
			pts = score.DefendScore
		}
		api.SetPlayerScore(playerID, pts)
	}
}

func (t *Teams) SetupClasses(api API, mapCfg MapConfig) {
	// Pack weapons are granted by the gamemode; class kits must stay empty.
	weapons := [6]int{0, 0, 0, 0, 0, 0}
	escortSkins := []int{0, 1, 2, 9}
	defendSkins := []int{9, 28, 47, 57}
	for i, sp := range mapCfg.EscortSpawns {
		if i >= 4 {
			break
		}
		skin := escortSkins[i%len(escortSkins)]
		api.AddPlayerClass(TeamEscort, 0xFF6EC6FF, skin, sp, 0, weapons)
	}
	for i, sp := range mapCfg.DefendSpawns {
		if i >= 4 {
			break
		}
		skin := defendSkins[i%len(defendSkins)]
		api.AddPlayerClass(TeamDefend, 0xFFFF6E6E, skin, sp, 0, weapons)
	}
	if len(mapCfg.EscortSpawns) > 0 {
		api.SetSpawnPos(mapCfg.EscortSpawns[0])
	}
}

// ServerOption constants mirrored for API (vcmpServerOptionUseClasses = 18).
const ServerOptionUseClasses = 18
const ServerOptionJoinMessages = 15
const ServerOptionDeathMessages = 16
const ServerOptionDisableDriveBy = 8
const ServerOptionFastSwitch = 6
const ServerOptionStuntBike = 14
const ServerOptionWallGlitch = 19
const ServerOptionDisableHeliBladeDamage = 21
