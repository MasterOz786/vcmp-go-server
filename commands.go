package main

import (
	"fmt"
	"strconv"
	"strings"
)

func resolvePlayer(arg string, partial bool) (int, string) {
	if arg == "" {
		return -1, "missing player"
	}
	if strings.HasPrefix(arg, "#") {
		id, err := strconv.Atoi(strings.TrimPrefix(arg, "#"))
		if err != nil || !bridgeIsConnected(id) {
			return -1, "player not found"
		}
		return id, ""
	}
	if id := bridgePlayerIDFromName(arg); id >= 0 && bridgeIsConnected(id) {
		return id, ""
	}
	if !partial {
		return -1, "player not found"
	}
	var matches []int
	lower := strings.ToLower(arg)
	for id := 0; id < MaxPlayers; id++ {
		if !bridgeIsConnected(id) {
			continue
		}
		if strings.Contains(strings.ToLower(bridgePlayerName(id)), lower) {
			matches = append(matches, id)
		}
	}
	switch len(matches) {
	case 0:
		return -1, "no such guy lives here"
	case 1:
		return matches[0], ""
	default:
		return -1, "multiple players match that name"
	}
}

func resolveVehicle(arg string) (int, string) {
	if strings.HasPrefix(arg, "#") {
		id, err := strconv.Atoi(strings.TrimPrefix(arg, "#"))
		if err != nil {
			return -1, "invalid vehicle id"
		}
		return id, ""
	}
	id, err := strconv.Atoi(arg)
	if err != nil {
		return -1, "usage: #<vehicle id>"
	}
	return id, ""
}

