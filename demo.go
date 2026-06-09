package main

import "fmt"

const (
	spawnX = -974.0
	spawnY = -106.0
	spawnZ = 11.2
)

type Demo struct {
	cfg Config
}

func newDemo(cfg Config) *Demo { return &Demo{cfg: cfg} }

func (d *Demo) register() {
	events.OnServerStart = d.onServerStart
	events.OnServerStop = d.onServerStop
	events.OnIncomingConnection = d.onIncomingConnection
	events.OnPlayerConnect = d.onPlayerConnect
	events.OnPlayerDisconnect = d.onPlayerDisconnect
	events.OnPlayerRequestClass = d.onPlayerRequestClass
	events.OnPlayerRequestSpawn = d.onPlayerRequestSpawn
	events.OnPlayerCommand = d.onPlayerCommand
	events.OnPlayerEnterVehicle = d.onPlayerEnterVehicle
	events.OnPlayerExitVehicle = d.onPlayerExitVehicle
	events.OnVehicleExplode = d.onVehicleExplode
}

func (d *Demo) onServerStart() FilterResult {
	bridgeSetServerName(d.cfg.ServerName)
	bridgeSetGameModeText(d.cfg.GameModeText)
	bridgeSetServerOption(ServerOptionJoinMessages, false)
	bridgeSetServerOption(ServerOptionDeathMessages, false)
	bridgeSetServerOption(ServerOptionUseClasses, true)
	bridgeSetSpawnPos(Vec3{X: spawnX, Y: spawnY, Z: spawnZ})
	bridgeAddPlayerClass(0, 0xFF6EC6FF, 0, Vec3{X: spawnX, Y: spawnY, Z: spawnZ}, 0, []int{0, 0, 0, 0, 0, 0})
	bridgeAddPlayerClass(0, 0xFFFFB86C, 9, Vec3{X: spawnX, Y: spawnY, Z: spawnZ}, 0, []int{0, 0, 0, 0, 0, 0})
	bridgeAddPlayerClass(0, 0xFF98E898, 10, Vec3{X: spawnX, Y: spawnY, Z: spawnZ}, 0, []int{0, 0, 0, 0, 0, 0})
	bridgeLog(fmt.Sprintf("[%s] demo server started (Kotlin-style commands enabled)", PluginName))
	return FilterAllow
}

func (d *Demo) onServerStop() {
	bridgeLog(fmt.Sprintf("[%s] demo server stopped", PluginName))
}

func (d *Demo) onIncomingConnection(name, password, ip string) string {
	_ = password
	_ = ip
	return name + "!"
}

func (d *Demo) onPlayerConnect(playerID int) {
	bridgeBroadcast(ColourYellowish, fmt.Sprintf("Player %s joined.", bridgePlayerName(playerID)))
}

func (d *Demo) onPlayerDisconnect(playerID int, reason DisconnectReason) {
	_ = reason
	clearPlayerPing(playerID)
	bridgeBroadcast(ColourYellowish, fmt.Sprintf("Player %s disconnected.", bridgePlayerName(playerID)))
}

func (d *Demo) onPlayerRequestClass(playerID int, offset int) FilterResult {
	_ = playerID
	_ = offset
	return FilterAllow
}

func (d *Demo) onPlayerRequestSpawn(playerID int) FilterResult {
	_ = playerID
	return FilterAllow
}

func (d *Demo) onPlayerCommand(playerID int, command string) FilterResult {
	return handleDemoCommand(playerID, command)
}

func (d *Demo) onPlayerEnterVehicle(playerID, vehicleID, slot int) {
	bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Player %s entered vehicle %d at slot %d.", bridgePlayerName(playerID), vehicleID, slot))
}

func (d *Demo) onPlayerExitVehicle(playerID, vehicleID int) {
	bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Player %s exited vehicle %d.", bridgePlayerName(playerID), vehicleID))
}

func (d *Demo) onVehicleExplode(vehicleID int) {
	bridgeBroadcast(ColourYellowish, fmt.Sprintf("Vehicle %d exploded.", vehicleID))
}
