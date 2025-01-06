package installer

import (
	"context"
	"fmt"
)

// Installer represents the main installer instance
type Installer struct {
	client   *SSHClient
	config   *Config
	progress *Progress
}

// NewInstaller creates a new installer instance
func NewInstaller(ip string, port int, username string, config *Config) (*Installer, error) {
	client, err := NewSSHClient(ip, port, username, config)
	if err != nil {
		return nil, err
	}

	return &Installer{
		client:   client,
		config:   config,
		progress: NewProgress(config),
	}, nil
}

// Install executes the full installation process
func (i *Installer) Install(ctx context.Context) error {
	phases := []struct {
		name     string
		commands []string
	}{
		{"Pre-installation", i.config.PreInstall},
		{"Installation", i.config.Installation},
		{"Post-installation", i.config.PostInstall},
	}

	for _, phase := range phases {
		if len(phase.commands) == 0 {
			continue
		}

		for _, cmd := range phase.commands {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			i.progress.Update(phase.name, cmd)
			output, err := i.client.ExecuteCommand(ctx, cmd)
			if err != nil {
				return fmt.Errorf("phase '%s' command '%s' failed: %v", phase.name, cmd, err)
			}

			if len(output.Output) > 0 {
				fmt.Printf("Command '%s' output:\n%s\n", cmd, output.Output)
			}
		}
	}

	return nil
}

// Close closes the SSH connection
func (i *Installer) Close() error {
	return i.client.Close()
}

// GetProgress returns the current progress tracker
func (i *Installer) GetProgress() *Progress {
	return i.progress
}
