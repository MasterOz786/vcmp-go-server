package safari

import "fmt"

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

// sendHydraCamPacket tells the client script to enable/disable or switch Hydra views.
// Modes: -1 off (left hydra), 0 pilot default, 1 chase, 2 side, 3 tactical.
// Camera math runs client-side in store/script to avoid server SetCamera lag.
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
		return
	}
	if mode >= HydraCamDefault {
		e.warnIfNoClientScript(playerID)
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
	e.api.Send(playerID, ColourCyan, fmt.Sprintf("Hydra view: %s (H or /hydraview)", hydraCameraName(s.HydraCameraMode)))
}

func (e *Engine) resetHydraCamera(playerID int) {
	s := e.teams.session(playerID)
	if s == nil {
		return
	}
	s.HydraCameraMode = HydraCamDefault
	e.sendHydraCamPacket(playerID, HydraCamOff, -1)
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
