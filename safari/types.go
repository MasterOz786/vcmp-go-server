package safari

const (
	MaxPlayers = 100
	MaxPack    = 3

	TeamEscort   = 1
	TeamDefend   = 2
	HydraModel   = 6460
	HydraMaxHP   = 1000.0

	HydraCamDefault   = 0
	HydraCamChase     = 1
	HydraCamSide      = 2
	HydraCamTactical  = 3
	HydraCamCount     = 4

	ColourWhite  uint32 = 0xFFFFFFFF
	ColourGreen  uint32 = 0xFF00FF00
	ColourRed    uint32 = 0xFFFF4040
	ColourYellow uint32 = 0xFFFFFF00
	ColourCyan   uint32 = 0xFF00FFFF
)

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

type RoundState int

const (
	RoundIdle RoundState = iota
	RoundActive
	RoundEnded
)

type PlayerStats struct {
	UID          string
	EscortPts    int
	DefendPts    int
	Marks        int
	RoundsPlayed int
	RoundsWon    int
}
