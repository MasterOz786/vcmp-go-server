package gameplay

import (
	"fmt"

	"github.com/masteroz/vcmp-go-server/safari/apidef"
)

type Teams struct {
	countEscort int
	countDefend int
	sessions    map[int]*apidef.PlayerSession
	connected   map[int]struct{}
}

func NewTeams() *Teams {
	return &Teams{
		sessions:  make(map[int]*apidef.PlayerSession),
		connected: make(map[int]struct{}),
	}
}

func (t *Teams) TrackConnect(playerID int) {
	t.connected[playerID] = struct{}{}
}

func (t *Teams) TrackDisconnect(playerID int) {
	delete(t.connected, playerID)
}

func (t *Teams) ConnectedCount() int {
	return len(t.connected)
}

func (t *Teams) ConnectedIDs() []int {
	ids := make([]int, 0, len(t.connected))
	for id := range t.connected {
		ids = append(ids, id)
	}
	return ids
}

func (t *Teams) Session(playerID int) *apidef.PlayerSession {
	return t.sessions[playerID]
}

func (t *Teams) ForEachSession(fn func(playerID int, s *apidef.PlayerSession)) {
	for playerID, s := range t.sessions {
		fn(playerID, s)
	}
}

func (t *Teams) Assign(api apidef.API, playerID, pack int) int {
	if s := t.sessions[playerID]; s != nil {
		return s.Team
	}
	if pack < 1 || pack > apidef.MaxPack {
		pack = 1
	}
	team := apidef.TeamEscort
	if t.countEscort > t.countDefend {
		team = apidef.TeamDefend
	}
	t.sessions[playerID] = apidef.NewPlayerSession(team, pack)
	api.SetPlayerTeam(playerID, team)
	if team == apidef.TeamEscort {
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
	if s.Team == apidef.TeamEscort {
		t.countEscort--
	} else if s.Team == apidef.TeamDefend {
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
	if pack < 1 || pack > apidef.MaxPack {
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
		s.RoundKills = 0
		s.RoundDeaths = 0
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
	case apidef.TeamEscort:
		return t.countEscort <= t.countDefend
	case apidef.TeamDefend:
		return t.countDefend <= t.countEscort
	}
	return true
}

func (t *Teams) RoleName(team int) string {
	if team == apidef.TeamEscort {
		return "Escort"
	}
	return "Defender"
}

func (t *Teams) Welcome(api apidef.API, playerID int) {
	team := t.Team(playerID)
	if team == 0 {
		return
	}
	colour := apidef.ColourGreen
	if team == apidef.TeamDefend {
		colour = apidef.ColourRed
	}
	api.Send(playerID, colour, fmt.Sprintf("Project Safari: you are on team %s. Use /pack 1|2 and /help.", t.RoleName(team)))
}

func (t *Teams) TeleportToSpawns(api apidef.API, mapCfg apidef.MapConfig) {
	spIdx := 0
	for playerID, s := range t.sessions {
		if !api.IsConnected(playerID) {
			continue
		}
		var spawns []apidef.Vec3
		switch s.Team {
		case apidef.TeamEscort:
			spawns = mapCfg.EscortSpawns
		case apidef.TeamDefend:
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

func (t *Teams) ApplyLoadouts(api apidef.API) {
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

func (t *Teams) SyncScores(api apidef.API, score apidef.Scoring) {
	for playerID, s := range t.sessions {
		if !api.IsConnected(playerID) {
			continue
		}
		pts := score.EscortScore
		if s.Team == apidef.TeamDefend {
			pts = score.DefendScore
		}
		api.SetPlayerScore(playerID, pts)
	}
}

func (t *Teams) SetupClasses(api apidef.API, mapCfg apidef.MapConfig) {
	// Pack weapons are granted by the gamemode; class kits must stay empty.
	weapons := [6]int{0, 0, 0, 0, 0, 0}
	escortSkins := []int{0, 1, 2, 9}
	defendSkins := []int{9, 28, 47, 57}
	for i, sp := range mapCfg.EscortSpawns {
		if i >= 4 {
			break
		}
		skin := escortSkins[i%len(escortSkins)]
		api.AddPlayerClass(apidef.TeamEscort, 0xFF6EC6FF, skin, sp, 0, weapons)
	}
	for i, sp := range mapCfg.DefendSpawns {
		if i >= 4 {
			break
		}
		skin := defendSkins[i%len(defendSkins)]
		api.AddPlayerClass(apidef.TeamDefend, 0xFFFF6E6E, skin, sp, 0, weapons)
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
