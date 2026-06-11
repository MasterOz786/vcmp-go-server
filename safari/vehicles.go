package safari

import "fmt"

const hydraVehicleArchive = "store/vehicles/v6460_t0_p1_Hydra.7z"

func (e *Engine) hydraModel() int {
	if e.cfg.HydraModel > 0 {
		return e.cfg.HydraModel
	}
	return HydraModel
}

func createHydraVehicle(api API, model, world int, pos Vec3, angle float32) int {
	return api.CreateVehicle(model, world, pos, angle, 0, 0)
}

func testHydraSpawnPos(mapCfg MapConfig, playerPos Vec3) Vec3 {
	if mapCfg.HydraStart.X != 0 || mapCfg.HydraStart.Y != 0 {
		return Vec3{
			X: mapCfg.HydraStart.X,
			Y: mapCfg.HydraStart.Y,
			Z: mapCfg.HydraStart.Z + 3,
		}
	}
	z := playerPos.Z + 12
	if z < 14 {
		z = 14.5
	}
	return Vec3{X: playerPos.X + 6, Y: playerPos.Y + 6, Z: z}
}

func formatHydraSpawnFailure(api API, model int) string {
	msg := fmt.Sprintf("Failed to spawn Hydra (model %d).", model)
	if detail := api.LastErrorString(); detail != "" {
		msg += " " + detail
	}
	msg += " Custom vehicle file must be " + hydraVehicleArchive + " (see forum.vc-mp.org topic 975)."
	return msg
}

func warnHydraModelMismatch(api API, playerID, vehicleID, expectedModel int) {
	if vehicleID < 0 {
		return
	}
	actual := api.VehicleModel(vehicleID)
	if actual == expectedModel {
		return
	}
	msg := fmt.Sprintf(
		"Hydra spawned as model %d (expected %d). Install %s and reconnect — otherwise it flies like a Maverick.",
		actual, expectedModel, hydraVehicleArchive,
	)
	api.Log(fmt.Sprintf("[safari] WARNING player %d vehicle %d: %s", playerID, vehicleID, msg))
	if playerID >= 0 && api.IsConnected(playerID) {
		api.Send(playerID, ColourRed, msg)
	}
}
