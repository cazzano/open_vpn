package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sync"
	"time"
)

func main() {
	// Get current username
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}
	
	username := currentUser.Username
	
	// Wait group to prevent the main program from exiting immediately
	var wg sync.WaitGroup
	wg.Add(1)
	
	// Run OpenVPN in a goroutine
	go func() {
		defer wg.Done()
		
		// Build the OpenVPN command
		sudoCmd := exec.Command("sudo", append([]string{"/usr/sbin/openvpn"}, 
			"--config", "/home/" + username + "/open_vpn/config.ovpn",
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
		
		// Optional: print the process ID
		fmt.Printf("OpenVPN started with PID: %d\n", sudoCmd.Process.Pid)
		
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
