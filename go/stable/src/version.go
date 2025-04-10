// version.go
package main

import (
	"fmt"
)

// DisplayVersion prints the version information of the project.
func DisplayVersion() {
	// Define the version and any other relevant information
	version := "1.0" // Update this to your actual version
	author := "cazzano" // Replace with your name or organization
	repository := "https://github.com/cazzano/open_vpn.git" // Replace with your repository link

	// Print the version information
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Author: %s\n", author)
	fmt.Printf("Repository: %s\n", repository)
}
