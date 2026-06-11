// Weapon display lines per pack (matches server/safari/loadouts.go).
PackCatalog <- {
	"1": {
		"ESCORT": [
			{ "title": "BREACHER", "weapons": "SHOTGUN\nCOLT .45\nMOLOTOV" },
			{ "title": "SUPPORT", "weapons": "M4\nSHOTGUN\nTEAR GAS" },
			{ "title": "DEMOLITION", "weapons": "RPG\nCOLT .45\nMOLOTOV" },
		],
		"DEFEND": [
			{ "title": "AA GUNNER", "weapons": "MINIGUN\nSHOTGUN" },
			{ "title": "SABOTEUR", "weapons": "RPG\nSHOTGUN\nTEAR GAS" },
			{ "title": "MARKSMAN", "weapons": "SNIPER\nRUGER\nSTUBBY SG" },
		],
	},
};

function packListForTeam(team) {
	if (team == Teams.ESCORT) {
		return PackCatalog["1"].ESCORT;
	}
	return PackCatalog["1"].DEFEND;
}
