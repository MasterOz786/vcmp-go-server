package safari

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func (e *Engine) HandleClientScriptData(playerID int, data []byte) {
	if !e.api.IsConnected(playerID) {
		return
	}
	if len(data) < 4 {
		e.api.Log(fmt.Sprintf("[safari-stream] player %d sent empty stream", playerID))
		return
	}
	r := NewStreamReader(data)
	pkt, err := r.ReadInt()
	if err != nil {
		e.api.Log(fmt.Sprintf("[safari-stream] player %d bad packet header: %v", playerID, err))
		return
	}
	switch pkt {
	case PacketHydraCamHello:
		e.markClientScriptReady(playerID)
	case PacketHydraCamCycle:
		e.cycleHydraCamera(playerID)
	case PacketSelectPack:
		pack, err := r.ReadInt()
		if err != nil {
			e.api.Log(fmt.Sprintf("[safari-stream] player %d SELECT_PACK read error: %v", playerID, err))
			return
		}
		e.handleSelectPack(playerID, int(pack))
	case PacketRequestShowPacks:
		e.handleRequestShowPacks(playerID)
	case PacketRequestRegisterUI:
		uid := e.api.PlayerUID(playerID)
		if uid != "" {
			e.maybePromptRegistration(playerID, uid)
		}
	case PacketRegister:
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
	s := NewStreamWriter()
	s.WriteInt(PacketShowRegister)
	if err := e.api.SendScriptData(playerID, s.Bytes()); err != nil {
		e.api.Log(fmt.Sprintf("[safari-stream] SHOW_REGISTER to %d failed: %v", playerID, err))
		return
	}
	e.api.Log(fmt.Sprintf("[safari-stream] SHOW_REGISTER sent to player %d", playerID))
}

func (e *Engine) SendHideRegister(playerID int) {
	if !e.api.IsConnected(playerID) {
		return
	}
	s := NewStreamWriter()
	s.WriteInt(PacketHideRegister)
	if err := e.api.SendScriptData(playerID, s.Bytes()); err != nil {
		e.api.Log(fmt.Sprintf("[safari-stream] HIDE_REGISTER to %d failed: %v", playerID, err))
	}
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
	s := e.teams.session(playerID)
	if s == nil {
		e.ensurePlayerSession(playerID)
		s = e.teams.session(playerID)
	}
	if s != nil {
		s.ClientScriptReady = true
	}
	e.api.Log(fmt.Sprintf("[safari] hydra camera client loaded for player %d (%s)", playerID, e.api.PlayerName(playerID)))
}

func (e *Engine) warnIfNoClientScript(playerID int) {
	s := e.teams.session(playerID)
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
