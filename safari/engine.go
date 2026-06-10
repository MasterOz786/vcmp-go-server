package safari

import (
	"fmt"
	"time"
)

type Engine struct {
	api        API
	db         *DBWorker
	cfg        Config
	mapCfg     MapConfig
	serverName string
	gameMode   string

	round   *Round
	teams   *Teams
	marking *Marking

	autostartEnabled bool
	tickCount        int
	loadoutRetries   map[int]int

	events chan Event
	stop   chan struct{}
	done   chan struct{}
}

func NewEngine(api API, db *DBWorker, cfg Config, mapCfg MapConfig, serverName, gameMode string) *Engine {
	autostart := cfg.AutoStartPlayers > 0
	return &Engine{
		api:              api,
		db:               db,
		cfg:              cfg,
		mapCfg:           mapCfg,
		serverName:       serverName,
		gameMode:         gameMode,
		round:            NewRound(),
		teams:            NewTeams(),
		marking:          NewMarking(cfg.MarkCooldownSec),
		autostartEnabled: autostart,
		events:           make(chan Event, 256),
		stop:             make(chan struct{}),
		done:             make(chan struct{}),
	}
}

func (e *Engine) Start() {
	go e.loop()
}

func (e *Engine) Stop() {
	close(e.stop)
	<-e.done
}

func (e *Engine) Enqueue(ev Event) {
	select {
	case e.events <- ev:
	default:
		e.api.Log("[safari] event queue full")
	}
}

func (e *Engine) configureServer() {
	e.api.SetServerOption(int(ServerOptionUseClasses), true)
	e.api.SetServerOption(int(ServerOptionJoinMessages), false)
	e.api.SetServerOption(int(ServerOptionDeathMessages), false)
	e.api.SetServerOption(int(ServerOptionDisableDriveBy), e.cfg.DisableDriveBy)
	e.api.SetServerOption(int(ServerOptionFastSwitch), e.cfg.FastSwitch)
	e.api.SetServerOption(int(ServerOptionStuntBike), e.cfg.StuntBike)
	e.api.SetServerOption(int(ServerOptionWallGlitch), e.cfg.WallGlitch)
	e.api.SetServerOption(int(ServerOptionDisableHeliBladeDamage), e.cfg.DisableHeliBladeDmg)
	e.teams.SetupClasses(e.api, e.mapCfg)
}

func (e *Engine) OnServerStart() {
	e.api.SetServerName(e.serverName)
	e.api.SetGameModeText(e.gameMode)
	e.configureServer()
	e.api.Log("[safari] server ready — Project Safari: Hydra Warfare")
	e.announce(MsgServerReady)
}

func (e *Engine) loop() {
	defer close(e.done)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-e.stop:
			return
		case <-ticker.C:
			e.handleEvent(NewTickEvent())
		case ev := <-e.events:
			e.handleEvent(ev)
		}
	}
}

func (e *Engine) handleEvent(ev Event) {
	switch ev.Type {
	case EvTick:
		e.onTick()
	case EvPlayerConnect:
		e.onConnect(ev.PlayerID)
	case EvPlayerDisconnect:
		e.onDisconnect(ev.PlayerID)
	case EvPlayerSpawn:
		e.onSpawn(ev.PlayerID)
	case EvPlayerDeath:
		e.onDeath(ev.PlayerID, ev.KillerID)
	case EvPlayerCommand:
		_ = e.HandleCommand(ev.PlayerID, ev.Command)
	case EvVehicleExplode:
		e.onVehicleExplode(ev.VehicleID)
	case EvRequestSpawn:
		ev.SpawnResult <- e.HandleRequestSpawn(ev.PlayerID)
	case EvClientScriptData:
		e.onClientScriptData(ev.PlayerID, ev.ScriptData)
	case EvVehicleUpdate:
		e.onVehicleUpdate(ev.VehicleID, ev.VehicleUpdateType)
	case EvVehicleRespawn:
		e.onVehicleRespawn(ev.VehicleID)
	case EvPickupPicked:
		e.onPickupPicked(ev.PickupID, ev.PlayerID)
	case EvCheckpointEntered:
		e.onCheckpointEntered(ev.CheckpointID, ev.PlayerID)
	case EvCheckpointExited:
		e.onCheckpointExited(ev.CheckpointID, ev.PlayerID)
	case EvPlayerKeyBind:
		e.onPlayerKeyBind(ev.PlayerID, ev.KeyBindID, ev.KeyBindReleased)
	case EvObjectShot:
		e.onObjectShot(ev.ObjectID, ev.PlayerID, ev.WeaponID)
	case EvObjectTouched:
		e.onObjectTouched(ev.ObjectID, ev.PlayerID)
	case EvPickupRespawn:
		e.onPickupRespawn(ev.PickupID)
	case EvEntityPoolChange:
		e.onEntityPoolChange(ev.EntityPool, ev.ObjectID, ev.EntityDeleted)
	case EvPlayerUpdate:
		e.onPlayerUpdate(ev.PlayerID, ev.VehicleUpdateType)
	case EvPlayerEnterVehicle:
		e.onPlayerEnterVehicle(ev.PlayerID, ev.VehicleID, ev.VehicleSlot)
	case EvPlayerExitVehicle:
		e.onPlayerExitVehicle(ev.PlayerID, ev.VehicleID)
	}
}

