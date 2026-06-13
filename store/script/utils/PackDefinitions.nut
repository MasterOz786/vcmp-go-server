// Weapon display lines per pack (matches safari/gameplay/loadouts.go).
PackCatalog <- {
	"1": {
		"ESCORT": [
			{ "title": "BREACHER", "weapons": "SHOTGUN\nCOLT .45\nMOLOTOV", "slot": 1 },
			{ "title": "SUPPORT", "weapons": "M4\nSHOTGUN\nTEAR GAS", "slot": 2 },
			{ "title": "DEMOLITION", "weapons": "RPG\nCOLT .45\nMOLOTOV", "slot": 3 },
		],
		"DEFEND": [
			{ "title": "AA GUNNER", "weapons": "MINIGUN\nSHOTGUN", "slot": 1 },
			{ "title": "SABOTEUR", "weapons": "RPG\nSHOTGUN\nTEAR GAS", "slot": 2 },
			{ "title": "MARKSMAN", "weapons": "SNIPER\nRUGER\nSTUBBY SG", "slot": 3 },
		],
	},
};

function packListForTeam(team) {
	if (team == Teams.ESCORT) {
		return PackCatalog["1"].ESCORT;
	}
	return PackCatalog["1"].DEFEND;
}
