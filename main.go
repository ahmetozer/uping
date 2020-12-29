package main

import (
	"fmt"
	"os"

	client "github.com/ahmetozer/uping/client"
	server "github.com/ahmetozer/uping/server"
)

var (
	helpText = `
Usage: oping mode -help
Modes:
	server: Run as server mode
	client: Run as client mode

`
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide program mode")
		fmt.Println(helpText)
		os.Exit(1)
	}

	switch os.Args[1] {

	case "client":
		client.Main(os.Args[2:])
	case "server":
		server.Main(os.Args[2:])
	default:
		fmt.Println("Unknown Program mode")
		fmt.Println(helpText)
		os.Exit(1)
	}
}
