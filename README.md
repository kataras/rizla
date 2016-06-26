Rizla builds, runs and monitors your Go Applications with ease.

[![Travis Widget]][Travis] [![Release Widget]][Release] [![Report Widget]][Report] [![License Widget]][License] [![Gitter Widget]][Gitter]

[Travis Widget]: https://img.shields.io/travis/kataras/rizla.svg?style=flat-square
[Travis]: http://travis-ci.org/kataras/rizla
[License Widget]: https://img.shields.io/badge/license-MIT%20%20License%20-E91E63.svg?style=flat-square
[License]: https://github.com/kataras/rizla/blob/master/LICENSE
[Release Widget]: https://img.shields.io/badge/release-v0.0.1-blue.svg?style=flat-square
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
)

func main() {
  // Create a new project by the main source file
  project := rizla.NewProject("C:/myproject/main.go")
  // Provide a Name which will be printed before the 'A change has detected, reloading now...'
  project.Name = "My super project"
  // Allow reload every 2 seconds or more no less
  project.AllowReloadAfter = 2*time.Second
  // Custom file matcher
  project.Matcher = func(filename string) bool {
	 return filename == "I_want_to_reload_only_when_this_file_changed.go"
  }

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
You can find answers by exploring [these questions](https://github.com/kataras/rizla/issues?q=label%3Aquestion).


Features
------------
- Super easy - is created for everyone.
- You can use it either as command line tool either as part of your project's source code!
- Multi-Monitoring - Supports monitoring of unlimited projects
- No forever loops with filepath.Walk, Rizla uses the Operating System's signals to fire a reload.


People
------------
If you'd like to discuss this package, or ask questions about it, feel free to [Chat]( https://gitter.im/kataras/rizla).

The author of rizla is [@kataras](https://github.com/kataras).


Versioning
------------

Current: **v0.0.1**


Read more about Semantic Versioning 2.0.0

 - http://semver.org/
 - https://en.wikipedia.org/wiki/Software_versioning
 - https://wiki.debian.org/UpstreamGuide#Releases_and_Versions


Third-Party Licenses
------------

Third-Party Licenses can be found [here](THIRDPARTY-LICENSE.md)


License
------------

This project is licensed under the MIT License.

License can be found [here](LICENSE).
