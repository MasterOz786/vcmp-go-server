package safari

import (
	"fmt"

	"github.com/masteroz/vcmp-go-server/safari/admin"
	"github.com/masteroz/vcmp-go-server/safari/gameplay"
)

func (e *Engine) isAllowlistedAdmin(playerID int) bool {
	return admin.IsAllowlisted(e.api, e.cfg.AdminNames, e.cfg.AdminUIDs, playerID)
}

func (e *Engine) isAdmin(playerID int) bool {
	return admin.IsPrivileged(e.api, e.cfg.AdminNames, e.cfg.AdminUIDs, playerID)
}

func (e *Engine) syncAdminPrivileges(playerID int) {
	admin.SyncAllowlist(e.api, e.cfg.AdminNames, e.cfg.AdminUIDs, playerID)
}

type scriptReloadKicker struct{ e *Engine }

func (k scriptReloadKicker) ConnectedIDs() []int     { return k.e.teams.ConnectedIDs() }
func (k scriptReloadKicker) IsConnected(id int) bool { return k.e.api.IsConnected(id) }
func (k scriptReloadKicker) Kick(id int) error       { return k.e.api.Kick(id) }
func (k scriptReloadKicker) Log(msg string)          { k.e.api.Log(msg) }

type serverHotReloader struct{ e *Engine }

func (h serverHotReloader) Log(msg string) { h.e.api.Log(msg) }
func (h serverHotReloader) Shutdown()    { h.e.api.Shutdown() }

func (e *Engine) ReloadFromDisk() error {
	cfg := LoadConfig()
	mapCfg, err := LoadMap(cfg.MapFile)
	if err != nil {
		return fmt.Errorf("map %s: %w", cfg.MapFile, err)
	}
	e.cfg = cfg
	e.mapCfg = mapCfg
	e.marking = gameplay.NewMarking(cfg.MarkCooldownSec)
	if cfg.AutoStartPlayers > 0 {
		e.autostartEnabled = true
	}
	e.configureServer()
	e.api.Log(fmt.Sprintf("[safari] reloaded config and map (%s)", cfg.MapFile))
	return nil
}

func (e *Engine) reloadClientScripts(adminID int) {
	n := admin.KickConnectedForScriptReload(scriptReloadKicker{e: e})
	e.api.Send(adminID, ColourGreen, fmt.Sprintf("Kicked %d player(s) — reconnect to load updated client scripts.", n))
}

func (e *Engine) scheduleServerHotReload() {
	admin.ScheduleServerHotReload(serverHotReloader{e: e})
}
