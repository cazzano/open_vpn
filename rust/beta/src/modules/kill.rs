use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;
use std::process::{Command, exit};
use std::thread::sleep;
use std::time::Duration;

#[derive(Serialize, Deserialize, Debug)]
struct PidInfo {
    pid: i32,
    // Add other fields that might be in the original PidInfo
}

fn check_sudo() -> Result<(), String> {
    let output = Command::new("id")
        .arg("-u")
        .output()
        .map_err(|e| format!("Failed to execute command: {}", e))?;

    let uid = String::from_utf8_lossy(&output.stdout)
        .trim()
        .parse::<u32>()
        .map_err(|e| format!("Failed to parse UID: {}", e))?;

    if uid == 0 {
        Ok(())
    } else {
        Err("This program requires sudo privileges".to_string())
    }
}

fn read_pid_info(config_path: &PathBuf) -> Result<PidInfo, String> {
    let pid_file = config_path.join("pid.json");
    let data =
        fs::read_to_string(&pid_file).map_err(|e| format!("Error reading PID file: {}", e))?;

    let pid_info: PidInfo =
        serde_json::from_str(&data).map_err(|e| format!("Error parsing PID file: {}", e))?;

    Ok(pid_info)
}

fn kill_vpn() {
    // Check sudo permissions first
    println!("Checking sudo permissions...");
    if let Err(e) = check_sudo() {
        println!("Error: {}", e);
        println!("Please run with sudo or enter your password when prompted");
        exit(1);
    }

    // Get config directory path
    let home_dir = match dirs::home_dir() {
        Some(path) => path,
        None => {
            println!("Error getting home directory");
            exit(1);
        }
    };

    let config_path = home_dir.join(".config").join("secret_vpn");
    let pid_file = config_path.join("pid.json");

    // Check if PID file exists
    if !pid_file.exists() {
        println!("No VPN process found (PID file does not exist)");
        exit(1);
    }

    // Read the PID file
    let pid_info = match read_pid_info(&config_path) {
        Ok(info) => info,
        Err(e) => {
            println!("Error: {}", e);
            println!("Is the VPN running?");

            // If PID file is corrupted, remove it
            if e.contains("Error parsing PID file") {
                let _ = fs::remove_file(&pid_file);
            }

            exit(1);
        }
    };

    // Try to kill the process
    println!(
        "Attempting to kill OpenVPN process (PID: {})...",
        pid_info.pid
    );

    // First try SIGTERM for graceful shutdown
    let status = Command::new("sudo")
        .args(["kill", "-TERM", &pid_info.pid.to_string()])
        .status();

    if let Err(e) = status {
        println!("Warning: SIGTERM failed: {}", e);
        exit(1);
    } else if !status.unwrap().success() {
        println!("Warning: SIGTERM failed, attempting force kill");

        // If SIGTERM fails, try SIGKILL
        let status = Command::new("sudo")
            .args(["kill", "-9", &pid_info.pid.to_string()])
            .status();

        if let Err(e) = status {
            println!("Error: Failed to kill process: {}", e);
            exit(1);
        } else if !status.unwrap().success() {
            println!("Error: Failed to kill process");
            exit(1);
        }
    }

    // Wait a moment to ensure the process is killed
    sleep(Duration::from_secs(2));

    // Verify the process is killed by checking if it exists
    let status = Command::new("ps")
        .args(["-p", &pid_info.pid.to_string()])
        .status();

    match status {
        Ok(exit_status) => {
            if exit_status.success() {
                println!("Warning: Process might still be running");
                println!("Please check the process status manually");
            } else {
                println!("OpenVPN process successfully terminated");
            }
        }
        Err(_) => {
            println!("OpenVPN process successfully terminated");
        }
    }

    // Remove the PID file
    match fs::remove_file(&pid_file) {
        Ok(_) => println!("PID file removed successfully"),
        Err(e) => println!("Warning: Could not remove PID file: {}", e),
    }

    println!("VPN shutdown complete");
}

fn main() {
    kill_vpn();
}
