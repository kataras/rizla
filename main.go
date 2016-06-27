//Package main depon builds, runs and monitors your Go Applications with ease.
//
//   depon main.go
//   depon C:/myprojects/project1/main.go C:/myprojects/project2/main.go C:/myprojects/project3/main.go
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kataras/depon/depon"
)

const (
	// Version of Depon command line tool
	Version = "0.0.1"
	// Name of Depon
	Name = "Depon"
	// Description of Depon
	Description = "Depon builds, runs and monitors your Go Applications with ease."
)

var helpTmpl = fmt.Sprintf(`NAME:
   %s - %s

USAGE:
   depon main.go
   depon C:/myprojects/project1/main.go C:/myprojects/project2/main.go C:/myprojects/project3/main.go

VERSION:
   %s
   `, Name, Description, Version)

func main() {
	argsLen := len(os.Args)

	if argsLen <= 1 {
		help(-1)
	} else if isArgHelp(os.Args[1]) {
		help(0)
	}

	args := os.Args[1:]
	for _, a := range args {
		if !strings.HasSuffix(a, ".go") {
			color.Red("Error: Please provide files with '.go' extension.\n")
			help(-1)
		} else if p, _ := filepath.Abs(a); !fileExists(p) {
			color.Red("Error: File " + p + " does not exists.\n")
			help(-1)
		}
	}

	for _, a := range args {
		p := depon.NewProject(a)
		depon.Add(p)
	}

	depon.Run()
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
