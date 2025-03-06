use serde::{Deserialize, Serialize};
use std::fs;
use std::os::unix::process::CommandExt;
use std::path::{Path, PathBuf};
use std::process::{Child, Command, exit};
use std::thread;
use std::time::{Duration, SystemTime};

// PIDInfo stores the process information
#[derive(Serialize, Deserialize, Debug)]
struct PidInfo {
    pid: i32,
    start_time: SystemTime,
}

fn check_sudo() -> Result<(), String> {
    let status = Command::new("sudo")
        .arg("-v")
        .status()
        .map_err(|e| format!("Failed to execute sudo command: {}", e))?;

    if status.success() {
        Ok(())
    } else {
        Err("Failed to validate sudo privileges".to_string())
    }
}

fn ensure_config_dir(config_path: &Path) -> Result<(), String> {
    fs::create_dir_all(config_path)
        .map_err(|e| format!("Error creating config directory: {}", e))?;

    // Change permissions to 0700 (rwx------)
    #[cfg(unix)]
    {
        use std::os::unix::fs::PermissionsExt;
        let permissions = fs::Permissions::from_mode(0o700);
        fs::set_permissions(config_path, permissions)
            .map_err(|e| format!("Failed to set directory permissions: {}", e))?;
    }

    Ok(())
}

fn save_pid(pid: i32, config_path: &Path) -> Result<(), String> {
    let pid_info = PidInfo {
        pid,
        start_time: SystemTime::now(),
    };

    let data = serde_json::to_string(&pid_info)
        .map_err(|e| format!("Error serializing PID info: {}", e))?;

    let pid_file = config_path.join("pid.json");
    fs::write(&pid_file, data).map_err(|e| format!("Error writing PID file: {}", e))?;

    // Set file permissions to 0600 (rw-------)
    #[cfg(unix)]
    {
        use std::os::unix::fs::PermissionsExt;
        let permissions = fs::Permissions::from_mode(0o600);
        fs::set_permissions(&pid_file, permissions)
            .map_err(|e| format!("Failed to set file permissions: {}", e))?;
    }

    Ok(())
}

fn check_existing_vpn(config_path: &Path) -> Result<bool, String> {
    let pid_file = config_path.join("pid.json");

    // Check if PID file exists
    if !pid_file.exists() {
        return Ok(false);
    }

    // Read the PID file
    let data =
        fs::read_to_string(&pid_file).map_err(|e| format!("Error reading PID file: {}", e))?;

    let pid_info: PidInfo = match serde_json::from_str(&data) {
        Ok(info) => info,
        Err(_) => {
            // If PID file is corrupted, remove it
            let _ = fs::remove_file(&pid_file);
            return Ok(false);
        }
    };

    // Check if the process exists and is running
    let output = Command::new("ps")
        .args(["-p", &pid_info.pid.to_string()])
        .output()
        .map_err(|e| format!("Failed to check process status: {}", e))?;

    if !output.status.success() {
        // Process is not running, clean up PID file
        let _ = fs::remove_file(&pid_file);
        return Ok(false);
    }

    // Process exists and is running
    Ok(true)
}

fn start_openvpn(username: &str, config_path: &Path) -> Result<Child, String> {
    let config_file = format!("/home/{}/.open_vpn/config.ovpn", username);
    let auth_file = "/etc/openvpn/auth.txt";

    // Build the OpenVPN command
    let mut cmd = Command::new("sudo");
    cmd.args([
        "/usr/sbin/openvpn",
        "--config",
        &config_file,
        "--auth-user-pass",
        auth_file,
    ]);

    // Don't capture stdout/stderr to let it run in background
    cmd.stdout(std::process::Stdio::null());
    cmd.stderr(std::process::Stdio::null());

    // Start the process
    let child = cmd
        .spawn()
        .map_err(|e| format!("Error starting OpenVPN command: {}", e))?;

    Ok(child)
}

pub fn main_vpn() {
    // Check sudo permissions first
    println!("Checking sudo permissions...");
    if let Err(e) = check_sudo() {
        println!("Error: This program requires sudo privileges");
        println!("Please run with sudo or enter your password when prompted");
        println!("Error details: {}", e);
        exit(1);
    }

    // Get current username
    let username = match std::env::var("USER") {
        Ok(name) => name,
        Err(_) => {
            println!("Error getting current user");
            exit(1);
        }
    };

    // Create config directory path
    let home_dir = match dirs::home_dir() {
        Some(path) => path,
        None => {
            println!("Error getting home directory");
            exit(1);
        }
    };

    let config_path = home_dir.join(".config").join("secret_vpn");

    // Ensure config directory exists
    if let Err(e) = ensure_config_dir(&config_path) {
        println!("Error creating config directory: {}", e);
        exit(1);
    }

    // Check for existing VPN process
    match check_existing_vpn(&config_path) {
        Ok(true) => {
            println!("VPN is already running");
            println!("Use './vpn stop' to stop the existing VPN before starting a new one");
            exit(1);
        }
        Ok(false) => {
            // Continue execution, VPN is not running
        }
        Err(e) => {
            println!("Error checking existing VPN process: {}", e);
            exit(1);
        }
    }

    // Run OpenVPN
    println!("Starting OpenVPN as root in the background...");
    match start_openvpn(&username, &config_path) {
        Ok(child) => {
            let pid = child.id() as i32;

            // Save the PID
            if let Err(e) = save_pid(pid, &config_path) {
                println!("Error saving PID: {}", e);
            } else {
                println!(
                    "OpenVPN started with PID: {} (saved to {})",
                    pid,
                    config_path.join("pid.json").display()
                );
            }

            // Process runs in background, we don't wait for it
            thread::spawn(move || {
                if let Err(e) = child.wait_with_output() {
                    println!("OpenVPN process exited with error: {}", e);
                }
            });
        }
        Err(e) => {
            println!("Error starting OpenVPN: {}", e);
            exit(1);
        }
    }

    // Give the thread a moment to start the process
    thread::sleep(Duration::from_secs(1));

    println!("OpenVPN launcher has started the process, continuing execution...");

    // Program exits here, but OpenVPN continues running in the background
}

fn main() {
    main_vpn();
}
