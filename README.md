Rizla builds, runs and monitors your Go Applications with ease.

[![Travis Widget]][Travis] [![Release Widget]][Release] [![Report Widget]][Report] [![License Widget]][License] [![Gitter Widget]][Gitter]

[Travis Widget]: https://img.shields.io/travis/kataras/rizla.svg?style=flat-square
[Travis]: http://travis-ci.org/kataras/rizla
[License Widget]: https://img.shields.io/badge/license-MIT%20%20License%20-E91E63.svg?style=flat-square
[License]: https://github.com/kataras/rizla/blob/master/LICENSE
[Release Widget]: https://img.shields.io/badge/release-v0.0.2-blue.svg?style=flat-square
[Release]: https://github.com/kataras/rizla/releases
[Gitter Widget]: https://img.shields.io/badge/chat-on%20gitter-00BCD4.svg?style=flat-square
[Gitter]: https://gitter.im/kataras/rizla
[Report Widget]: https://img.shields.io/badge/report%20card-A%2B-F44336.svg?style=flat-square
[Report]: http://goreportcard.com/report/kataras/rizla
[Language Widget]: https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square
[Language]: http://golang.org
[Platform Widget]: https://img.shields.io/badge/platform-Any--OS-gray.svg?style=flat-square

# Get Started

```bash
$ rizla main.go #single project monitoring
$ rizla C:/myprojects/project1/main.go C:/myprojects/project2/main.go #multi projects monitoring
```

Want to use it from your project's source code? easy
```sh
$ cat from_code_simple.go
```
```go
package main

import (
	"github.com/kataras/rizla/rizla"
)

func main() {
  // Build, run & start monitoring the projects
  rizla.Run("C:/iris-project/main.go", "C:/otherproject/main.go")
}
```

```sh
$ cat from_code_pro.go
```
```go
package main

import (
	"github.com/kataras/rizla/rizla"
	"time"
	"os"
)

func main() {
  // Optional, set the messages/logs output destination of our application,
  // let's set them to their defaults
  rizla.Out = os.Stdout
  rizla.Err = os.Stderr

  // Create a new project by the main source file
  project := rizla.NewProject("C:/myproject/main.go")

  // The below are optional

  // Provide a Name which will be printed before the 'A change has detected, reloading now...'
  project.Name = "My super project"
  // Allow reload every 3 seconds or more no less
  project.AllowReloadAfter = time.Duration(3) * time.Second

  // Custom subdirectory matcher, for the watcher, return true to include this folder to the watcher
  // the default adds all subdirectories to the watcher, except ".git", "node_modules", "vendor"
  //
  // NOTE: This also executes on runtime if a new folder added, so you calculate and possible 'future' subdirectories.
  project.Watcher = func(absolutePath string) bool {
     return absolutePath != "THIS_SUBDIRECTORY_SHOULD_BE_IGNORED_FROM_THE_WATCHER"
  }

  // Custom file matcher on runtime (file change), return true to reload when a file with this name changed
  project.Matcher = func(filename string) bool {
	 return filename == "I_want_to_reload_only_when_this_file_changed.go"
  }
  // Add arguments, these will be used from the executable file
  project.Args = []string{"-myargument","the value","-otherargument","a value"}
  // Set custom callback when a change to this project is happening
  project.OnChange = func(name string) {
    println("my project's source code "+name+" has been changed, the rizla will care take of the app reloading!!!!!'")
  }

  // End of optional

  // Add the project to the rizla container
  rizla.Add(project)
  //  Build, run & start monitoring the project(s)
  rizla.Run()
}
```


> That's all!

Installation
------------
The only requirement is the [Go Programming Language](https://golang.org/dl)

`$ go get -u github.com/kataras/rizla`

FAQ
------------
Ask questions and get real-time answer from the [Chat](https://gitter.im/kataras/rizla).


Features
------------
- Super easy - is created for everyone.
- You can use it either as command line tool either as part of your project's source code!
- Multi-Monitoring - Supports monitoring of unlimited projects
- No forever loops with filepath.Walk, rizla uses the Operating System's signals to fire a reload.


People
------------
If you'd like to discuss this package, or ask questions about it, feel free to [Chat]( https://gitter.im/kataras/rizla).

The author of rizla is [@kataras](https://github.com/kataras).


Versioning
------------

Current: **v0.0.2**

[HISTORY](https://github.com/kataras/rizla/blob/master/HISTORY.md) file is your best friend!

Read more about Semantic Versioning 2.0.0

 - http://semver.org/
 - https://en.wikipedia.org/wiki/Software_versioning
 - https://wiki.debian.org/UpstreamGuide#Releases_and_Versions


Todo
------------

 - [ ] Tests
 - [ ] Provide full examples.

Third-Party Licenses
------------

Third-Party Licenses can be found [here](THIRDPARTY-LICENSE.md)


License
------------

This project is licensed under the MIT License.

License can be found [here](LICENSE).
