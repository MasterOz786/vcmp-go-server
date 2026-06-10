package safari

import (
	"fmt"
	"strings"
)

func (e *Engine) cmdTestHydra(playerID int, args []string) CommandResult {
	e.api.Send(playerID, ColourWhite, "Test Hydra: working...")

	if len(args) > 0 && strings.EqualFold(args[0], "stop") {
		e.stopTestHydra(playerID)
		e.api.Send(playerID, ColourGreen, "Test Hydra removed.")
		return CommandResult{Handled: true, Deny: true}
	}

	if !e.api.IsSpawned(playerID) {
		if err := e.api.ForceSpawn(playerID); err != nil {
			e.api.Send(playerID, ColourRed, "Spawn first, then run /testhydra.")
			return CommandResult{Handled: true, Deny: true}
		}
	}

	e.stopTestHydra(playerID)

	models := e.hydraVehicleModels()
	pos := testHydraSpawnPos(e.mapCfg, e.api.PlayerPosition(playerID))
	world := e.mapCfg.World
	if e.api.IsConnected(playerID) {
		world = e.api.PlayerWorld(playerID)
	}

	vid, modelUsed := createHydraVehicle(e.api, models, world, pos, e.mapCfg.HydraAngle)
	if vid < 0 {
		e.api.Send(playerID, ColourRed, formatHydraSpawnFailure(e.api, models))
		return CommandResult{Handled: true, Deny: true}
	}
	e.api.SetVehicleHealth(vid, HydraMaxHP)
	if world != e.mapCfg.World {
		_ = e.api.SetVehicleWorld(vid, world)
	}

	s := e.teams.session(playerID)
	if s == nil {
		e.ensurePlayerSession(playerID)
		s = e.teams.session(playerID)
	}
	if s != nil {
		s.TestHydraVehicleID = vid
		s.HydraCameraMode = HydraCamDefault
	}

	e.api.PutPlayerInVehicle(playerID, vid, 0)
	e.resetHydraCamera(playerID)

	label := "Hydra"
	if modelUsed != HydraModel {
		label = fmt.Sprintf("aircraft (model %d)", modelUsed)
	}
	e.api.Send(playerID, ColourGreen, fmt.Sprintf("Test %s ready — you are in the pilot seat.", label))
	e.api.Send(playerID, ColourCyan, "Press V or /hydraview to cycle camera views.")
	if modelUsed != HydraModel {
		e.api.Send(playerID, ColourYellow, "Using fallback helicopter — add Hydra (520) custom vehicle pack for the real Hydra.")
	}
	e.api.Log(fmt.Sprintf("[safari] test hydra spawned for player %d vehicle=%d model=%d", playerID, vid, modelUsed))
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) stopTestHydra(playerID int) {
	s := e.teams.session(playerID)
	if s == nil || s.TestHydraVehicleID < 0 {
		return
	}
	if e.api.PlayerVehicleID(playerID) == s.TestHydraVehicleID {
		_ = e.api.RemoveFromVehicle(playerID)
	}
	e.api.DeleteVehicle(s.TestHydraVehicleID)
	s.TestHydraVehicleID = -1
	e.resetHydraCamera(playerID)
}
