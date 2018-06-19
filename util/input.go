package util

import (
	"flag"
	"github.com/DingHub/rGo/util/security"
	"fmt"
	"os"
)

var Input = struct {
	User     string
	Password string
	Cmds     string
	Sh       string
	Duration int64
}{}

func init() {
	const (
		tUsage = "the whole time expected (second, default 1),\n" +
			"infact, rGo will execute commands with mutiple threads,\n" +
			"and the whole time will be the max one for all the threads"
		uUsage = "if there is no \"user\" field of some record in config.json,\n" +
			"value of this option will be used"
		pUsage = "if there is no \"paasword\" or \"key_files\" field of some records in config.json,\n" +
			"we can use this option and then enter the paasword"
	)
	var (
		h = false
		help = false
		useEnteredPwd = false
	)
	flag.Int64Var(&Input.Duration, "t", 1, tUsage)
	flag.StringVar(&Input.Cmds, "c", "", "commands (can be separated with \";\")")
	flag.StringVar(&Input.Sh, "sh", "", "shell script to be executed")
	flag.StringVar(&Input.User, "u", "", uUsage)
	flag.BoolVar(&useEnteredPwd, "p", false, pUsage)
	flag.BoolVar(&h, "h", false, "show help info")
	flag.BoolVar(&help, "help", false, "show help info")
	flag.Parse()

	if h || help {
		printHelpInfo()
		os.Exit(0)
	}

	if useEnteredPwd {
		prompt := "Please enter the password"
		Input.Password = security.EnterWithPrompt(prompt)
	}
}

const helpInfo = `Usage of rGo:
  -c string
	commands (can be separated with ";")
  -sh string
	shell script to be executed
  -t int
	the whole time expected (second, default 1),
	infact, rGo will execute commands with mutiple threads,
	and the whole time will be the max one for all the threads

  -u string
	if there is no "user" field of some record in config.json,
	value of this option will be used
  -p
	if there is no "paasword" or "key_files" field of some records in config.json,
	we can use this option and then enter the paasword`
func printHelpInfo()  {
	fmt.Println(helpInfo)
}
