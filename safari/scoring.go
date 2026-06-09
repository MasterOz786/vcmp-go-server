package safari

const (
	PointsCheckpoint   = 50
	PointsMark         = 15
	PointsHydraMinute  = 10
	PointsHydraDestroy = 100
	PointsKill         = 0
)

type Scoring struct {
	EscortScore int
	DefendScore int
}

func (s *Scoring) AddEscort(pts int) {
	s.EscortScore += pts
}

func (s *Scoring) AddDefend(pts int) {
	s.DefendScore += pts
}

func (s *Scoring) WinnerByScore() int {
	if s.EscortScore >= s.DefendScore {
		return TeamEscort
	}
	return TeamDefend
}
