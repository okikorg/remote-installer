package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"remote-installer/installer"
	"syscall"

	"github.com/fatih/color"
)

func main() {
	ipAddress := flag.String("ip", "", "Remote server IP address")
	port := flag.Int("port", 22, "SSH port number")
	configPath := flag.String("config", "", "Path to configuration YAML file")
	flag.Parse()

	if *ipAddress == "" || *configPath == "" {
		color.Red("Both IP address and config path are required")
		os.Exit(1)
	}

	// Load configuration
	config, err := installer.LoadConfig(*configPath)
	if err != nil {
		color.Red("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Create installer
	inst, err := installer.NewInstaller(*ipAddress, *port, config)
	if err != nil {
		color.Red("Failed to create installer: %v", err)
		os.Exit(1)
	}
	defer inst.Close()

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		color.Yellow("\nReceived interrupt signal, gracefully shutting down...")
		cancel()
	}()

	// Start installation
	color.Cyan("\n🚀 Starting installation process...")
	if err := inst.Install(ctx); err != nil {
		if ctx.Err() != nil {
			color.Yellow("\n⚠️ Installation interrupted")
		} else {
			color.Red("\n❌ Installation failed: %v", err)
		}
		os.Exit(1)
	}

	color.Green("\n✨ Installation completed successfully!")
}
