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
	DBPath              string `json:"db_path"`
	MapFile             string `json:"map_file"`
	RoundMinutes        int    `json:"round_minutes"`
	AutoStartPlayers    int    `json:"auto_start_players"`
	MarkCooldownSec     int    `json:"mark_cooldown_sec"`
	StatusBroadcastSec  int    `json:"status_broadcast_sec"`
	WeaponCheckSec      int    `json:"weapon_check_sec"`
	DisableDriveBy      bool   `json:"disable_driveby"`
	FastSwitch          bool   `json:"fast_switch"`
	StuntBike           bool   `json:"stunt_bike"`
	WallGlitch          bool   `json:"wallglitch"`
	DisableHeliBladeDmg bool   `json:"disable_heli_blade_damage"`
	LobbySpawn          *Vec3  `json:"lobby_spawn"`
}

type MapConfig struct {
	HydraStart   Vec3   `json:"hydra_start"`
	HydraAngle   float32 `json:"hydra_angle"`
	World        int    `json:"world"`
	Waypoints    []Vec3 `json:"waypoints"`
	EscortSpawns []Vec3 `json:"escort_spawns"`
	DefendSpawns []Vec3 `json:"defend_spawns"`
}

func DefaultConfig() Config {
	return Config{
		DBPath:             defaultDBPath,
		MapFile:            defaultMapFile,
		RoundMinutes:       15,
		AutoStartPlayers:   2,
		MarkCooldownSec:    30,
		StatusBroadcastSec: 30,
		WeaponCheckSec:     5,
		DisableDriveBy:     true,
		FastSwitch:         true,
		StuntBike:          true,
		WallGlitch:         false,
		DisableHeliBladeDmg: true,
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
	if cfg.StatusBroadcastSec <= 0 {
		cfg.StatusBroadcastSec = 30
	}
	if cfg.WeaponCheckSec <= 0 {
		cfg.WeaponCheckSec = 5
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
	if m.HydraAngle == 0 {
		m.HydraAngle = 90
	}
	return m, nil
}

func (c Config) LobbyPosition(mapCfg MapConfig) Vec3 {
	if c.LobbySpawn != nil {
		return *c.LobbySpawn
	}
	if len(mapCfg.EscortSpawns) > 0 {
		return mapCfg.EscortSpawns[0]
	}
	return Vec3{}
}
