package safari

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CommandResult struct {
	Handled bool
	Deny    bool
}

func normalizeCommand(raw string) (string, bool) {
	cmd := strings.TrimSpace(raw)
	if cmd == "" {
		return "", false
	}
	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}
	return cmd, true
}

func (e *Engine) HandleCommand(playerID int, raw string) CommandResult {
	cmd, ok := normalizeCommand(raw)
	if !ok {
		return CommandResult{}
	}
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return CommandResult{}
	}
	name := strings.ToLower(parts[0])
	args := parts[1:]

	switch name {
	case "/help":
		e.sendHelp(playerID)
		return CommandResult{Handled: true, Deny: true}
	case "/startsafari":
		if !e.isAdmin(playerID) {
			e.api.Send(playerID, ColourRed, "Admin only.")
			return CommandResult{Handled: true, Deny: true}
		}
		if e.round.State == RoundActive {
			e.api.Send(playerID, ColourYellow, "Round already active.")
			return CommandResult{Handled: true, Deny: true}
		}
		e.startRound()
		return CommandResult{Handled: true, Deny: true}
	case "/stopsafari":
		if !e.isAdmin(playerID) {
			e.api.Send(playerID, ColourRed, "Admin only.")
			return CommandResult{Handled: true, Deny: true}
		}
		if e.round.State == RoundActive {
			e.endRound(TeamDefend, "Round stopped by admin.")
		}
		return CommandResult{Handled: true, Deny: true}
	case "/pausesafari":
		if !e.isAdmin(playerID) {
			e.api.Send(playerID, ColourRed, "Admin only.")
			return CommandResult{Handled: true, Deny: true}
		}
		if e.round.State != RoundActive {
			e.api.Send(playerID, ColourYellow, "No active round.")
			return CommandResult{Handled: true, Deny: true}
		}
		e.TogglePause()
		return CommandResult{Handled: true, Deny: true}
	case "/autostart":
		if !e.isAdmin(playerID) {
			e.api.Send(playerID, ColourRed, "Admin only.")
			return CommandResult{Handled: true, Deny: true}
		}
		if len(args) == 0 {
			state := "off"
			if e.AutostartEnabled() {
				state = "on"
			}
			e.api.Send(playerID, ColourWhite, "Autostart is "+state+".")
			return CommandResult{Handled: true, Deny: true}
		}
		switch strings.ToLower(args[0]) {
		case "on", "1", "true":
			e.SetAutostart(true)
		case "off", "0", "false":
			e.SetAutostart(false)
		default:
			e.api.Send(playerID, ColourYellow, "Usage: /autostart on|off")
		}
		return CommandResult{Handled: true, Deny: true}
	case "/pack":
		return e.cmdPack(playerID, args)
	case "/mark":
		return e.cmdMark(playerID, args)
	case "/status":
		e.cmdStatus(playerID)
		return CommandResult{Handled: true, Deny: true}
	case "/stats":
		e.cmdStats(playerID)
		return CommandResult{Handled: true, Deny: true}
	case "/register":
		e.cmdRegister(playerID)
		return CommandResult{Handled: true, Deny: true}
	case "/testhydra":
		return e.cmdTestHydra(playerID, args)
	case "/hydraview":
		e.cycleHydraCamera(playerID)
		return CommandResult{Handled: true, Deny: true}
	case "/getpos":
		e.cmdGetPos(playerID, args)
		return CommandResult{Handled: true, Deny: true}
	case "/lobby":
		e.cmdLobby(playerID)
		return CommandResult{Handled: true, Deny: true}
	case "/scoreboard":
		return e.cmdScoreboard(playerID, args)
	case "/leaderboard", "/leaderboards":
		return e.cmdLeaderboard(playerID)
	case "/gotopos":
		return e.cmdGoToPos(playerID, args)
	case "/getvehicle":
		return e.cmdGetVehicle(playerID, args)
	case "/wep":
		return e.cmdWep(playerID, args)
	case "/reload":
		return e.cmdReload(playerID, args)
	default:
		if strings.HasPrefix(name, "/") && len(name) > 1 {
			e.api.Send(playerID, ColourYellow, "Unknown command. Try /help")
			return CommandResult{Handled: true, Deny: true}
		}
		return CommandResult{}
	}
}

