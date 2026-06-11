package safari

import (
	"fmt"
	"math"
)

const hydraCameraUpdateInterval = float32(0.05) // 20 Hz — smooth enough without per-frame lag

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

func hydraCameraOffsets(api API, vehicleID, mode int) (camPos, lookAt Vec3) {
	pos := api.VehiclePos(vehicleID)
	lookAt = Vec3{X: pos.X, Y: pos.Y, Z: pos.Z + 1}
	rot := api.VehicleRotationEuler(vehicleID)
	heading := rot.Z * math.Pi / 180
	sinH := float32(math.Sin(float64(heading)))
	cosH := float32(math.Cos(float64(heading)))

	switch mode {
	case HydraCamChase:
		camPos = Vec3{
			X: pos.X - sinH*38,
			Y: pos.Y + cosH*38,
			Z: pos.Z + 14,
		}
	case HydraCamSide:
		camPos = Vec3{
			X: pos.X + cosH*28,
			Y: pos.Y + sinH*28,
			Z: pos.Z + 10,
		}
	case HydraCamTactical:
		camPos = Vec3{X: pos.X, Y: pos.Y, Z: pos.Z + 55}
	default:
		camPos = pos
	}
	return camPos, lookAt
}

func (e *Engine) sendHydraCamPacket(playerID, mode, vehicleID int) {
	if !e.api.IsConnected(playerID) {
		return
	}
	s := NewStreamWriter()
	s.WriteInt(PacketHydraCam)
	s.WriteInt(int32(mode))
	s.WriteInt(int32(vehicleID))
	if err := e.api.SendScriptData(playerID, s.Bytes()); err != nil {
		e.api.Log(fmt.Sprintf("[safari] hydra cam stream to %d failed: %v", playerID, err))
	}
}

func (e *Engine) applyHydraCamera(playerID, vehicleID, mode int) {
	if mode <= HydraCamDefault {
		_ = e.api.RestoreCamera(playerID)
		return
	}
	if vehicleID < 0 {
		return
	}
	camPos, lookAt := hydraCameraOffsets(e.api, vehicleID, mode)
	_ = e.api.SetCamera(playerID, camPos, lookAt)
}

func (e *Engine) updateHydraCameras() {
	for playerID, s := range e.teams.sessions {
		if !e.api.IsConnected(playerID) || s.HydraCameraMode <= HydraCamDefault {
			continue
		}
		vid := e.playerHydraVehicleID(playerID)
		if vid < 0 || e.api.PlayerVehicleID(playerID) != vid {
			continue
		}
		e.applyHydraCamera(playerID, vid, s.HydraCameraMode)
	}
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
	e.sendHydraCamPacket(playerID, s.HydraCameraMode, vid)
	e.applyHydraCamera(playerID, vid, s.HydraCameraMode)
	e.api.Send(playerID, ColourCyan, fmt.Sprintf("Hydra view: %s (H or /hydraview)", hydraCameraName(s.HydraCameraMode)))
}

func (e *Engine) resetHydraCamera(playerID int) {
	s := e.teams.session(playerID)
	if s == nil {
		return
	}
	s.HydraCameraMode = HydraCamDefault
	e.sendHydraCamPacket(playerID, HydraCamOff, -1)
	_ = e.api.RestoreCamera(playerID)
}

func (e *Engine) syncHydraCamera(playerID, vehicleID int) {
	s := e.teams.session(playerID)
	if s == nil {
		return
	}
	mode := s.HydraCameraMode
	if mode < HydraCamDefault {
		mode = HydraCamDefault
	}
	e.sendHydraCamPacket(playerID, mode, vehicleID)
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
