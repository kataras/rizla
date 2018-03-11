//Package rizla contains the source code of the rizla project
package rizla

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/kataras/golog"
)

const (
	isWindows = runtime.GOOS == "windows"
	isMac     = runtime.GOOS == "darwin"
	goExt     = ".go"
)

var (
	// Out is the logger which prints watcher errors and information.
	//
	// Change its printer's Output (io.Writer) by `Out.SetOutput(io.Writer)`.
	Out = golog.New().SetOutput(os.Stdout)

	projects []*Project

	pathSeparator = string(os.PathSeparator)

	stopChan = make(chan bool, 1)

	fsWatcher Watcher
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

var errUnexpected = errors.New("unexpected error!!! Please post an issue here: https://github.com/kataras/rizla/issues")

// RunWith starts the repeat of the build-run-watch-reload task of all projects
// receives optional parameters which can be the main source file
// of the project(s) you want to add, they can work nice with .Add(project) also, so dont worry use it.
//
// First receiver is the type of watcher
// second (optional) parameter(s) are the directories of the projects.
//    it's optional because they can be added with the .Add(NewProject) before the RunWith.
//
func RunWith(watcher Watcher, sources map[string][]string, delayOnDetect time.Duration) {
	// Author's notes: Because rizla's Run is not allowed to be called more than once
	// the whole package works as it is, so the watcher here
	// is CHANGING THE UNEXPORTED PACKGE VARIABLE 'fsWatcher'.
	// We don't export the 'fsWatcher' directly because this may cause issues
	// if user tries to change it while it runs.
	fsWatcher = watcher

	if len(sources) > 0 {
		for programFile, args := range sources {
			project := NewProject(programFile, args...)
			project.AllowRunAfter = delayOnDetect
			Add(project)
		}
	}

	for _, p := range projects {
		// go build
		if err := buildProject(p); err != nil {
			p.Err.Errorf(err.Error())
			continue
		}

		// exec run the builded program
		if err := runProject(p); err != nil {
			p.Err.Errorf(err.Error())
			continue
		}

	}

	watcher.OnError(func(err error) {
		Out.Errorf(err.Error())
	})

	watcher.OnChange(func(p *Project, filename string) {
		if time.Now().After(p.lastChange.Add(p.AllowReloadAfter)) {
			if match := p.Matcher(filename); !match {
				return
			}

			if p.AllowRunAfter > 0 {

				// Note that here, at "AllowRunAfter", maybe a lot of re-builds
				// at the same time if the user saved without
				// "AllowReloadAfter" configured, so every of the changes
				// are allowed in short period of time.
				// As a solution if AllowReloadAfter is not configured we will
				// configure it here, we can't do it before the first try
				// because it will wrong to wait "x" time to the first change detect allow.
				// We could make it with 1-buf channel as well.
				if p.AllowReloadAfter == 0 {
					p.lastChange = time.Now()
					p.AllowReloadAfter = p.AllowRunAfter
					// if minus := 250 * time.Millisecond; p.AllowRunAfter > minus {
					// 	p.AllowReloadAfter = p.AllowRunAfter - minus
					// }
				}
				time.Sleep(p.AllowRunAfter)
			}

			p.lastChange = time.Now()
			p.OnReload(filename)

			// kill previous running instance
			err := killProcess(p.proc, p.AppName)
			if err != nil {
				p.Err.Errorf("kill: %v", err)
				return
			}

			// go build
			err = buildProject(p)
			if err != nil {
				p.Err.Errorf(err.Error())
				return
			}

			// exec run the builded program
			err = runProject(p)
			if err != nil {
				p.Err.Errorf("failed to run the project: %v", err)
				return
			}

			p.OnReloaded(filename)

		}
	})

	watcher.Loop()
}

// Run same as RunWith but runs with the default file system watcher
// which is the fsnotify (watch over file system's signals) or the last used with RunWith.
//
// It's a map of main files and their arguments, if any.
func Run(sources map[string][]string) {
	if fsWatcher != nil {
		// if user already called RunWith before, the watcher is saved on the 'fsWatcher' variable,
		// use that instead.
		RunWith(fsWatcher, sources, 0)
		return
	}

	RunWith(newSignalWatcher(), sources, 0)
}

// Stop any projects are watched by the RunWith/Run method, this function should be call when you call the Run inside a goroutine.
func Stop() {
	if fsWatcher != nil {
		fsWatcher.Stop()
	}
}

func isDirectory(fullname string) bool {
	if info, err := os.Stat(fullname); err == nil && info.IsDir() {
		return true
	}
	return false
}

func buildProject(p *Project) error {

	// relative := p.MainFile[len(p.dir)+1:len(p.MainFile)-3] + goExt
	goBuild := exec.Command("go", "build", ".")
	goBuild.Dir = p.dir
	goBuild.Stdout = p.Out.Printer.Output
	goBuild.Stderr = p.Err.Printer.Output
	return goBuild.Run()
}

func runProject(p *Project) error {

	// buildProject := p.MainFile[len(p.dir) : len(p.MainFile)-3] // with prepended slash

	buildProject := filepath.Base(p.dir)

	if isWindows {
		buildProject += ".exe"
	}

	// runCmd := exec.Command("."+buildProject, p.Args...)

	runCmd := exec.Command("./"+buildProject, p.Args...)
	runCmd.Dir = p.dir

	if p.DisableProgramRerunOutput && p.i > 0 && p.proc != nil {
		// if already ran once succesfuly, we don't need to printout the output of the program, because we will have big outputs if the program has banner (like Iris :))
	} else {
		runCmd.Stdout = p.Out.Printer.Output
	}

	runCmd.Stderr = p.Err.Printer.Output

	// Moved to exec.Command's second argument instead:
	// if p.Args != nil && len(p.Args) > 0 {
	// 	runCmd.Args = p.Args[0 : len(p.Args)-1]
	// }

	if err := runCmd.Start(); err != nil {
		return err
	}
	p.proc = runCmd.Process
	return nil
}

func killProcess(proc *os.Process, appName string) (err error) {
	if proc == nil {
		return nil
	}

	if !isMac {
		err = proc.Release()
		if err != nil {
			return nil // to prevent throw an error if the proc is not yet started correctly (= previous build error)
		}
	}

	if (isWindows || isMac) && proc.Pid <= 0 {
		return nil
	}

	err = proc.Kill()
	if err == nil {
		_, err = proc.Wait()
	} else {
		// force kill, sometimes proc.Kill or Signal(os.Kill) doesn't kills
		if isWindows {
			err = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(proc.Pid)).Run()
			if err != nil && err.Error() == "exit status 128" {
				err = nil // skip that stupid error here.
			}

			// err = exec.Command("taskkill", "/im", appName+".exe").Run()
		} else if isMac {
			err = exec.Command("killall", "-KILL", strconv.Itoa(proc.Pid)).Run()
		} else {
			// err = exec.Command("kill", "-INT", "-"+strconv.Itoa(proc.Pid)).Run()
			err = exec.Command("pkill", "-SIGINT", appName).Run()
		}
	}
	proc = nil
	return
}
