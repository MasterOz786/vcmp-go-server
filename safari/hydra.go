package safari

import (
	"fmt"
	"math"
)

const (
	hydraStepMeters   = 5.0
	checkpointRadiusM = 15.0
)

type HydraState int

const (
	HydraIdle HydraState = iota
	HydraPatrol
	HydraDestroyed
	HydraObjectiveReached
)

type Hydra struct {
	State      HydraState
	VehicleID  int
	Waypoints  []Vec3
	Index      int
	lastMinute uint64
}

func NewHydra() *Hydra {
	return &Hydra{State: HydraIdle, Index: 0, VehicleID: -1}
}

func (h *Hydra) Spawn(api API, mapCfg MapConfig, model int) int {
	pos := mapCfg.HydraStart
	if len(mapCfg.Waypoints) > 0 && mapCfg.HydraStart.X == 0 && mapCfg.HydraStart.Y == 0 {
		pos = mapCfg.Waypoints[0]
	}
	h.Waypoints = mapCfg.Waypoints
	h.Index = 0
	h.VehicleID = createHydraVehicle(api, model, mapCfg.World, pos, mapCfg.HydraAngle)
	if h.VehicleID >= 0 {
		api.SetVehicleHealth(h.VehicleID, HydraMaxHP)
		h.State = HydraPatrol
	} else {
		h.State = HydraIdle
		api.Log(formatHydraSpawnFailure(api, model))
	}
	return h.VehicleID
}

func (h *Hydra) Destroy(api API) {
	if h.VehicleID >= 0 {
		api.DeleteVehicle(h.VehicleID)
		h.VehicleID = -1
	}
	if h.State == HydraPatrol {
		h.State = HydraIdle
	}
}

func (h *Hydra) Health(api API) float32 {
	if h.VehicleID < 0 {
		return 0
	}
	return api.VehicleHealth(h.VehicleID)
}

// TickResult carries hydra tick outcomes for the engine to announce.
type TickResult struct {
	Checkpoint      bool
	EscortWin       bool
	DefendWin       bool
	CheckpointMsg   string
	HydraMinuteMsg  string
}

func (h *Hydra) Tick(api API, score *Scoring, nowMs uint64) TickResult {
	var res TickResult
	if h.State != HydraPatrol || h.VehicleID < 0 {
		return res
	}
	hp := api.VehicleHealth(h.VehicleID)
	if hp <= 0 {
		score.AddDefend(PointsHydraDestroy)
		h.State = HydraDestroyed
		res.DefendWin = true
		return res
	}

	if len(h.Waypoints) == 0 {
		return res
	}

	if h.Index >= len(h.Waypoints) {
		h.State = HydraObjectiveReached
		res.EscortWin = true
		return res
	}

	target := h.Waypoints[h.Index]
	pos := api.VehiclePos(h.VehicleID)
	if dist(pos, target) <= checkpointRadiusM {
		h.Index++
		score.AddEscort(PointsCheckpoint)
		res.Checkpoint = true
		if h.Index >= len(h.Waypoints) {
			h.State = HydraObjectiveReached
			res.EscortWin = true
			res.CheckpointMsg = fmt.Sprintf("Final checkpoint reached! Escort +%d", PointsCheckpoint)
		} else {
			res.CheckpointMsg = fmt.Sprintf("Checkpoint %d/%d reached! Escort +%d", h.Index, len(h.Waypoints), PointsCheckpoint)
		}
		return res
	}

	h.moveToward(api, pos, target)

	if h.lastMinute == 0 {
		h.lastMinute = nowMs
	} else if nowMs-h.lastMinute >= 60000 {
		h.lastMinute = nowMs
		score.AddEscort(PointsHydraMinute)
		res.HydraMinuteMsg = fmt.Sprintf("Hydra holding route (+%d escort)", PointsHydraMinute)
	}
	return res
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
