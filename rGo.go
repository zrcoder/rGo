package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/zrcoder/rGo/util"
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
				config.User = util.Input.User
			}
			if config.Password == "" && len(config.KeyFiles) == 0 {
				config.Password = util.Input.Password
			}
			cmds := util.Input.Cmds
			headMsg := fmt.Sprintf("%s%s:", config, cmds)
			if util.Input.Sh != "" {
				shContent, err := ioutil.ReadFile(util.Input.Sh)
				if err != nil {
					logger.Println(err)
					os.Exit(1)
				}
				headMsg = fmt.Sprintf("%s%s:", config, util.Input.Sh)
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
	time.Sleep(time.Second * time.Duration(util.Input.Duration))
}

func localExec(headMsg, cmds string) {
	result, err := cmd.Run(cmds)
	if err != nil {
		logger.Println(headMsg, err.Error())
		return
	}
	logger.Printf("%s\n%s\n", headMsg, result)
}

func remoteExec(sshClient *ssh.Client, headMsg, cmds string) {
	stdout, stderr, err := sshClient.Run(cmds)
	defer sshClient.Close()

	headMsg = strings.TrimRight(headMsg, lineSep) + lineSep
	if stderr != "" {
		headMsg += "stderr: " + stderr + lineSep
	}
	if err != nil {
		headMsg += "error: " + err.Error() + lineSep
	} else {
		headMsg += strings.TrimRight(stdout, lineSep) + lineSep
	}
	logger.Println(headMsg)
}
