package safari

import "fmt"

type Teams struct {
	countEscort int
	countDefend int
	assignments map[int]int
	packs       map[int]int
}

func NewTeams() *Teams {
	return &Teams{
		assignments: make(map[int]int),
		packs:       make(map[int]int),
	}
}

func (t *Teams) Assign(api API, playerID int) int {
	if team, ok := t.assignments[playerID]; ok {
		return team
	}
	team := TeamEscort
	if t.countEscort > t.countDefend {
		team = TeamDefend
	}
	t.assignments[playerID] = team
	t.packs[playerID] = 1
	api.SetPlayerTeam(playerID, team)
	if team == TeamEscort {
		t.countEscort++
	} else {
		t.countDefend++
	}
	return team
}

func (t *Teams) Remove(playerID int) {
	team, ok := t.assignments[playerID]
	if !ok {
		return
	}
	if team == TeamEscort {
		t.countEscort--
	} else {
		t.countDefend--
	}
	delete(t.assignments, playerID)
	delete(t.packs, playerID)
}

func (t *Teams) Team(playerID int) int {
	return t.assignments[playerID]
}

func (t *Teams) SetPack(playerID, pack int) bool {
	if pack < 1 || pack > 2 {
		return false
	}
	if _, ok := t.assignments[playerID]; !ok {
		return false
	}
	t.packs[playerID] = pack
	return true
}

func (t *Teams) Pack(playerID int) int {
	if p, ok := t.packs[playerID]; ok {
		return p
	}
	return 1
}

func (t *Teams) RoleName(team int) string {
	if team == TeamEscort {
		return "Escort"
	}
	return "Defender"
}

func (t *Teams) Welcome(api API, playerID int) {
	team := t.Assign(api, playerID)
	colour := ColourGreen
	if team == TeamDefend {
		colour = ColourRed
	}
	api.Send(playerID, colour, fmt.Sprintf("Project Safari: you are on team %s. Use /pack 1|2 and /help.", t.RoleName(team)))
}

func (t *Teams) SetupClasses(api API, mapCfg MapConfig) {
	api.SetServerOption(int(ServerOptionUseClasses), true)
	api.SetServerOption(int(ServerOptionJoinMessages), false)
	api.SetServerOption(int(ServerOptionDeathMessages), false)

	weapons := [6]int{WeaponShotgun, 50, 0, 0, 0, 0}
	for i, sp := range mapCfg.EscortSpawns {
		if i >= 4 {
			break
		}
		api.AddPlayerClass(TeamEscort, 0xFF6EC6FF, 0, sp, 0, weapons)
	}
	for i, sp := range mapCfg.DefendSpawns {
		if i >= 4 {
			break
		}
		api.AddPlayerClass(TeamDefend, 0xFFFF6E6E, 9, sp, 0, weapons)
	}
	if len(mapCfg.EscortSpawns) > 0 {
		api.SetSpawnPos(mapCfg.EscortSpawns[0])
	}
}

// ServerOption constants mirrored for API (vcmpServerOptionUseClasses = 18).
const ServerOptionUseClasses = 18
const ServerOptionJoinMessages = 15
const ServerOptionDeathMessages = 16
