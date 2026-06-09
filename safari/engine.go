package safari

import (
	"strings"
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

	events chan Event
	stop   chan struct{}
	done   chan struct{}
}

func NewEngine(api API, db *DBWorker, cfg Config, mapCfg MapConfig, serverName, gameMode string) *Engine {
	return &Engine{
		api:        api,
		db:         db,
		cfg:        cfg,
		mapCfg:     mapCfg,
		serverName: serverName,
		gameMode:   gameMode,
		round:      NewRound(),
		teams:      NewTeams(),
		marking:    NewMarking(cfg.MarkCooldownSec),
		events:     make(chan Event, 256),
		stop:       make(chan struct{}),
		done:       make(chan struct{}),
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

func (e *Engine) OnServerStart() {
	e.api.SetServerName(e.serverName)
	e.api.SetGameModeText(e.gameMode)
	e.teams.SetupClasses(e.api, e.mapCfg)
	e.api.Log("[safari] server ready — Project Safari: Hydra Warfare")
	e.api.Broadcast(ColourCyan, "Project Safari loaded. /help for commands. Admin: /startsafari")
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
	case EvPlayerCommand:
		_ = e.HandleCommand(ev.PlayerID, ev.Command)
	case EvVehicleExplode:
		e.onVehicleExplode(ev.VehicleID)
	case EvRequestSpawn:
		allow := e.round.State == RoundActive || e.round.State == RoundIdle
		ev.SpawnResult <- allow
	}
}

func (e *Engine) onTick() {
	if e.round.MaybeReset() {
		e.api.Broadcast(ColourWhite, "Ready for next round. /startsafari")
	}
	if e.round.State != RoundActive {
		if e.round.State == RoundIdle {
			e.maybeAutoStart()
		}
		return
	}

	if ended, winner, reason := e.round.CheckTimer(e.api); ended {
		e.endRound(winner, reason)
		return
	}

	nowMs := e.api.ServerTimeMs()
	_, escortWin, defendWin, msg := e.round.Hydra.Tick(e.api, &e.round.Score, nowMs)
	if msg != "" {
		e.api.Broadcast(ColourYellow, msg)
	}
	if defendWin {
		e.endRound(TeamDefend, "Hydra destroyed!")
		return
	}
	if escortWin {
		e.endRound(TeamEscort, "Hydra reached its objective!")
	}
}

func (e *Engine) maybeAutoStart() {
	if e.cfg.AutoStartPlayers <= 0 {
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
	e.round.Start(e.api, e.mapCfg, e.cfg.RoundMinutes)
}

func (e *Engine) endRound(winnerTeam int, reason string) {
	e.round.End(e.api, winnerTeam, reason)
	e.persistRound(winnerTeam)
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

func (e *Engine) onConnect(playerID int) {
	uid := e.api.PlayerUID(playerID)
	name := e.api.PlayerName(playerID)
	if uid != "" {
		e.db.Enqueue(upsertPlayerJob{uid: uid, name: name})
	}
	e.teams.Welcome(e.api, playerID)
}

func (e *Engine) onDisconnect(playerID int) {
	e.teams.Remove(playerID)
	e.marking.ClearPlayer(playerID)
}

func (e *Engine) onSpawn(playerID int) {
	team := e.teams.Team(playerID)
	if team == 0 {
		team = e.teams.Assign(e.api, playerID)
	}
	pack := e.teams.Pack(playerID)
	ApplyLoadout(e.api, playerID, team, pack)
	score := e.round.Score.EscortScore
	if team == TeamDefend {
		score = e.round.Score.DefendScore
	}
	e.api.SetPlayerScore(playerID, score)
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

func (e *Engine) HandleCommandSync(playerID int, cmd string) bool {
	cmd = strings.TrimSpace(cmd)
	if !strings.HasPrefix(cmd, "/") {
		return false
	}
	res := e.HandleCommand(playerID, cmd)
	return res.Deny
}
