//go:build !windows
// +build !windows

package server

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// startPTY starts a PTY session on Unix systems.
func startPTY(cmd *exec.Cmd) (io.ReadWriteCloser, error) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, fmt.Errorf("pty.Start failed: %w", err)
	}
	return ptmx, nil
}

// startPTYWindows is not used on Unix - stub for build compatibility.
func startPTYWindows(shell string, args []string) (io.ReadWriteCloser, error) {
	return nil, fmt.Errorf("startPTYWindows not supported on Unix")
}

// resizePTY resizes the PTY window.
func resizePTY(ptmx io.ReadWriteCloser, cols, rows uint16) error {
	f, ok := ptmx.(*os.File)
	if !ok {
		return fmt.Errorf("invalid pty type")
	}
	return pty.Setsize(f, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	})
}