func (e *Engine) sendHelp(playerID int) {
	e.api.Send(playerID, ColourCyan, "=== Project Safari: Hydra Warfare ===")
	lines := []string{
		"/pack 1|2|3 - choose loadout (keys 1-3)",
		"/mark [name] - Escort: designate target",
		"/testhydra - spawn test Hydra and warp in",
		"/testhydra stop - remove your test Hydra",
		"/hydraview or H - cycle Hydra camera",
		"/getpos [name] - your position (or another player)",
		"/lobby - warp to the lobby spawn",
		"/scoreboard - show live round HUD",
		"/scoreboard end - round stats overlay (P to close)",
		"/scoreboard hide - hide scoreboard UI",
		"/leaderboard or /leaderboards - toggle leaderboard UI",
		"/status - round info",
		"/stats - your persistent stats",
		"/register - open registration window",
		"P - weapon pack picker (also closes UI)",
		"/startsafari - admin: start round",
		"/stopsafari - admin: stop round",
		"/pausesafari - admin: pause/resume round",
		"/autostart on|off - admin: toggle autostart",
		"/gotopos <x> <y> <z> - admin: warp to coordinates",
		"/getvehicle <model id> - admin: vehicle info by model (no args lists all)",
		"/wep <id> [ammo] - admin: give weapon",
		"/reload - admin: reload config/map from disk",
		"/reload scripts - admin: kick all (refresh client script)",
		"/reload server - admin: rebuild plugin and restart server",
	}
	for _, l := range lines {
		e.api.Send(playerID, ColourYellow, l)
	}
}

func (e *Engine) cmdPack(playerID int, args []string) CommandResult {
	if len(args) != 1 {
		e.api.Send(playerID, ColourYellow, fmt.Sprintf("Usage: /pack 1-%d", MaxPack))
		return CommandResult{Handled: true, Deny: true}
	}
	pack, err := strconv.Atoi(args[0])
	if err != nil || pack < 1 || pack > MaxPack {
		e.api.Send(playerID, ColourYellow, fmt.Sprintf("Pack must be 1 to %d.", MaxPack))
		return CommandResult{Handled: true, Deny: true}
	}
	if e.round.State == RoundActive && e.teams.HasSpawnedThisRound(playerID) {
		e.api.Send(playerID, ColourYellow, "Cannot change pack after spawning this round.")
		return CommandResult{Handled: true, Deny: true}
	}
	e.ensurePlayerSession(playerID)
	if e.teams.Team(playerID) == 0 {
		e.api.Send(playerID, ColourRed, "You are not assigned to a team yet.")
		return CommandResult{Handled: true, Deny: true}
	}
	e.ApplyPack(playerID, pack)
	team := e.teams.Team(playerID)
	var name string
	if team == TeamEscort {
		name = EscortPacks()[pack].Name
	} else {
		name = DefendPacks()[pack].Name
	}
	e.api.Send(playerID, ColourGreen, fmt.Sprintf("Loadout equipped: %s", name))
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) cmdMark(playerID int, args []string) CommandResult {
	if e.round.State != RoundActive {
		e.announceTo(playerID, MsgMarkFail, "No active round.")
		return CommandResult{Handled: true, Deny: true}
	}
	target := strings.Join(args, " ")
	ok, msg := e.marking.TryMark(e.api, e.db, playerID, target, &e.round.Score)
	if !ok {
		e.announceTo(playerID, MsgMarkFail, msg)
	} else {
		e.announce(MsgMarkSuccess, msg)
		e.teams.SyncScores(e.api, e.round.Score)
		e.BroadcastScoreboard()
	}
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) cmdStatus(playerID int) {
	switch e.round.State {
	case RoundIdle:
		e.api.Send(playerID, ColourWhite, "Safari idle. Waiting for round start.")
	case RoundActive:
		e.api.Send(playerID, ColourCyan, e.formatStatusLine())
	case RoundEnded:
		e.api.Send(playerID, ColourWhite, fmt.Sprintf("Round ended: %s", e.round.EndReason))
	}
}

