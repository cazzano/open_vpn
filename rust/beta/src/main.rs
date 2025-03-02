use std::env;
use std::process::{Command, Stdio};
use std::thread;
use std::time::Duration;

fn main() {
    // Get current username
    let username = match env::var("USER") {
        Ok(user) => user,
        Err(_) => {
            eprintln!("Error getting current user");
            return;
        }
    };
    
    println!("Starting OpenVPN as root in the background...");
    
    // Spawn a new thread to run OpenVPN in the background
    thread::spawn(move || {
        // Build the OpenVPN command with sudo
        let sudo_process = Command::new("sudo")
            .arg("/usr/sbin/openvpn")
            .arg("--config")
            .arg(format!("/home/{}/open_vpn/config.ovpn", username))
            .arg("--auth-user-pass")
            .arg("/etc/openvpn/auth.txt")
            .stdout(Stdio::null())
            .stderr(Stdio::null())
            .spawn();
            
        match sudo_process {
            Ok(child) => {
                println!("OpenVPN started with PID: {}", child.id());
                
                // We deliberately don't wait for the child process
                // This allows it to continue running in the background
            },
            Err(e) => {
                eprintln!("Error starting OpenVPN command: {}", e);
            }
        }
    });
    
    // Give the thread a moment to start the process
    thread::sleep(Duration::from_secs(1));
    
    println!("OpenVPN launcher has started the process, continuing execution...");
    
    // The main program will exit here, but OpenVPN will continue running in the background
    // If you want the program to stay running, uncomment the line below:
    // std::thread::park(); // This will block the main thread indefinitely
}
