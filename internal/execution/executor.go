// Package execution provides interfaces and implementations for running shell commands.
// This file defines the Executor interface, which is the main contract for running commands.
package execution

// Executor is the interface that all command runners must implement.
// An interface is like a contract or promise - any type that implements
// this interface promises it can run commands.
//
// Why use an interface?
// 1. We can swap implementations (local, remote, Docker, etc.)
// 2. We can test code that uses Executor without running real commands
// 3. We can add new execution methods without changing existing code
type Executor interface {
	// Execute runs a command and returns the result.
	// It takes an ExecutionContext (what to run and how) and returns
	// an ExecutionResult (what happened when we ran it).
	//
	// Example:
	//   ctx := ExecutionContext{Command: "echo hello"}
	//   result := executor.Execute(ctx)
	//   // result.Stdout will be "hello\n"
	//   // result.ExitCode will be 0
	Execute(ctx ExecutionContext) ExecutionResult
}