func (e *Engine) cmdStats(playerID int) {
	uid := e.api.PlayerUID(playerID)
	if uid == "" {
		e.api.Send(playerID, ColourRed, "Could not load your UID.")
		return
	}
	st, ok := e.db.CachedStats(uid)
	if !ok {
		e.db.PrefetchStats(uid)
		e.api.Send(playerID, ColourYellow, "Stats loading — try again in a moment.")
		return
	}
	e.api.Send(playerID, ColourWhite, fmt.Sprintf(
		"Stats — Escort pts: %d | Defend pts: %d | Marks: %d | Rounds: %d | Wins: %d",
		st.EscortPts, st.DefendPts, st.Marks, st.RoundsPlayed, st.RoundsWon,
	))
}

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "0:00"
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

func (e *Engine) cmdRegister(playerID int) {
	uid := e.api.PlayerUID(playerID)
	if uid == "" {
		e.api.Send(playerID, ColourRed, "Could not load your UID.")
		return
	}
	registered, err := e.db.IsRegistered(uid)
	if err != nil {
		e.api.Send(playerID, ColourRed, "Could not check registration status.")
		return
	}
	if registered {
		e.api.Send(playerID, ColourYellow, "You are already registered.")
		e.SendHideRegister(playerID)
		return
	}
	e.promptRegistration(playerID)
}

func (e *Engine) cmdLobby(playerID int) {
	lobby := e.cfg.LobbyPosition(e.mapCfg)
	if lobby.X == 0 && lobby.Y == 0 && lobby.Z == 0 {
		e.api.Send(playerID, ColourRed, "Lobby spawn is not configured.")
		return
	}
	if !e.api.IsSpawned(playerID) {
		if err := e.api.ForceSpawn(playerID); err != nil {
			e.api.Send(playerID, ColourYellow, "Spawn first, then use /lobby.")
			return
		}
	}
	if e.api.PlayerVehicleID(playerID) >= 0 {
		_ = e.api.RemoveFromVehicle(playerID)
	}
	if err := e.api.SetPlayerPosition(playerID, lobby); err != nil {
		e.api.Send(playerID, ColourRed, "Could not warp to lobby.")
		return
	}
	e.api.Send(playerID, ColourGreen, fmt.Sprintf(
		"Warped to lobby: %.2f, %.2f, %.2f",
		lobby.X, lobby.Y, lobby.Z,
	))
}

func (e *Engine) cmdLeaderboard(playerID int) CommandResult {
	e.ToggleLobbyLeaderboard(playerID)
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) cmdScoreboard(playerID int, args []string) CommandResult {
	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "hide", "off", "close":
			e.SendScoreboardHide(playerID)
			e.SendHideRoundStatsTo(playerID)
			e.api.Send(playerID, ColourGreen, "Scoreboard hidden.")
			return CommandResult{Handled: true, Deny: true}
		case "end", "stats", "round":
			winner, reason := e.roundEndWinnerAndReason()
			e.SendRoundEndStatsTo(playerID, winner, reason)
			e.api.Send(playerID, ColourGreen, "Round scoreboard opened (P to close).")
			return CommandResult{Handled: true, Deny: true}
		}
	}

	forceState := scoreboardStateAuto
	if e.round.State == RoundIdle {
		forceState = scoreboardStatePreview
	}
	e.SendScoreboardTo(playerID, forceState)
	e.api.Send(playerID, ColourGreen, "Live scoreboard shown.")
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) cmdGoToPos(playerID int, args []string) CommandResult {
	if !e.isAdmin(playerID) {
		e.api.Send(playerID, ColourRed, "Admin only.")
		return CommandResult{Handled: true, Deny: true}
	}
	if len(args) < 3 {
		e.api.Send(playerID, ColourYellow, "Usage: /gotopos <x> <y> <z>")
		return CommandResult{Handled: true, Deny: true}
	}
	x, err := parseCoord(args[0])
	if err != nil {
		e.api.Send(playerID, ColourYellow, "X must be a number.")
		return CommandResult{Handled: true, Deny: true}
	}
	y, err := parseCoord(args[1])
	if err != nil {
		e.api.Send(playerID, ColourYellow, "Y must be a number.")
		return CommandResult{Handled: true, Deny: true}
	}
	z, err := parseCoord(args[2])
	if err != nil {
		e.api.Send(playerID, ColourYellow, "Z must be a number.")
		return CommandResult{Handled: true, Deny: true}
	}
	if !e.api.IsSpawned(playerID) {
		if err := e.api.ForceSpawn(playerID); err != nil {
			e.api.Send(playerID, ColourYellow, "Spawn first, then use /gotopos.")
			return CommandResult{Handled: true, Deny: true}
		}
	}
	if e.api.PlayerVehicleID(playerID) >= 0 {
		_ = e.api.RemoveFromVehicle(playerID)
	}
	pos := Vec3{X: x, Y: y, Z: z}
	if err := e.api.SetPlayerPosition(playerID, pos); err != nil {
		e.api.Send(playerID, ColourRed, "Could not warp to position.")
		return CommandResult{Handled: true, Deny: true}
	}
	e.api.Send(playerID, ColourGreen, fmt.Sprintf(
		"Warped to %.2f, %.2f, %.2f",
		pos.X, pos.Y, pos.Z,
	))
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) activeVehicleIDs() []int {
	const maxScan = 200
	var ids []int
	for id := 0; id < maxScan; id++ {
		if e.api.VehicleExists(id) {
			ids = append(ids, id)
		}
	}
	return ids
}

