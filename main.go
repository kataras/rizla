//Package main contains the source code for goo binary
package main

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	// Version of Goo command line tool
	Version = "0.0.1"
	// Name of Goo
	Name = "Goo"
	// Description og Goo
	Description = "Builds, runs and watches your Go Applications with ease.\nSupports multi-projects also."
)

func main() {
	argsLen := len(os.Args)

	if argsLen <= 1 {
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
	println("Help asked")
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
