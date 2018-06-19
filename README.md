# rGo
rGo means "Ready? Go!" or "remote go" ~<br>
It can execute commands(or shell script) on multiple remote hosts automatically

## Usage

### common
Edit config.json to config the remote hosts. For example:<br>
```
[
  {"host":"1.1.1.1"},
  {"host":"1.2.3.4"},
  {"host":"4.3.2.1"}
]
```
ready? let's go <br>
```
./rGo -u root -p -c "cd /tmp; touch abc"
```
and a promtp appears:
```
Please enter the password:
```
after entered the password (which will not be showed),<br>
rGo will login to each host and execute your commands.<br>

### Another example
config.json:
```
[
  {"host":"localhost"},
  {"host":"1.1.1.1"},
  {"host":"1.2.3.4", "user":"abc", "password":"xxx"},
  {"host":"2.2.2.2", "user":"ttt", "password":"vvv"},
  {"host":"4.3.2.1", "user":"kkk", "key_files":["~/.ssh/id_rsa"]}
]
```
If host is 'localhost' or '127.0.0.1', commands will be executed locally<br>
'password' and 'key_files' are optional, just for authentication
```
./rGo -c "echo Hello; ls ~"
```
or you can execute a shell script with the "-sh" option, like:
```
./rGo -sh test.sh
```

## Help info
You can type this line:
```
./rGo --help
```
to get a breif help info:
```
Usage of rGo:
  -c string
	commands (can be separated with ";")
  -sh string
	shell script to be executed
  -t int
	the expected whole duration (second, default 1),
	infact, rGo will execute commands with mutiple threads,
	and the whole time will be the max one for all the threads

  -u string
	if there is no "user" field of some record in config.json,
	value of this option will be used
  -p
	if there is no "paasword" or "key_files" field of some records in config.json,
	we can use this option and then enter the paasword
```
