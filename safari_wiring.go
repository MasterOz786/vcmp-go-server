package main

import (
	"github.com/masteroz/vcmp-go-server/safari"
)

func (p *Plugin) register() {
	if p.engine == nil {
		return
	}

	events.OnServerStart = func() FilterResult {
		p.engine.OnServerStart()
		return FilterAllow
	}

	events.OnServerStop = func() {
		p.shutdown()
	}

	events.OnPlayerConnect = func(playerID int) {
		p.engine.Enqueue(safari.NewConnectEvent(playerID))
	}

	events.OnPlayerDisconnect = func(playerID int, _ DisconnectReason) {
		p.engine.Enqueue(safari.NewDisconnectEvent(playerID))
	}

	events.OnPlayerRequestSpawn = func(playerID int) FilterResult {
		if p.engine.HandleRequestSpawn(playerID) {
			return FilterAllow
		}
		return FilterDeny
	}

	events.OnPlayerSpawn = func(playerID int) {
		p.engine.Enqueue(safari.NewSpawnEvent(playerID))
	}

	events.OnPlayerCommand = func(playerID int, command string) FilterResult {
		if p.engine.HandleCommandSync(playerID, command) {
			return FilterDeny
		}
		return FilterAllow
	}

	events.OnVehicleExplode = func(vehicleID int) {
		p.engine.Enqueue(safari.NewVehicleExplodeEvent(vehicleID))
	}
}
