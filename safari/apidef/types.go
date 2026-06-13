package apidef

const (
	MaxPlayers = 100
	MaxPack    = 3
	MaxSkin    = 4 // spawn-screen class slots per team (indices 0–3)

	TeamEscort = 1
	TeamDefend = 2
	TeamNone   = 255 // VC:MP unassigned — used for lobby skin-only classes

	HydraModel = 6460
	HydraMaxHP = 1000.0

	HydraCamDefault  = 0
	HydraCamChase    = 1
	HydraCamSide     = 2
	HydraCamTactical = 3
	HydraCamCount    = 4
	HydraCamOff      = -1

	ColourWhite  uint32 = 0xFFFFFFFF
	ColourGreen  uint32 = 0xFF00FF00
	ColourRed    uint32 = 0xFFFF4040
	ColourYellow uint32 = 0xFFFFFF00
	ColourCyan   uint32 = 0xFF00FFFF

	PointsCheckpoint   = 50
	PointsMark         = 15
	PointsHydraMinute  = 10
	PointsHydraDestroy = 100
	PointsKill         = 0
)

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

type MapConfig struct {
	LobbySpawn   *Vec3
	HydraStart   Vec3
	HydraAngle   float32
	World        int
	Waypoints    []Vec3
	EscortSpawns []Vec3
	DefendSpawns []Vec3
}

type RoundState int

const (
	RoundIdle RoundState = iota
	RoundActive
	RoundEnded
)

type Scoring struct {
	EscortScore int
	DefendScore int
}

func (s *Scoring) AddEscort(pts int) { s.EscortScore += pts }
func (s *Scoring) AddDefend(pts int)  { s.DefendScore += pts }

func (s *Scoring) WinnerByScore() int {
	if s.EscortScore >= s.DefendScore {
		return TeamEscort
	}
	return TeamDefend
}

type PlayerSession struct {
	Team                int
	Pack                int
	SkinIndex           int // spawn-screen class index 0–3
	SpawnIndex          int
	HasSpawnedThisRound bool
	HydraCameraMode     int
	TestHydraVehicleID  int
	ClientScriptReady   bool
	RoundKills          int
	RoundDeaths         int
	LeaderboardVisible  bool
}

func NewPlayerSession(team, pack int) *PlayerSession {
	if pack < 1 {
		pack = 1
	}
	return &PlayerSession{
		Team:               team,
		Pack:               pack,
		SkinIndex:          0,
		HydraCameraMode:    HydraCamDefault,
		TestHydraVehicleID: -1,
	}
}
