package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func readPIDInfo(configPath string) (*PIDInfo, error) {
	pidFile := filepath.Join(configPath, "pid.json")
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return nil, fmt.Errorf("error reading PID file: %v", err)
	}

	var pidInfo PIDInfo
	if err := json.Unmarshal(data, &data); err != nil {
		return nil, fmt.Errorf("error parsing PID file: %v", err)
	}

	return &pidInfo, nil
}

func killVPN() {
	// Check sudo permissions first
	fmt.Println("Checking sudo permissions...")
	if err := checkSudo(); err != nil {
		fmt.Println("Error: This program requires sudo privileges")
		fmt.Println("Please run with sudo or enter your password when prompted")
		os.Exit(1)
	}

	// Get config directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	configPath := filepath.Join(homeDir, ".config", "secret_vpn")
	pidFile := filepath.Join(configPath, "pid.json")

	// Check if PID file exists
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		fmt.Println("No VPN process found (PID file does not exist)")
		os.Exit(1)
	}

	// Read the PID file
	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println("Error reading PID file:", err)
		fmt.Println("Is the VPN running?")
		os.Exit(1)
	}

	var pidInfo PIDInfo
	if err := json.Unmarshal(data, &pidInfo); err != nil {
		fmt.Println("Error parsing PID file:", err)
		// If PID file is corrupted, remove it
		os.Remove(pidFile)
		os.Exit(1)
	}

	// Check if the process exists
	process, err := os.FindProcess(pidInfo.PID)
	if err != nil {
		fmt.Printf("Process with PID %d not found\n", pidInfo.PID)
		// Clean up the PID file
		os.Remove(pidFile)
		os.Exit(1)
	}

	// Try to kill the process
	fmt.Printf("Attempting to kill OpenVPN process (PID: %d)...\n", pidInfo.PID)
	
	// First try SIGTERM for graceful shutdown
	cmd := exec.Command("sudo", "kill", "-TERM", fmt.Sprintf("%d", pidInfo.PID))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: SIGTERM failed, attempting force kill: %v\n", err)
		// If SIGTERM fails, try SIGKILL
		cmd = exec.Command("sudo", "kill", "-9", fmt.Sprintf("%d", pidInfo.PID))
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error: Failed to kill process: %v\n", err)
			os.Exit(1)
		}
	}

	// Wait a moment to ensure the process is killed
	time.Sleep(2 * time.Second)

	// Verify the process is killed by sending signal 0
	if err := process.Signal(syscall.Signal(0)); err == nil {
		fmt.Println("Warning: Process might still be running")
		fmt.Println("Please check the process status manually")
	} else {
		fmt.Println("OpenVPN process successfully terminated")
	}

	// Remove the PID file
	if err := os.Remove(pidFile); err != nil {
		fmt.Printf("Warning: Could not remove PID file: %v\n", err)
	} else {
		fmt.Println("PID file removed successfully")
	}

	fmt.Println("VPN shutdown complete")
}
