package safari

import (
	"encoding/json"
	"os"
)

const (
	defaultConfigFile = "safari.json"
	defaultMapFile    = "safari_maps/patrol_default.json"
	defaultDBPath     = "safari.db"
)

type Config struct {
	DBPath           string `json:"db_path"`
	MapFile          string `json:"map_file"`
	RoundMinutes     int    `json:"round_minutes"`
	AutoStartPlayers int    `json:"auto_start_players"`
	MarkCooldownSec  int    `json:"mark_cooldown_sec"`
}

type MapConfig struct {
	HydraStart    Vec3   `json:"hydra_start"`
	HydraAngle    float32 `json:"hydra_angle"`
	World         int    `json:"world"`
	Waypoints     []Vec3 `json:"waypoints"`
	EscortSpawns  []Vec3 `json:"escort_spawns"`
	DefendSpawns  []Vec3 `json:"defend_spawns"`
}

type SpawnPoint struct {
	Pos   Vec3
	Angle float32
}

func DefaultConfig() Config {
	return Config{
		DBPath:           defaultDBPath,
		MapFile:          defaultMapFile,
		RoundMinutes:     15,
		AutoStartPlayers: 2,
		MarkCooldownSec:  30,
	}
}

func LoadConfig() Config {
	cfg := DefaultConfig()
	data, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		return cfg
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig()
	}
	if cfg.DBPath == "" {
		cfg.DBPath = defaultDBPath
	}
	if cfg.MapFile == "" {
		cfg.MapFile = defaultMapFile
	}
	if cfg.RoundMinutes <= 0 {
		cfg.RoundMinutes = 15
	}
	if cfg.AutoStartPlayers <= 0 {
		cfg.AutoStartPlayers = 2
	}
	if cfg.MarkCooldownSec <= 0 {
		cfg.MarkCooldownSec = 30
	}
	return cfg
}

func LoadMap(path string) (MapConfig, error) {
	var m MapConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return m, err
	}
	if m.World == 0 {
		m.World = 0
	}
	if m.HydraAngle == 0 {
		m.HydraAngle = 90
	}
	return m, nil
}
