package rizla

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kataras/golog"
)

const minimumAllowReloadAfter = time.Duration(2) * time.Second

// DefaultDisableProgramRerunOutput a long name but, it disables the output of the program's 'messages' after the first successfully run for each of the projects
// the project iteral can be override this value.
// set to true to disable the program's output when reloads
var DefaultDisableProgramRerunOutput = false

// MatcherFunc returns whether the file should be watched for the reload
type MatcherFunc func(string) bool

// DefaultGoMatcher is the default Matcher for the Project iteral
func DefaultGoMatcher(fullname string) bool {
	return (filepath.Ext(fullname) == goExt) ||
		(!isWindows && strings.Contains(fullname, goExt))
}

// DefaultWatcher is the default Watcher for the Project iteral
// allows all subdirs except .git, node_modules and vendor
func DefaultWatcher(abs string) bool {
	base := filepath.Base(abs)
	// by-default ignore .git folder, node_modules, vendor and any hidden files.
	return !(base == ".git" || base == "node_modules" || base == "vendor" || base == ".")
}

// OnReloadScripts simple file names which will execute a script, i.e `./on_reload.sh` or `./on_reload.bat` or even `service supervisor restart`
// on windows, it will just execute that based on the operating system, nothing crazy here,
// they are filled by the cli but they can be customized by the source as well.
//
// If contains whitespaces, after the first whitespace they are the command's flags (if not a script file).
var OnReloadScripts []string

// DefaultOnReload fired when file has changed and reload going to happens
func DefaultOnReload(p *Project) func(string) {
	return func(string) {
		fromproject := ""
		if p.Name != "" {
			fromproject = "From project '" + p.Name + "': "
		}
		p.Out.Infof("%sA change has been detected, reloading now...", fromproject)

		if len(OnReloadScripts) > 0 {
			p.Out.Infof("%sExecuting commands from %s before restart...", fromproject, strings.Join(OnReloadScripts, ", "))
			for _, s := range OnReloadScripts {

				// The below should work for things like
				// service supervisor restart
				nameAndFlags := strings.Split(s, " ")
				var args []string
				name := nameAndFlags[0]
				if len(nameAndFlags) > 1 {
					args = nameAndFlags[1:]
				}

				cmd := exec.Command(name, args...)
				cmd.Stderr = p.Err.Printer.Output
				cmd.Stdout = p.Out.Printer.Output
				if err := cmd.Run(); err != nil {
					p.Out.Errorf("%s%s run: %v", fromproject, s, err)
					os.Exit(1)
				}
				cmd.Wait() // ignore error.
			}
		}
	}
}

// var rdy = []byte("ready!\n")

// DefaultOnReloaded fired when reload has been finished.
// Defaults to noOp.
func DefaultOnReloaded(p *Project) func(string) {
	return func(string) {
		// p.Out.Printer.Output.Write(rdy)
	}
}

// Project the struct which contains the necessary fields to watch and reload(rerun) a go project
type Project struct {
	// optional Name for the project
	Name string
	// MainFile is the absolute path of the go project's main file source.
	MainFile string
	// The application's name, usually is the MainFile withotu the extension.
	// At the future we may provide a way for custom naming which will be used on the "go build -o" flag.
	AppName string
	Args    []string
	// The Output destination (sent by rizla and your program)
	Out *golog.Logger
	// The Err Output destination (sent on rizla errors and your program's errors)
	Err *golog.Logger
	// Watcher accepts subdirectories by the watcher
	// executes before the watcher starts,
	// if return true, then this (absolute) subdirectory is watched by watcher
	// the default accepts all subdirectories but ignores the ".git", "node_modules" and "vendor"
	Watcher MatcherFunc
	Matcher MatcherFunc
	// AllowReloadAfter skip reload on file changes that made too fast from the last reload
	// minimum allowed duration is 3 seconds.
	AllowReloadAfter time.Duration
	// AllowRunAfter it accepts the file changes
	// but it waits "x" duration for the reload to happen.
	AllowRunAfter time.Duration
	// OnReload fires when when file has been changed and rizla is going to reload the project
	// the parameter is the changed file name
	OnReload func(string)
	// OnReloaded fires when rizla finish with the reload
	// the parameter is the changed file name
	OnReloaded func(string)
	// DisableRuntimeDir set to true to disable adding subdirectories into the watcher, when a folder created at runtime
	// set to true to disable the program's output when reloads
	// defaults to false
	DisableRuntimeDir bool
	// DisableProgramRerunOutput a long name but, it disables the output of the program's 'messages' after the first successfully run
	// defaults to false
	DisableProgramRerunOutput bool

	dir string
	// proc the system Process of a running instance (if any)
	proc *os.Process
	// when the last change was made
	lastChange time.Time
	// i%2 ==0 if windows, then the reload is allowed
	i int
}

// NewProject returns a simple project iteral which doesn't needs argument parameters
// and has the default file matcher ( which is valid if you want to reload only on .Go files).
//
// You can change all of its fields before the .Run function.
func NewProject(mainfile string, args ...string) *Project {
	if mainfile == "" {
		mainfile = "main.go"
	}
	appName := mainfile[0 : len(mainfile)-len(goExt)]
	mainfile, _ = filepath.Abs(mainfile)

	dir := filepath.Dir(mainfile)

	p := &Project{
		MainFile:                  mainfile,
		AppName:                   appName,
		Args:                      args,
		Out:                       golog.New().SetOutput(os.Stdout),
		Err:                       golog.New().SetOutput(os.Stderr),
		Watcher:                   DefaultWatcher,
		Matcher:                   DefaultGoMatcher,
		AllowReloadAfter:          minimumAllowReloadAfter,
		DisableProgramRerunOutput: DefaultDisableProgramRerunOutput,
		dir:        dir,
		lastChange: time.Now(),
	}

	p.OnReload = DefaultOnReload(p)
	p.OnReloaded = DefaultOnReloaded(p)
	return p
}
