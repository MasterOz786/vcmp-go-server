package safari

// PlayerSession holds per-player round scratch state (Flag Raids player.setData equivalent).
type PlayerSession struct {
	Team                int
	Pack                int
	SpawnIndex          int
	HasSpawnedThisRound bool
	HydraCameraMode     int
	TestHydraVehicleID  int
	ClientScriptReady   bool
	RoundKills          int
	RoundDeaths         int
	LeaderboardVisible  bool
}

func newPlayerSession(team, pack int) *PlayerSession {
	if pack < 1 {
		pack = 1
	}
	return &PlayerSession{
		Team:                team,
		Pack:                pack,
		HydraCameraMode:     HydraCamDefault,
		TestHydraVehicleID:  -1,
	}
}
