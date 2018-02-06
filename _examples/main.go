package main

import (
	"flag"
	"fmt"
)

func main() {
	host := flag.String("host", "", "the host")
	flag.Parse()
	fmt.Printf("The 'host' argument is: %v\n", *host)
}
