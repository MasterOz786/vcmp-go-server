package safari

import "github.com/masteroz/vcmp-go-plugin/vcmp"

// VCMPAPI implements API using github.com/masteroz/vcmp-go-plugin/vcmp.
type VCMPAPI struct{}

func toSafariVec3(p vcmp.Vec3) Vec3 {
	return Vec3{X: p.X, Y: p.Y, Z: p.Z}
}

func fromSafariVec3(p Vec3) vcmp.Vec3 {
	return vcmp.Vec3{X: p.X, Y: p.Y, Z: p.Z}
}

func (VCMPAPI) Log(msg string) { vcmp.API.Server.Log(msg) }

func (VCMPAPI) Broadcast(colour uint32, msg string) { vcmp.API.Server.Broadcast(colour, msg) }

func (VCMPAPI) Send(playerID int, colour uint32, msg string) {
	vcmp.API.Player.SendMessage(playerID, colour, msg)
}

func (VCMPAPI) IsConnected(playerID int) bool { return vcmp.API.Player.IsConnected(playerID) }

func (VCMPAPI) IsSpawned(playerID int) bool { return vcmp.API.Player.IsSpawned(playerID) }

func (VCMPAPI) IsAdmin(playerID int) bool { return vcmp.API.Player.IsAdmin(playerID) }

func (VCMPAPI) PlayerName(playerID int) string { return vcmp.API.Player.Name(playerID) }

func (VCMPAPI) PlayerIDFromName(name string) int { return vcmp.API.Player.IDFromName(name) }

func (VCMPAPI) PlayerUID(playerID int) string { return vcmp.API.Player.UID(playerID) }

func (VCMPAPI) PlayerTeam(playerID int) int { return vcmp.API.Player.Team(playerID) }

func (VCMPAPI) SetPlayerTeam(playerID, team int) { vcmp.API.Player.SetTeam(playerID, team) }

func (VCMPAPI) SetPlayerScore(playerID, score int) { vcmp.API.Player.SetScore(playerID, score) }

func (VCMPAPI) GetPlayerScore(playerID int) int { return vcmp.API.Player.Score(playerID) }

func (VCMPAPI) SetPlayerPosition(playerID int, pos Vec3) error {
	return vcmp.API.Player.SetPosition(playerID, fromSafariVec3(pos))
}

func (VCMPAPI) SetServerName(name string) { vcmp.API.Server.SetName(name) }

func (VCMPAPI) SetGameModeText(text string) { vcmp.API.Server.SetGameModeText(text) }

func (VCMPAPI) SetServerOption(option int, on bool) {
	vcmp.API.Server.SetOption(vcmp.ServerOption(option), on)
}

func (VCMPAPI) SetSpawnPos(pos Vec3) {
	vcmp.API.Server.SetSpawnPosition(fromSafariVec3(pos))
}

func (VCMPAPI) AddPlayerClass(teamID int, colour uint32, model int, pos Vec3, angle float32, weapons [6]int) {
	vcmp.API.Server.AddPlayerClass(teamID, colour, model, fromSafariVec3(pos), angle, weapons[:]...)
}

func (VCMPAPI) CreateVehicle(model, world int, pos Vec3, angle float32, c1, c2 int) int {
	return vcmp.API.Vehicle.Create(model, world, fromSafariVec3(pos), angle, c1, c2)
}

func (VCMPAPI) DeleteVehicle(vehicleID int) { vcmp.API.Vehicle.Delete(vehicleID) }

func (VCMPAPI) VehiclePos(vehicleID int) Vec3 {
	return toSafariVec3(vcmp.API.Vehicle.Position(vehicleID))
}

func (VCMPAPI) VehicleHealth(vehicleID int) float32 { return vcmp.API.Vehicle.Health(vehicleID) }

func (VCMPAPI) SetVehiclePosition(vehicleID int, pos Vec3) {
	vcmp.API.Vehicle.SetPosition(vehicleID, fromSafariVec3(pos), false)
}

func (VCMPAPI) SetVehicleHealth(vehicleID int, health float32) {
	vcmp.API.Vehicle.SetHealth(vehicleID, health)
}

func (VCMPAPI) RemoveAllWeapons(playerID int) { vcmp.API.Player.RemoveAllWeapons(playerID) }

func (VCMPAPI) GiveWeapon(playerID, weaponID, ammo int) {
	vcmp.API.Player.GiveWeapon(playerID, weaponID, ammo)
}

func (VCMPAPI) WeaponAtSlot(playerID, slot int) int {
	return vcmp.API.Player.WeaponAtSlot(playerID, slot)
}

func (VCMPAPI) RemoveWeapon(playerID, weaponID int) error {
	return vcmp.API.Player.RemoveWeapon(playerID, weaponID)
}

func (VCMPAPI) ServerTimeMs() uint64 { return vcmp.API.Server.Time() }

func (VCMPAPI) SendScriptData(playerID int, data []byte) error {
	return vcmp.API.Player.SendScriptData(playerID, data)
}
