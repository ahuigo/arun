package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/crypto/ssh/terminal"
)

// Engine ...
type Engine struct {
	config    *config
	Args      *Args
	logger    *logger
	watcher   *fsnotify.Watcher
	debugMode bool

	eventCh       chan string
	watcherStopCh chan bool
	binStopCh     chan bool
	binStopEndCh  chan bool
	exitCh        chan bool

	mu         sync.RWMutex
	binRunning bool
	watchers   uint

	ll sync.Mutex // lock for logger
}

// NewEngine ...
func NewEngine(cfgPath string, debugMode bool) (*Engine, error) {
	var err error
	cfg, err := initConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	args := parseArgv()
	cfg.Build.Bin = args.Cmd
	fmt.Printf("%#v", args.Cmd)

	logger := newLogger(cfg)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Engine{
		config:        cfg,
		Args:          args,
		logger:        logger,
		watcher:       watcher,
		debugMode:     debugMode,
		eventCh:       make(chan string, 1000),
		watcherStopCh: make(chan bool, 10),
		binStopCh:     make(chan bool),
		binStopEndCh:  make(chan bool, 1),
		exitCh:        make(chan bool),
		binRunning:    false,
		watchers:      0,
	}, nil
}

// Run run run
func (e *Engine) Run() {
	e.mainDebug("CWD: %s", e.config.Root)
	e.runnerLog("root: %s", e.config.Root)

	var err error
	if err = e.checkRunEnv(); err != nil {
		os.Exit(1)
	}
	if err = e.watching(e.config.Root); err != nil {
		os.Exit(1)
	}

	e.start()
	e.cleanup()
}

func (e *Engine) checkRunEnv() error {
	p := e.config.tmpPath()
	if _, err := os.Stat(p); os.IsNotExist(err) {
		e.runnerLog("mkdir %s", p)
		if err := os.Mkdir(p, 0755); err != nil {
			e.runnerLog("failed to mkdir, error: %s", err.Error())
			return err
		}
	}
	return nil
}

func (e *Engine) watching(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// NOTE: path is absolute
		if info != nil && !info.IsDir() {
			return nil
		}
		// exclude tmp dir
		if e.isTmpDir(path) {
			e.watcherLog("!exclude %s", e.config.rel(path))
			return filepath.SkipDir
		}
		// exclude hidden directories like .git, .idea, etc.
		if isHiddenDirectory(path) {
			return filepath.SkipDir
		}
		// exclude user specified directories
		if e.isExcludeDir(path) {
			e.watcherLog("!exclude %s", e.config.rel(path))
			return filepath.SkipDir
		}
		isIn, walkDir := e.checkIncludeDir(path)
		if !walkDir {
			e.watcherLog("!exclude %s", e.config.rel(path))
			return filepath.SkipDir
		}
		if isIn {
			return e.watchDir(path)
		}
		return nil
	})
}

func (e *Engine) watchDir(path string) error {
	if err := e.watcher.Add(path); err != nil {
		e.watcherLog("failed to watching %s, error: %s", path, err.Error())
		return err
	}
	e.watcherLog("watching %s", e.config.rel(path))

	go func() {
		e.withLock(func() {
			e.watchers++
		})
		defer func() {
			e.withLock(func() {
				e.watchers--
			})
		}()

		for {
			select {
			case <-e.watcherStopCh:
				return
			case ev := <-e.watcher.Events:
				e.mainDebug("event: %+v", ev)
				if !validEvent(ev) {
					break
				}
				if isDir(ev.Name) {
					e.watchNewDir(ev.Name, removeEvent(ev))
					break
				}
				if e.isExcludeFile(ev.Name) {
					break
				}
				if !e.isIncludeExt(ev.Name) {
					break
				}
				e.watcherDebug("%s has changed", e.config.rel(ev.Name))
				e.eventCh <- ev.Name
			case err := <-e.watcher.Errors:
				e.watcherLog("error: %s", err.Error())
			}
		}
	}()
	return nil
}

