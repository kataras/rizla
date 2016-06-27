package rizla

import (
	"os"
	"path/filepath"

	"strings"
	"time"
)

const minimumAllowReloadAfter = time.Duration(3) * time.Second

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
	return !(base == ".git" || base == "node_modules" || base == "vendor")
}

// Project the struct which contains the necessary fields to watch and reload(rerun) a go project
type Project struct {
	// optional Name for the project
	Name string
	// MainFile is the absolute path of the go project's main file source.
	MainFile string
	Args     []string

	// Watcher accepts subdirectories by the watcher
	// executes before the watcher starts,
	// if return true, then this (absolute) subdirectory is watched by watcher
	// the default accepts all subdirectories but ignores the ".git", "node_modules" and "vendor"
	Watcher MatcherFunc

	Matcher MatcherFunc
	// OnChange call something when this project's source code has changed and rizla going to reload
	OnChange func(string)
	// AllowReloadAfter skip reload on file changes that made too fast from the last reload
	// minimum allowed duration is 3 seconds.
	AllowReloadAfter time.Duration

	// DisableRuntimeDir set to true to disable adding subdirectories into the watcher, when a folder created at runtime
	// defaults to false
	DisableRuntimeDir bool

	dir string
	// proc the system Process of a running instance (if any)
	proc *os.Process
	// when the last change was made
	lastChange time.Time
}

// NewProject returns a simple project iteral which doesn't needs argument parameters
// and has the default file matcher ( which is valid if you want to reload only on .Go files).
//
// You can change all of its fields before the .Run function.
func NewProject(mainfile string) *Project {
	if mainfile == "" {
		mainfile = "main.go"
	}
	mainfile, _ = filepath.Abs(mainfile)

	dir := filepath.Dir(mainfile)

	return &Project{
		MainFile:         mainfile,
		Watcher:          DefaultWatcher,
		Matcher:          DefaultGoMatcher,
		AllowReloadAfter: minimumAllowReloadAfter,
		dir:              dir,
		lastChange:       time.Now(),
	}
}
