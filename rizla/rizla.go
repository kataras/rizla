//Package rizla contains the source code of the rizla project
package rizla

import (
	"os"
	"os/exec"
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
		p.prepare()
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
	errBuild       = errors.New("Failed to build the program. Trace: %s\n")
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

		// go build
		err := buildProject(p)
		if err != nil {
			dangerf(errBuild.Format(err.Error()).Error())
			continue
		}

		// exec run the builded program
		err = runProject(p)
		if err != nil {
			dangerf(errRun.Format(err.Error()).Error())
			continue
		}

		// add to the watcher
		// add its root folder
		if werr = watcher.Add(p.dir); werr != nil {
			dangerf("\n" + werr.Error() + "\n")
		}

		// add subdirs also
		for _, subdir := range p.subdirs {
			if werr = watcher.Add(subdir); werr != nil {
				dangerf("\n" + werr.Error() + "\n")
			}
		}

	}
	hasStoppedManually := false

	// if something bad happens and program exits, show an unexpecter error message
	defer func() {
		if !hasStoppedManually {
			dangerf(errUnexpected.Error())
		}
	}()

	stopChan <- false

	// run the watcher
	for {
		select {
		case stop := <-stopChan:
			if stop {
				for _, p := range projects {
					killProcess(p.proc)
				}
				hasStoppedManually = true
				watcher.Close()

				break
			}

		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				filename := event.Name

				for _, p := range projects {
					//this is received two times, the last time is the real changed file (at least on windows(?)), so
					p.winEvtCount++
					if time.Now().After(p.lastChange.Add(p.AllowReloadAfter)) {

						if p.winEvtCount%2 == 0 || !isWindows { // this 'hack' works for windows & linux but I dont know if works for osx too, we can wait for issue reports here.
							if p.Matcher(filename) {
								// call the user defined change callback
								if p.OnChange != nil {
									p.OnChange()
								}

								fromproject := ""
								if p.Name != "" {
									fromproject = "From project '" + p.Name + "': "
								}
								infof("\n%sA change has been detected, reloading now...", fromproject)
								p.lastChange = time.Now()
								// kill previous running instance
								err := killProcess(p.proc)
								if err != nil {
									dangerf(err.Error())
									continue
								}
								// go build
								err = buildProject(p)
								if err != nil {
									dangerf(errBuild.Format(err.Error()).Error())
									continue
								}

								// exec run the builded program
								err = runProject(p)
								if err != nil {
									dangerf(errRun.Format(err.Error()).Error())
									continue
								}
								successf("ready!\n")

							}

						}
					}
				}

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
