package gameplay

import (
	"fmt"
	"math"

	"github.com/masteroz/vcmp-go-server/safari/apidef"
)

const (
	hydraStepMeters   = 5.0
	checkpointRadiusM = 15.0
)

const HydraVehicleArchive = "store/vehicles/v6460_t0_p1_Hydra.7z"

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
	Waypoints  []apidef.Vec3
	Index      int
	lastMinute uint64
}

func NewHydra() *Hydra {
	return &Hydra{State: HydraIdle, Index: 0, VehicleID: -1}
}

func CreateHydraVehicle(api apidef.API, model, world int, pos apidef.Vec3, angle float32) int {
	return api.CreateVehicle(model, world, pos, angle, 0, 0)
}

func TestHydraSpawnPos(mapCfg apidef.MapConfig, playerPos apidef.Vec3) apidef.Vec3 {
	if mapCfg.HydraStart.X != 0 || mapCfg.HydraStart.Y != 0 {
		return apidef.Vec3{
			X: mapCfg.HydraStart.X,
			Y: mapCfg.HydraStart.Y,
			Z: mapCfg.HydraStart.Z + 3,
		}
	}
	z := playerPos.Z + 12
	if z < 14 {
		z = 14.5
	}
	return apidef.Vec3{X: playerPos.X + 6, Y: playerPos.Y + 6, Z: z}
}

func FormatHydraSpawnFailure(api apidef.API, model int) string {
	msg := fmt.Sprintf("Failed to spawn Hydra (model %d).", model)
	if detail := api.LastErrorString(); detail != "" {
		msg += " " + detail
	}
	msg += " Custom vehicle file must be " + HydraVehicleArchive + " (see forum.vc-mp.org topic 975)."
	return msg
}

func WarnHydraModelMismatch(api apidef.API, playerID, vehicleID, expectedModel int) {
	if vehicleID < 0 {
		return
	}
	actual := api.VehicleModel(vehicleID)
	if actual == expectedModel {
		return
	}
	msg := fmt.Sprintf(
		"Hydra spawned as model %d (expected %d). Install %s and reconnect — otherwise it flies like a Maverick.",
		actual, expectedModel, HydraVehicleArchive,
	)
	api.Log(fmt.Sprintf("[safari] WARNING player %d vehicle %d: %s", playerID, vehicleID, msg))
	if playerID >= 0 && api.IsConnected(playerID) {
		api.Send(playerID, apidef.ColourRed, msg)
	}
}

func (h *Hydra) Spawn(api apidef.API, mapCfg apidef.MapConfig, model int) int {
	pos := mapCfg.HydraStart
	if len(mapCfg.Waypoints) > 0 && mapCfg.HydraStart.X == 0 && mapCfg.HydraStart.Y == 0 {
		pos = mapCfg.Waypoints[0]
	}
	h.Waypoints = mapCfg.Waypoints
	h.Index = 0
	h.VehicleID = CreateHydraVehicle(api, model, mapCfg.World, pos, mapCfg.HydraAngle)
	if h.VehicleID >= 0 {
		api.SetVehicleHealth(h.VehicleID, apidef.HydraMaxHP)
		WarnHydraModelMismatch(api, -1, h.VehicleID, model)
		h.State = HydraPatrol
	} else {
		h.State = HydraIdle
		api.Log(FormatHydraSpawnFailure(api, model))
	}
	return h.VehicleID
}

func (h *Hydra) Destroy(api apidef.API) {
	if h.VehicleID >= 0 {
		api.DeleteVehicle(h.VehicleID)
		h.VehicleID = -1
	}
	if h.State == HydraPatrol {
		h.State = HydraIdle
	}
}

func (h *Hydra) Health(api apidef.API) float32 {
	if h.VehicleID < 0 {
		return 0
	}
	return api.VehicleHealth(h.VehicleID)
}

// TickResult carries hydra tick outcomes for the engine to announce.
type TickResult struct {
	Checkpoint     bool
	EscortWin      bool
	DefendWin      bool
	CheckpointMsg  string
	HydraMinuteMsg string
}

func (h *Hydra) Tick(api apidef.API, score *apidef.Scoring, nowMs uint64) TickResult {
	var res TickResult
	if h.State != HydraPatrol || h.VehicleID < 0 {
		return res
	}
	hp := api.VehicleHealth(h.VehicleID)
	if hp <= 0 {
		score.AddDefend(apidef.PointsHydraDestroy)
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
		score.AddEscort(apidef.PointsCheckpoint)
		res.Checkpoint = true
		if h.Index >= len(h.Waypoints) {
			h.State = HydraObjectiveReached
			res.EscortWin = true
			res.CheckpointMsg = fmt.Sprintf("Final checkpoint reached! Escort +%d", apidef.PointsCheckpoint)
		} else {
			res.CheckpointMsg = fmt.Sprintf("Checkpoint %d/%d reached! Escort +%d", h.Index, len(h.Waypoints), apidef.PointsCheckpoint)
		}
		return res
	}

	h.moveToward(api, pos, target)

	if h.lastMinute == 0 {
		h.lastMinute = nowMs
	} else if nowMs-h.lastMinute >= 60000 {
		h.lastMinute = nowMs
		score.AddEscort(apidef.PointsHydraMinute)
		res.HydraMinuteMsg = fmt.Sprintf("Hydra holding route (+%d escort)", apidef.PointsHydraMinute)
	}
	return res
}

func (h *Hydra) moveToward(api apidef.API, pos, target apidef.Vec3) {
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
	newPos := apidef.Vec3{
		X: pos.X + dx*ratio,
		Y: pos.Y + dy*ratio,
		Z: pos.Z + dz*ratio,
	}
	api.SetVehiclePosition(h.VehicleID, newPos)
}

func dist(a, b apidef.Vec3) float32 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	dz := b.Z - a.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}
