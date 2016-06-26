//Package goo contains the source code of the goo project
package goo

import (
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
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

//Run starts the repeat of the build-run-watch-reload task
func Run() {
	color.Output = colorable.NewColorable(Out)

	for _, p := range projects {

		watcher, werr := fsnotify.NewWatcher()
		if werr != nil {
			color.Red(werr.Error())
			return
		}

		go func() {
			var lastChange = time.Now()
			var i = 0
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
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
		}()
		for _, subdir := range p.compiledDirectories {
			if werr = watcher.Add(subdir); werr != nil {
				color.Red(werr.Error())
			}
		}

	}
}
