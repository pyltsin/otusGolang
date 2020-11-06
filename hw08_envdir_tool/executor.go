package main

import (
	"context"
	"os"
	"os/exec"
	"time"
)

const (
	defaultErrorCode = 1
	successCode      = 0
	minimalCmdLen    = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var resultCommand *exec.Cmd
	cmdLen := len(cmd)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	switch {
	case cmdLen == minimalCmdLen:
		resultCommand = exec.CommandContext(ctx, cmd[0]) //nolint:gosec
	default:
		resultCommand = exec.CommandContext(ctx, cmd[0], cmd[1:]...) //nolint:gosec
	}

	resultCommand.Stdout = os.Stdout
	resultCommand.Stdin = os.Stdin
	resultCommand.Stderr = os.Stderr

	for envName, envVal := range env {
		if envVal == "" {
			err := os.Unsetenv(envName)
			if err != nil {
				return defaultErrorCode
			}
		}
		_ = os.Setenv(envName, envVal)
	}
	resultCommand.Env = os.Environ()
	if err := resultCommand.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
		return defaultErrorCode
	}
	return successCode
}
