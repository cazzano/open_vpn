package main

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Println("Usage: ./main <command>")
	fmt.Println("Commands:")
	fmt.Println("  init     Initialize VPN configuration")
	fmt.Println("  start    Start the VPN connection")
	fmt.Println("  stop     Stop the VPN connection")
	fmt.Println("  help     Show this help message")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
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
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
