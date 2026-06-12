package admin

import (
	"fmt"
	"strings"
)

// PlayerInfo reads identity fields used for allowlist checks.
type PlayerInfo interface {
	PlayerName(playerID int) string
	PlayerUID(playerID int) string
	IsAdmin(playerID int) bool
	SetAdmin(playerID int, admin bool)
	Log(msg string)
}

func normalizeName(name string) string {
	return strings.TrimSpace(name)
}

func IsAllowlisted(info PlayerInfo, names, uids []string, playerID int) bool {
	name := normalizeName(info.PlayerName(playerID))
	for _, allowed := range names {
		if strings.EqualFold(name, normalizeName(allowed)) {
			return true
		}
	}
	uid := strings.TrimSpace(info.PlayerUID(playerID))
	if uid == "" {
		return false
	}
	for _, allowed := range uids {
		if uid == strings.TrimSpace(allowed) {
			return true
		}
	}
	return false
}

func IsPrivileged(info PlayerInfo, names, uids []string, playerID int) bool {
	if info.IsAdmin(playerID) {
		return true
	}
	return IsAllowlisted(info, names, uids, playerID)
}

func SyncAllowlist(info PlayerInfo, names, uids []string, playerID int) {
	if !IsAllowlisted(info, names, uids, playerID) {
		return
	}
	if info.IsAdmin(playerID) {
		return
	}
	info.SetAdmin(playerID, true)
	info.Log(fmt.Sprintf(
		"[safari] granted admin to %q (uid=%s)",
		info.PlayerName(playerID),
		info.PlayerUID(playerID),
	))
}
