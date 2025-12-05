// Package execution provides interfaces and implementations for running shell commands.
// This package abstracts command execution so we can test without running real commands,
// and potentially add remote execution capabilities in the future.
package execution

// ExecutionContext holds all the information needed to run a command.
// Think of it like filling out a form before you ask someone to do something:
// - What do you want me to do? (Command)
// - Where should I do it? (WorkingDir)
// - How long should I wait? (TimeoutSeconds)
// - What extra info might I need? (Environment)
type ExecutionContext struct {
	// Command is the shell command to execute (e.g., "echo hello world").
	// This is the main instruction we want to run.
	Command string

	// WorkingDir is the directory where the command should run.
	// If empty, the command runs in the current directory.
	// Like telling someone "do this in the kitchen" vs "do this in the garage".
	WorkingDir string

	// TimeoutSeconds is how long to wait for the command to finish.
	// If the command takes longer than this, we stop it.
	// 0 means no timeout (wait forever).
	TimeoutSeconds int

	// Environment is a map of extra environment variables to set.
	// Environment variables are like settings that programs can read.
	// For example: {"DEBUG": "true", "LOG_LEVEL": "verbose"}
	Environment map[string]string
}

// ExecutionResult holds everything that happened when we ran a command.
// It's like a report card for the command:
// - What did it say? (Stdout, Stderr)
// - Did it succeed? (ExitCode)
// - Did something go wrong with running it? (Error)
type ExecutionResult struct {
	// Stdout is the normal output from the command.
	// This is what you'd see printed to the screen on success.
	// For "echo hello", stdout would be "hello".
	Stdout string

	// Stderr is the error output from the command.
	// Programs write error messages and warnings here.
	// Even successful commands might write to stderr.
	Stderr string

	// ExitCode is the numeric code the command returned when it finished.
	// 0 means success (everything went well).
	// Any other number usually means something went wrong.
	// It's like a grade: 0 is perfect, anything else is a problem.
	ExitCode int

	// Error is set if something went wrong with running the command itself.
	// This is different from the command failing - this means we couldn't
	// even run it (like if the command doesn't exist).
	Error error
}
