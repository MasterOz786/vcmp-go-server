package persist

type PlayerStats struct {
	UID          string
	EscortPts    int
	DefendPts    int
	Marks        int
	RoundsPlayed int
	RoundsWon    int
}

type LeaderboardEntry struct {
	Name   string
	Points int
	Marks  int
	Wins   int
}

type RoundPlayerRecord struct {
	UID        string
	Team       int
	EscortPts  int
	DefendPts  int
	MarksAdded int
}
