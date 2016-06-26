package goo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	isWindows = runtime.GOOS == "windows"
	goExt     = ".go"
)

var (
	workingDir    string
	pathSeparator = string(os.PathSeparator)
)

func init() {
	if d, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		workingDir = d
	}
}

// MatcherFunc returns whether the file should be watched for the reload
type MatcherFunc func(string) bool

// DefaultMatcher is the default Matcher for the Project iteral
func DefaultMatcher(fullname string) bool {
	return filepath.Ext(fullname) == goExt || (!isWindows && strings.Contains(fullname, goExt))
}

// Project the struct which contains the necessary fields to watch and reload(rerun) a go project
type Project struct {
	// MainFile is the absolute path of the go project's main file source.
	MainFile          string
	compiledDirectory string
	// compiledDirectories contains all subdirectories from the Directory, this field is actually used
	compiledDirectories []string
	Args                map[string]string
	Matcher             MatcherFunc
}

func (p *Project) prepare() {
	if p.MainFile == "" {
		p.MainFile = "main.go"
	}
	if !filepath.IsAbs(p.MainFile) {
		p.MainFile = workingDir + pathSeparator + p.MainFile
	}

	p.compiledDirectory = filepath.Dir(p.MainFile)

	subfiles, err := ioutil.ReadDir(p.compiledDirectory)
	if err != nil {
		panic(err)
	}

	for _, subfile := range subfiles {
		if subfile.IsDir() {
			if abspath, err := filepath.Abs(p.compiledDirectory + pathSeparator + subfile.Name()); err == nil {
				p.compiledDirectories = append(p.compiledDirectories, abspath)
			}
		}
	}

	if p.Matcher == nil {
		p.Matcher = DefaultMatcher
	}
}
