package safari

import (
	"fmt"
	"os"

	"github.com/masteroz/vcmp-go-server/safari/gameplay"
	"github.com/masteroz/vcmp-go-server/safari/persist"
)

// Engine runs all gameplay on the VC:MP callback thread (Flag Raids style).
// Only SQLite I/O uses a background worker with in-memory caches.
type Engine struct {
	api        API
	db         *persist.DBWorker
	cfg        Config
	mapCfg     MapConfig
	serverName string
	gameMode   string

	round   *gameplay.Round
	teams   *gameplay.Teams
	marking *gameplay.Marking

	loadoutRetries    map[int]int
	autostartEnabled  bool
	tickCount         int
	secondsAccum      float32
	hydraCameraAccum  float32
}

func NewEngine(api API, db *persist.DBWorker, cfg Config, mapCfg MapConfig, serverName, gameMode string) *Engine {
	autostart := cfg.AutoStartPlayers > 0
	return &Engine{
		api:              api,
		db:               db,
		cfg:              cfg,
		mapCfg:           mapCfg,
		serverName:       serverName,
		gameMode:         gameMode,
		round:            gameplay.NewRound(),
		teams:            gameplay.NewTeams(),
		marking:          gameplay.NewMarking(cfg.MarkCooldownSec),
		autostartEnabled: autostart,
	}
}

func (e *Engine) configureServer() {
	e.api.SetServerOption(int(gameplay.ServerOptionUseClasses), true)
	e.api.SetServerOption(int(gameplay.ServerOptionJoinMessages), false)
	e.api.SetServerOption(int(gameplay.ServerOptionDeathMessages), false)
	e.api.SetServerOption(int(gameplay.ServerOptionDisableDriveBy), e.cfg.DisableDriveBy)
	e.api.SetServerOption(int(gameplay.ServerOptionFastSwitch), e.cfg.FastSwitch)
	e.api.SetServerOption(int(gameplay.ServerOptionStuntBike), e.cfg.StuntBike)
	e.api.SetServerOption(int(gameplay.ServerOptionWallGlitch), e.cfg.WallGlitch)
	e.api.SetServerOption(int(gameplay.ServerOptionDisableHeliBladeDamage), e.cfg.DisableHeliBladeDmg)
	e.teams.SetupClasses(e.api, e.mapCfg)
	lobby := e.cfg.LobbyPosition(e.mapCfg)
	if lobby.X != 0 || lobby.Y != 0 || lobby.Z != 0 {
		e.api.SetSpawnPos(lobby)
	}
}

func (e *Engine) OnServerStart() {
	e.api.SetServerName(e.serverName)
	e.api.SetGameModeText(e.gameMode)
	e.configureServer()
	if _, err := os.Stat("store/script/main.nut"); err != nil {
		e.api.Log("[safari] WARNING: store/script/main.nut not found — Hydra camera will not work until you add it next to server64.exe")
	} else {
		e.api.Log("[safari] store/script/main.nut found — clients download it on connect (F8 shows load message)")
	}
	if _, err := os.Stat(gameplay.HydraVehicleArchive); err != nil {
		e.api.Log("[safari] WARNING: " + gameplay.HydraVehicleArchive + " not found — Hydra will spawn as a default helicopter until you add the custom vehicle")
	} else {
		e.api.Log("[safari] custom Hydra vehicle found at " + gameplay.HydraVehicleArchive)
	}
	e.api.Log("[safari] server ready — Project Safari: Hydra Warfare (direct callbacks)")
	go func() {
		e.refreshLeaderboardCache()
	}()
	e.announce(MsgServerReady)
}

// OnServerFrame drives throttled Hydra camera updates and the 1 Hz game tick.
func (e *Engine) OnServerFrame(elapsed float32) {
	e.hydraCameraAccum += elapsed
	if e.hydraCameraAccum >= hydraCameraUpdateInterval {
		e.hydraCameraAccum -= hydraCameraUpdateInterval
		e.updateHydraCameras()
	}

	e.secondsAccum += elapsed
	if e.secondsAccum < 1.0 {
		return
	}
	e.secondsAccum = 0
	e.onTick()
}

