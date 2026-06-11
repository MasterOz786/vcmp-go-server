//go:build !windows

package safari

import (
	"os/exec"
	"syscall"
)

func configureHotReloadCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
