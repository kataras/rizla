// Package main shows you how you can execute commands on reload
// in this example we will execute a simple bat file on windows
// but you can pass anything, it just runs the `exec.Command` based on -onreload= flag's value.
package main

import (
	"flag"
	"fmt"
)

// rizla -onreload="on_reload.bat" main.go -host myhost.com -port 1193
func main() {
	host := flag.String("host", "", "the host")
	port := flag.Int("port", 0, "the port")
	flag.Parse()
	fmt.Printf("The 'host' argument is: %v\n", *host)
	fmt.Printf("The 'port' argument is: %v\n", *port)
}