// OnServerFrame runs pending loadout retries every server frame (faster than the 1s tick).
func (e *Engine) OnServerFrame() {
	e.processPendingLoadouts()
}

func (e *Engine) onTick() {
	e.tickCount++
	e.processPendingLoadouts()

	if e.round.MaybeReset() {
		e.announce(MsgRoundReset)
	}
	if e.round.State != RoundActive {
		if e.round.State == RoundIdle {
			e.maybeAutoStart()
		}
		return
	}

	if e.round.Paused {
		return
	}

	if ended, winner, reason := e.round.CheckTimer(); ended {
		e.endRound(winner, reason)
		return
	}

	if e.cfg.StatusBroadcastSec > 0 && e.tickCount%e.cfg.StatusBroadcastSec == 0 {
		e.announce(MsgStatusBroadcast, e.formatStatusLine())
	}

	if e.cfg.WeaponCheckSec > 0 && e.tickCount%e.cfg.WeaponCheckSec == 0 {
		e.enforceAllWeapons()
	}

	nowMs := e.api.ServerTimeMs()
	tick := e.round.Hydra.Tick(e.api, &e.round.Score, nowMs)
	if tick.CheckpointMsg != "" {
		e.announce(MsgCheckpoint, tick.CheckpointMsg)
		e.teams.SyncScores(e.api, e.round.Score)
	}
	if tick.HydraMinuteMsg != "" {
		e.announce(MsgHydraMinute, tick.HydraMinuteMsg)
		e.teams.SyncScores(e.api, e.round.Score)
	}
	if tick.DefendWin {
		e.endRound(TeamDefend, "Hydra destroyed!")
		return
	}
	if tick.EscortWin {
		e.endRound(TeamEscort, "Hydra reached its objective!")
	}
}

func (e *Engine) enforceAllWeapons() {
	for i := 0; i < MaxPlayers; i++ {
		if !e.api.IsConnected(i) {
			continue
		}
		team := e.teams.Team(i)
		if team == 0 {
			continue
		}
		EnforceAllowed(e.api, i, team, e.teams.Pack(i))
	}
}

func (e *Engine) maybeAutoStart() {
	if !e.autostartEnabled || e.cfg.AutoStartPlayers <= 0 {
		return
	}
	n := 0
	for i := 0; i < MaxPlayers; i++ {
		if e.api.IsConnected(i) {
			n++
		}
	}
	if n >= e.cfg.AutoStartPlayers {
		e.startRound()
	}
}

func (e *Engine) startRound() {
	if e.round.State == RoundActive {
		return
	}
	if e.round.Start(e.api, e.mapCfg, e.cfg.RoundMinutes, e.teams, e.marking) {
		e.announce(MsgRoundStart, e.round.RoundMinutesStr())
	}
}

func (e *Engine) endRound(winnerTeam int, reason string) {
	e.round.End(e.api, winnerTeam, reason)
	colour := "green"
	teamName := "Escort"
	if winnerTeam == TeamDefend {
		colour = "red"
		teamName = "Defenders"
	}
	e.announce(MsgRoundEnd, colour, teamName, reason,
		strItoa(e.round.Score.EscortScore), strItoa(e.round.Score.DefendScore))
	e.persistRound(winnerTeam)
}

func strItoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [16]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func (e *Engine) persistRound(winnerTeam int) {
	escortN, defendN := 0, 0
	for i := 0; i < MaxPlayers; i++ {
		if !e.api.IsConnected(i) {
			continue
		}
		switch e.teams.Team(i) {
		case TeamEscort:
			escortN++
		case TeamDefend:
			defendN++
		}
	}
	escortShare, defendShare := 0, 0
	if escortN > 0 {
		escortShare = e.round.Score.EscortScore / escortN
	}
	if defendN > 0 {
		defendShare = e.round.Score.DefendScore / defendN
	}

	var records []RoundPlayerRecord
	for i := 0; i < MaxPlayers; i++ {
		if !e.api.IsConnected(i) {
			continue
		}
		uid := e.api.PlayerUID(i)
		if uid == "" {
			continue
		}
		team := e.teams.Team(i)
		rec := RoundPlayerRecord{UID: uid, Team: team}
		if team == TeamEscort {
			rec.EscortPts = escortShare
		} else if team == TeamDefend {
			rec.DefendPts = defendShare
		}
		records = append(records, rec)
	}
	e.db.Enqueue(recordRoundJob{
		players:      records,
		winnerTeam:   winnerTeam,
		escortScore:  e.round.Score.EscortScore,
		defendScore:  e.round.Score.DefendScore,
		survivedSecs: e.round.SurvivedSecs(),
	})
}

func (e *Engine) preferredPack(playerID int) int {
	uid := e.api.PlayerUID(playerID)
	if uid == "" {
		return 1
	}
	pack, err := e.db.GetPreferredPack(uid)
	if err != nil {
		e.api.Log(fmt.Sprintf("[safari] preferred pack lookup failed for %s: %v", uid, err))
		return 1
	}
	if pack < 1 || pack > 2 {
		return 1
	}
	return pack
}

func (e *Engine) ensurePlayerSession(playerID int) int {
	pack := e.preferredPack(playerID)
	if e.teams.Team(playerID) == 0 {
		return e.teams.Assign(e.api, playerID, pack)
	}
	_ = e.teams.SetPack(playerID, pack)
	return e.teams.Team(playerID)
}

func (e *Engine) applyPlayerLoadout(playerID int) {
	team := e.teams.Team(playerID)
	if team == 0 {
		return
	}
	pack := e.teams.Pack(playerID)
	ApplyLoadout(e.api, playerID, team, pack)
	EnforceAllowed(e.api, playerID, team, pack)
}

// SchedulePlayerLoadout applies the saved pack immediately and retries until
// the full kit is present (covers /kill respawns and class weapon races).
func (e *Engine) SchedulePlayerLoadout(playerID int) {
	if e.loadoutRetries == nil {
		e.loadoutRetries = make(map[int]int)
	}
	e.ensurePlayerSession(playerID)
	if e.api.IsSpawned(playerID) {
		e.applyPlayerLoadout(playerID)
	}
	const retries = 5
	if e.loadoutRetries[playerID] < retries {
		e.loadoutRetries[playerID] = retries
	}
}

func (e *Engine) processPendingLoadouts() {
	for playerID, retries := range e.loadoutRetries {
		if !e.api.IsConnected(playerID) {
			delete(e.loadoutRetries, playerID)
			continue
		}
		if !e.api.IsSpawned(playerID) {
			continue
		}
		team := e.teams.Team(playerID)
		if team == 0 {
			continue
		}
		pack := e.teams.Pack(playerID)
		if LoadoutComplete(e.api, playerID, team, pack) {
			delete(e.loadoutRetries, playerID)
			continue
		}
		e.applyPlayerLoadout(playerID)
		if retries <= 1 {
			delete(e.loadoutRetries, playerID)
		} else {
			e.loadoutRetries[playerID] = retries - 1
		}
	}
}

func (e *Engine) onConnect(playerID int) {
	uid := e.api.PlayerUID(playerID)
	name := e.api.PlayerName(playerID)
	if uid != "" {
		e.db.Enqueue(upsertPlayerJob{uid: uid, name: name})
	}
	e.ensurePlayerSession(playerID)
	e.teams.Welcome(e.api, playerID)
	if e.round.State == RoundIdle {
		lobby := e.cfg.LobbyPosition(e.mapCfg)
		if lobby.X != 0 || lobby.Y != 0 || lobby.Z != 0 {
			_ = e.api.SetPlayerPosition(playerID, lobby)
		}
	}
	e.SchedulePlayerLoadout(playerID)
}

func (e *Engine) onDisconnect(playerID int) {
	delete(e.loadoutRetries, playerID)
	e.teams.Remove(playerID)
	e.marking.ClearPlayer(playerID)
}

func (e *Engine) onSpawn(playerID int) {
	e.SchedulePlayerLoadout(playerID)
	e.teams.MarkSpawned(playerID)
	team := e.teams.Team(playerID)
	score := e.round.Score.EscortScore
	if team == TeamDefend {
		score = e.round.Score.DefendScore
	}
	e.api.SetPlayerScore(playerID, score)
}

