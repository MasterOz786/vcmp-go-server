package main

import (
	"fmt"

	"github.com/masteroz/vcmp-go-server/safari"
)

type Plugin struct {
	engine *safari.Engine
	store  *safari.Store
	db     *safari.DBWorker
}

func newPlugin(cfg Config) *Plugin {
	safariCfg := safari.LoadConfig()
	mapCfg, err := safari.LoadMap(safariCfg.MapFile)
	if err != nil {
		bridgeLog(fmt.Sprintf("[safari] map load failed (%s): %v — using defaults", safariCfg.MapFile, err))
		mapCfg = defaultSafariMap()
	}

	store, err := safari.OpenStore(safariCfg.DBPath)
	if err != nil {
		bridgeLog(fmt.Sprintf("[safari] database open failed: %v", err))
		return &Plugin{}
	}

	db := safari.NewDBWorker(store, 128)
	db.Start()

	gameMode := cfg.GameModeText
	if gameMode == "" {
		gameMode = "Project Safari: Hydra Warfare"
	}

	engine := safari.NewEngine(safariBridge{}, db, safariCfg, mapCfg, cfg.ServerName, gameMode)
	engine.Start()

	return &Plugin{engine: engine, store: store, db: db}
}

func (p *Plugin) shutdown() {
	if p.engine != nil {
		p.engine.Stop()
	}
	if p.db != nil {
		p.db.Stop()
	}
	if p.store != nil {
		_ = p.store.Close()
	}
}

func defaultSafariMap() safari.MapConfig {
	return safari.MapConfig{
		HydraStart:   safari.Vec3{X: -974, Y: -106, Z: 11.2},
		HydraAngle:   90,
		World:        0,
		Waypoints:    []safari.Vec3{{X: -900, Y: -50, Z: 11.2}, {X: -800, Y: 0, Z: 11.2}},
		EscortSpawns: []safari.Vec3{{X: -980, Y: -110, Z: 11.2}},
		DefendSpawns: []safari.Vec3{{X: -850, Y: -20, Z: 11.2}},
	}
}
