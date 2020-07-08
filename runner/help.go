package runner

import (
	"fmt"
	"os"
	"strings"
)

// Args todo
type Args struct {
	Debug   string
	Cmd     string
	Verbose bool
	Help    bool
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
		args.Cmd = strings.Join(os.Args[i:], " ")
	} else {
		args.Help = true
	}
	return &args
}

const helpDesc = `
Usage:
	watchdo [options] command arguments......
options:
	-d <item> 	Debug Item
	-h			Help
	-v <level>	Verbose Level
	-s			Keep Silent, without log output
	-l <seconds> loop execute(run command every n seconds)
example:
	watchdo 
`

// Help help
func Help() {
	fmt.Print(helpDesc)
}
