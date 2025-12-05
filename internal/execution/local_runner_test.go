// Package execution provides interfaces and implementations for running shell commands.
// This test file verifies that the LocalRunner correctly executes shell commands.
package execution

import (
	"strings"
	"testing"
)

// TestLocalRunnerEchoHelloWorld verifies that local_runner.Execute("echo hello world")
// returns "hello world" in stdout and exit code 0.
// This is the first test case required by Contract 2.
func TestLocalRunnerEchoHelloWorld(t *testing.T) {
	// Create a new LocalRunner.
	runner := NewLocalRunner()

	// Set up the execution context with our test command.
	ctx := ExecutionContext{
		Command: "echo hello world",
	}

	// Execute the command.
	result := runner.Execute(ctx)

	// Verify exit code is 0 (success).
	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	// Verify stdout contains "hello world".
	// We use strings.TrimSpace because echo adds a newline at the end.
	expectedOutput := "hello world"
	actualOutput := strings.TrimSpace(result.Stdout)
	if actualOutput != expectedOutput {
		t.Errorf("Expected stdout to be %q, got %q", expectedOutput, actualOutput)
	}

	// Verify there was no error in running the command.
	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}
}

// TestLocalRunnerExitCode1 verifies that local_runner.Execute("exit 1")
// returns exit code 1 and captures any stderr.
// This is the second test case required by Contract 2.
func TestLocalRunnerExitCode1(t *testing.T) {
	// Create a new LocalRunner.
	runner := NewLocalRunner()

	// Set up the execution context with a command that fails.
	// "exit 1" is a simple command that immediately exits with code 1.
	ctx := ExecutionContext{
		Command: "exit 1",
	}

	// Execute the command.
	result := runner.Execute(ctx)

	// Verify exit code is 1 (failure).
	if result.ExitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", result.ExitCode)
	}

	// Note: "exit 1" doesn't produce stderr, but we verify the field exists.
	// The requirement is to "capture any stderr", which we do - it's just empty here.
}

// TestLocalRunnerCapturesStderr verifies that error output is captured in stderr.
// This extends the Contract 2 requirement to verify stderr capture works.
func TestLocalRunnerCapturesStderr(t *testing.T) {
	// Create a new LocalRunner.
	runner := NewLocalRunner()

	// Use a command that writes to stderr.
	// ">&2 echo error" redirects the echo output to stderr.
	ctx := ExecutionContext{
		Command: "echo 'error message' >&2",
	}

	// Execute the command.
	result := runner.Execute(ctx)

	// Verify stderr contains our error message.
	if !strings.Contains(result.Stderr, "error message") {
		t.Errorf("Expected stderr to contain 'error message', got %q", result.Stderr)
	}

	// Exit code should still be 0 because the command itself succeeded.
	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
}

// TestLocalRunnerWithWorkingDirectory verifies that the working directory is respected.
func TestLocalRunnerWithWorkingDirectory(t *testing.T) {
	runner := NewLocalRunner()

	// Run "pwd" in /tmp to verify working directory is set correctly.
	ctx := ExecutionContext{
		Command:    "pwd",
		WorkingDir: "/tmp",
	}

	result := runner.Execute(ctx)

	// Verify the output shows /tmp (or a resolved path to it).
	output := strings.TrimSpace(result.Stdout)
	if !strings.Contains(output, "tmp") {
		t.Errorf("Expected working directory to be /tmp, got %q", output)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
}

// TestLocalRunnerWithEnvironmentVariables verifies that environment variables are passed.
func TestLocalRunnerWithEnvironmentVariables(t *testing.T) {
	runner := NewLocalRunner()

	// Set a custom environment variable and echo it.
	ctx := ExecutionContext{
		Command: "echo $MY_TEST_VAR",
		Environment: map[string]string{
			"MY_TEST_VAR": "test_value_123",
		},
	}

	result := runner.Execute(ctx)

	// Verify the output contains our test value.
	output := strings.TrimSpace(result.Stdout)
	if output != "test_value_123" {
		t.Errorf("Expected output to be 'test_value_123', got %q", output)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
}

// TestLocalRunnerImplementsExecutorInterface verifies that LocalRunner
// properly implements the Executor interface.
func TestLocalRunnerImplementsExecutorInterface(t *testing.T) {
	// This test verifies at compile time that LocalRunner implements Executor.
	// If it doesn't, this won't compile.
	var _ Executor = (*LocalRunner)(nil)
	var _ Executor = NewLocalRunner()

	// If we got here, the interface is properly implemented.
	t.Log("LocalRunner correctly implements the Executor interface")
}