func handleDemoCommand(playerID int, raw string) FilterResult {
	parts := strings.Fields(strings.TrimSpace(raw))
	if len(parts) == 0 {
		return FilterDeny
	}
	cmd := strings.ToLower(strings.TrimPrefix(parts[0], "/"))
	args := parts[1:]

	switch cmd {
	case "renamed":
		bridgeSendClientMessage(playerID, ColourResponse, "This command has a custom name that differs from method name.")
	case "lottery":
		if len(args) < 1 {
			bridgeSendClientMessage(playerID, ColourResponse, "Usage: /lottery <number>")
			return FilterDeny
		}
		n, err := strconv.Atoi(args[0])
		if err != nil {
			bridgeSendClientMessage(playerID, ColourResponse, "Enter a valid number.")
			return FilterDeny
		}
		bridgeSendClientMessage(playerID, ColourResponse, fmt.Sprintf("Good job on finally entering the correct parameters. Oh and %d is not a winning number", n))
	case "finddefault":
		if len(args) < 1 {
			return FilterDeny
		}
		target, msg := resolvePlayer(args[0], false)
		if msg != "" {
			bridgeSendClientMessage(playerID, ColourResponse, msg)
			return FilterDeny
		}
		bridgeSendClientMessage(playerID, ColourResponse, fmt.Sprintf("Player '%s' found.", bridgePlayerName(target)))
	case "findnoerror":
		if len(args) < 1 {
			return FilterDeny
		}
		target, _ := resolvePlayer(args[0], false)
		if target < 0 {
			bridgeSendClientMessage(playerID, ColourResponse, "No such guy lives here.")
		} else {
			bridgeSendClientMessage(playerID, ColourResponse, fmt.Sprintf("Yep, found %s.", bridgePlayerName(target)))
		}
	case "findpartial":
		if len(args) < 1 {
			return FilterDeny
		}
		target, msg := resolvePlayer(args[0], true)
		if msg != "" {
			bridgeSendClientMessage(playerID, ColourResponse, msg)
			return FilterDeny
		}
		bridgeSendClientMessage(playerID, ColourResponse, fmt.Sprintf("The dude called %s seems close enough.", bridgePlayerName(target)))
	case "getservername":
		if !requireAdmin(playerID) {
			return FilterDeny
		}
		bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Server name: %s", bridgeGetServerName()))
	case "setservername":
		if !requireAdmin(playerID) {
			return FilterDeny
		}
		if len(args) < 1 {
			return FilterDeny
		}
		name := strings.Join(args, " ")
		bridgeSetServerName(name)
		bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Server name changed to: %s", name))
	case "reload":
		if !requireAdmin(playerID) {
			return FilterDeny
		}
		bridgeSendClientMessage(playerID, ColourYellowish, "Reload is not available in the native Go plugin (unlike the Java plugin). Restart the server instead.")
	case "pingme":
		if startPlayerPing(playerID) {
			bridgeSendClientMessage(playerID, ColourTimer, "Pinging you every 5 seconds.")
		} else {
			bridgeSendClientMessage(playerID, ColourTimer, "You are already being pinged.")
		}
	case "stopping":
		if stopPlayerPing(playerID) {
			bridgeSendClientMessage(playerID, ColourTimer, "Stopping your pings.")
		} else {
			bridgeSendClientMessage(playerID, ColourTimer, "You are not being pinged.")
		}
	case "createvehicle":
		if len(args) < 1 {
			return FilterDeny
		}
		model, err := strconv.Atoi(args[0])
		if err != nil {
			return FilterDeny
		}
		pos := bridgePlayerPos(playerID)
		pos.X += 5
		world := bridgePlayerWorld(playerID)
		vehicleID := bridgeCreateVehicle(model, world, pos, 0, 1, 1)
		if vehicleID >= 0 {
			bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Vehicle %d created! YAY!", vehicleID))
		} else {
			bridgeSendClientMessage(playerID, ColourYellowish, "Could not create vehicle.")
		}
	case "getworld":
		bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Your world is %d.", bridgePlayerWorld(playerID)))
	case "setworld":
		if len(args) < 1 {
			return FilterDeny
		}
		world, err := strconv.Atoi(args[0])
		if err != nil {
			return FilterDeny
		}
		bridgeSetPlayerWorld(playerID, world)
		bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Set your world to %d.", world))
	case "getplayervehicle":
		if len(args) < 1 {
			return FilterDeny
		}
		target, msg := resolvePlayer(args[0], false)
		if msg != "" {
			bridgeSendClientMessage(playerID, ColourYellowish, msg)
			return FilterDeny
		}
		vid := bridgePlayerVehicleID(target)
		if vid >= 0 {
			bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Player %s is in vehicle %d.", bridgePlayerName(target), vid))
		} else {
			bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Player %s is not in a vehicle.", bridgePlayerName(target)))
		}
	case "getvehiclehealth":
		if len(args) < 1 {
			return FilterDeny
		}
		vid, msg := resolveVehicle(args[0])
		if msg != "" {
			bridgeSendClientMessage(playerID, ColourYellowish, msg)
			return FilterDeny
		}
		bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Vehicle %d health: %.2f.", vid, bridgeVehicleHealth(vid)))
	case "putplayerinvehicle":
		if len(args) < 5 {
			bridgeSendClientMessage(playerID, ColourYellowish, "Usage: /putplayerinvehicle <player> <vehicle> <slot> <makeroom 0|1> <warp 0|1>")
			return FilterDeny
		}
		target, pmsg := resolvePlayer(args[0], false)
		if pmsg != "" {
			bridgeSendClientMessage(playerID, ColourYellowish, pmsg)
			return FilterDeny
		}
		vid, vmsg := resolveVehicle(args[1])
		if vmsg != "" {
			bridgeSendClientMessage(playerID, ColourYellowish, vmsg)
			return FilterDeny
		}
		slot, _ := strconv.Atoi(args[2])
		makeRoom := args[3] == "1"
		warp := args[4] == "1"
		bridgePutInVehicle(target, vid, slot, makeRoom, warp)
		bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Sending player %s into vehicle %d.", bridgePlayerName(target), vid))
	case "getvehicleoccupant":
		if len(args) < 2 {
			return FilterDeny
		}
		vid, vmsg := resolveVehicle(args[0])
		if vmsg != "" {
			bridgeSendClientMessage(playerID, ColourYellowish, vmsg)
			return FilterDeny
		}
		slot, err := strconv.Atoi(args[1])
		if err != nil {
			return FilterDeny
		}
		occupant := bridgeVehicleOccupant(vid, slot)
		if occupant < 0 {
			bridgeSendClientMessage(playerID, ColourYellowish, "That vehicle slot is not occupied.")
		} else {
			bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("That slot is occupied by %s (%d).", bridgePlayerName(occupant), occupant))
		}
	case "getvehicleposition":
		if len(args) < 1 {
			return FilterDeny
		}
		vid, vmsg := resolveVehicle(args[0])
		if vmsg != "" {
			bridgeSendClientMessage(playerID, ColourYellowish, vmsg)
			return FilterDeny
		}
		pos := bridgeVehiclePos(vid)
		bridgeSendClientMessage(playerID, ColourYellowish, fmt.Sprintf("Vehicle %d is at (%.2f, %.2f, %.2f).", vid, pos.X, pos.Y, pos.Z))
	case "breakcar":
		if len(args) < 1 {
			return FilterDeny
		}
		vid, vmsg := resolveVehicle(args[0])
		if vmsg != "" {
			bridgeSendClientMessage(playerID, ColourYellowish, vmsg)
			return FilterDeny
		}
		bridgeBreakVehicle(vid)
		bridgeSendClientMessage(playerID, ColourYellowish, "Broke that car.")
	case "gibadmin":
		bridgeSetAdmin(playerID, true)
		bridgeSendClientMessage(playerID, ColourWhite, "Gabe admin")
	case "help":
		sendHelp(playerID)
	default:
		bridgeSendClientMessage(playerID, ColourYellow, "Unknown command. Type /help")
	}
	return FilterDeny
}

func requireAdmin(playerID int) bool {
	if bridgeIsAdmin(playerID) {
		return true
	}
	bridgeSendClientMessage(playerID, ColourYellowish, "You need to be an admin for this (/gibadmin).")
	return false
}

func sendHelp(playerID int) {
	lines := []string{
		"--- Go Demo Commands (Kotlin port) ---",
		"/renamed /lottery /finddefault /findnoerror /findpartial",
		"/pingme /stopping",
		"/createvehicle /getworld /setworld /getplayervehicle",
		"/getvehiclehealth /putplayerinvehicle /getvehicleoccupant",
		"/getvehicleposition /breakcar /gibadmin",
		"Admin: /getservername /setservername /reload",
	}
	for _, line := range lines {
		bridgeSendClientMessage(playerID, ColourWhite, line)
	}
}
