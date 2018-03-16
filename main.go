//Package main rizla builds, runs and monitors your Go Applications with ease.
//
//   rizla main.go
//   rizla C:/myprojects/project1/main.go C:/myprojects/project2/main.go C:/myprojects/project3/main.go
//   rizla -walk main.go [if -walk then rizla uses the stdlib's filepath.Walk method instead of file system's signals]
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/rizla/rizla"
)

const (
	// Version of rizla command line tool
	Version = "0.1.1"
	// Name of rizla
	Name = "Rizla"
	// Description of rizla
	Description = "Rizla builds, runs and monitors your Go Applications with ease."
)

const delayArgName = "-delay"

func getDelayFromArg(arg string) (time.Duration, bool) {
	if strings.HasPrefix(arg, delayArgName) || strings.HasPrefix(arg, delayArgName[1:]) {
		// [-]delay=5s
		// [-]delay 5s

		if equalIdx := strings.IndexRune(arg, '='); equalIdx > 0 {
			delayStr := arg[equalIdx+1:]
			d, _ := time.ParseDuration(delayStr)
			return d, true
		}

		if spaceIdx := strings.IndexRune(arg, ' '); spaceIdx > 0 {
			delayStr := arg[spaceIdx+1:]
			d, _ := time.ParseDuration(delayStr)
			return d, true
		}

	}

	return 0, false
}

const onReloadArg = "-onreload"

func getOnReloadArg(arg string) (string, bool) {
	if strings.HasPrefix(arg, onReloadArg) || strings.HasPrefix(arg, onReloadArg[1:]) {
		// first equality ofc...
		// [-]onreload=cmd1,cmd2,file.sh
		// [-]onreload cmd1,cmd2,file.sh

		if equalIdx := strings.IndexRune(arg, '='); equalIdx > 0 {
			src := arg[equalIdx+1:]
			return src, true
		}

		if spaceIdx := strings.IndexRune(arg, ' '); spaceIdx > 0 {
			src := arg[spaceIdx+1:]
			return src, true
		}

	}

	return "", false
}

var helpTmpl = fmt.Sprintf(`NAME:
   %s - %s

USAGE:
   rizla main.go
   rizla C:/myprojects/project1/main.go C:/myprojects/project2/main.go C:/myprojects/project3/main.go
   rizla -walk main.go [if -walk then rizla uses the stdlib's filepath.Walk method instead of file system's signals]
   rizla -delay=5s main.go [if delay > 0 then it delays the reload, also note that it accepts the first change but the rest of changes every "delay"]
   rizla -onreload="service supervisor restart" main.go or rizla -onreload="cmd /C echo Hello World!" main.go
VERSION:
   %s
   `, Name, Description, Version)

func main() {
	argsLen := len(os.Args)

	errorf := golog.New().SetOutput(os.Stderr).Errorf

	if argsLen <= 1 {
		help(-1)
	} else if isArgHelp(os.Args[1]) {
		help(0)
	}

	args := os.Args[1:]
	programFiles := make(map[string][]string, 0) // key = main file, value = args.
	fsWatcher, _ := rizla.WatcherFromFlag("signal")

	var lastProgramFile string
	var delayOnDetect time.Duration

	for i, a := range args {
		// if main files with arguments aren't passed yet,
		// then the argument(s) should refer to the rizla tool and not the
		// external programs.
		if lastProgramFile == "" {
			// The first argument must be the method type of the file system's watcher.
			// if -w,-walk,walk then
			//   asks to use the stdlib's filepath.walk method instead of the operating system's signal.
			//   It's only usage is when the user's IDE overrides the os' signals.
			// otherwise
			//   use the fsnotify's operating system's file system's signals.
			if watcher, ok := rizla.WatcherFromFlag(a); ok {
				fsWatcher = watcher
				continue
			}

			if delay, ok := getDelayFromArg(a); ok {
				delayOnDetect = delay
				continue
			}

			if onReloadSources, ok := getOnReloadArg(a); ok {
				rizla.OnReloadScripts = strings.Split(onReloadSources, ",")
			}
		}

		// it's main.go or any go main program
		if strings.HasSuffix(a, ".go") {
			programFiles[a] = []string{}
			lastProgramFile = a
			continue
		}

		if lastProgramFile != "" && len(args) > i+1 {
			programFiles[lastProgramFile] = append(programFiles[lastProgramFile], args[i:]...) // note that: the executable argument (1st arg) is set-ed by the exec.Command on `runProject`.
			continue
		}
	}

	// no program files given
	if len(programFiles) == 0 {
		errorf("please provide a *.go file.\n")
		help(-1)
		return
	}

	// check if given program files exist
	for programFile := range programFiles {
		// the argument is not the first  given is *.go but doesn't exists on user's disk
		if p, _ := filepath.Abs(programFile); !fileExists(p) {
			errorf("file " + p + " does not exists.\n")
			help(-1)
			return
		}
	}

	rizla.RunWith(fsWatcher, programFiles, delayOnDetect)
}

func help(code int) {
	os.Stdout.WriteString(helpTmpl)
	os.Exit(code)
}

func isArgHelp(s string) bool {
	return s == "help" || s == "-h" || s == "-help"
}

func fileExists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}
