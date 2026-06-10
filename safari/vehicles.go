package safari

import "fmt"

// DefaultHydraModelFallbacks: Hydra (custom), then vanilla VC helicopters.
var DefaultHydraModelFallbacks = []int{520, 487, 422, 446}

func mergeHydraModels(primary int, fallbacks []int) []int {
	if primary <= 0 {
		primary = HydraModel
	}
	seen := map[int]bool{primary: true}
	out := []int{primary}
	for _, m := range fallbacks {
		if m <= 0 || seen[m] {
			continue
		}
		seen[m] = true
		out = append(out, m)
	}
	for _, m := range DefaultHydraModelFallbacks {
		if m <= 0 || seen[m] {
			continue
		}
		seen[m] = true
		out = append(out, m)
	}
	return out
}

func (e *Engine) hydraVehicleModels() []int {
	return mergeHydraModels(e.cfg.HydraModel, e.cfg.HydraModelFallbacks)
}

func createHydraVehicle(api API, models []int, world int, pos Vec3, angle float32) (vehicleID, modelUsed int) {
	for _, model := range models {
		vid := api.CreateVehicle(model, world, pos, angle, 1, 1)
		if vid >= 0 {
			return vid, model
		}
	}
	return -1, 0
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

func formatHydraSpawnFailure(api API, models []int) string {
	msg := fmt.Sprintf("Failed to spawn aircraft (tried models %v).", models)
	if detail := api.LastErrorString(); detail != "" {
		msg += " " + detail
	}
	msg += " Hydra (520) needs custom vehicle files on the server; Maverick (487) is used as fallback when present."
	return msg
}
