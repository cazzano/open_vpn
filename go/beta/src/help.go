// help.go
package main

import (
	"fmt"
)

// PrintUsage prints the usage information for the application.
func PrintUsage() {
	fmt.Println("Usage: svpn <command>")
	fmt.Println("Commands:")
	fmt.Println("  init     Initialize VPN configuration")
	fmt.Println("  start    Start the VPN connection")
	fmt.Println("  stop     Stop the VPN connection")
	fmt.Println("  --h      Show this help message")
	fmt.Println("  --v      Display version information")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  1. svpn init     It will initialize the VPN configs")
	fmt.Println("  2. svpn start    It will start the VPN in background")
	fmt.Println("  3. svpn stop     It will stop the VPN in background")
}
