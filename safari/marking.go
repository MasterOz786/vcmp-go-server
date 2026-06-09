package safari

import (
	"fmt"
	"strings"
	"time"
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

func (m *Marking) TryMark(api API, db *DBWorker, senderID int, targetName string, score *Scoring) (bool, string) {
	if api.PlayerTeam(senderID) != TeamEscort {
		return false, "Only Escort team can mark targets."
	}
	if until, ok := m.cooldowns[senderID]; ok && time.Now().Before(until) {
		left := int(time.Until(until).Seconds())
		return false, fmt.Sprintf("Mark on cooldown (%ds left).", left)
	}
	targetName = strings.TrimSpace(targetName)
	if targetName == "" {
		return false, "Usage: /mark [player name]"
	}
	targetID := api.PlayerIDFromName(targetName)
	if targetID < 0 || !api.IsConnected(targetID) {
		return false, "Target not found."
	}
	if api.PlayerTeam(targetID) != TeamDefend {
		return false, "You can only mark Defenders."
	}
	m.cooldowns[senderID] = time.Now().Add(time.Duration(m.cooldownSec) * time.Second)
	score.AddEscort(PointsMark)
	uid := api.PlayerUID(senderID)
	if uid != "" {
		db.Enqueue(addMarkJob{uid: uid})
	}
	name := api.PlayerName(targetID)
	api.Broadcast(ColourGreen, fmt.Sprintf("[ESCORT] Target marked: %s (+15)", name))
	return true, ""
}
