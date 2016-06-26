//Package main Rizla builds, runs and watches your Go Applications with ease.
//
//   rizla main.go
//   rizla C:/myprojects/project1/main.go C:/myprojects/project2/main.go C:/myprojects/project3/main.go
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Version of Rizla command line tool
	Version = "0.0.1"
	// Name of Rizla
	Name = "Rizla"
	// Description of Rizla
	Description = "Builds, runs and watches your Go Applications with ease."
)

var helpTmpl = fmt.Sprintf(`NAME:
   %s - %s

USAGE:
   rizla main.go
   rizla C:/myprojects/project1/main.go C:/myprojects/project2/main.go C:/myprojects/project3/main.go

VERSION:
   %s
   `, Name, Description, Version)

func main() {
	argsLen := len(os.Args)

	if argsLen <= 1 {
		help()
	} else if argsLen == 2 && os.Args[1] == "help" {
		help()
	}
	args := os.Args[1:]
	for _, a := range args {
		if !strings.HasSuffix(a, ".go") {
			pleaseGo()
		} else if p, _ := filepath.Abs(a); !fileExists(p) {
			fileDoesntFound(p)
		}
	}

	println(strings.Join(args, ","))
}

func fileExists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}

func help() {
	println(helpTmpl)
	os.Exit(-1)
}

func pleaseGo() {
	println("Please provide files with .go")
	os.Exit(-1)
}

func fileDoesntFound(arg string) {
	println("file " + arg + " doesnt found")
	os.Exit(-1)
}
