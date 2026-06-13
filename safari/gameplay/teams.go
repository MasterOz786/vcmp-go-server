package gameplay

import (
	"fmt"

	"github.com/masteroz/vcmp-go-server/safari/apidef"
)

// Spawn-screen skin models per team (class index 0–1). Taxi driver skins.
var (
	EscortSkins = []int{133, 134}
	DefendSkins = []int{133, 134}
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

func (t *Teams) AllowClassRequest(playerID, classOffset int, roundActive bool) bool {
	// VC:MP passes a navigation offset on the spawn screen — never block browsing.
	_ = playerID
	_ = classOffset
	_ = roundActive
	return true
}

func (t *Teams) AllowSpawnRequest(playerID int, roundActive bool) bool {
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

func TeamAndSkinFromClass(classID int) (team, skin int) {
	if classID < 0 {
		return 0, 0
	}
	if classID >= apidef.MaxSkin {
		return apidef.TeamDefend, clampSkinIndex(classID - apidef.MaxSkin)
	}
	return apidef.TeamEscort, clampSkinIndex(classID)
}

func (t *Teams) SyncFromSpawnScreen(api apidef.API, playerID int, classHint int) {
	s := t.sessions[playerID]
	if s == nil {
		return
	}

	classID := classHint
	if classHint == -1 || classHint == 1 {
		classID = s.SkinIndex
		if s.Team == apidef.TeamDefend {
			classID += apidef.MaxSkin
		}
		classID += classHint
		if classID < 0 {
			classID = apidef.MaxSkin*2 - 1
		}
		if classID >= apidef.MaxSkin*2 {
			classID = 0
		}
	} else if classHint < 0 {
		classID = api.PlayerClass(playerID)
	}

	team := api.PlayerTeam(playerID)
	skin := clampSkinIndex(classID)

	if classHint == -1 || classHint == 1 || (classHint >= 0 && classHint < apidef.MaxSkin*2) {
		team, skin = TeamAndSkinFromClass(classID)
	} else if team != apidef.TeamEscort && team != apidef.TeamDefend {
		team, skin = TeamAndSkinFromClass(classID)
	} else if classID >= apidef.MaxSkin {
		skin = clampSkinIndex(classID - apidef.MaxSkin)
	}

	if team != apidef.TeamEscort && team != apidef.TeamDefend {
		return
	}

	if s.Team != team {
		switch s.Team {
		case apidef.TeamEscort:
			t.countEscort--
		case apidef.TeamDefend:
			t.countDefend--
		}
		switch team {
		case apidef.TeamEscort:
			t.countEscort++
		case apidef.TeamDefend:
			t.countDefend++
		}
		s.Team = team
		api.SetPlayerTeam(playerID, team)
	}
	s.SkinIndex = skin
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

func SpawnScreenCamera(lobby apidef.Vec3) (camPos, lookAt apidef.Vec3) {
	return apidef.Vec3{
			X: lobby.X - 5,
			Y: lobby.Y + 4,
			Z: lobby.Z + 1,
		}, apidef.Vec3{
			X: lobby.X + 10,
			Y: lobby.Y - 12,
			Z: lobby.Z + 2.5,
		}
}

func (t *Teams) SetupClasses(api apidef.API, lobby apidef.Vec3, angle float32) {
	// Pack weapons are granted by the gamemode; class kits must stay empty.
	// All class previews use the lobby position so skins render on the spawn screen.
	weapons := [6]int{0, 0, 0, 0, 0, 0}
	if lobby.X == 0 && lobby.Y == 0 && lobby.Z == 0 {
		api.Log("[safari] WARNING: lobby spawn unset — spawn-screen skins may not appear")
	}
	for i := 0; i < apidef.MaxSkin; i++ {
		skin := EscortSkins[i%len(EscortSkins)]
		api.AddPlayerClass(apidef.TeamEscort, 0xFFFFE448, skin, lobby, angle, weapons)
	}
	for i := 0; i < apidef.MaxSkin; i++ {
		skin := DefendSkins[i%len(DefendSkins)]
		api.AddPlayerClass(apidef.TeamDefend, 0xFFFF78AF, skin, lobby, angle, weapons)
	}
	api.SetSpawnPos(lobby)
	camPos, lookAt := SpawnScreenCamera(lobby)
	api.SetSpawnCameraPosition(camPos)
	api.SetSpawnCameraLookAt(lookAt)
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
