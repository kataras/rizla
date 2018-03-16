## 0.1.1

Add `-onreload` to execute commands on reload through the `rizla` cli.

Example:

**on_reload.bat**

```sh
@echo off
echo Hello, custom script can goes here before reload, i.e build https://github.comkataras/bindata things!
```

**main.go**

```go
// Package main shows you how you can execute commands on reload
// in this example we will execute a simple bat file on windows
// but you can pass anything, it just runs the `exec.Command` based on -onreload= flag's value.
package main

import (
	"flag"
	"fmt"
)

func main() {
	host := flag.String("host", "", "the host")
	port := flag.Int("port", 0, "the port")
	flag.Parse()
	fmt.Printf("The 'host' argument is: %v\n", *host)
	fmt.Printf("The 'port' argument is: %v\n", *port)
}

```

```sh
$ rizla -onreload="on_reload.bat" main.go -host myhost.com -port 1193
```

## 0.1.0

Rizla drops support for multi `main.go:func main()` programs in the same directory. It still accepts a filename with `main.go` but it depends on the directory now (as all examples already shown) in order to be able to run and watch projects with multiple `.go` files in the project's root directory, this is very useful especially when the project depends on libraries like `go-bindata` with a result of `.go` file in the root directory. For most cases that will not change anything. If you used to have many go programs with `func main()` in the same root directory please consider that this is not idiomatic and you must change this habit, the sooner the better.

> `main.go` can be any filename contains the `func main(){ ... }` of course.

## 0.0.9

1. `rizla.Run()` -> `rizla.Run(map[string][]string)`. Run now accepts a `sources map[string][]string`, can be nil if projects added manually previously, the `key string` is the program filepath with `.go` extension and the `values []string` are any optional arguments that the program excepts to be passed. Therefore use `Run(nil)` if you used `rizla.Run()` before.

2. `rizla.RunWith(watcher rizla.Watcher, programFile string)` -> `rizla.RunWith(watcher rizla.Watcher, sources map[string][]string, delayOnDetect time.Duration)`.

3. At `rizla#Project` the property `AllowRunAfter time.Duration` added.

4. New `-delay` cli flag added, as requested by @scorpnode at https://github.com/kataras/rizla/issues/14

## 0.0.8

Support flags as requested at [#13](https://github.com/kataras/rizla/issues/13) by @Zeno-Code.

## 0.0.6 -> 0.0.7

Rizla uses the operating system's signals to fire a change because it is the fastest way and it consumes the minimal CPU.
But as the [feature request](https://github.com/kataras/rizla/issues/6) explains, some IDEs overrides the Operating System's signals, so I needed to change the things a bit in order to allow
looping-every-while and compare the file(s) with its modtime for new changes while in the same time keep the default as it's.

- **NEW**: Add a common interface for file system watchers in order to accoblish select between two methods of scanning for file changes.
    - file system's signals (default)
    - `filepath.Walk` (using the `-walk` flag)

### When to enable `-walk`?
When the default method doesn't works for your IDE's save method.

### How to enable `-walk`?
- If you're command line user: `rizla -walk main.go` , just append the `-walk` flag.
- If you use rizla behind your source code then use the `rizla.RunWith(rizla.WatcherFromFlag("-flag"))` instead of `rizla.Run()`.


## 0.0.4 -> 0.0.5 and 0.0.6

- **Fix**: Reload more than one time on Mac

## 0.0.3 -> 0.0.4

- **Added**: global `DefaultDisableProgramRerunOutput` and per-project `DisableProgramRerunOutput` option, to disable the program's output when reloads.
- **Fix** two-times reload on windows
- **Fix** fail to re-run when previous build error, issue [here](https://github.com/kataras/rizla/issues/1)

## 0.0.2 -> 0.0.3
- Change: `rizla.Out & rizla.Err` moved to the Project iteral (`project.Out` & `project.Err`), each project can have its own output now, and are type of *Printer
- Change/NEW: `project.OnChange` removed, new  `project.OnReload(func(string))` & `project.OnReloaded(func(string))` , with these you can change the output message when a file has changed and also when project reloaded, see the [project.go](https://github.com/kataras/rizla/blob/master/project.go) for more.

## 0.0.1 -> 0.0.2

- A lot of underline code improvements & fixes
- New: `project.Watcher(string) bool`
- New: `project.OnChange(string)`
- New: Allow watching new directories in runtime
- New: Rizla accepts all fs os events as change events
- Fix: Not watching third-level subdirectories, now watch everything except ` ".git", "node_modules", "vendor"` (you can change this behavior with the `project.Watcher`)
- Maybe I'm missing something, just upgrade with `go get -u github.com/kataras/rizla` and have fun :)
