package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmtrek/air/runner"
)

var (
	cfgPath     string
	debugMode   bool
	showVersion bool
)

func init() {
	debugMode = true
	// flag.StringVar(&cfgPath, "c", "", "config path")
	// flag.BoolVar(&debugMode, "d", false, "debug mode")
	// flag.BoolVar(&showVersion, "v", false, "show version")
	// flag.Parse()
}

func main() {
	fmt.Printf(`
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v%s // live reload for Go apps, with Go%s

`, airVersion, goVersion)

	if showVersion {
		return
	}

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
