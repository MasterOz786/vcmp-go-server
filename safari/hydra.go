package safari

import (
	"fmt"
	"math"
	"strings"

	"github.com/masteroz/vcmp-go-server/safari/apidef"
	"github.com/masteroz/vcmp-go-server/safari/clientscript"
	"github.com/masteroz/vcmp-go-server/safari/gameplay"
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

func hydraCameraOffsets(api API, vehicleID, mode int) (camPos, lookAt Vec3) {
	pos := api.VehiclePos(vehicleID)
	lookAt = Vec3{X: pos.X, Y: pos.Y, Z: pos.Z + 1}
	rot := api.VehicleRotationEuler(vehicleID)
	heading := rot.Z * math.Pi / 180
	sinH := float32(math.Sin(float64(heading)))
	cosH := float32(math.Cos(float64(heading)))

	switch mode {
	case HydraCamChase:
		camPos = Vec3{X: pos.X - sinH*38, Y: pos.Y + cosH*38, Z: pos.Z + 14}
	case HydraCamSide:
		camPos = Vec3{X: pos.X + cosH*28, Y: pos.Y + sinH*28, Z: pos.Z + 10}
	case HydraCamTactical:
		camPos = Vec3{X: pos.X, Y: pos.Y, Z: pos.Z + 55}
	default:
		camPos = pos
	}
	return camPos, lookAt
}

func (e *Engine) hydraModel() int {
	if e.cfg.HydraModel > 0 {
		return e.cfg.HydraModel
	}
	return HydraModel
}

func (e *Engine) sendHydraCamPacket(playerID, mode, vehicleID int) {
	e.sendScriptPacket(playerID, clientscript.HydraCam(mode, vehicleID), "hydra cam stream")
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
	e.teams.ForEachSession(func(playerID int, s *apidef.PlayerSession) {
		if !e.api.IsConnected(playerID) || s.HydraCameraMode <= HydraCamDefault {
			return
		}
		vid := e.playerHydraVehicleID(playerID)
		if vid < 0 || e.api.PlayerVehicleID(playerID) != vid {
			return
		}
		e.applyHydraCamera(playerID, vid, s.HydraCameraMode)
	})
}

func (e *Engine) cycleHydraCamera(playerID int) {
	s := e.teams.Session(playerID)
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
	s := e.teams.Session(playerID)
	if s == nil {
		return
	}
	s.HydraCameraMode = HydraCamDefault
	e.sendHydraCamPacket(playerID, HydraCamOff, -1)
	_ = e.api.RestoreCamera(playerID)
}

func (e *Engine) syncHydraCamera(playerID, vehicleID int) {
	s := e.teams.Session(playerID)
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
	s := e.teams.Session(playerID)
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
	found := false
	e.teams.ForEachSession(func(_ int, s *apidef.PlayerSession) {
		if s.TestHydraVehicleID == vehicleID {
			found = true
		}
	})
	return found
}

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

	model := e.hydraModel()
	pos := gameplay.TestHydraSpawnPos(e.mapCfg, e.api.PlayerPosition(playerID))
	world := e.mapCfg.World
	if e.api.IsConnected(playerID) {
		world = e.api.PlayerWorld(playerID)
	}

	vid := gameplay.CreateHydraVehicle(e.api, model, world, pos, e.mapCfg.HydraAngle)
	if vid < 0 {
		e.api.Send(playerID, ColourRed, gameplay.FormatHydraSpawnFailure(e.api, model))
		return CommandResult{Handled: true, Deny: true}
	}
	e.api.SetVehicleHealth(vid, HydraMaxHP)
	if world != e.mapCfg.World {
		_ = e.api.SetVehicleWorld(vid, world)
	}

	s := e.teams.Session(playerID)
	if s == nil {
		e.ensurePlayerSession(playerID)
		s = e.teams.Session(playerID)
	}
	if s != nil {
		s.TestHydraVehicleID = vid
		s.HydraCameraMode = HydraCamDefault
	}

	e.api.PutPlayerInVehicle(playerID, vid, 0)
	gameplay.WarnHydraModelMismatch(e.api, playerID, vid, model)
	e.syncHydraCamera(playerID, vid)

	e.api.Send(playerID, ColourGreen, "Test Hydra ready — you are in the pilot seat.")
	e.api.Send(playerID, ColourCyan, "Press H or /hydraview to cycle camera views.")
	e.api.Log(fmt.Sprintf("[safari] test hydra spawned for player %d vehicle=%d model=%d", playerID, vid, model))
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) stopTestHydra(playerID int) {
	s := e.teams.Session(playerID)
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
