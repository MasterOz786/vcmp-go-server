package safari

import (
	"fmt"
	"math"
)

var hydraCameraNames = []string{
	"Pilot (default)",
	"Chase cam",
	"Side orbit",
	"Tactical overhead",
}

func hydraCameraName(mode int) string {
	if mode >= 0 && mode < len(hydraCameraNames) {
		return hydraCameraNames[mode]
	}
	return "Unknown"
}

func (e *Engine) cycleHydraCamera(playerID int) {
	s := e.teams.session(playerID)
	if s == nil {
		e.api.Send(playerID, ColourYellow, "You are not in a session.")
		return
	}
	vid := e.playerHydraVehicleID(playerID)
	if vid < 0 {
		e.api.Send(playerID, ColourYellow, "You must be in a Hydra. Use /testhydra first.")
		return
	}
	s.HydraCameraMode = (s.HydraCameraMode + 1) % HydraCamCount
	e.applyHydraCamera(playerID, vid, s.HydraCameraMode)
	e.api.Send(playerID, ColourCyan, fmt.Sprintf("Hydra view: %s (/hydraview or V)", hydraCameraName(s.HydraCameraMode)))
}

func (e *Engine) resetHydraCamera(playerID int) {
	s := e.teams.session(playerID)
	if s == nil {
		return
	}
	s.HydraCameraMode = HydraCamDefault
	_ = e.api.RestoreCamera(playerID)
}

func (e *Engine) applyHydraCamera(playerID, vehicleID, mode int) {
	if mode == HydraCamDefault {
		_ = e.api.RestoreCamera(playerID)
		return
	}
	camPos, lookAt := hydraCameraOffsets(e.api, vehicleID, mode)
	_ = e.api.SetCamera(playerID, camPos, lookAt)
}

func (e *Engine) updateHydraCameras() {
	for playerID, s := range e.teams.sessions {
		if !e.api.IsConnected(playerID) || s.HydraCameraMode == HydraCamDefault {
			continue
		}
		vid := e.playerHydraVehicleID(playerID)
		if vid < 0 {
			continue
		}
		e.applyHydraCamera(playerID, vid, s.HydraCameraMode)
	}
}

func hydraCameraOffsets(api API, vehicleID, mode int) (camPos, lookAt Vec3) {
	pos := api.VehiclePos(vehicleID)
	lookAt = pos
	rot := api.VehicleRotationEuler(vehicleID)
	heading := rot.Z * math.Pi / 180
	sinH := float32(math.Sin(float64(heading)))
	cosH := float32(math.Cos(float64(heading)))

	switch mode {
	case HydraCamChase:
		dist, height := 38.0, 14.0
		camPos = Vec3{
			X: pos.X - sinH*float32(dist),
			Y: pos.Y + cosH*float32(dist),
			Z: pos.Z + float32(height),
		}
	case HydraCamSide:
		dist, height := 28.0, 10.0
		camPos = Vec3{
			X: pos.X + cosH*float32(dist),
			Y: pos.Y + sinH*float32(dist),
			Z: pos.Z + float32(height),
		}
	case HydraCamTactical:
		camPos = Vec3{X: pos.X, Y: pos.Y, Z: pos.Z + 55}
	default:
		camPos = pos
	}
	return camPos, lookAt
}

func (e *Engine) playerHydraVehicleID(playerID int) int {
	if vid := e.api.PlayerVehicleID(playerID); vid >= 0 {
		if e.isHydraVehicle(vid) {
			return vid
		}
	}
	s := e.teams.session(playerID)
	if s != nil && s.TestHydraVehicleID >= 0 && e.api.VehicleHealth(s.TestHydraVehicleID) > 0 {
		return s.TestHydraVehicleID
	}
	return -1
}

func (e *Engine) isHydraVehicle(vehicleID int) bool {
	if vehicleID < 0 {
		return false
	}
	if e.round.Hydra.VehicleID == vehicleID {
		return true
	}
	for _, s := range e.teams.sessions {
		if s.TestHydraVehicleID == vehicleID {
			return true
		}
	}
	return false
}
