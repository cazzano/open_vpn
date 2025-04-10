package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		initVPN()
	case "start":
		main_vpn()
	case "stop":
		killVPN()
	case "--h":
		PrintUsage()
	case "--v":
		DisplayVersion()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		PrintUsage()
		os.Exit(1)
	}
}
