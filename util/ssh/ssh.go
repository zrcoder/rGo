package ssh

/*
package ssh implements a set of functions to interact with remote servers using SSH.

Thanks for:
https://github.com/mrrooijen/simplessh
https://github.com/wingedpig/loom
*/

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	envHome         = "HOME"
	defaultKeySufix = "/.ssh/id_rsa"

	netUnix = "unix"
	netTcp  = "tcp"

	defaultPortSufix = ":22"

	envSshAuthSock = "SSH_AUTH_SOCK"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	// The file names of additional key files to use for authentication (~/.ssh/id_rsa is defaulted).
	// RSA (PKCS#1), DSA (OpenSSL), and ECDSA private keys are supported.
	KeyFiles []string `json:"key_files"`
}

func (c Config) String() string {
	return fmt.Sprintf("[%s@%s]", c.User, c.Host)
}

type Client struct {
	sshClient *ssh.Client
}

func NewClient(config Config) (*Client, error) {
	clientConfig := &ssh.ClientConfig{
		User: config.User,
	}
	if config.Password != "" {
		password := ssh.Password(config.Password)
		clientConfig.Auth = append(clientConfig.Auth, password)
	}
	if len(config.KeyFiles) > 0 {
		if keys, err := parseKeys(config.KeyFiles); err == nil {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(keys...))
		}
	} else {
		keyPath := os.Getenv(envHome) + defaultKeySufix
		if key, err := parseKey(keyPath); err == nil {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(key))
		}
	}
	if !strings.Contains(config.Host, ":") {
		config.Host += defaultPortSufix
	}

	sock, err := net.Dial(netUnix, os.Getenv(envSshAuthSock))
	if err == nil {
		signers, err := agent.NewClient(sock).Signers()
		if err == nil {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signers...))
		}
	}
	sshClient, err := ssh.Dial(netTcp, config.Host, clientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{sshClient}, nil
}

func parseKeys(keyPaths []string) ([]ssh.Signer, error) {
	var keys []ssh.Signer
	for _, keyPath := range keyPaths {
		if key, err := parseKey(keyPath); err == nil {
			keys = append(keys, key)
		} else {
			return keys, err
		}
	}
	return keys, nil
}

func parseKey(file string) (ssh.Signer, error) {
	keyBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(keyBytes)
}

func (client *Client) Run(command string) (string, string, error) {
	session, err := client.sshClient.NewSession()
	if err != nil {
		return "", "", nil
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(command)
	return stdout.String(), stderr.String(), err
}

// Close closes the underlying network connection
func (client *Client) Close() {
	client.sshClient.Close()
}
