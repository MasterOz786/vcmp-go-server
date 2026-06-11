package safari

import "fmt"

type LeaderboardEntry struct {
	Name   string
	Points int
	Marks  int
	Wins   int
}

const (
	LobbyLeaderboardHideOverlay = -1
	LobbyLeaderboardWorldOnly   = 0
	LobbyLeaderboardWithOverlay = 1
)

func (e *Engine) refreshLeaderboardCache() {
	<-e.db.RefreshLeaderboardsAsync()
}

func (e *Engine) buildLobbyLeaderboardPayload(mode int) []byte {
	s := NewStreamWriter()
	s.WriteInt(PacketLobbyLeaderboard)
	s.WriteInt(int32(mode))
	if mode < 0 {
		return s.Bytes()
	}

	escort, defend := e.leaderboardRows()
	s.WriteInt(int32(len(escort)))
	for _, row := range escort {
		s.WriteString(row.Name)
		s.WriteInt(int32(row.Points))
		s.WriteInt(int32(row.Marks))
		s.WriteInt(int32(row.Wins))
	}
	s.WriteInt(int32(len(defend)))
	for _, row := range defend {
		s.WriteString(row.Name)
		s.WriteInt(int32(row.Points))
		s.WriteInt(int32(row.Marks))
		s.WriteInt(int32(row.Wins))
	}
	return s.Bytes()
}

func (e *Engine) leaderboardRows() (escort, defend []LeaderboardEntry) {
	escort, defend, ok := e.db.Leaderboards()
	if !ok {
		go e.refreshLeaderboardCache()
	}
	if len(escort) == 0 && len(defend) == 0 {
		escort, defend = e.fallbackLeaderboardRows()
	}
	return escort, defend
}

func (e *Engine) fallbackLeaderboardRows() (escort, defend []LeaderboardEntry) {
	for _, id := range e.teams.ConnectedIDs() {
		if !e.api.IsConnected(id) {
			continue
		}
		name := e.api.PlayerName(id)
		if name == "" {
			continue
		}
		pts, marks, wins := 0, 0, 0
		if uid := e.api.PlayerUID(id); uid != "" {
			if st, ok := e.db.CachedStats(uid); ok {
				marks = st.Marks
				wins = st.RoundsWon
				switch e.teams.Team(id) {
				case TeamEscort:
					pts = st.EscortPts
				case TeamDefend:
					pts = st.DefendPts
				}
			}
		}
		row := LeaderboardEntry{Name: name, Points: pts, Marks: marks, Wins: wins}
		switch e.teams.Team(id) {
		case TeamEscort:
			escort = append(escort, row)
		case TeamDefend:
			defend = append(defend, row)
		default:
			escort = append(escort, row)
		}
	}
	return escort, defend
}

func (e *Engine) SendLobbyLeaderboard(playerID int, mode int) {
	if !e.api.IsConnected(playerID) {
		return
	}
	_ = e.api.SendScriptData(playerID, e.buildLobbyLeaderboardPayload(mode))
}

func (e *Engine) SendHideLobbyLeaderboard(playerID int) {
	e.SendLobbyLeaderboard(playerID, LobbyLeaderboardHideOverlay)
}

func (e *Engine) BroadcastLobbyLeaderboardWorld() {
	payload := e.buildLobbyLeaderboardPayload(LobbyLeaderboardWorldOnly)
	for _, id := range e.teams.ConnectedIDs() {
		if e.api.IsConnected(id) {
			_ = e.api.SendScriptData(id, payload)
		}
	}
}

func (e *Engine) RefreshLobbyLeaderboard() {
	e.refreshLeaderboardCache()
}

func (e *Engine) ToggleLobbyLeaderboard(playerID int) {
	e.ensurePlayerSession(playerID)
	sess := e.teams.session(playerID)
	if sess == nil {
		return
	}
	if sess.LeaderboardVisible {
		e.SendHideLobbyLeaderboard(playerID)
		sess.LeaderboardVisible = false
		e.api.Send(playerID, ColourGreen, "Leaderboard overlay closed. 3D boards remain at lobby.")
		return
	}
	go e.refreshLeaderboardCache()
	e.SendLobbyLeaderboard(playerID, LobbyLeaderboardWithOverlay)
	sess.LeaderboardVisible = true
	escort, defend := e.leaderboardRows()
	e.api.Send(playerID, ColourGreen, fmt.Sprintf(
		"Leaderboard open — %d escort, %d defend entries.",
		len(escort), len(defend),
	))
}

