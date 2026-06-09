package safari

const (
	MaxPlayers = 100

	TeamEscort   = 1
	TeamDefend   = 2
	HydraModel   = 520
	HydraMaxHP   = 1000.0

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