func (e *Engine) onTick() {
	e.tickCount++
	e.retryPendingLoadouts()
	e.syncPrefetchedPacks()

	if e.round.MaybeReset() {
		e.announce(MsgRoundReset)
		e.BroadcastScoreboard()
		e.BroadcastHideRoundStats()
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

	e.BroadcastScoreboard()

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

func (e *Engine) syncPrefetchedPacks() {
	for _, playerID := range e.teams.ConnectedIDs() {
		if !e.api.IsConnected(playerID) {
			continue
		}
		uid := e.api.PlayerUID(playerID)
		if uid == "" {
			continue
		}
		pack, ok := e.db.CachedPreferredPack(uid)
		if !ok || pack == e.teams.Pack(playerID) {
			continue
		}
		if !e.teams.SetPack(playerID, pack) {
			continue
		}
		if e.api.IsSpawned(playerID) {
			e.applyPlayerLoadout(playerID)
		}
	}
}

func (e *Engine) maybeAutoStart() {
	if !e.autostartEnabled || e.cfg.AutoStartPlayers <= 0 {
		return
	}
	if e.teams.ConnectedCount() >= e.cfg.AutoStartPlayers {
		e.startRound()
	}
}

func (e *Engine) startRound() {
	if e.round.State == RoundActive {
		return
	}
	if e.round.Start(e.api, e.mapCfg, e.cfg.RoundMinutes, e.teams, e.marking, e.hydraModel()) {
		e.announce(MsgRoundStart, e.round.RoundMinutesStr())
		e.BroadcastScoreboard()
		for _, playerID := range e.teams.ConnectedIDs() {
			if e.api.IsConnected(playerID) {
				e.SendShowPacks(playerID)
			}
		}
	}
}

func (e *Engine) endRound(winnerTeam int, reason string) {
	e.round.End(e.api, winnerTeam, reason)
	e.BroadcastRoundEndStats(winnerTeam, reason)
	colour := "green"
	teamName := "Escort"
	if winnerTeam == TeamDefend {
		colour = "red"
		teamName = "Defenders"
	}
	e.announce(MsgRoundEnd, colour, teamName, reason,
		strItoa(e.round.Score.EscortScore), strItoa(e.round.Score.DefendScore))
	e.persistRound(winnerTeam)
	go func() {
		e.refreshLeaderboardCache()
		e.BroadcastLobbyLeaderboardWorld()
	}()
	e.BroadcastScoreboard()
	for _, id := range e.teams.ConnectedIDs() {
		if uid := e.api.PlayerUID(id); uid != "" {
			e.db.InvalidateStats(uid)
			e.db.PrefetchStats(uid)
		}
	}
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
	for _, i := range e.teams.ConnectedIDs() {
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

	var records []persist.RoundPlayerRecord
	for _, i := range e.teams.ConnectedIDs() {
		if !e.api.IsConnected(i) {
			continue
		}
		uid := e.api.PlayerUID(i)
		if uid == "" {
			continue
		}
		team := e.teams.Team(i)
		rec := persist.RoundPlayerRecord{UID: uid, Team: team}
		if team == TeamEscort {
			rec.EscortPts = escortShare
		} else if team == TeamDefend {
			rec.DefendPts = defendShare
		}
		records = append(records, rec)
	}
	e.db.EnqueueRecordRound(
		records,
		winnerTeam,
		e.round.Score.EscortScore,
		e.round.Score.DefendScore,
		e.round.SurvivedSecs(),
	)
}

func (e *Engine) initialPack(uid string) int {
	if uid == "" {
		return 1
	}
	if pack, ok := e.db.CachedPreferredPack(uid); ok {
		if pack < 1 || pack > MaxPack {
			return 1
		}
		return pack
	}
	e.db.PrefetchPreferredPack(uid)
	return 1
}

func (e *Engine) ensurePlayerSession(playerID int) int {
	if e.teams.Team(playerID) != 0 {
		return e.teams.Team(playerID)
	}
	uid := e.api.PlayerUID(playerID)
	pack := e.initialPack(uid)
	return e.teams.Assign(e.api, playerID, pack)
}

func (e *Engine) scheduleLoadoutRetry(playerID int) {
	if e.loadoutRetries == nil {
		e.loadoutRetries = make(map[int]int)
	}
	const retries = 3
	if e.loadoutRetries[playerID] < retries {
		e.loadoutRetries[playerID] = retries
	}
}

func (e *Engine) retryPendingLoadouts() {
	for playerID, retries := range e.loadoutRetries {
		if !e.api.IsConnected(playerID) {
			delete(e.loadoutRetries, playerID)
			continue
		}
		if !e.api.IsSpawned(playerID) {
			continue
		}
		// Never strip/re-grant weapons while seated — causes stutter in vehicles (Hydra).
		if e.api.PlayerVehicleID(playerID) >= 0 {
			continue
		}
		team := e.teams.Team(playerID)
		if team == 0 {
			delete(e.loadoutRetries, playerID)
			continue
		}
		pack := e.teams.Pack(playerID)
		if gameplay.LoadoutComplete(e.api, playerID, team, pack) {
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

func (e *Engine) enforcePlayerWeapons(playerID int) {
	team := e.teams.Team(playerID)
	if team == 0 {
		return
	}
	gameplay.EnforceAllowed(e.api, playerID, team, e.teams.Pack(playerID))
}

func (e *Engine) applyPlayerLoadout(playerID int) {
	team := e.teams.Team(playerID)
	if team == 0 {
		return
	}
	pack := e.teams.Pack(playerID)
	gameplay.ApplyLoadout(e.api, playerID, team, pack)
	gameplay.EnforceAllowed(e.api, playerID, team, pack)
}

func (e *Engine) OnConnect(playerID int) {
	e.teams.TrackConnect(playerID)
	uid := e.api.PlayerUID(playerID)
	name := e.api.PlayerName(playerID)
	if uid != "" {
		e.db.EnqueueUpsertPlayer(uid, name)
		e.db.PrefetchPreferredPack(uid)
		e.db.PrefetchStats(uid)
		e.db.PrefetchRegistered(uid)
	}
	e.ensurePlayerSession(playerID)
	e.syncAdminPrivileges(playerID)
	e.teams.Welcome(e.api, playerID)
	if e.round.State != RoundIdle {
		e.BroadcastScoreboard()
	}
	if e.round.State == RoundIdle {
		lobby := e.cfg.LobbyPosition(e.mapCfg)
		if lobby.X != 0 || lobby.Y != 0 || lobby.Z != 0 {
			_ = e.api.SetPlayerPosition(playerID, lobby)
		}
	}
	mode := LobbyLeaderboardWorldOnly
	if sess := e.teams.Session(playerID); sess != nil && sess.LeaderboardVisible {
		mode = LobbyLeaderboardWithOverlay
	}
	e.SendLobbyLeaderboard(playerID, mode)

}

func (e *Engine) maybePromptRegistration(playerID int, uid string) {
	registered, err := e.db.IsRegistered(uid)
	if err != nil {
		e.api.Log(fmt.Sprintf("[safari] registration lookup failed for %s: %v", uid, err))
	}
	if registered {
		return
	}
	e.promptRegistration(playerID)
}

func (e *Engine) OnDisconnect(playerID int) {
	e.stopTestHydra(playerID)
	delete(e.loadoutRetries, playerID)
	e.teams.TrackDisconnect(playerID)
	e.teams.Remove(playerID)
	e.marking.ClearPlayer(playerID)
}

func (e *Engine) OnSpawn(playerID int) {
	e.ensurePlayerSession(playerID)
	if !e.api.IsSpawned(playerID) {
		return
	}
	if !e.teams.HasSpawnedThisRound(playerID) {
		e.applyPlayerLoadout(playerID)
		e.scheduleLoadoutRetry(playerID)
		e.teams.MarkSpawned(playerID)
	}
	team := e.teams.Team(playerID)
	score := e.round.Score.EscortScore
	if team == TeamDefend {
		score = e.round.Score.DefendScore
	}
	e.api.SetPlayerScore(playerID, score)
}

func (e *Engine) OnDeath(playerID, killerID int) {
	e.teams.AdvanceSpawn(playerID)
	e.scheduleLoadoutRetry(playerID)
	if s := e.teams.Session(playerID); s != nil {
		s.RoundDeaths++
	}
	if killerID >= 0 && e.api.IsConnected(killerID) {
		if ks := e.teams.Session(killerID); ks != nil {
			ks.RoundKills++
		}
	}
	victim := e.api.PlayerName(playerID)
	var msg string
	if killerID < 0 || !e.api.IsConnected(killerID) {
		msg = victim + " committed suicide!"
	} else {
		msg = victim + " was killed by " + e.api.PlayerName(killerID) + "!"
	}
	e.announce(MsgKillFeed, msg)
}

func (e *Engine) OnVehicleExplode(vehicleID int) {
	if e.round.State != RoundActive {
		return
	}
	if vehicleID == e.round.Hydra.VehicleID {
		e.endRound(TeamDefend, "Hydra exploded!")
	}
}

func (e *Engine) OnPlayerKeyBind(playerID, bindID int, released bool) {
	if released {
		return
	}
	switch bindID {
	case 1:
		_ = e.HandleCommand(playerID, "/pack 1")
	case 2:
		_ = e.HandleCommand(playerID, "/pack 2")
	case 3:
		_ = e.HandleCommand(playerID, "/pack 3")
	}
}

func (e *Engine) HandlePickupPickAttempt(pickupID, playerID int) bool {
	_ = pickupID
	if e.round.State == RoundActive {
		e.enforcePlayerWeapons(playerID)
	}
	return true
}

func (e *Engine) OnPickupPicked(pickupID, playerID int) {
	_ = pickupID
	if e.round.State == RoundActive {
		e.enforcePlayerWeapons(playerID)
	}
}

func (e *Engine) OnPlayerEnterVehicle(playerID, vehicleID, slot int) {
	if !e.isHydraVehicle(vehicleID) {
		return
	}
	s := e.teams.Session(playerID)
	if s != nil {
		s.HydraCameraMode = HydraCamDefault
	}
	e.syncHydraCamera(playerID, vehicleID)
	e.api.Send(playerID, ColourCyan, "Hydra ready — press H or /hydraview to cycle camera.")
}

func (e *Engine) OnPlayerExitVehicle(playerID, vehicleID int) {
	if e.isHydraVehicle(vehicleID) {
		e.resetHydraCamera(playerID)
	}
}

func (e *Engine) HandleEnterVehicleRequest(playerID, vehicleID, slot int) bool {
	s := e.teams.Session(playerID)
	if s != nil && s.TestHydraVehicleID == vehicleID {
		return true
	}
	if e.round.State == RoundActive && vehicleID == e.round.Hydra.VehicleID {
		return e.teams.Team(playerID) == TeamEscort
	}
	_ = slot
	return true
}

func (e *Engine) HandleRequestSpawn(playerID int) bool {
	_ = playerID
	if e.round.State == RoundEnded {
		return false
	}
	return true
}

func (e *Engine) HandleRequestClass(playerID, classIndex int) bool {
	return e.teams.AllowClassRequest(playerID, classIndex, e.round.State == RoundActive)
}

func (e *Engine) TogglePause() bool {
	paused := e.round.TogglePause()
	if paused {
		e.announce(MsgPause)
	} else {
		e.announce(MsgUnpause)
	}
	e.BroadcastScoreboard()
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

func (e *Engine) ApplyPack(playerID, pack int) {
	e.ensurePlayerSession(playerID)
	if !e.teams.SetPack(playerID, pack) {
		return
	}
	if uid := e.api.PlayerUID(playerID); uid != "" {
		e.db.SavePreferredPack(uid, pack)
	}
	if e.api.IsSpawned(playerID) {
		e.applyPlayerLoadout(playerID)
	}
}
