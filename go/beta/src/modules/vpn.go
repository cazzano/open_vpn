package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sync"
	"time"
)

// PIDInfo stores the process information
type PIDInfo struct {
	PID       int       `json:"pid"`
	StartTime time.Time `json:"start_time"`
}

func ensureConfigDir(configPath string) error {
	return os.MkdirAll(configPath, 0700)
}

func savePID(pid int, configPath string) error {
	pidInfo := PIDInfo{
		PID:       pid,
		StartTime: time.Now(),
	}

	data, err := json.Marshal(pidInfo)
	if err != nil {
		return fmt.Errorf("error marshaling PID info: %v", err)
	}

	pidFile := filepath.Join(configPath, "pid.json")
	err = os.WriteFile(pidFile, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing PID file: %v", err)
	}

	return nil
}

func main() {
	// Get current username
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	username := currentUser.Username

	// Create config directory path
	configPath := filepath.Join(currentUser.HomeDir, ".config", "secret_vpn")

	// Ensure config directory exists
	if err := ensureConfigDir(configPath); err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}

	// Wait group to prevent the main program from exiting immediately
	var wg sync.WaitGroup
	wg.Add(1)

	// Run OpenVPN in a goroutine
	go func() {
		defer wg.Done()

		// Build the OpenVPN command
		sudoCmd := exec.Command("sudo", append([]string{"/usr/sbin/openvpn"},
			"--config", "/home/"+username+"/open_vpn/config.ovpn",
			"--auth-user-pass", "/etc/openvpn/auth.txt")...)

		// Redirect stdout and stderr to the background
		sudoCmd.Stdout = nil
		sudoCmd.Stderr = nil

		fmt.Println("Starting OpenVPN as root in the background...")

		// Execute the command
		err := sudoCmd.Start()
		if err != nil {
			fmt.Println("Error starting OpenVPN command:", err)
			return
		}

		// Get and save the PID
		pid := sudoCmd.Process.Pid
		if err := savePID(pid, configPath); err != nil {
			fmt.Println("Error saving PID:", err)
		} else {
			fmt.Printf("OpenVPN started with PID: %d (saved to %s)\n", pid, filepath.Join(configPath, "pid.json"))
		}

		// Process runs in background, we don't wait for it
		// But we can optionally report its status if needed
		go func() {
			err := sudoCmd.Wait()
			if err != nil {
				fmt.Println("OpenVPN process exited with error:", err)
			}
		}()
	}()

	// Give the goroutine a moment to start the process
	time.Sleep(1 * time.Second)

	fmt.Println("OpenVPN launcher has started the process, continuing execution...")

	// If you want the program to exit after launching OpenVPN,
	// you can comment out the wg.Wait() below
	// Uncomment it if you want the program to stay running

	// wg.Wait()

	// The program will exit here, but OpenVPN will continue running in the background
}
