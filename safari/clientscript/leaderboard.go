package clientscript

import "github.com/masteroz/vcmp-go-server/safari/stream"

type LeaderboardRow struct {
	Name   string
	Points int
	Marks  int
	Wins   int
}

func LobbyLeaderboard(mode int, escort, defend []LeaderboardRow) []byte {
	w := stream.NewWriter()
	w.WriteInt(stream.PacketLobbyLeaderboard)
	w.WriteInt(int32(mode))
	if mode < 0 {
		return w.Bytes()
	}
	writeRows := func(rows []LeaderboardRow) {
		w.WriteInt(int32(len(rows)))
		for _, row := range rows {
			w.WriteString(row.Name)
			w.WriteInt(int32(row.Points))
			w.WriteInt(int32(row.Marks))
			w.WriteInt(int32(row.Wins))
		}
	}
	writeRows(escort)
	writeRows(defend)
	return w.Bytes()
}
