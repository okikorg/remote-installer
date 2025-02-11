package installer

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHClient wraps ssh.Client with additional functionality
type SSHClient struct {
	*ssh.Client
	config *Config
}

// NewSSHClient creates a new SSH client connection
func NewSSHClient(ip string, port int, username string, config *Config) (*SSHClient, error) {
	key, err := os.ReadFile(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		if err.Error() == "ssh: no key found" {
			// If the error is "no key found", the private key might be in the wrong format
			return nil, fmt.Errorf("invalid private key format. Please ensure you're using a valid SSH private key (RSA, ED25519, etc)")
		}
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote server: %v", err)
	}

	return &SSHClient{
		Client: client,
		config: config,
	}, nil
}