func (e *Engine) vehiclesByModel(modelID int) []int {
	var matches []int
	for _, id := range e.activeVehicleIDs() {
		if e.api.VehicleModel(id) == modelID {
			matches = append(matches, id)
		}
	}
	return matches
}

func (e *Engine) vehicleSummary(vehicleID int) string {
	pos := e.api.VehiclePos(vehicleID)
	rot := e.api.VehicleRotationEuler(vehicleID)
	model := e.api.VehicleModel(vehicleID)
	hp := e.api.VehicleHealth(vehicleID)
	return fmt.Sprintf(
		"model %d (slot %d) | HP %.0f | pos %.2f, %.2f, %.2f | rot %.2f, %.2f, %.2f",
		model, vehicleID, hp, pos.X, pos.Y, pos.Z, rot.X, rot.Y, rot.Z,
	)
}

func (e *Engine) activeModelIDs() []int {
	seen := map[int]bool{}
	var models []int
	for _, id := range e.activeVehicleIDs() {
		m := e.api.VehicleModel(id)
		if !seen[m] {
			seen[m] = true
			models = append(models, m)
		}
	}
	return models
}

func (e *Engine) resolveVehicleQuery(arg string) (slots []int, label, errMsg string) {
	modelID, err := strconv.Atoi(strings.TrimSpace(arg))
	if err != nil || modelID < 0 {
		return nil, "", "Invalid model id. Usage: /getvehicle <model id>"
	}
	if matches := e.vehiclesByModel(modelID); len(matches) > 0 {
		return matches, fmt.Sprintf("model %d", modelID), ""
	}
	active := e.activeModelIDs()
	if len(active) == 0 {
		return nil, "", fmt.Sprintf("Model %d not found. No active server vehicles.", modelID)
	}
	return nil, "", fmt.Sprintf("Model %d not found. Active models: %v", modelID, active)
}

func (e *Engine) cmdGetVehicle(playerID int, args []string) CommandResult {
	if !e.isAdmin(playerID) {
		e.api.Send(playerID, ColourRed, "Admin only.")
		return CommandResult{Handled: true, Deny: true}
	}
	if len(args) < 1 {
		ids := e.activeVehicleIDs()
		if len(ids) == 0 {
			e.api.Send(playerID, ColourRed, "No active server vehicles. Start a round or /testhydra first.")
			return CommandResult{Handled: true, Deny: true}
		}
		e.api.Send(playerID, ColourCyan, fmt.Sprintf("Active server vehicles (%d):", len(ids)))
		for _, id := range ids {
			e.api.Send(playerID, ColourWhite, e.vehicleSummary(id))
		}
		e.api.Send(playerID, ColourYellow, "Usage: /getvehicle <model id>")
		return CommandResult{Handled: true, Deny: true}
	}

	slots, label, errMsg := e.resolveVehicleQuery(args[0])
	if errMsg != "" {
		e.api.Send(playerID, ColourRed, errMsg)
		return CommandResult{Handled: true, Deny: true}
	}

	if len(slots) > 1 {
		e.api.Send(playerID, ColourCyan, fmt.Sprintf("Matched %d vehicles for %s:", len(slots), label))
	}
	for _, id := range slots {
		summary := e.vehicleSummary(id)
		if label != "" && len(slots) == 1 {
			summary += " (" + label + ")"
		}
		e.api.Send(playerID, ColourCyan, summary)
	}
	return CommandResult{Handled: true, Deny: true}
}

