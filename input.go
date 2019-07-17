package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zrcoder/rGo/util/security"
)

var Input struct {
	User     string
	Password string
	Cmds     string
	Sh       string
}

func init() {
	const (
		uUsage = "if there is no \"user\" field of some record in config.json,\n" +
			"value of this option will be used"
		pUsage = "if there is no \"paasword\" or \"key_files\" field of some records in config.json,\n" +
			"we can use this option and then enter the paasword"
	)
	var (
		help          = false
		useEnteredPwd = false
	)
	flag.StringVar(&Input.Cmds, "c", "", "commands (can be separated with \";\")")
	flag.StringVar(&Input.Sh, "sh", "", "shell script file to be executed")
	flag.StringVar(&Input.User, "u", "", uUsage)
	flag.BoolVar(&useEnteredPwd, "p", false, pUsage)
	flag.BoolVar(&help, "help", false, "show help info")
	flag.Parse()

	if help || (Input.Cmds == "" && Input.Sh == "") {
		printHelpInfo()
		os.Exit(0)
	}

	if useEnteredPwd {
		prompt := "Please enter the password"
		Input.Password = security.EnterWithPrompt(prompt)
	}
}

const helpInfo = `  ---------------------------------------------------------------------------------------------
  rGo can execute commands(or a shell script) on multiple remote hosts almost at the same time
  We must firstly config the hosts in the file config.json, and then run rGo
  ---------------------------------------------------------------------------------------------

  Usage of rGo:

  -c string
	commands (can be separated with ";")
  -sh string
	shell script file to be executed

  -u string
	if there is no "user" field of some record in config.json,
	value of this flag will be used
  -p
	if there is no "paasword" or "key_files" field of some records in config.json,
	we can use this flag and then enter the paasword
`

func printHelpInfo() {
	fmt.Println(helpInfo)
}
