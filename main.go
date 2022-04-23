package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
    "io/ioutil"
	"path"
	"runtime"
	"strings"


	"github.com/ahuigo/arun/runner"
)

func getVersion() string {
	_, filename, _, _ := runtime.Caller(0)
	versionFile := path.Dir(filename) + "/version"
	version, _ := ioutil.ReadFile(versionFile)
	return strings.TrimSpace(string(version))
}

var (
	cfgPath     string
	debugMode   bool
	goVersion   =  runtime.Version()
	VERSION = getVersion()
)

func init() {
	debugMode = true

}

func main() {
	fmt.Printf(`
     ___
    /   |  _______  ______
   / /| | / ___/ / / / __ \
  / ___ |/ /  / /_/ / / / /
 /_/  |_/_/   \__,_/_/ /_/  %s // live reload for any command, with %s
`, VERSION, goVersion)


	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var err error
	r, err := runner.NewEngine(cfgPath, debugMode)
	if err != nil {
		log.Fatal(err)
		return
	}
	if r.Args.Debug != "" {
		fmt.Println("[debug]", r.Args.Debug)
	}
	if r.Args.Help {
		runner.Help()
		return
	}
	go func() {
		<-sigs
		r.Stop()
	}()

	defer func() {
		if e := recover(); e != nil {
			log.Fatalf("PANIC: %+v", e)
		}
	}()

	r.Run()
}
