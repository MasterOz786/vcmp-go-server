package safari

// VC:MP / GTA Vice City weapon IDs (NOT GTA SA). See wiki.vc-mp.org/wiki/Weapons
const (
	WeaponTearGas        = 14
	WeaponMolotov        = 15
	WeaponColt45         = 17
	WeaponPython         = 18
	WeaponShotgun        = 19
	WeaponStubbyShotgun  = 21
	WeaponTec9           = 22
	WeaponM4             = 26
	WeaponRuger          = 27 // VC has no silenced pistol; marksman secondary
	WeaponSniper         = 28
	WeaponRPG            = 30
	WeaponMinigun        = 33
)

type WeaponGrant struct {
	ID   int
	Ammo int
}

type PackDef struct {
	Name    string
	Weapons []WeaponGrant
}

func EscortPacks() map[int]PackDef {
	return map[int]PackDef{
		1: {Name: "Escort Breacher", Weapons: []WeaponGrant{
			{WeaponShotgun, 80},
			{WeaponColt45, 120},
			{WeaponMolotov, 4},
		}},
		2: {Name: "Escort Support", Weapons: []WeaponGrant{
			{WeaponM4, 200},
			{WeaponShotgun, 60},
			{WeaponTearGas, 4},
		}},
		3: {Name: "Escort Demolition", Weapons: []WeaponGrant{
			{WeaponRPG, 4},
			{WeaponColt45, 100},
			{WeaponMolotov, 2},
		}},
	}
}

func DefendPacks() map[int]PackDef {
	return map[int]PackDef{
		1: {Name: "Defender AA Gunner", Weapons: []WeaponGrant{
			{WeaponMinigun, 200},
			{WeaponShotgun, 60},
		}},
		2: {Name: "Defender Saboteur", Weapons: []WeaponGrant{
			{WeaponRPG, 4},
			{WeaponShotgun, 50},
			{WeaponTearGas, 4},
		}},
		3: {Name: "Defender Marksman", Weapons: []WeaponGrant{
			{WeaponSniper, 40},
			{WeaponRuger, 80},
			{WeaponStubbyShotgun, 50},
		}},
	}
}

func AllowedWeaponIDs() map[int]bool {
	ids := map[int]bool{}
	for _, p := range EscortPacks() {
		for _, w := range p.Weapons {
			ids[w.ID] = true
		}
	}
	for _, p := range DefendPacks() {
		for _, w := range p.Weapons {
			ids[w.ID] = true
		}
	}
	return ids
}

func clampPack(pack int) int {
	if pack < 1 {
		return 1
	}
	if pack > MaxPack {
		return MaxPack
	}
	return pack
}

func LoadoutComplete(api API, playerID, team, pack int) bool {
	pack = clampPack(pack)
	var packs map[int]PackDef
	if team == TeamEscort {
		packs = EscortPacks()
	} else {
		packs = DefendPacks()
	}
	def, ok := packs[pack]
	if !ok {
		def = packs[1]
	}
	for _, w := range def.Weapons {
		found := false
		for slot := 0; slot <= 12; slot++ {
			if api.WeaponAtSlot(playerID, slot) == w.ID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func ApplyLoadout(api API, playerID, team, pack int) {
	pack = clampPack(pack)
	var packs map[int]PackDef
	if team == TeamEscort {
		packs = EscortPacks()
	} else {
		packs = DefendPacks()
	}
	def, ok := packs[pack]
	if !ok {
		def = packs[1]
	}
	api.RemoveAllWeapons(playerID)
	for _, w := range def.Weapons {
		api.GiveWeapon(playerID, w.ID, w.Ammo)
	}
}

func allowedWeaponIDsForPack(team, pack int) map[int]bool {
	pack = clampPack(pack)
	var packs map[int]PackDef
	if team == TeamEscort {
		packs = EscortPacks()
	} else {
		packs = DefendPacks()
	}
	def, ok := packs[pack]
	if !ok {
		def = packs[1]
	}
	ids := make(map[int]bool)
	for _, w := range def.Weapons {
		ids[w.ID] = true
	}
	return ids
}

// EnforceAllowed removes weapons not permitted for the player's pack.
func EnforceAllowed(api API, playerID, team, pack int) {
	allowed := allowedWeaponIDsForPack(team, pack)
	for slot := 0; slot <= 12; slot++ {
		wid := api.WeaponAtSlot(playerID, slot)
		if wid <= 0 {
			continue
		}
		if !allowed[wid] {
			_ = api.RemoveWeapon(playerID, wid)
		}
	}
}
