package rizla

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"strings"
	"time"
)

const minimumAllowReloadAfter = time.Duration(1) * time.Second

// MatcherFunc returns whether the file should be watched for the reload
type MatcherFunc func(string) bool

// DefaultGoMatcher is the default Matcher for the Project iteral
func DefaultGoMatcher(fullname string) bool {
	return filepath.Ext(fullname) == goExt || (!isWindows && strings.Contains(fullname, goExt))
}

// Project the struct which contains the necessary fields to watch and reload(rerun) a go project
type Project struct {
	// optional Name for the project
	Name string
	// MainFile is the absolute path of the go project's main file source.
	MainFile string
	Args     []string
	Matcher  MatcherFunc
	// OnChange call something when this project's source code has changed and rizla going to reload
	OnChange func()
	// AllowReloadAfter skip reload on file changes that made too fast from the last reload
	// minimum duration is 1 second.
	AllowReloadAfter time.Duration
	dir              string
	// subdirs contains all dir from the directory
	subdirs []string
	// proc the system Process of a running instance (if any)
	proc *os.Process
	// when the last change was made
	lastChange time.Time
	// Used only on windows, winEvtCount ever this is not an odd  number then the event is valid
	winEvtCount int
}

// NewProject returns a simple project iteral which doesn't needs argument parameters
// and has the default file matcher ( which is valid if you want to reload only on .Go files).
//
// You can change all of its fields before the .Run function.
func NewProject(mainfile string) *Project {
	return &Project{MainFile: mainfile}
}

func (p *Project) prepare() {
	if p.Matcher == nil {
		p.Matcher = DefaultGoMatcher
	}

	if p.AllowReloadAfter < minimumAllowReloadAfter {
		p.AllowReloadAfter = minimumAllowReloadAfter
	}

	if p.MainFile == "" {
		p.MainFile = "main.go"
	}
	if !filepath.IsAbs(p.MainFile) {
		p.MainFile, _ = filepath.Abs(p.MainFile)
	}

	p.dir = filepath.Dir(p.MainFile)

	subfiles, err := ioutil.ReadDir(p.dir)
	if err != nil {
		panic(err)
	}

	for _, subfile := range subfiles {
		if subfile.IsDir() {
			path := p.dir + pathSeparator + subfile.Name()
			p.subdirs = append(p.subdirs, path)
		}
	}

}
