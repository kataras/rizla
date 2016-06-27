## 0.0.1 -> 0.0.2

- A lot of underline code improvements & fixes
- New: `project.Watcher(string) bool`
- New: `project.OnChange(string)`
- New: Allow watching new directories in runtime
- New: Rizla accepts all fs os events as change events
- Fix: Not watching third-level subdirectories, now watch everything except ` ".git", "node_modules", "vendor"` (you can change this behavior with the `project.Watcher`)
- Maybe I'm missing something, just upgrade with `go get -u github.com/kataras/rizla` and have fun :)
