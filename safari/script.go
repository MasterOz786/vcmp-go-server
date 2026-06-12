package safari

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/masteroz/vcmp-go-server/safari/clientscript"
	"github.com/masteroz/vcmp-go-server/safari/persist"
	"github.com/masteroz/vcmp-go-server/safari/stream"
)

const (
	scoreboardStateAuto    int32 = -1
	scoreboardStateHidden  int32 = 0
	scoreboardStatePreview int32 = 4

	LobbyLeaderboardHideOverlay = -1
	LobbyLeaderboardWorldOnly   = 0
	LobbyLeaderboardWithOverlay = 1
)

func (e *Engine) sendScriptData(playerID int, payload []byte) error {
	if !e.api.IsConnected(playerID) {
		return nil
	}
	return e.api.SendScriptData(playerID, payload)
}

func (e *Engine) broadcastScript(payload []byte) {
	for _, id := range e.teams.ConnectedIDs() {
		if e.api.IsConnected(id) {
			_ = e.api.SendScriptData(id, payload)
		}
	}
}

func (e *Engine) sendScriptPacket(playerID int, payload []byte, logLabel string) {
	if err := e.sendScriptData(playerID, payload); err != nil && logLabel != "" {
		e.api.Log(fmt.Sprintf("[safari-stream] %s to %d failed: %v", logLabel, playerID, err))
	}
}

func readScriptPacket(data []byte) (*stream.Buf, int32, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("empty stream")
	}
	r := stream.NewReader(data)
	pkt, err := r.ReadInt()
	if err != nil {
		return nil, 0, err
	}
	return r, pkt, nil
}

func (e *Engine) HandleClientScriptData(playerID int, data []byte) {
	if !e.api.IsConnected(playerID) {
		return
	}
	r, pkt, err := readScriptPacket(data)
	if err != nil {
		e.api.Log(fmt.Sprintf("[safari-stream] player %d bad packet: %v", playerID, err))
		return
	}
	switch pkt {
	case stream.PacketHydraCamHello:
		e.markClientScriptReady(playerID)
	case stream.PacketHydraCamCycle:
		e.cycleHydraCamera(playerID)
	case stream.PacketSelectPack:
		pack, err := r.ReadInt()
		if err != nil {
			e.api.Log(fmt.Sprintf("[safari-stream] player %d SELECT_PACK read error: %v", playerID, err))
			return
		}
		e.handleSelectPack(playerID, int(pack))
	case stream.PacketRequestShowPacks:
		e.handleRequestShowPacks(playerID)
	case stream.PacketRequestRegisterUI:
		uid := e.api.PlayerUID(playerID)
		if uid != "" {
			e.maybePromptRegistration(playerID, uid)
		}
	case stream.PacketRegister:
		password, err := r.ReadString()
		if err != nil {
			e.api.Send(playerID, ColourRed, "Registration failed: invalid stream payload.")
			e.api.Log(fmt.Sprintf("[safari-stream] player %d REGISTER read error: %v", playerID, err))
			return
		}
		e.completeRegistration(playerID, password)
	default:
		e.api.Log(fmt.Sprintf("[safari-stream] player %d unhandled packet %d", playerID, pkt))
	}
}

func (e *Engine) SendShowRegister(playerID int) {
	if !e.api.IsConnected(playerID) {
		return
	}
	e.sendScriptPacket(playerID, clientscript.ShowRegister(), "SHOW_REGISTER")
	e.api.Log(fmt.Sprintf("[safari-stream] SHOW_REGISTER sent to player %d", playerID))
}

func (e *Engine) SendHideRegister(playerID int) {
	e.sendScriptPacket(playerID, clientscript.HideRegister(), "HIDE_REGISTER")
}

func (e *Engine) promptRegistration(playerID int) {
	e.api.Send(playerID, ColourCyan, "Please register your account using the window or /register.")
	e.SendShowRegister(playerID)
}

func (e *Engine) completeRegistration(playerID int, password string) {
	password = strings.TrimSpace(password)
	if password == "" {
		e.api.Send(playerID, ColourRed, "Registration failed: password cannot be empty.")
		return
	}
	uid := e.api.PlayerUID(playerID)
	if uid == "" {
		e.api.Send(playerID, ColourRed, "Registration failed: could not read your UID.")
		return
	}
	name := e.api.PlayerName(playerID)
	registered, err := e.db.IsRegistered(uid)
	if err != nil {
		e.api.Log(fmt.Sprintf("[safari-stream] register lookup error for %s: %v", uid, err))
		e.api.Send(playerID, ColourRed, "Registration failed: database error.")
		return
	}
	if registered {
		e.api.Send(playerID, ColourYellow, "This account is already registered.")
		e.SendHideRegister(playerID)
		return
	}
	hash := hashPassword(password)
	if err := e.db.RegisterAccount(uid, name, hash); err != nil {
		e.api.Log(fmt.Sprintf("[safari-stream] register save error for %s: %v", uid, err))
		e.api.Send(playerID, ColourRed, "Registration failed: could not save account.")
		return
	}
	e.SendHideRegister(playerID)
	e.api.Send(playerID, ColourGreen, "Account registered successfully. Welcome to Project Safari!")
	e.api.Log(fmt.Sprintf("[safari-stream] player %d (%s) registered", playerID, name))
}

