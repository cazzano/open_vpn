package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

// PIDInfo stores the process information
type PIDInfo struct {
	PID       int       `json:"pid"`
	StartTime time.Time `json:"start_time"`
}

func readPIDFile(configPath string) (*PIDInfo, error) {
	pidFile := filepath.Join(configPath, "pid.json")
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return nil, fmt.Errorf("error reading PID file: %v", err)
	}

	var pidInfo PIDInfo
	if err := json.Unmarshal(data, &pidInfo); err != nil {
		return nil, fmt.Errorf("error parsing PID file: %v", err)
	}

	return &pidInfo, nil
}

func killProcess(pid int) error {
	// First try to kill the OpenVPN process directly
	sudoCmd := exec.Command("sudo", "kill", fmt.Sprintf("%d", pid))
	if err := sudoCmd.Run(); err != nil {
		return fmt.Errorf("error killing process with sudo: %v", err)
	}
	return nil
}

func cleanupPIDFile(configPath string) error {
	pidFile := filepath.Join(configPath, "pid.json")
	if err := os.Remove(pidFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing PID file: %v", err)
	}
	return nil
}

func main() {
	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	// Get config directory path
	configPath := filepath.Join(currentUser.HomeDir, ".config", "secret_vpn")

	// Read PID file
	pidInfo, err := readPIDFile(configPath)
	if err != nil {
		fmt.Println("Error reading PID file:", err)
		os.Exit(1)
	}

	fmt.Printf("Found OpenVPN process with PID: %d\n", pidInfo.PID)

	// Try to kill the process
	if err := killProcess(pidInfo.PID); err != nil {
		fmt.Println("Error killing process:", err)
		os.Exit(1)
	}

	// Clean up the PID file
	if err := cleanupPIDFile(configPath); err != nil {
		fmt.Println("Error cleaning up PID file:", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully terminated OpenVPN process (PID: %d)\n", pidInfo.PID)
}
