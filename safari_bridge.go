package main

import "github.com/masteroz/vcmp-go-server/safari"

type safariBridge struct{}

func (safariBridge) Log(msg string) { bridgeLog(msg) }

func (safariBridge) Broadcast(colour uint32, msg string) { bridgeBroadcast(colour, msg) }

func (safariBridge) Send(playerID int, colour uint32, msg string) {
	bridgeSendClientMessage(playerID, colour, msg)
}

func (safariBridge) IsConnected(playerID int) bool { return bridgeIsConnected(playerID) }

func (safariBridge) IsAdmin(playerID int) bool { return bridgeIsAdmin(playerID) }

func (safariBridge) PlayerName(playerID int) string { return bridgePlayerName(playerID) }

func (safariBridge) PlayerIDFromName(name string) int { return bridgePlayerIDFromName(name) }

func (safariBridge) PlayerUID(playerID int) string { return bridgePlayerUID(playerID) }

func (safariBridge) PlayerTeam(playerID int) int { return bridgePlayerTeam(playerID) }

func (safariBridge) SetPlayerTeam(playerID, team int) { bridgeSetPlayerTeam(playerID, team) }

func (safariBridge) SetPlayerScore(playerID, score int) { bridgeSetPlayerScore(playerID, score) }

func (safariBridge) GetPlayerScore(playerID int) int { return bridgeGetPlayerScore(playerID) }

func (safariBridge) SetServerName(name string) { bridgeSetServerName(name) }

func (safariBridge) SetGameModeText(text string) { bridgeSetGameModeText(text) }

func (safariBridge) SetServerOption(option int, on bool) {
	bridgeSetServerOption(ServerOption(option), on)
}

func (safariBridge) SetSpawnPos(pos safari.Vec3) {
	bridgeSetSpawnPos(Vec3{X: pos.X, Y: pos.Y, Z: pos.Z})
}

func (safariBridge) AddPlayerClass(teamID int, colour uint32, model int, pos safari.Vec3, angle float32, weapons [6]int) {
	bridgeAddPlayerClass(teamID, colour, model, Vec3{X: pos.X, Y: pos.Y, Z: pos.Z}, angle, weapons[:])
}

func (safariBridge) CreateVehicle(model, world int, pos safari.Vec3, angle float32, c1, c2 int) int {
	return bridgeCreateVehicle(model, world, Vec3{X: pos.X, Y: pos.Y, Z: pos.Z}, angle, c1, c2)
}

func (safariBridge) DeleteVehicle(vehicleID int) { bridgeDeleteVehicle(vehicleID) }

func (safariBridge) VehiclePos(vehicleID int) safari.Vec3 {
	p := bridgeVehiclePos(vehicleID)
	return safari.Vec3{X: p.X, Y: p.Y, Z: p.Z}
}

func (safariBridge) VehicleHealth(vehicleID int) float32 { return bridgeVehicleHealth(vehicleID) }

func (safariBridge) SetVehiclePosition(vehicleID int, pos safari.Vec3) {
	bridgeSetVehiclePosition(vehicleID, Vec3{X: pos.X, Y: pos.Y, Z: pos.Z})
}

func (safariBridge) SetVehicleHealth(vehicleID int, health float32) {
	bridgeSetVehicleHealth(vehicleID, health)
}

func (safariBridge) RemoveAllWeapons(playerID int) { bridgeRemoveAllWeapons(playerID) }

func (safariBridge) GiveWeapon(playerID, weaponID, ammo int) {
	bridgeGiveWeapon(playerID, weaponID, ammo)
}

func (safariBridge) ServerTimeMs() uint64 { return bridgeServerTimeMs() }
