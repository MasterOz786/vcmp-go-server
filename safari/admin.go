package safari

import (
	"fmt"
	"strings"
)

func normalizeAdminName(name string) string {
	return strings.TrimSpace(name)
}

func (e *Engine) isAllowlistedAdmin(playerID int) bool {
	name := normalizeAdminName(e.api.PlayerName(playerID))
	for _, allowed := range e.cfg.AdminNames {
		if strings.EqualFold(name, normalizeAdminName(allowed)) {
			return true
		}
	}
	uid := strings.TrimSpace(e.api.PlayerUID(playerID))
	if uid == "" {
		return false
	}
	for _, allowed := range e.cfg.AdminUIDs {
		if uid == strings.TrimSpace(allowed) {
			return true
		}
	}
	return false
}

func (e *Engine) isAdmin(playerID int) bool {
	if e.api.IsAdmin(playerID) {
		return true
	}
	return e.isAllowlistedAdmin(playerID)
}

func (e *Engine) syncAdminPrivileges(playerID int) {
	if !e.isAllowlistedAdmin(playerID) {
		return
	}
	if e.api.IsAdmin(playerID) {
		return
	}
	e.api.SetAdmin(playerID, true)
	e.api.Log(fmt.Sprintf(
		"[safari] granted admin to %q (uid=%s)",
		e.api.PlayerName(playerID),
		e.api.PlayerUID(playerID),
	))
}
