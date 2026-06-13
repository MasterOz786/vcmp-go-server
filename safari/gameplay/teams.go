package gameplay

import (
	"fmt"

	"github.com/masteroz/vcmp-go-server/safari/apidef"
)

// Spawn-screen skin models per team (class index 0–3).
var (
	EscortSkins = []int{0, 1, 2, 9}
	DefendSkins = []int{9, 28, 47, 57}
)

func clampSkinIndex(idx int) int {
	if idx < 0 {
		return 0
	}
	if idx >= apidef.MaxSkin {
		return apidef.MaxSkin - 1
	}
	return idx
}

func SkinModel(team, skinIndex int) int {
	skinIndex = clampSkinIndex(skinIndex)
	switch team {
	case apidef.TeamEscort:
		return EscortSkins[skinIndex%len(EscortSkins)]
	case apidef.TeamDefend:
		return DefendSkins[skinIndex%len(DefendSkins)]
	default:
		return 0
	}
}

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

func (t *Teams) Assign(api apidef.API, playerID, pack, preferredTeam int) int {
	if s := t.sessions[playerID]; s != nil {
		return s.Team
	}
	if pack < 1 || pack > apidef.MaxPack {
		pack = 1
	}
	team := t.pickTeam(preferredTeam)
	t.sessions[playerID] = apidef.NewPlayerSession(team, pack)
	api.SetPlayerTeam(playerID, team)
	if team == apidef.TeamEscort {
		t.countEscort++
	} else {
		t.countDefend++
	}
	return team
}

func (t *Teams) pickTeam(preferredTeam int) int {
	if preferredTeam == apidef.TeamEscort || preferredTeam == apidef.TeamDefend {
		if t.wouldBalanceAfterJoin(preferredTeam) {
			return preferredTeam
		}
	}
	if t.countEscort > t.countDefend {
		return apidef.TeamDefend
	}
	return apidef.TeamEscort
}

func (t *Teams) wouldBalanceAfterJoin(team int) bool {
	escort, defend := t.countEscort, t.countDefend
	switch team {
	case apidef.TeamEscort:
		escort++
		return escort <= defend
	case apidef.TeamDefend:
		defend++
		return defend <= escort
	}
	return false
}

func (t *Teams) CanSwitchTo(playerID, newTeam int) bool {
	if newTeam != apidef.TeamEscort && newTeam != apidef.TeamDefend {
		return false
	}
	cur := t.Team(playerID)
	if cur == newTeam {
		return false
	}
	escort, defend := t.countEscort, t.countDefend
	switch cur {
	case apidef.TeamEscort:
		escort--
	case apidef.TeamDefend:
		defend--
	}
	switch newTeam {
	case apidef.TeamEscort:
		escort++
		return escort <= defend
	case apidef.TeamDefend:
		defend++
		return defend <= escort
	}
	return false
}

func (t *Teams) SwitchTeam(api apidef.API, playerID, newTeam int) bool {
	if newTeam != apidef.TeamEscort && newTeam != apidef.TeamDefend {
		return false
	}
	s := t.sessions[playerID]
	if s == nil {
		return false
	}
	cur := s.Team
	if cur == newTeam {
		return false
	}
	switch cur {
	case apidef.TeamEscort:
		t.countEscort--
	case apidef.TeamDefend:
		t.countDefend--
	default:
		return false
	}
	switch newTeam {
	case apidef.TeamEscort:
		t.countEscort++
	case apidef.TeamDefend:
		t.countDefend++
	}
	s.Team = newTeam
	s.HasSpawnedThisRound = false
	api.SetPlayerTeam(playerID, newTeam)
	return true
}

func (t *Teams) ForceSwitchTeam(api apidef.API, playerID, newTeam int) bool {
	if newTeam != apidef.TeamEscort && newTeam != apidef.TeamDefend {
		return false
	}
	s := t.sessions[playerID]
	if s == nil {
		return false
	}
	cur := s.Team
	if cur == newTeam {
		return false
	}
	switch cur {
	case apidef.TeamEscort:
		t.countEscort--
	case apidef.TeamDefend:
		t.countDefend--
	default:
		return false
	}
	switch newTeam {
	case apidef.TeamEscort:
		t.countEscort++
	case apidef.TeamDefend:
		t.countDefend++
	}
	s.Team = newTeam
	s.HasSpawnedThisRound = false
	api.SetPlayerTeam(playerID, newTeam)
	return true
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

func (t *Teams) SetSkinIndex(playerID, skinIndex int) bool {
	skinIndex = clampSkinIndex(skinIndex)
	s := t.sessions[playerID]
	if s == nil {
		return false
	}
	s.SkinIndex = skinIndex
	return true
}

func (t *Teams) SkinIndex(playerID int) int {
	if s := t.sessions[playerID]; s != nil {
		return clampSkinIndex(s.SkinIndex)
	}
	return 0
}

func (t *Teams) ApplySkin(api apidef.API, playerID int) {
	team := t.Team(playerID)
	if team == 0 {
		return
	}
	skin := SkinModel(team, t.SkinIndex(playerID))
	_ = api.SetPlayerSkin(playerID, skin)
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
	if classIndex < 0 || classIndex >= apidef.MaxSkin {
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
	api.Send(playerID, colour, fmt.Sprintf(
		"Project Safari: you are on team %s. Press P or /pack 1-%d for loadout; /skin 1-%d or spawn screen for skin; /switch to change team. /help for commands.",
		t.RoleName(team), apidef.MaxPack, apidef.MaxSkin,
	))
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
	for i, sp := range mapCfg.EscortSpawns {
		if i >= apidef.MaxSkin {
			break
		}
		skin := EscortSkins[i%len(EscortSkins)]
		api.AddPlayerClass(apidef.TeamEscort, 0xFF6EC6FF, skin, sp, 0, weapons)
	}
	for i, sp := range mapCfg.DefendSpawns {
		if i >= apidef.MaxSkin {
			break
		}
		skin := DefendSkins[i%len(DefendSkins)]
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
