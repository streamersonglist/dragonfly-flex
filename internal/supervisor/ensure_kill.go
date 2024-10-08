// Copied from https://github.com/fly-apps/postgres-flex/blob/master/internal/supervisor/ensure_kill_linux.go
package supervisor

import (
	"os/exec"
	"syscall"
)

func ensureKill(cmd *exec.Cmd) {
	cmd.SysProcAttr.Pdeathsig = syscall.SIGKILL
}
