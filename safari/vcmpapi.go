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

func (VCMPAPI) SetAdmin(playerID int, admin bool) { vcmp.API.Player.SetAdmin(playerID, admin) }

func (VCMPAPI) Kick(playerID int) error { return vcmp.API.Player.Kick(playerID) }

func (VCMPAPI) Shutdown() { vcmp.API.Server.Shutdown() }

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

func (VCMPAPI) PlayerPosition(playerID int) Vec3 {
	return toSafariVec3(vcmp.API.Player.Position(playerID))
}

func (VCMPAPI) PlayerVehicleID(playerID int) int {
	return vcmp.API.Player.VehicleID(playerID)
}

func (VCMPAPI) ForceSpawn(playerID int) error {
	return vcmp.API.Player.ForceSpawn(playerID)
}

func (VCMPAPI) PutPlayerInVehicle(playerID, vehicleID, slot int) {
	vcmp.API.Player.PutInVehicle(playerID, vehicleID, slot, true, true)
}

func (VCMPAPI) RemoveFromVehicle(playerID int) error {
	return vcmp.API.Player.RemoveFromVehicle(playerID)
}

func (VCMPAPI) SetCamera(playerID int, pos, lookAt Vec3) error {
	return vcmp.API.Player.SetCamera(playerID, fromSafariVec3(pos), fromSafariVec3(lookAt))
}

func (VCMPAPI) RestoreCamera(playerID int) error {
	return vcmp.API.Player.RestoreCamera(playerID)
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

func (VCMPAPI) VehicleExists(vehicleID int) bool {
	return vcmp.API.Vehicle.Exists(vehicleID)
}

func (VCMPAPI) VehiclePos(vehicleID int) Vec3 {
	return toSafariVec3(vcmp.API.Vehicle.Position(vehicleID))
}

func (VCMPAPI) VehicleModel(vehicleID int) int {
	return vcmp.API.Vehicle.Model(vehicleID)
}

func (VCMPAPI) VehicleRotationEuler(vehicleID int) Vec3 {
	rot, err := vcmp.API.Vehicle.RotationEuler(vehicleID)
	if err != nil {
		return Vec3{}
	}
	return toSafariVec3(rot)
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

func (VCMPAPI) SetWeapon(playerID, weaponID, ammo int) error {
	return vcmp.API.Player.SetWeapon(playerID, weaponID, ammo)
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

func (VCMPAPI) LastErrorString() string {
	if err := vcmp.API.Server.LastError(); err != nil {
		return err.Error()
	}
	return ""
}

func (VCMPAPI) PlayerWorld(playerID int) int {
	return vcmp.API.Player.World(playerID)
}

func (VCMPAPI) SetVehicleWorld(vehicleID, world int) error {
	return vcmp.API.Vehicle.SetWorld(vehicleID, world)
}
