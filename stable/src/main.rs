use std::env;
use std::io::{self, Write};
use std::path::Path;
use std::process::Command;

fn main() -> io::Result<()> {
    // Get home directory for the current user
    let home_dir = env::var("HOME").expect("Failed to get HOME directory");

    // Define paths
    let config_path = format!("{}/open_vpn/config.ovpn", home_dir);
    let auth_path = "/etc/openvpn/auth.txt";

    // Check if config file exists
    if !Path::new(&config_path).exists() {
        eprintln!("Error: OpenVPN config file not found at: {}", config_path);
        return Err(io::Error::new(
            io::ErrorKind::NotFound,
            "Config file not found",
        ));
    }

    // Check if auth file exists
    if !Path::new(auth_path).exists() {
        eprintln!("Error: Auth file not found at: {}", auth_path);
        return Err(io::Error::new(
            io::ErrorKind::NotFound,
            "Auth file not found",
        ));
    }

    println!("Starting OpenVPN with config: {}", config_path);

    // Create command
    let mut cmd = Command::new("openvpn");
    cmd.arg("--config")
        .arg(&config_path)
        .arg("--auth-user-pass")
        .arg(auth_path);

    // Execute the command
    let status = cmd.status()?;

    if status.success() {
        println!("OpenVPN exited successfully");
        Ok(())
    } else {
        eprintln!("OpenVPN exited with error code: {:?}", status.code());
        Err(io::Error::new(
            io::ErrorKind::Other,
            "OpenVPN command failed",
        ))
    }
}