func (e *Engine) onClientScriptData(playerID int, data []byte) {
	if len(data) == 0 {
		return
	}
	// Reserved for Safari client packet routing (loadout UI, marks, HUD sync).
	_ = playerID
}

func (e *Engine) onVehicleUpdate(vehicleID, updateType int) {
	if e.round.State != RoundActive || vehicleID != e.round.Hydra.VehicleID {
		return
	}
	_ = updateType
}

func (e *Engine) onVehicleRespawn(vehicleID int) {
	if e.round.State != RoundActive || vehicleID != e.round.Hydra.VehicleID {
		return
	}
	e.api.Log("[safari] hydra vehicle respawned")
}

func (e *Engine) HandlePickupPickAttempt(pickupID, playerID int) bool {
	_ = pickupID
	if e.round.State == RoundActive {
		e.SchedulePlayerLoadout(playerID)
	}
	return true
}

func (e *Engine) onPickupPicked(pickupID, playerID int) {
	_ = pickupID
	if e.round.State == RoundActive {
		e.SchedulePlayerLoadout(playerID)
	}
}

func (e *Engine) onCheckpointEntered(checkpointID, playerID int) {
	_ = checkpointID
	_ = playerID
}

func (e *Engine) onCheckpointExited(checkpointID, playerID int) {
	_ = checkpointID
	_ = playerID
}

func (e *Engine) onPlayerKeyBind(playerID, bindID int, released bool) {
	if released {
		return
	}
	switch bindID {
	case 1:
		_ = e.HandleCommand(playerID, "/pack 1")
	case 2:
		_ = e.HandleCommand(playerID, "/pack 2")
	}
}

func (e *Engine) onObjectShot(objectID, playerID, weaponID int) {
	_, _, _ = objectID, playerID, weaponID
}

func (e *Engine) onObjectTouched(objectID, playerID int) {
	_, _ = objectID, playerID
}

func (e *Engine) onPickupRespawn(pickupID int) {
	_ = pickupID
}

func (e *Engine) onEntityPoolChange(pool, entityID int, deleted bool) {
	_, _, _ = pool, entityID, deleted
}

func (e *Engine) onPlayerUpdate(playerID, updateType int) {
	_, _ = playerID, updateType
}

func (e *Engine) onPlayerEnterVehicle(playerID, vehicleID, slot int) {
	_, _, _ = playerID, vehicleID, slot
}

func (e *Engine) onPlayerExitVehicle(playerID, vehicleID int) {
	_, _ = playerID, vehicleID
}

func (e *Engine) HandleEnterVehicleRequest(playerID, vehicleID, slot int) bool {
	_ = playerID
	if e.round.State == RoundActive && vehicleID == e.round.Hydra.VehicleID {
		return e.teams.Team(playerID) == TeamEscort
	}
	return true
}

func (e *Engine) onDeath(playerID, killerID int) {
	e.teams.AdvanceSpawn(playerID)
	e.SchedulePlayerLoadout(playerID)
	victim := e.api.PlayerName(playerID)
	var msg string
	if killerID < 0 || !e.api.IsConnected(killerID) {
		msg = victim + " committed suicide!"
	} else {
		msg = victim + " was killed by " + e.api.PlayerName(killerID) + "!"
	}
	e.announce(MsgKillFeed, msg)
}

func (e *Engine) onVehicleExplode(vehicleID int) {
	if e.round.State != RoundActive {
		return
	}
	if vehicleID == e.round.Hydra.VehicleID {
		e.endRound(TeamDefend, "Hydra exploded!")
	}
}

func (e *Engine) HandleRequestSpawn(playerID int) bool {
	if e.round.State == RoundEnded {
		return false
	}
	return true
}

func (e *Engine) HandleRequestClass(playerID, classIndex int) bool {
	return e.teams.AllowClassRequest(playerID, classIndex, e.round.State == RoundActive)
}

func (e *Engine) HandleCommandSync(playerID int, cmd string) bool {
	res := e.HandleCommand(playerID, cmd)
	return res.Deny
}

func (e *Engine) TogglePause() bool {
	paused := e.round.TogglePause()
	if paused {
		e.announce(MsgPause)
	} else {
		e.announce(MsgUnpause)
	}
	return paused
}

func (e *Engine) SetAutostart(on bool) {
	e.autostartEnabled = on
	if on {
		e.announce(MsgAutostart, "enabled")
	} else {
		e.announce(MsgAutostart, "disabled")
	}
}

func (e *Engine) AutostartEnabled() bool {
	return e.autostartEnabled
}
