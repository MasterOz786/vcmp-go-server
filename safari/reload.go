package safari

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func (e *Engine) ReloadFromDisk() error {
	cfg := LoadConfig()
	mapCfg, err := LoadMap(cfg.MapFile)
	if err != nil {
		return fmt.Errorf("map %s: %w", cfg.MapFile, err)
	}
	e.cfg = cfg
	e.mapCfg = mapCfg
	e.marking = NewMarking(cfg.MarkCooldownSec)
	if cfg.AutoStartPlayers > 0 {
		e.autostartEnabled = true
	}
	e.configureServer()
	e.api.Log(fmt.Sprintf("[safari] reloaded config and map (%s)", cfg.MapFile))
	return nil
}

func (e *Engine) reloadClientScripts(adminID int) {
	n := 0
	for _, id := range e.teams.ConnectedIDs() {
		if !e.api.IsConnected(id) {
			continue
		}
		if err := e.api.Kick(id); err != nil {
			e.api.Log(fmt.Sprintf("[safari] kick failed for player %d: %v", id, err))
			continue
		}
		n++
	}
	e.api.Send(adminID, ColourGreen, fmt.Sprintf("Kicked %d player(s) — reconnect to load updated client scripts.", n))
}

func (e *Engine) scheduleServerHotReload() {
	wd, err := os.Getwd()
	if err != nil {
		e.api.Log("[safari] hot reload: getwd failed: " + err.Error())
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		script := filepath.Join(wd, "tools", "hotreload.ps1")
		cmd = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", script)
	default:
		script := filepath.Join(wd, "tools", "hotreload.sh")
		cmd = exec.Command("bash", script)
	}
	cmd.Dir = wd
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr.Setpgid = true
	}

	if err := cmd.Start(); err != nil {
		e.api.Log("[safari] hot reload launcher failed: " + err.Error())
		return
	}

	go func() {
		time.Sleep(750 * time.Millisecond)
		e.api.Shutdown()
	}()
}
