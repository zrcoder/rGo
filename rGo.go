package main

import (
	"fmt"
	"encoding/json"
	"time"
	"github.com/DingHub/rGo/util"
	"github.com/DingHub/rGo/util/ssh"
	"github.com/DingHub/rGo/util/cmd"
	"strings"
	"os"
	"io/ioutil"
)

func main() {
	configData, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var configs []ssh.Config
	err = json.Unmarshal(configData, &configs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, config := range configs {
		if config.User == "" {
			config.User = util.Input.User
		}
		if config.Password == "" && len(config.KeyFiles) == 0 {
			config.Password = util.Input.Password
		}
		cmds := util.Input.Cmds
		headMsg := fmt.Sprintf("%s%s", config, cmds)
		if util.Input.Sh != "" {
			shContent, err := ioutil.ReadFile(util.Input.Sh)
			if err != nil {
+				fmt.Println(err)
+				os.Exit(1)
 			}
			headMsg = fmt.Sprintf("%s%s", config, util.Input.Sh)
			cmds = string(shContent)
		}
		if strings.ToLower(config.Host) == "localhost" || config.Host == "127.0.0.1" {
			go localExec(headMsg, cmds)
			continue
		}
		sshClient, err := ssh.NewClient(&config)
		if err != nil {
			fmt.Printf("%s%s\n%s\n", config, util.Input.Cmds, err.Error())
			continue
		}
		go remoteExec(sshClient, headMsg, cmds)
	}
	time.Sleep(time.Second * time.Duration(util.Input.Duration))
}

func localExec(headMsg, cmds string) {
	result, err := cmd.Run(cmds)
	if err != nil {
		fmt.Printf("%s\n%s\n", headMsg, err.Error())
		return
	}
	fmt.Printf("%s\n%s\n", headMsg, result)
}

func remoteExec(sshClient *ssh.Client, headMsg, cmds string) {
	result, stderr, err := sshClient.Run(cmds)
	defer sshClient.Close()
	if err != nil {
		fmt.Printf("%s\n%s\n%s\n", headMsg, stderr, err.Error())
		return
	}
	fmt.Printf("%s\n%s\n", headMsg, result)
}
