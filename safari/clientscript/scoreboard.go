package clientscript

import "github.com/masteroz/vcmp-go-server/safari/stream"

type ScoreboardData struct {
	EscortScore int32
	DefendScore int32
	State       int32
	Mins        int32
	Secs        int32
	HydraHP     float32
	CPIdx       int32
	CPTotal     int32
}

func Scoreboard(d ScoreboardData) []byte {
	w := stream.NewWriter()
	w.WriteInt(stream.PacketScoreboard)
	w.WriteInt(d.EscortScore)
	w.WriteInt(d.DefendScore)
	w.WriteInt(d.State)
	w.WriteInt(d.Mins)
	w.WriteInt(d.Secs)
	w.WriteFloat(d.HydraHP)
	w.WriteInt(d.CPIdx)
	w.WriteInt(d.CPTotal)
	return w.Bytes()
}

type RoundEndPlayer struct {
	Name   string
	Team   int
	Points int
	Kills  int
	Deaths int
}

func RoundEndStats(winnerTeam, escortScore, defendScore int, reason string, players []RoundEndPlayer) []byte {
	w := stream.NewWriter()
	w.WriteInt(stream.PacketRoundEndStats)
	w.WriteInt(int32(winnerTeam))
	w.WriteInt(int32(escortScore))
	w.WriteInt(int32(defendScore))
	w.WriteString(reason)
	w.WriteInt(int32(len(players)))
	for _, p := range players {
		w.WriteString(p.Name)
		w.WriteInt(int32(p.Team))
		w.WriteInt(int32(p.Points))
		w.WriteInt(int32(p.Kills))
		w.WriteInt(int32(p.Deaths))
	}
	return w.Bytes()
}

func HideRoundEndStats() []byte {
	w := stream.NewWriter()
	w.WriteInt(stream.PacketRoundEndStats)
	w.WriteInt(-1)
	return w.Bytes()
}
