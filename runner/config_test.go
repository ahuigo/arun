package runner

import (
	"runtime"
	"testing"
)

const (
	bin = `./tmp/main`
	cmd = "go build -o ./tmp/main ."
)

func getConfig() config {
	build := cfgBuild{
		Cmd:         "go build -o ./tmp/main .",
		Log:         "build-errors.log",
		IncludeExt:  []string{"go", "tpl", "tmpl", "html"},
		ExcludeDir:  []string{"assets", "tmp", "vendor"},
		Delay:       1000,
		StopOnError: true,
	}
	if runtime.GOOS == "windows" {
		build.Cmd = cmd
	}

	return config{
		Root:   ".",
		TmpDir: "tmp",
		Build:  build,
	}
}

func TestBinCmdPath(t *testing.T) {

	var err error

	c := getConfig()
	err = c.preprocess()
	if err != nil {
		t.Fatal(err)
	}
}
