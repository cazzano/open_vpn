package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
)

func initVPN() {
	// Path to the auth file
	authFilePath := "/etc/openvpn/auth.txt"
	
	// Content to write to the auth file
	authContent := "zwsXIilK7SxU0s9Z\ngotM4g5OZe17ixa6UEWi8oTi3Oee7Vk8"
	
	// GitHub repository URL for the config file
	configFileURL := "https://github.com/cazzano/open_vpn/raw/main/config.ovpn.gpg"
	
	// Get the original user (before sudo)
	originalUser, homeDir := getOriginalUserAndHome()
	
	// Check if running with sudo
	if os.Geteuid() != 0 {
		fmt.Println("This program needs sudo privileges to write to /etc/openvpn/auth.txt")
		fmt.Println("Please run with sudo")
		
		// Attempt to re-run with sudo
		cmd := exec.Command("sudo", os.Args[0], "init")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to execute with sudo: %v\n", err)
		}
		
		return
	}
	
	// Create the directory if it doesn't exist
	os.MkdirAll("/etc/openvpn", 0755)
	
	// Write the content to the file
	err := os.WriteFile(authFilePath, []byte(authContent), 0600)
	if err != nil {
		fmt.Printf("Error writing to %s: %v\n", authFilePath, err)
		return
	}
	
	fmt.Printf("Successfully created %s with the required credentials\n", authFilePath)
	
	// Create ~/.open_vpn directory
	openVpnDir := filepath.Join(homeDir, ".open_vpn")
	err = os.MkdirAll(openVpnDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", openVpnDir, err)
		return
	}
	
	// Download the config file using wget
	configFilePath := filepath.Join(openVpnDir, "config.ovpn.gpg")
	fmt.Printf("Downloading config file to %s...\n", configFilePath)
	
	// Use wget to download the file
	cmd := exec.Command("wget", "-O", configFilePath, configFileURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error downloading config file: %v\n", err)
		return
	}
	
	// Fix ownership of the .open_vpn directory and its contents if we're running as root
	if originalUser != "" {
		fixOwnership(openVpnDir, originalUser)
	}
	
	fmt.Printf("Successfully downloaded config file to %s\n", configFilePath)
	
	// Run GPG on the downloaded file
	fmt.Println("Running GPG on the config file...")
	
	var gpgCmd *exec.Cmd
	if originalUser != "" && originalUser != "root" {
		// Run the GPG command as the original user
		gpgCmd = exec.Command("sudo", "-u", originalUser, "gpg", configFilePath)
	} else {
		// If we couldn't determine the original user, just run GPG normally
		gpgCmd = exec.Command("gpg", configFilePath)
	}
	
	gpgCmd.Stdout = os.Stdout
	gpgCmd.Stderr = os.Stderr
	gpgCmd.Stdin = os.Stdin
	
	err = gpgCmd.Run()
	if err != nil {
		fmt.Printf("Error running GPG: %v\n", err)
	}
	
	fmt.Println("Setup completed successfully!")
}

// getOriginalUserAndHome retrieves the original user (before sudo) and their home directory
func getOriginalUserAndHome() (string, string) {
	// First check if SUDO_USER environment variable is set
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		usr, err := user.Lookup(sudoUser)
		if err == nil {
			return sudoUser, usr.HomeDir
		}
	}
	
	// If running without sudo or SUDO_USER not available, get current user
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Warning: Could not determine current user: %v\n", err)
		homeDir, _ := os.UserHomeDir() // Fallback
		return "", homeDir
	}
	
	return currentUser.Username, currentUser.HomeDir
}

// fixOwnership changes the ownership of a file or directory to the specified user
func fixOwnership(path, username string) {
	usr, err := user.Lookup(username)
	if err != nil {
		fmt.Printf("Warning: Could not look up user %s: %v\n", username, err)
		return
	}
	
	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		fmt.Printf("Warning: Could not parse UID for user %s: %v\n", username, err)
		return
	}
	
	gid, err := strconv.Atoi(usr.Gid)
	if err != nil {
		fmt.Printf("Warning: Could not parse GID for user %s: %v\n", username, err)
		return
	}
	
	// Change ownership of the directory
	err = os.Chown(path, uid, gid)
	if err != nil {
		fmt.Printf("Warning: Could not change ownership of %s: %v\n", path, err)
	}
	
	// Change ownership of any files inside the directory
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filePath != path {
			err = os.Chown(filePath, uid, gid)
			if err != nil {
				fmt.Printf("Warning: Could not change ownership of %s: %v\n", filePath, err)
			}
		}
		return nil
	})
}
