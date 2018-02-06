package main

import (
	"flag"
	"fmt"
)

// rizla main.go -host myhost.com -port 1193
func main() {
	host := flag.String("host", "", "the host")
	port := flag.Int("port", 0, "the port")
	flag.Parse()
	fmt.Printf("The 'host' argument is: %v\n", *host)
	fmt.Printf("The 'port' argument is: %v\n", *port)
}
