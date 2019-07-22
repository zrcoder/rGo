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
	flag.StringVar(&Input.Cmds, "c", "", "")
	flag.StringVar(&Input.Sh, "sh", "", "")
	flag.StringVar(&Input.User, "u", "", "")
	useEnteredPwd := flag.Bool("p", false, "")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, helpInfo)
	}

	flag.Parse()

	if Input.Cmds == "" && Input.Sh == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *useEnteredPwd {
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
