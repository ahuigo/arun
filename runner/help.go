package runner

import (
	"fmt"
	"os"
)

// Args todo
type Args struct {
	Debug      string
	Cmd        []string
	Verbose    bool
	Help       bool
	IgnoreDirs []string
}

func parseArgv() *Args {
	args := Args{Verbose: true}
	argsLen := len(os.Args)
	i := 1
	for ; i < argsLen; i++ {
		v := os.Args[i]
		if v == "" {
			continue
		}
		if v[0] == '-' {
			switch v {
			case "-h":
				args.Help = true
			case "-i":
				i++
				if i < argsLen {
					args.IgnoreDirs = append(args.IgnoreDirs, os.Args[i])
				}
			case "-d":
				i++
				if i < argsLen {
					args.Debug = os.Args[i]
				}
			}
			continue
		} else {
			break
		}
	}
	if i < argsLen {
		args.Cmd = os.Args[i:]
	} else {
		args.Help = true
	}
	return &args
}

const helpDesc = `
Usage:
	arun [options] command arguments......
options:
	-d <item>	Debug Item
	-h			Help
	-i 			Specify Ignore Directory
	-v <level>	Verbose Level
	-s			Keep Silent, without log output
example:
	arun go run main.go
`

// -l <seconds> loop execute(run command every n seconds)

// Help help
func Help() {
	fmt.Print(helpDesc)
}