func (e *Engine) watchNewDir(dir string, removeDir bool) {
	if e.isTmpDir(dir) {
		return
	}
	if isHiddenDirectory(dir) || e.isExcludeDir(dir) {
		e.watcherLog("!exclude %s", e.config.rel(dir))
		return
	}
	if removeDir {
		if err := e.watcher.Remove(dir); err != nil {
			e.watcherLog("failed to stop watching %s, error: %s", dir, err.Error())
		}
		return
	}
	go func(dir string) {
		if err := e.watching(dir); err != nil {
			e.watcherLog("failed to watching %s, error: %s", dir, err.Error())
		}
	}(dir)
}

// Endless loop and never return
func (e *Engine) start() {
	firstRunCh := make(chan bool, 1)
	firstRunCh <- true
	e.binStopEndCh <- true

	for {
		var filename string

		e.mainLog("firstCh...")
		select {
		case <-e.exitCh:
			return
		case filename = <-e.eventCh:
			time.Sleep(e.config.buildDelay())
			e.flushEvents()
			if !e.isIncludeExt(filename) {
				continue
			}
			e.mainLog("%s has changed", e.config.rel(filename))
		case <-firstRunCh:
			// go down
			break
		}

		e.withLock(func() {
			e.mainDebug("binRuning: %v", e.binRunning)
			if e.binRunning {
				e.binStopCh <- true
			}
		})
		e.mainLog("start build run ...")

		<-e.binStopEndCh
		go e.buildRun()
	}
}

func (e *Engine) buildRun() {
	// e.buildLog("failed to build, error: %s", err.Error())
	// e.writeBuildErrorLog(err.Error())
	if err := e.runBin(); err != nil {
		e.runnerLog("failed to run, error: %s", err.Error())
	}
}

func (e *Engine) flushEvents() {
	for {
		select {
		case <-e.eventCh:
			e.mainDebug("flushing events")
		default:
			return
		}
	}
}

func (e *Engine) runBin() error {
	var err error
	// e.config.Build.Bin = "./echo1.sh abc"
	e.runnerLog("startCmd bin:%s", e.config.Build.Bin)
	cmd, stdout, stderr, err := e.startCmd(e.config.Build.Bin)
	if err != nil {
		return err
	}
	e.withLock(func() {
		e.binRunning = true
	})

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	go func() {
		io.Copy(os.Stdout, stdout)
		io.Copy(os.Stderr, stderr)
	}()

	go func(cmd *exec.Cmd, stdout io.ReadCloser, stderr io.ReadCloser) {
		<-e.binStopCh
		e.mainDebug("trying to kill cmd %+v", cmd.Args)
		defer func() {
			stdout.Close()
			stderr.Close()
		}()

		var err error
		pid, err := e.killCmd(cmd)
		if err != nil {
			e.mainDebug("failed to kill PID %d, error: %s", pid, err.Error())
			if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				e.mainDebug("Process not exited!!!!")
				os.Exit(1)
			}
		} else {
			e.mainDebug("cmd killed, pid: %d", pid)
		}
		e.binStopEndCh <- true
		e.withLock(func() {
			e.binRunning = false
		})
	}(cmd, stdout, stderr)
	// err = cmd.Wait()
	return nil
}

func (e *Engine) cleanup() {
	e.mainLog("cleaning...")
	defer e.mainLog("see you again~")

	e.mainLog("set binStopCh <- true ...")
	e.withLock(func() {
		if e.binRunning {
			e.binStopCh <- true
		}
	})

	e.withLock(func() {
		for i := 0; i < int(e.watchers); i++ {
			e.watcherStopCh <- true
		}
	})

	var err error
	if err = e.watcher.Close(); err != nil {
		e.mainLog("failed to close watcher, error: %s", err.Error())
	}

	if e.config.Misc.CleanOnExit {
		e.mainLog("deleting %s", e.config.tmpPath())
		if err = os.RemoveAll(e.config.tmpPath()); err != nil {
			e.mainLog("failed to delete tmp dir, err: %+v", err)
		}
	}
}

// Stop the air
func (e *Engine) Stop() {
	e.exitCh <- true
}
