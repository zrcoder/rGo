package ssh
/*
package ssh implements a set of functions to interact with remote servers using SSH.

Thanks for:
https://github.com/mrrooijen/simplessh
https://github.com/wingedpig/loom
*/

import (
	"path/filepath"
	"fmt"
	"io/ioutil"
	"os"
	"bytes"
	"strings"
	"strconv"
)

// Put copies one or more local files to the remote host, using scp. localfiles can
// contain wildcards, and remotefile can be either a directory or a file.
func (client *Client) Put(localfiles string, remotefile string) error {
	files, err := filepath.Glob(localfiles)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no files match %s", localfiles)
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

		session, err := client.sshClient.NewSession()
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
	session, err := client.sshClient.NewSession()
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
