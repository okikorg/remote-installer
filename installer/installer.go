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

// Install runs the installation process
func (i *Installer) Install(ctx context.Context) error {
	defer i.progress.Stop()

	phases := []struct {
		name     string
		commands []string
	}{
		{"Pre-installation", i.config.PreInstall},
		{"Installation", i.config.Installation},
		{"Post-installation", i.config.PostInstall},
	}

	for _, phase := range phases {
		for _, cmd := range phase.commands {
			select {
			case <-ctx.Done():
				i.progress.Error(ctx.Err())
				return ctx.Err()
			default:
				i.progress.Update(phase.name, cmd)
				output, err := i.client.ExecuteCommand(ctx, cmd)
				if err != nil {
					i.progress.Error(err)
					return fmt.Errorf("phase '%s' command '%s' failed: %v", phase.name, cmd, err)
				}

				if len(output.Output) > 0 {
					// Only print output if debug is enabled or there's an error
					if i.config.Debug {
						fmt.Printf("\n📝 Command '%s' output:\n%s\n", cmd, output.Output)
					}
				}
			}
		}
	}

	i.progress.Success()
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
