package ssh
/*
package ssh implements a set of functions to interact with remote servers using SSH.

Thanks for:
https://github.com/mrrooijen/simplessh
https://github.com/wingedpig/loom
*/

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"net"
	"golang.org/x/crypto/ssh/agent"
	"strings"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

// Config contains ssh and other configuration data needed for all the public functions in loom.
type Config struct {
	// The user name used in SSH connections.
	User string `json:"user"`

	// Password for SSH connections. This is optional. If the user has an ~/.ssh/id_rsa keyfile,
	// that will also be tried. In addition, other key files can be specified.
	Password string `json:"password"`

	// The machine:port to connect to.
	Host string `json:"host"`

	// The file names of additional key files to use for authentication (~/.ssh/id_rsa is defaulted).
	// RSA (PKCS#1), DSA (OpenSSL), and ECDSA private keys are supported.
	KeyFiles []string `json:"key_files"`
}

func (c Config) String() string  {
	return fmt.Sprintf("[%s@%s]", c.User, c.Host)
}

type Client struct {
	SshClient *ssh.Client
}

func NewClient(config Config) (client *Client, err error) {
	clientConfig := &ssh.ClientConfig{
		User: config.User,
	}
	if config.Password != "" {
		password := ssh.Password(config.Password)
		clientConfig.Auth = append(clientConfig.Auth, password)
	}
	if len(config.KeyFiles) == 0 {
		keyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
		key, err := parseKey(keyPath)
		if err == nil {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(key))
		}
	} else {
		keys, err := parseKeys(config.KeyFiles)
		if err == nil {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(keys...))
		}
	}
	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err == nil {
		signers, err := agent.NewClient(sock).Signers()
		if err == nil {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signers...))
		}
	}
	if !strings.Contains(config.Host, ":") {
		config.Host += ":22"
	}
	sshClient, err := ssh.Dial("tcp", config.Host, clientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{sshClient}, nil
}

func (client *Client) Run(command string) (stdoutStr string, stderrStr string, err error) {
	session, err := client.SshClient.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(command)
	stdoutStr = stdout.String()
	stderrStr = stderr.String()
	return
}

func (client *Client) Close() {
	client.SshClient.Close()
}

func parseKeys(keyPaths []string) (keys []ssh.Signer, err error) {
	for _, keyPath := range keyPaths {
		if key, err := parseKey(keyPath); err == nil {
			keys = append(keys, key)
		} else {
			return keys, err
		}
	}
	return
}

// reads in a keyfile containing a private key and parses it.
func parseKey(file string) (key ssh.Signer, err error) {
	keyBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return
	}
	return
}

// Put copies one or more local files to the remote host, using scp. localfiles can
// contain wildcards, and remotefile can be either a directory or a file.
func (client *Client) Put(localfiles string, remotefile string) error {
	files, err := filepath.Glob(localfiles)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("No files match %s", localfiles)
	}
	for _, localfile := range files {
		contents, err := ioutil.ReadFile(localfile)
		if err != nil {
			return err
		}
		// get the local file mode bits
		fi, err := os.Stat(localfile)
		if err != nil {
			return err
		}
		// the file mode bits are the 9 least significant bits of Mode()
		mode := fi.Mode() & 1023

		session, err := client.SshClient.NewSession()
		if err != nil {
			return err
		}
		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Stderr = &stderrBuf

		w, _ := session.StdinPipe()

		_, lfile := filepath.Split(localfile)
		err = session.Start("/usr/bin/scp -qrt " + remotefile)
		if err != nil {
			w.Close()
			session.Close()
			return err
		}
		fmt.Fprintf(w, "C%04o %d %s\n", mode, len(contents), lfile /*remotefile*/)
		w.Write(contents)
		fmt.Fprint(w, "\x00")
		w.Close()

		err = session.Wait()
		if err != nil {
			session.Close()
			return err
		}
		session.Close()
	}
	return nil
}

// Get copies the file from the remote host to the local host, using scp. Wildcards are not currently supported.
func (client *Client) Get(remotefile string, localfile string) error {
	session, err := client.SshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// TODO: Handle wildcards in remotefile

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	w, _ := session.StdinPipe()

	err = session.Start("/usr/bin/scp -qrf " + remotefile)
	if err != nil {
		w.Close()
		return err
	}
	// TODO: better error checking than just firing and forgetting these nulls.
	fmt.Fprintf(w, "\x00")
	fmt.Fprintf(w, "\x00")
	fmt.Fprintf(w, "\x00")
	fmt.Fprintf(w, "\x00")
	fmt.Fprintf(w, "\x00")
	fmt.Fprintf(w, "\x00")

	err = session.Wait()
	if err != nil {
		return err
	}

	stdout := stdoutBuf.String()
	//stderr := stderrBuf.String()

	// first line of stdout contains file information
	fields := strings.SplitN(stdout, "\n", 2)
	mode, _ := strconv.ParseInt(fields[0][1:5], 8, 32)

	// need to generate final local file name
	// localfile could be a directory or a filename
	// if it's a directory, we need to append the remotefile filename
	// if it doesn't exist, we assume file
	var lfile string
	_, rfile := filepath.Split(remotefile)
	l := len(localfile)
	if localfile[l-1] == '/' {
		localfile = localfile[:l-1]
	}
	fi, err := os.Stat(localfile)
	if err != nil || fi.IsDir() == false {
		lfile = localfile
	} else if fi.IsDir() == true {
		lfile = localfile + "/" + rfile
	}
	// there's a trailing 0 in the file that we need to nuke
	l = len(fields[1])
	err = ioutil.WriteFile(lfile, []byte(fields[1][:l-1]), os.FileMode(mode))
	if err != nil {
		return err
	}
	return nil
}
