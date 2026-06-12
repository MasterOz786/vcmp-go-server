package admin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Kicker disconnects players so they reconnect with updated client scripts.
type Kicker interface {
	ConnectedIDs() []int
	IsConnected(playerID int) bool
	Kick(playerID int) error
	Log(msg string)
}

func KickConnectedForScriptReload(k Kicker) int {
	n := 0
	for _, id := range k.ConnectedIDs() {
		if !k.IsConnected(id) {
			continue
		}
		if err := k.Kick(id); err != nil {
			k.Log(fmt.Sprintf("[safari] kick failed for player %d: %v", id, err))
			continue
		}
		n++
	}
	return n
}

// HotReloader rebuilds the plugin and restarts the server process.
type HotReloader interface {
	Log(msg string)
	Shutdown()
}

func ScheduleServerHotReload(h HotReloader) {
	wd, err := os.Getwd()
	if err != nil {
		h.Log("[safari] hot reload: getwd failed: " + err.Error())
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
	configureHotReloadCmd(cmd)

	if err := cmd.Start(); err != nil {
		h.Log("[safari] hot reload launcher failed: " + err.Error())
		return
	}

	go func() {
		time.Sleep(750 * time.Millisecond)
		h.Shutdown()
	}()
}
