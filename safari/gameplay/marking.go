package gameplay

import (
	"strconv"
	"strings"
	"time"

	"github.com/masteroz/vcmp-go-server/safari/apidef"
	"github.com/masteroz/vcmp-go-server/safari/persist"
)

type Marking struct {
	cooldownSec int
	cooldowns   map[int]time.Time
}

func NewMarking(cooldownSec int) *Marking {
	return &Marking{
		cooldownSec: cooldownSec,
		cooldowns:   make(map[int]time.Time),
	}
}

func (m *Marking) ClearPlayer(playerID int) {
	delete(m.cooldowns, playerID)
}

func (m *Marking) ResetAll() {
	m.cooldowns = make(map[int]time.Time)
}

func (m *Marking) TryMark(api apidef.API, db *persist.DBWorker, senderID int, targetName string, score *apidef.Scoring) (bool, string) {
	if api.PlayerTeam(senderID) != apidef.TeamEscort {
		return false, "Only Escort team can mark targets."
	}
	if until, ok := m.cooldowns[senderID]; ok && time.Now().Before(until) {
		left := int(time.Until(until).Seconds())
		return false, "Mark on cooldown (" + strconv.Itoa(left) + "s left)."
	}
	targetName = strings.TrimSpace(targetName)
	if targetName == "" {
		return false, "Usage: /mark [player name]"
	}
	targetID := api.PlayerIDFromName(targetName)
	if targetID < 0 || !api.IsConnected(targetID) {
		return false, "Target not found."
	}
	if api.PlayerTeam(targetID) != apidef.TeamDefend {
		return false, "You can only mark Defenders."
	}
	m.cooldowns[senderID] = time.Now().Add(time.Duration(m.cooldownSec) * time.Second)
	score.AddEscort(apidef.PointsMark)
	uid := api.PlayerUID(senderID)
	if uid != "" {
		db.EnqueueMark(uid)
	}
	name := api.PlayerName(targetID)
	return true, name
}
