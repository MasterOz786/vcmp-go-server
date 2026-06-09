package main

type Events struct {
	OnServerStart func() FilterResult
	OnServerStop  func()
	OnServerFrame func(elapsed float32)

	OnIncomingConnection func(name string, password string, ip string) string

	OnPlayerConnect    func(playerID int)
	OnPlayerDisconnect func(playerID int, reason DisconnectReason)
	OnPlayerRequestClass func(playerID int, offset int) FilterResult
	OnPlayerRequestSpawn func(playerID int) FilterResult
	OnPlayerSpawn      func(playerID int)
	OnPlayerDeath      func(playerID, killerID int, reason int, bodyPart int)
	OnPlayerMessage    func(playerID int, message string) FilterResult
	OnPlayerCommand    func(playerID int, command string) FilterResult

	OnPlayerEnterVehicle func(playerID, vehicleID, slot int)
	OnPlayerExitVehicle  func(playerID, vehicleID int)
	OnVehicleExplode     func(vehicleID int)
}

var events Events
