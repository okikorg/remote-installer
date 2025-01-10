package installer

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// CommandOutput represents the result of a command execution
type CommandOutput struct {
	Command  string
	Output   []byte
	ExitCode int
	Error    error
}

// ExecuteCommand executes a single command and returns its output
func (c *SSHClient) ExecuteCommand(ctx context.Context, cmd string) (*CommandOutput, error) {
	session, err := c.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	execCmd := cmd
	if len(c.config.Environment) > 0 {
		envStrings := make([]string, 0, len(c.config.Environment))
		for key, value := range c.config.Environment {
			envStrings = append(envStrings, fmt.Sprintf("%s=%s", key, value))
		}
		execCmd = fmt.Sprintf("export %s && %s", strings.Join(envStrings, " "), cmd)
	}

	output, err := session.CombinedOutput(execCmd)
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode = exitErr.ExitStatus()
		}
	}

	return &CommandOutput{
		Command:  execCmd,
		Output:   output,
		ExitCode: exitCode,
		Error:    err,
	}, nil
}
