//Package goo contains the source code of the goo project
package goo

import (
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/iris-contrib/errors"
	"github.com/mattn/go-colorable"
)

var (
	// Out The logger output for all projects
	Out = os.Stdout
	// Err The logger output for errors for all projects
	Err = os.Stderr
	// In The input for all projects
	In       = os.Stdin
	projects []*Project
)

// Add receives a Project and adds it to the projects
func Add(p *Project) {
	p.prepare()
	projects = append(projects, p)
}

// Reset clears the current projects, doesn't stop them if running
func Reset() {
	projects = projects[0:0]
}

var (
	errInvalidArgs = errors.New("Invalid arguments [%s], type -h to get assistant")
	errInvalidExt  = errors.New("%s is not a go program")
	errUnexpected  = errors.New("Unexpected error!!! Please post an issue here: https://github.com/kataras/goo/issues")
	errBuild       = errors.New("\n Failed to build the %s iris program. Trace: %s")
	errRun         = errors.New("\n Failed to run the %s iris program. Trace: %s")
)

//Run starts the repeat of the build-run-watch-reload task
func Run() {
	color.Output = colorable.NewColorable(Out)

	watcher, werr := fsnotify.NewWatcher()
	if werr != nil {
		color.Red(werr.Error())
		return
	}

	for _, p := range projects {
		for _, subdir := range p.compiledDirectories {
			if werr = watcher.Add(subdir); werr != nil {
				color.Red(werr.Error())
			}
		}
	}

	defer func() {
		color.Red(errUnexpected.Error())
	}()

	var lastChange = time.Now()
	var i = 0
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				filename := event.Name
				for _, p := range projects {
					if p.Matcher(filename) {

					}
				}
				//this is received two times, the last time is the real changed file (at least on windows(?)), so
				i++
				if i%2 == 0 || !isWindows { // this 'hack' works for windows & linux but I dont know if works for osx too, we can wait for issue reports here.
					if time.Now().After(lastChange.Add(time.Duration(1) * time.Second)) {
						lastChange = time.Now()
						//TODO
					}
				}

			}
		case err := <-watcher.Errors:
			color.Red(err.Error())
		}
	}

}

func goBuildProject(p *Project) error {
	goBuild := exec.Command("go", "build", p.MainFile)
	goBuild.Stdout = Out
	goBuild.Stderr = Err
	if err := goBuild.Run(); err != nil {
		return err
	}
	return nil
}

func goRunProject(p *Project) error {

	execFilename := p.MainFile[len(p.compiledDirectory) : len(p.MainFile)-3]
	if isWindows {
		execFilename += ".exe"
	}

	runCmd := exec.Command("." + pathSeparator + execFilename)
	runCmd.Dir = p.compiledDirectory
	runCmd.Stdout = Out
	runCmd.Stderr = Err
	runCmd.Stdin = In

	if err := runCmd.Start(); err != nil {
		return err
	}
	p.proc = runCmd.Process
	return nil
}
