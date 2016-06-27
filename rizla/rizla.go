//Package rizla contains the source code of the rizla project
package rizla

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/iris-contrib/errors"
	"github.com/mattn/go-colorable"
)

const (
	isWindows = runtime.GOOS == "windows"
	goExt     = ".go"
)

var (
	// Out The logger output for all projects
	Out = os.Stdout
	// Err The logger output for errors for all projects
	Err = os.Stderr

	projects []*Project

	pathSeparator = string(os.PathSeparator)

	stopChan = make(chan bool, 1)
)

// Add project(s) to the container
func Add(proj ...*Project) {
	for _, p := range proj {
		projects = append(projects, p)
	}
}

// RemoveAll clears the current projects, doesn't stop them if running
func RemoveAll() {
	projects = make([]*Project, 0)
}

// Len how much projects have  been added so far
func Len() int {
	return len(projects)
}

var (
	errInvalidArgs = errors.New("Invalid arguments [%s], type -h to get assistant\n")
	errUnexpected  = errors.New("Unexpected error!!! Please post an issue here: https://github.com/kataras/rizla/issues\n")
	errBuild       = errors.New("Failed to build the program.\n")
	errRun         = errors.New("Failed to run the the program. Trace: %s\n")
)

// newPrinter returns a new colorable printer
func newPrinter() *color.Color {
	color.Output = colorable.NewColorable(Out)
	return color.New()
}

// Run starts the repeat of the build-run-watch-reload task of all projects
// receives optional parameters which can be the main source file of the project(s) you want to add, they can work nice with .Add(project) also, so dont worry use it.
func Run(sources ...string) {
	if len(sources) > 0 {
		for _, s := range sources {
			Add(NewProject(s))
		}
	}

	printer := newPrinter()

	dangerf := func(format string, a ...interface{}) {
		printer.Add(color.FgRed)
		printer.Printf(format, a...)
	}

	infof := func(format string, a ...interface{}) {
		printer.Add(color.FgCyan)
		printer.Printf(format, a...)
	}

	successf := func(format string, a ...interface{}) {
		printer.Add(color.FgGreen)
		printer.Printf(format, a...)
	}

	watcher, werr := fsnotify.NewWatcher()
	if werr != nil {
		dangerf(werr.Error())
		return
	}

	for _, p := range projects {

		// add to the watcher first in order to watch changes and re-builds if the first build has fallen

		// add its root folder first
		if werr = watcher.Add(p.dir); werr != nil {
			dangerf("\n" + werr.Error() + "\n")
		}

		visitFn := func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				// check if this subdir is allowed
				if p.Watcher(path) {
					if werr = watcher.Add(path); werr != nil {
						dangerf("\n" + werr.Error() + "\n")
					}
				} else {
					return filepath.SkipDir
				}

			}
			return nil
		}

		if err := filepath.Walk(p.dir, visitFn); err != nil {
			panic(err)
		}

		// go build
		if err := buildProject(p); err != nil {
			dangerf(errBuild.Error())
			continue
		}

		// exec run the builded program
		if err := runProject(p); err != nil {
			dangerf(errRun.Error())
			continue
		}

	}
	hasStoppedManually := false

	defer func() {
		watcher.Close()
		for _, p := range projects {
			killProcess(p.proc)
		}
		if !hasStoppedManually {
			// if something bad happens and program exits, show an unexpected error message
			dangerf(errUnexpected.Error())
		}
	}()

	stopChan <- false

	// run the watcher
	for {
		select {
		case stop := <-stopChan:
			if stop {
				hasStoppedManually = true
				return
			}

		case event := <-watcher.Events:
			filename := event.Name
			for _, p := range projects {
				go func() {

					if time.Now().After(p.lastChange.Add(p.AllowReloadAfter)) {
						p.lastChange = time.Now()

						isDir := false
						match := p.Matcher(filename)
						if !p.DisableRuntimeDir { //we don't check if !match because the folder maybe be: myfolder.go
							isDir = isDirectory(filename)
						}

						if match || isDir && p.Watcher(filename) {
							if isDir {
								if werr = watcher.Add(filename); werr != nil {
									dangerf("\n" + werr.Error() + "\n")
								}
							}

							fromproject := ""
							if p.Name != "" {
								fromproject = "From project '" + p.Name + "': "
							}
							infof("\n%sA change has been detected, reloading now...", fromproject)

							// kill previous running instance
							err := killProcess(p.proc)
							if err != nil {
								dangerf(err.Error())
								return
							}
							// go build
							err = buildProject(p)
							if err != nil {
								dangerf(errBuild.Error())
								return
							}

							// exec run the builded program
							err = runProject(p)
							if err != nil {
								dangerf(errRun.Format(err.Error()).Error())
								return
							}
							successf("ready!\n")

							// call the user defined change callback
							if p.OnChange != nil {
								p.OnChange(filename)
							}

						}
					}
				}()
			}
		case err := <-watcher.Errors:
			if !hasStoppedManually {
				dangerf(err.Error())
			}
		}
	}

}

// Stop any projects are watched by the Run method, this function should be call when you call the Run inside a goroutine.
func Stop() {
	stopChan <- true
}

func isDirectory(fullname string) bool {
	if info, err := os.Stat(fullname); err == nil && info.IsDir() {
		return true
	}
	return false
}

func buildProject(p *Project) error {
	relative := p.MainFile[len(p.dir)+1:len(p.MainFile)-3] + goExt
	goBuild := exec.Command("go", "build", relative)
	goBuild.Dir = p.dir
	goBuild.Stdout = Out
	goBuild.Stderr = Err
	if err := goBuild.Run(); err != nil {
		return err
	}
	return nil
}

func runProject(p *Project) error {

	buildProject := p.MainFile[len(p.dir) : len(p.MainFile)-3] // with prepended slash
	if isWindows {
		buildProject += ".exe"
	}

	runCmd := exec.Command("." + buildProject)
	runCmd.Dir = p.dir
	runCmd.Stdout = Out
	runCmd.Stderr = Err

	if p.Args != nil && len(p.Args) > 0 {
		runCmd.Args = p.Args[0 : len(p.Args)-1]
	}

	if err := runCmd.Start(); err != nil {
		return err
	}
	p.proc = runCmd.Process
	return nil
}

func killProcess(proc *os.Process) (err error) {
	if proc == nil {
		return nil
	}
	err = proc.Kill()
	if err == nil {
		_, err = proc.Wait()
	} else {
		// force kill, sometimes proc.Kill or Signal(os.Kill) doesn't kills
		if isWindows {
			err = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(proc.Pid)).Run()
		} else {
			err = exec.Command("kill", "-INT", "-"+strconv.Itoa(proc.Pid)).Run()
		}
	}
	proc = nil
	return
}
