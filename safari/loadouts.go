package safari

// GTA VC weapon IDs used for Safari loadouts.
const (
	WeaponColt45  = 22
	WeaponShotgun = 24
	WeaponRPG     = 35
	WeaponTearGas = 17
	WeaponMolotov = 18
	WeaponMinigun = 38
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
		1: {Name: "Escort Shotgun", Weapons: []WeaponGrant{
			{WeaponShotgun, 80},
			{WeaponColt45, 100},
		}},
		2: {Name: "Escort Utility", Weapons: []WeaponGrant{
			{WeaponShotgun, 60},
			{WeaponTearGas, 3},
			{WeaponMolotov, 3},
		}},
	}
}

func DefendPacks() map[int]PackDef {
	return map[int]PackDef{
		1: {Name: "Defender RPG", Weapons: []WeaponGrant{
			{WeaponShotgun, 60},
			{WeaponRPG, 3},
		}},
		2: {Name: "Defender AA", Weapons: []WeaponGrant{
			{WeaponShotgun, 60},
			{WeaponMinigun, 150},
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

func ApplyLoadout(api API, playerID, team, pack int) {
	if pack < 1 || pack > 2 {
		pack = 1
	}
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
	if pack < 1 || pack > 2 {
		pack = 1
	}
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
