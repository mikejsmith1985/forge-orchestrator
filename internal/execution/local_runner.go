// Package execution provides interfaces and implementations for running shell commands.
// This file implements the LocalRunner, which executes commands on the local machine.
package execution

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// LocalRunner runs commands on the local machine using Go's os/exec package.
// It implements the Executor interface, meaning it can be used anywhere
// an Executor is needed.
type LocalRunner struct{}

// NewLocalRunner creates a new LocalRunner.
// It's a simple constructor, but having it makes the code more consistent
// and allows us to add initialization logic later if needed.
func NewLocalRunner() *LocalRunner {
	return &LocalRunner{}
}

// Execute runs a shell command on the local machine and captures all output.
// It uses the system's default shell (bash on Linux/Mac) to interpret the command.
//
// How it works:
// 1. Create a shell process with the command
// 2. Set up buffers to capture stdout and stderr
// 3. Run the command and wait for it to finish
// 4. Return everything that happened in an ExecutionResult
func (l *LocalRunner) Execute(ctx ExecutionContext) ExecutionResult {
	// Create a context for timeout management.
	// If TimeoutSeconds is 0, we use a background context (no timeout).
	var cmdContext context.Context
	var cancel context.CancelFunc

	if ctx.TimeoutSeconds > 0 {
		// Create a context that will automatically cancel after the timeout.
		// This is like setting a timer - if the command isn't done when
		// the timer goes off, we stop it.
		cmdContext, cancel = context.WithTimeout(
			context.Background(),
			time.Duration(ctx.TimeoutSeconds)*time.Second,
		)
		defer cancel()
	} else {
		cmdContext = context.Background()
	}

	// Create the command using bash to interpret the shell command.
	// We use "bash -c" so that shell features like pipes and redirects work.
	// For example: "echo hello | grep h" needs bash to work correctly.
	cmd := exec.CommandContext(cmdContext, "bash", "-c", ctx.Command)

	// Set the working directory if specified.
	// This is where the command will run from.
	if ctx.WorkingDir != "" {
		cmd.Dir = ctx.WorkingDir
	}

	// Set up environment variables if any were provided.
	// We start with the current environment and add the extras.
	if len(ctx.Environment) > 0 {
		// Get the current environment as a starting point.
		env := cmd.Environ()
		// Add each extra variable.
		for key, value := range ctx.Environment {
			env = append(env, key+"="+value)
		}
		cmd.Env = env
	}

	// Create buffers to capture stdout and stderr.
	// A buffer is like a container that collects text as the command runs.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command and wait for it to finish.
	// This blocks (waits) until the command is done.
	err := cmd.Run()

	// Prepare the result.
	result := ExecutionResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	// Determine the exit code.
	// If the command ran at all, we can get the exit code from it.
	// If it failed to run entirely, we set exit code to -1.
	if err != nil {
		// Try to get the exit code from the error.
		// This works if the command ran but returned a non-zero exit code.
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			// The command failed to run at all (e.g., command not found).
			result.ExitCode = -1
			result.Error = err
		}
	} else {
		// Command ran successfully - exit code is 0.
		result.ExitCode = 0
	}

	return result
}
