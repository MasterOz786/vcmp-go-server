package safari

// API is the gamemode-facing surface. Production uses VCMPAPI (vcmp-go-plugin).
type API interface {
	Log(msg string)
	Broadcast(colour uint32, msg string)
	Send(playerID int, colour uint32, msg string)

	IsConnected(playerID int) bool
	IsSpawned(playerID int) bool
	IsAdmin(playerID int) bool
	PlayerName(playerID int) string
	PlayerIDFromName(name string) int
	PlayerUID(playerID int) string
	PlayerTeam(playerID int) int
	SetPlayerTeam(playerID, team int)
	SetPlayerScore(playerID, score int)
	GetPlayerScore(playerID int) int
	SetPlayerPosition(playerID int, pos Vec3) error

	SetServerName(name string)
	SetGameModeText(text string)
	SetServerOption(option int, on bool)
	SetSpawnPos(pos Vec3)
	AddPlayerClass(teamID int, colour uint32, model int, pos Vec3, angle float32, weapons [6]int)

	CreateVehicle(model, world int, pos Vec3, angle float32, c1, c2 int) int
	DeleteVehicle(vehicleID int)
	VehiclePos(vehicleID int) Vec3
	VehicleHealth(vehicleID int) float32
	SetVehiclePosition(vehicleID int, pos Vec3)
	SetVehicleHealth(vehicleID int, health float32)

	RemoveAllWeapons(playerID int)
	GiveWeapon(playerID, weaponID, ammo int)
	WeaponAtSlot(playerID, slot int) int
	RemoveWeapon(playerID, weaponID int) error

	ServerTimeMs() uint64

	SendScriptData(playerID int, data []byte) error
}
