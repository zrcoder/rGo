package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/zrcoder/rGo/util/ssh"
	"github.com/zrcoder/rGo/util/cmd"
)

const lineSep = "\n"

var logger = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)

func main() {
	configData, err := ioutil.ReadFile("config.json")
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}
	var configs []ssh.Config
	err = json.Unmarshal(configData, &configs)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	for _, config := range configs {
		go func() {
			if config.User == "" {
				config.User = Input.User
			}
			if config.Password == "" && len(config.KeyFiles) == 0 {
				config.Password = Input.Password
			}
			cmds := Input.Cmds
			headMsg := fmt.Sprintf("%s%s:", config, cmds)
			if Input.Sh != "" {
				shContent, err := ioutil.ReadFile(Input.Sh)
				if err != nil {
					logger.Println(err)
					os.Exit(1)
				}
				headMsg = fmt.Sprintf("%s%s:", config, Input.Sh)
				cmds = string(shContent)
			}
			if strings.ToLower(config.Host) == "localhost" || config.Host == "127.0.0.1" {
				localExec(headMsg, cmds)
			} else {
				sshClient, err := ssh.NewClient(config)
				if err != nil {
					logger.Printf("%s\n%s\n", headMsg, err.Error())
				} else {
					remoteExec(sshClient, headMsg, cmds)
				}
			}
		}()
	}
	time.Sleep(time.Second * time.Duration(Input.Duration))
}

func localExec(headMsg, cmds string) {
	stdout, stderr, err := cmd.Run(cmds)
	printResult(headMsg, stdout, stderr, err)
}

func remoteExec(sshClient *ssh.Client, headMsg, cmds string) {
	stdout, stderr, err := sshClient.Run(cmds)
	printResult(headMsg, stdout, stderr, err)
	sshClient.Close()
}

func printResult(headMsg, stdout, stderr string, err error) {
	result := strings.TrimRight(headMsg, lineSep) + lineSep
	if stderr != "" {
		result += "stderr: " + stderr + lineSep
	}
	if err != nil {
		result += "error: " + err.Error() + lineSep
	} else {
		result += strings.TrimRight(stdout, lineSep) + lineSep
	}
	logger.Println(result)
}