func (e *Engine) markClientScriptReady(playerID int) {
	if !e.api.IsConnected(playerID) {
		return
	}
	s := e.teams.Session(playerID)
	if s == nil {
		e.ensurePlayerSession(playerID)
		s = e.teams.Session(playerID)
	}
	if s != nil {
		s.ClientScriptReady = true
	}
	e.api.Log(fmt.Sprintf("[safari] hydra camera client loaded for player %d (%s)", playerID, e.api.PlayerName(playerID)))
}

func (e *Engine) warnIfNoClientScript(playerID int) {
	s := e.teams.Session(playerID)
	if s != nil && s.ClientScriptReady {
		return
	}
	e.api.Send(playerID, ColourYellow,
		"Hydra camera needs the client script (store/script/main.nut). Reconnect after store sync; press F8 and look for [safari] hydra camera client loaded.")
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func (e *Engine) SendShowPacks(playerID int) {
	if !e.api.IsConnected(playerID) {
		return
	}
	team := e.teams.Team(playerID)
	if team == 0 {
		e.ensurePlayerSession(playerID)
		team = e.teams.Team(playerID)
	}
	payload := clientscript.ShowPacks(team, e.teams.Pack(playerID))
	e.sendScriptPacket(playerID, payload, "SHOW_PACKS")
}

func (e *Engine) SendHidePacks(playerID int) {
	e.sendScriptData(playerID, clientscript.HidePacks())
}

func (e *Engine) SendPackFeedback(playerID int, message string) {
	e.sendScriptData(playerID, clientscript.PackFeedback(message))
}

func (e *Engine) handleSelectPack(playerID, pack int) {
	if pack < 1 || pack > MaxPack {
		e.SendPackFeedback(playerID, fmt.Sprintf("Pack must be 1 to %d.", MaxPack))
		return
	}
	if e.round.State == RoundActive && e.teams.HasSpawnedThisRound(playerID) {
		e.SendPackFeedback(playerID, "Cannot change pack after spawning this round.")
		return
	}
	e.ensurePlayerSession(playerID)
	if e.teams.Team(playerID) == 0 {
		e.SendPackFeedback(playerID, "You are not assigned to a team yet.")
		return
	}
	e.ApplyPack(playerID, pack)
	e.SendPackFeedback(playerID, "")
	team := e.teams.Team(playerID)
	var name string
	if team == TeamEscort {
		name = EscortPacks()[pack].Name
	} else {
		name = DefendPacks()[pack].Name
	}
	e.api.Send(playerID, ColourGreen, fmt.Sprintf("Loadout equipped: %s", name))
}

func (e *Engine) handleRequestShowPacks(playerID int) {
	if e.round.State == RoundActive && e.teams.HasSpawnedThisRound(playerID) {
		msg := "Cannot change pack after spawning this round."
		e.api.Send(playerID, ColourYellow, msg)
		e.SendPackFeedback(playerID, msg)
		return
	}
	e.SendShowPacks(playerID)
}

func (e *Engine) scoreboardState(forceState int32) int32 {
	if forceState >= 0 {
		return forceState
	}
	switch e.round.State {
	case RoundActive:
		if e.round.Paused {
			return 3
		}
		return 1
	case RoundEnded:
		return 2
	default:
		return 0
	}
}

func (e *Engine) buildScoreboardPacket(forceState int32) []byte {
	state := e.scoreboardState(forceState)
	left := e.round.TimeLeft()
	mins := int32(left / time.Minute)
	secs := int32(left.Seconds()) % 60

	hydraHP := float32(0)
	cpIdx, cpTotal := int32(0), int32(0)
	if e.round.State == RoundActive && e.round.Hydra.VehicleID >= 0 {
		hydraHP = e.round.Hydra.Health(e.api)
		cpIdx = int32(e.round.Hydra.Index)
		cpTotal = int32(len(e.round.Hydra.Waypoints))
	}

	return clientscript.Scoreboard(clientscript.ScoreboardData{
		EscortScore: int32(e.round.Score.EscortScore),
		DefendScore: int32(e.round.Score.DefendScore),
		State:       state,
		Mins:        mins,
		Secs:        secs,
		HydraHP:     hydraHP,
		CPIdx:       cpIdx,
		CPTotal:     cpTotal,
	})
}

func (e *Engine) SendScoreboardTo(playerID int, forceState int32) {
	e.sendScriptData(playerID, e.buildScoreboardPacket(forceState))
}

func (e *Engine) SendScoreboardHide(playerID int) {
	e.SendScoreboardTo(playerID, scoreboardStateHidden)
}

func (e *Engine) BroadcastScoreboard() {
	e.broadcastScript(e.buildScoreboardPacket(scoreboardStateAuto))
}

func (e *Engine) roundEndPlayers() []clientscript.RoundEndPlayer {
	escortN, defendN := 0, 0
	for _, id := range e.teams.ConnectedIDs() {
		if !e.api.IsConnected(id) {
			continue
		}
		switch e.teams.Team(id) {
		case TeamEscort:
			escortN++
		case TeamDefend:
			defendN++
		}
	}
	escortShare, defendShare := 0, 0
	if escortN > 0 {
		escortShare = e.round.Score.EscortScore / escortN
	}
	if defendN > 0 {
		defendShare = e.round.Score.DefendScore / defendN
	}

	var players []clientscript.RoundEndPlayer
	for _, id := range e.teams.ConnectedIDs() {
		if !e.api.IsConnected(id) {
			continue
		}
		s := e.teams.Session(id)
		if s == nil {
			continue
		}
		pts := escortShare
		if s.Team == TeamDefend {
			pts = defendShare
		}
		players = append(players, clientscript.RoundEndPlayer{
			Name:   e.api.PlayerName(id),
			Team:   s.Team,
			Points: pts,
			Kills:  s.RoundKills,
			Deaths: s.RoundDeaths,
		})
	}
	return players
}

func (e *Engine) buildRoundEndStatsPayload(winnerTeam int, reason string) []byte {
	return clientscript.RoundEndStats(
		winnerTeam,
		e.round.Score.EscortScore,
		e.round.Score.DefendScore,
		reason,
		e.roundEndPlayers(),
	)
}

func (e *Engine) roundEndWinnerAndReason() (int, string) {
	switch e.round.State {
	case RoundActive:
		return e.round.Score.WinnerByScore(), "Current standings"
	case RoundEnded:
		if e.round.WinnerTeam != 0 {
			return e.round.WinnerTeam, e.round.EndReason
		}
		return e.round.Score.WinnerByScore(), e.round.EndReason
	default:
		return TeamEscort, "Waiting for round start"
	}
}

func (e *Engine) SendRoundEndStatsTo(playerID int, winnerTeam int, reason string) {
	e.sendScriptData(playerID, e.buildRoundEndStatsPayload(winnerTeam, reason))
}

func (e *Engine) SendHideRoundStatsTo(playerID int) {
	e.sendScriptData(playerID, clientscript.HideRoundEndStats())
}

func (e *Engine) BroadcastRoundEndStats(winnerTeam int, reason string) {
	e.broadcastScript(e.buildRoundEndStatsPayload(winnerTeam, reason))
}

func (e *Engine) BroadcastHideRoundStats() {
	e.broadcastScript(clientscript.HideRoundEndStats())
}

func (e *Engine) refreshLeaderboardCache() {
	<-e.db.RefreshLeaderboardsAsync()
}

func toLeaderboardRows(rows []persist.LeaderboardEntry) []clientscript.LeaderboardRow {
	out := make([]clientscript.LeaderboardRow, len(rows))
	for i, row := range rows {
		out[i] = clientscript.LeaderboardRow{
			Name:   row.Name,
			Points: row.Points,
			Marks:  row.Marks,
			Wins:   row.Wins,
		}
	}
	return out
}

func (e *Engine) buildLobbyLeaderboardPayload(mode int) []byte {
	if mode < 0 {
		return clientscript.LobbyLeaderboard(mode, nil, nil)
	}
	escort, defend := e.leaderboardRows()
	return clientscript.LobbyLeaderboard(mode, toLeaderboardRows(escort), toLeaderboardRows(defend))
}

func (e *Engine) leaderboardRows() (escort, defend []persist.LeaderboardEntry) {
	escort, defend, ok := e.db.Leaderboards()
	if !ok {
		go e.refreshLeaderboardCache()
	}
	if len(escort) == 0 && len(defend) == 0 {
		escort, defend = e.fallbackLeaderboardRows()
	}
	return escort, defend
}

func (e *Engine) fallbackLeaderboardRows() (escort, defend []persist.LeaderboardEntry) {
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
		row := persist.LeaderboardEntry{Name: name, Points: pts, Marks: marks, Wins: wins}
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
	e.sendScriptData(playerID, e.buildLobbyLeaderboardPayload(mode))
}

func (e *Engine) SendHideLobbyLeaderboard(playerID int) {
	e.SendLobbyLeaderboard(playerID, LobbyLeaderboardHideOverlay)
}

func (e *Engine) BroadcastLobbyLeaderboardWorld() {
	e.broadcastScript(e.buildLobbyLeaderboardPayload(LobbyLeaderboardWorldOnly))
}

func (e *Engine) RefreshLobbyLeaderboard() {
	e.refreshLeaderboardCache()
}

func (e *Engine) ToggleLobbyLeaderboard(playerID int) {
	e.ensurePlayerSession(playerID)
	sess := e.teams.Session(playerID)
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
