package safari

import (
	"github.com/masteroz/vcmp-go-server/safari/apidef"
	"github.com/masteroz/vcmp-go-server/safari/gameplay"
	"github.com/masteroz/vcmp-go-server/safari/persist"
)

type API = apidef.API

const (
	MaxPlayers = apidef.MaxPlayers
	MaxPack    = apidef.MaxPack
	MaxSkin    = apidef.MaxSkin

	TeamEscort = apidef.TeamEscort
	TeamDefend = apidef.TeamDefend

	HydraModel = apidef.HydraModel
	HydraMaxHP = apidef.HydraMaxHP

	HydraCamDefault  = apidef.HydraCamDefault
	HydraCamChase    = apidef.HydraCamChase
	HydraCamSide     = apidef.HydraCamSide
	HydraCamTactical = apidef.HydraCamTactical
	HydraCamCount    = apidef.HydraCamCount
	HydraCamOff      = apidef.HydraCamOff

	ColourWhite  = apidef.ColourWhite
	ColourGreen  = apidef.ColourGreen
	ColourRed    = apidef.ColourRed
	ColourYellow = apidef.ColourYellow
	ColourCyan   = apidef.ColourCyan

	PointsCheckpoint   = apidef.PointsCheckpoint
	PointsMark         = apidef.PointsMark
	PointsHydraMinute  = apidef.PointsHydraMinute
	PointsHydraDestroy = apidef.PointsHydraDestroy
)

type RoundState = apidef.RoundState

const (
	RoundIdle    = apidef.RoundIdle
	RoundActive  = apidef.RoundActive
	RoundEnded   = apidef.RoundEnded
)

type Scoring = apidef.Scoring

type Store = persist.Store
type DBWorker = persist.DBWorker
type PlayerStats = persist.PlayerStats
type LeaderboardEntry = persist.LeaderboardEntry
type RoundPlayerRecord = persist.RoundPlayerRecord

var (
	OpenStore    = persist.OpenStore
	NewDBWorker  = persist.NewDBWorker
	EscortPacks  = gameplay.EscortPacks
	DefendPacks  = gameplay.DefendPacks
)
