package safari

// HandlePickupPickAttemptSync runs on the VC:MP callback thread (must not block on the event queue).
func (e *Engine) HandlePickupPickAttemptSync(pickupID, playerID int) bool {
	return e.HandlePickupPickAttempt(pickupID, playerID)
}

func (e *Engine) HandleEnterVehicleRequestSync(playerID, vehicleID, slot int) bool {
	return e.HandleEnterVehicleRequest(playerID, vehicleID, slot)
}
