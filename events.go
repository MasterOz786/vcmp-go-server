package main

type Events struct {
	OnServerStart func() FilterResult
	OnServerStop  func()
	OnServerFrame func(elapsed float32)

	OnPluginCommand      func(commandID uint32, message string) FilterResult
	OnIncomingConnection func(name string, password string, ip string) string
	OnClientScriptData   func(playerID int, data []byte)

	OnPlayerConnect    func(playerID int)
	OnPlayerDisconnect func(playerID int, reason DisconnectReason)
	OnPlayerRequestClass func(playerID int, offset int) FilterResult
	OnPlayerRequestSpawn func(playerID int) FilterResult
	OnPlayerSpawn      func(playerID int)
	OnPlayerDeath      func(playerID, killerID int, reason int, bodyPart BodyPart)
	OnPlayerUpdate     func(playerID int, updateType PlayerUpdate)

	OnPlayerRequestEnterVehicle func(playerID, vehicleID, slot int) FilterResult
	OnPlayerEnterVehicle        func(playerID, vehicleID, slot int)
	OnPlayerExitVehicle         func(playerID, vehicleID int)

	OnPlayerNameChange     func(playerID int, oldName, newName string)
	OnPlayerStateChange    func(playerID int, oldState, newState PlayerState)
	OnPlayerActionChange   func(playerID int, oldAction, newAction int)
	OnPlayerOnFireChange   func(playerID int, isOnFire bool)
	OnPlayerCrouchChange   func(playerID int, isCrouching bool)
	OnPlayerGameKeysChange func(playerID int, oldKeys, newKeys uint32)
	OnPlayerBeginTyping    func(playerID int)
	OnPlayerEndTyping      func(playerID int)
	OnPlayerAwayChange     func(playerID int, isAway bool)

	OnPlayerMessage        func(playerID int, message string) FilterResult
	OnPlayerCommand        func(playerID int, command string) FilterResult
	OnPlayerPrivateMessage func(playerID, targetPlayerID int, message string) FilterResult

	OnPlayerKeyBindDown  func(playerID, bindID int)
	OnPlayerKeyBindUp    func(playerID, bindID int)
	OnPlayerSpectate     func(playerID, targetPlayerID int)
	OnPlayerCrashReport  func(playerID int, report string)
	OnPlayerModuleList   func(playerID int, list string)

	OnVehicleUpdate  func(vehicleID int, updateType VehicleUpdate)
	OnVehicleExplode func(vehicleID int)
	OnVehicleRespawn func(vehicleID int)

	OnObjectShot    func(objectID, playerID, weaponID int)
	OnObjectTouched func(objectID, playerID int)

	OnPickupPickAttempt func(pickupID, playerID int) FilterResult
	OnPickupPicked      func(pickupID, playerID int)
	OnPickupRespawn     func(pickupID int)

	OnCheckpointEntered func(checkpointID, playerID int)
	OnCheckpointExited  func(checkpointID, playerID int)

	OnEntityPoolChange func(entityType EntityPool, entityID int, isDeleted bool)

	OnServerPerformanceReport func(descriptions []string, times []uint64)
}

var events Events
