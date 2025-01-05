package installer

import (
	"fmt"
	"strings"
)

// Progress tracks installation progress
type Progress struct {
	TotalSteps    int
	CurrentStep   int
	CurrentPhase  string
	CurrentAction string
}

// NewProgress creates a new progress tracker
func NewProgress(config *Config) *Progress {
	total := len(config.PreInstall) + len(config.Installation) + len(config.PostInstall)
	return &Progress{
		TotalSteps:    total,
		CurrentStep:   0,
		CurrentPhase:  "Initializing",
		CurrentAction: "Starting installation",
	}
}

// Update updates the progress state
func (p *Progress) Update(phase, action string) {
	p.CurrentStep++
	p.CurrentPhase = phase
	p.CurrentAction = action
}

// GetProgressBar returns a string representation of the progress bar
func (p *Progress) GetProgressBar() string {
	width := 30
	filled := int(float64(p.CurrentStep) / float64(p.TotalSteps) * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := int(float64(p.CurrentStep) / float64(p.TotalSteps) * 100)
	return fmt.Sprintf("[%s] %d%%", bar, percentage)
}
