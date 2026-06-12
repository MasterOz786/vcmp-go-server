//go:build !windows

package admin

import (
	"os/exec"
	"syscall"
)

func configureHotReloadCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
