//go:build windows
// +build windows

package server

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/UserExistsError/conpty"
)

// startPTY is not used on Windows - we use startPTYWindows instead.
func startPTY(cmd *exec.Cmd) (io.ReadWriteCloser, error) {
	return nil, fmt.Errorf("startPTY not supported on Windows, use startPTYWindows")
}

// startPTYWindows starts a PTY session with a specific shell and arguments on Windows.
func startPTYWindows(shell string, args []string) (io.ReadWriteCloser, error) {
	// Build command line
	commandLine := shell
	if len(args) > 0 {
		commandLine += " " + strings.Join(args, " ")
	}

	cpty, err := conpty.Start(commandLine)
	if err != nil {
		return nil, fmt.Errorf("conpty start failed for %s: %w", commandLine, err)
	}
	return cpty, nil
}

// resizePTY resizes the PTY window.
func resizePTY(ptmx io.ReadWriteCloser, cols, rows uint16) error {
	cpty, ok := ptmx.(*conpty.ConPty)
	if !ok {
		return fmt.Errorf("invalid pty type")
	}
	return cpty.Resize(int(cols), int(rows))
}
