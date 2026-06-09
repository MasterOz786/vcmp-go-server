package safari

import (
	"fmt"
	"math"
)

const (
	hydraStepMeters     = 5.0
	checkpointRadiusM   = 15.0
)

type Hydra struct {
	VehicleID   int
	Waypoints   []Vec3
	Index       int
	lastMinute  uint64
}

func NewHydra() *Hydra {
	return &Hydra{Index: 0}
}

func (h *Hydra) Spawn(api API, mapCfg MapConfig) int {
	pos := mapCfg.HydraStart
	if len(mapCfg.Waypoints) > 0 && mapCfg.HydraStart.X == 0 && mapCfg.HydraStart.Y == 0 {
		pos = mapCfg.Waypoints[0]
	}
	h.Waypoints = mapCfg.Waypoints
	h.Index = 0
	h.VehicleID = api.CreateVehicle(HydraModel, mapCfg.World, pos, mapCfg.HydraAngle, 1, 1)
	if h.VehicleID >= 0 {
		api.SetVehicleHealth(h.VehicleID, HydraMaxHP)
	}
	return h.VehicleID
}

func (h *Hydra) Destroy(api API) {
	if h.VehicleID >= 0 {
		api.DeleteVehicle(h.VehicleID)
		h.VehicleID = -1
	}
}

func (h *Hydra) Health(api API) float32 {
	if h.VehicleID < 0 {
		return 0
	}
	return api.VehicleHealth(h.VehicleID)
}

func (h *Hydra) Tick(api API, score *Scoring, nowMs uint64) (checkpoint bool, escortWin bool, defendWin bool, msg string) {
	if h.VehicleID < 0 {
		return false, false, false, ""
	}
	hp := api.VehicleHealth(h.VehicleID)
	if hp <= 0 {
		score.AddDefend(PointsHydraDestroy)
		return false, false, true, "Hydra destroyed! Defenders win."
	}

	if len(h.Waypoints) == 0 {
		return false, false, false, ""
	}

	if h.Index >= len(h.Waypoints) {
		return false, true, false, "Hydra reached its objective! Escort wins."
	}

	target := h.Waypoints[h.Index]
	pos := api.VehiclePos(h.VehicleID)
	if dist(pos, target) <= checkpointRadiusM {
		h.Index++
		score.AddEscort(PointsCheckpoint)
		if h.Index >= len(h.Waypoints) {
			return true, true, false, fmt.Sprintf("Final checkpoint reached! Escort +%d", PointsCheckpoint)
		}
		return true, false, false, fmt.Sprintf("Checkpoint %d/%d reached! Escort +%d", h.Index, len(h.Waypoints), PointsCheckpoint)
	}

	h.moveToward(api, pos, target)

	if h.lastMinute == 0 {
		h.lastMinute = nowMs
	} else if nowMs-h.lastMinute >= 60000 {
		h.lastMinute = nowMs
		score.AddEscort(PointsHydraMinute)
		msg = fmt.Sprintf("Hydra holding route (+%d escort)", PointsHydraMinute)
	}
	return false, false, false, msg
}

func (h *Hydra) moveToward(api API, pos, target Vec3) {
	dx := target.X - pos.X
	dy := target.Y - pos.Y
	dz := target.Z - pos.Z
	d := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
	if d < 0.01 {
		return
	}
	step := float32(hydraStepMeters)
	if step > d {
		step = d
	}
	ratio := step / d
	newPos := Vec3{
		X: pos.X + dx*ratio,
		Y: pos.Y + dy*ratio,
		Z: pos.Z + dz*ratio,
	}
	api.SetVehiclePosition(h.VehicleID, newPos)
}

func dist(a, b Vec3) float32 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	dz := b.Z - a.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}
