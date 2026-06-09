package main

import (
	"github.com/masteroz/vcmp-go-plugin/vcmp"
	"github.com/masteroz/vcmp-go-server/safari"
)

func (p *Plugin) register() {
	if p.engine == nil {
		return
	}

	vcmp.Events.OnServerStart = func() vcmp.FilterResult {
		p.engine.OnServerStart()
		return vcmp.FilterAllow
	}

	vcmp.Events.OnServerStop = func() {
		p.shutdown()
	}

	vcmp.Events.OnPlayerConnect = func(playerID int) {
		p.engine.Enqueue(safari.NewConnectEvent(playerID))
	}

	vcmp.Events.OnPlayerDisconnect = func(playerID int, _ vcmp.DisconnectReason) {
		p.engine.Enqueue(safari.NewDisconnectEvent(playerID))
	}

	vcmp.Events.OnPlayerRequestSpawn = func(playerID int) vcmp.FilterResult {
		if p.engine.HandleRequestSpawn(playerID) {
			return vcmp.FilterAllow
		}
		return vcmp.FilterDeny
	}

	vcmp.Events.OnPlayerSpawn = func(playerID int) {
		p.engine.Enqueue(safari.NewSpawnEvent(playerID))
	}

	vcmp.Events.OnPlayerCommand = func(playerID int, command string) vcmp.FilterResult {
		if p.engine.HandleCommandSync(playerID, command) {
			return vcmp.FilterDeny
		}
		return vcmp.FilterAllow
	}

	vcmp.Events.OnVehicleExplode = func(vehicleID int) {
		p.engine.Enqueue(safari.NewVehicleExplodeEvent(vehicleID))
	}
}
