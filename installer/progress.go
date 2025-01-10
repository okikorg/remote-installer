package installer

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

// Progress tracks installation progress
type Progress struct {
	TotalSteps    int
	CurrentStep   int
	CurrentPhase  string
	CurrentAction string
	spinner       *spinner.Spinner
}

// NewProgress creates a new progress tracker
func NewProgress(config *Config) *Progress {
	total := len(config.PreInstall) + len(config.Installation) + len(config.PostInstall)
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Initializing..."
	s.Color("cyan")
	s.Start()

	return &Progress{
		TotalSteps:    total,
		CurrentStep:   0,
		CurrentPhase:  "Initializing",
		CurrentAction: "Starting installation",
		spinner:       s,
	}
}

// Update updates the progress state
func (p *Progress) Update(phase, action string) {
	p.CurrentStep++
	p.CurrentPhase = phase
	p.CurrentAction = action

	// Update spinner message
	progress := p.GetProgressBar()
	p.spinner.Suffix = fmt.Sprintf(" %s | %s: %s", progress, phase, action)
}

// GetProgressBar returns a string representation of the progress bar
func (p *Progress) GetProgressBar() string {
	width := 30
	filled := int(float64(p.CurrentStep) / float64(p.TotalSteps) * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := int(float64(p.CurrentStep) / float64(p.TotalSteps) * 100)
	return fmt.Sprintf("[%s] %d%%", bar, percentage)
}

// Stop stops the spinner
func (p *Progress) Stop() {
	if p.spinner != nil {
		p.spinner.Stop()
	}
}

// Success marks the progress as complete with a success message
func (p *Progress) Success() {
	if p.spinner != nil {
		p.spinner.FinalMSG = fmt.Sprintf("✅ Installation completed successfully!\n")
		p.spinner.Stop()
	}
}

// Error marks the progress as failed with an error message
func (p *Progress) Error(err error) {
	if p.spinner != nil {
		p.spinner.FinalMSG = fmt.Sprintf("❌ Installation failed: %v\n", err)
		p.spinner.Stop()
	}
}
