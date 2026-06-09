package safari

// API is implemented by the main plugin bridge (VC:MP natives).
type API interface {
	Log(msg string)
	Broadcast(colour uint32, msg string)
	Send(playerID int, colour uint32, msg string)

	IsConnected(playerID int) bool
	IsAdmin(playerID int) bool
	PlayerName(playerID int) string
	PlayerIDFromName(name string) int
	PlayerUID(playerID int) string
	PlayerTeam(playerID int) int
	SetPlayerTeam(playerID, team int)
	SetPlayerScore(playerID, score int)
	GetPlayerScore(playerID int) int

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

	ServerTimeMs() uint64
}