func parseCoord(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

func (e *Engine) cmdGetPos(playerID int, args []string) {
	targetID := playerID
	label := e.api.PlayerName(playerID)
	if len(args) > 0 {
		name := strings.Join(args, " ")
		targetID = e.api.PlayerIDFromName(name)
		if targetID < 0 || !e.api.IsConnected(targetID) {
			e.api.Send(playerID, ColourRed, name+" was not found.")
			return
		}
		label = e.api.PlayerName(targetID)
	}
	if !e.api.IsSpawned(targetID) {
		e.api.Send(playerID, ColourYellow, label+" is not spawned.")
		return
	}
	pos := e.api.PlayerPosition(targetID)
	e.api.Send(playerID, ColourCyan, fmt.Sprintf(
		"%s position: %.2f, %.2f, %.2f",
		label, pos.X, pos.Y, pos.Z,
	))
}

func (e *Engine) cmdReload(playerID int, args []string) CommandResult {
	if !e.isAdmin(playerID) {
		e.api.Send(playerID, ColourRed, "Admin only.")
		return CommandResult{Handled: true, Deny: true}
	}

	target := "config"
	if len(args) > 0 {
		target = strings.ToLower(args[0])
	}

	switch target {
	case "config", "cfg", "map":
		if err := e.ReloadFromDisk(); err != nil {
			e.api.Send(playerID, ColourRed, "Reload failed: "+err.Error())
			return CommandResult{Handled: true, Deny: true}
		}
		e.api.Broadcast(ColourCyan, "Safari config/map reloaded.")
		e.api.Send(playerID, ColourGreen, "Reloaded safari.json and map from disk.")
	case "scripts", "script", "client":
		e.api.Broadcast(ColourYellow, "Client scripts updated — reconnect to reload.")
		e.reloadClientScripts(playerID)
	case "server", "restart":
		e.api.Broadcast(ColourYellow, "Server reloading — reconnect in a few seconds.")
		e.api.Log("[safari] hot reload requested by " + e.api.PlayerName(playerID))
		e.scheduleServerHotReload()
		e.api.Send(playerID, ColourGreen, "Server restart scheduled.")
	default:
		e.api.Send(playerID, ColourYellow, "Usage: /reload [config|scripts|server]")
	}
	return CommandResult{Handled: true, Deny: true}
}

func (e *Engine) cmdWep(playerID int, args []string) CommandResult {
	if !e.isAdmin(playerID) {
		e.api.Send(playerID, ColourRed, "Admin only.")
		return CommandResult{Handled: true, Deny: true}
	}
	if len(args) < 1 {
		e.api.Send(playerID, ColourYellow, "Usage: /wep <weapon id> [ammo]")
		return CommandResult{Handled: true, Deny: true}
	}
	if !e.api.IsSpawned(playerID) {
		e.api.Send(playerID, ColourYellow, "You must be spawned.")
		return CommandResult{Handled: true, Deny: true}
	}
	weaponID, err := strconv.Atoi(args[0])
	if err != nil || weaponID <= 0 {
		e.api.Send(playerID, ColourYellow, "Weapon id must be a positive number.")
		return CommandResult{Handled: true, Deny: true}
	}
	ammo := 5000
	if len(args) >= 2 {
		ammo, err = strconv.Atoi(args[1])
		if err != nil || ammo < 0 {
			e.api.Send(playerID, ColourYellow, "Ammo must be zero or greater.")
			return CommandResult{Handled: true, Deny: true}
		}
	}
	e.api.GiveWeapon(playerID, weaponID, ammo)
	e.api.Send(playerID, ColourGreen, fmt.Sprintf("Weapon %d granted (%d ammo).", weaponID, ammo))
	return CommandResult{Handled: true, Deny: true}
}
