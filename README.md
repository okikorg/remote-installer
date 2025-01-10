# Remote Installer

A Go library for automating remote server installations and configurations via SSH. This package makes it easy to execute a series of commands on remote servers with proper error handling and progress tracking.

## Installation

```bash
go get github.com/okikorg/remote-installer
```

## Quick Start

1. Create a configuration file (`config.yml`):

```yaml
private_key_path: "~/.ssh/id_rsa"

pre_install_commands:
  - "sudo apt-get update"
  - "sudo apt-get upgrade -y"

installation_commands:
  - "sudo apt-get install -y nginx"
  - "sudo apt-get install -y docker.io"

post_install_commands:
  - "sudo systemctl start nginx"
  - "sudo systemctl start docker"

environment_variables:
  DEBIAN_FRONTEND: "noninteractive"
```

2. Use the library in your code:

```go
package main

import (
    "context"
    "log"

    "github.com/okikorg/remote-installer/installer"
)

func main() {
    // Load configuration
    config, err := installer.LoadConfig("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // Create installer
    inst, err := installer.NewInstaller("192.168.1.100", 22, "root", config)
    if err != nil {
        log.Fatal(err)
    }
    defer inst.Close()

    // Run installation
    ctx := context.Background()
    if err := inst.Install(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Features

- SSH-based remote command execution
- Configurable pre-installation, installation, and post-installation commands
- Environment variable support
- Progress tracking
- Context-based cancellation
- Robust error handling
- Automatic reconnection on connection loss

## Requirements

- Go 1.21 or later
- SSH access to target server
- Valid SSH private key

## Configuration Options

The YAML configuration file supports the following options:

- `private_key_path`: Path to SSH private key
- `pre_install_commands`: Commands to run before main installation
- `installation_commands`: Main installation commands
- `post_install_commands`: Commands to run after installation
- `environment_variables`: Environment variables to se
