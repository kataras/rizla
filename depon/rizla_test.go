package depon

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestProjectAddRemoveLen(t *testing.T) {
	if Len() > 0 {
		t.Fatalf("Init projects should be len 0 but got %d", Len())
	}

	Add(NewProject("main.go"))
	if Len() != 1 {
		t.Fatalf("Expected projects' length to be 1 but got %d", Len())
	}

	RemoveAll()
	if Len() > 0 {
		t.Fatalf("Expected projects' length to be 0 but got %d", Len())
	}

	Add(NewProject("main.go"))
	Add(NewProject("main.go"))
	if Len() != 2 {
		t.Fatalf("Expected projects' length to be 2 but got %d", Len())
	}
	RemoveAll()
}

func TestPrinter(t *testing.T) {
	// I tried to use memfs and some other libraries for virtual file system but the colorable  doesn't accept these, if anyone can help send me a message on the chat...
	logggerpath := "mytestlogger.txt"

	logger, ferr := os.Create(logggerpath)
	defer func() {
		logger.Close()
		Out.Close()
		os.RemoveAll(logggerpath)
	}()

	if ferr != nil {
		t.Fatal(ferr)
	}
	//set the global output for the printer
	Out = logger
	// get the printer
	printer := newPrinter()

	s := "Hello"
	printer.Print(s)

	contents, err := ioutil.ReadFile(logggerpath)
	if err != nil || len(contents) == 0 && err == io.EOF {
		t.Fatalf("While trying to read from the logger %s", err.Error())
	} else {
		if len(contents) != len(s) || string(contents) != s {
			t.Fatalf("Logger reads but the its contents are not valid, expected len bytes %d but got %d, expected %s but got %s", len(s), len(contents), s, string(contents))
		}
	}
}
