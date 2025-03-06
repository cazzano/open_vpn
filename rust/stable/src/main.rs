use std::env;
use std::process::exit;

// Import the two modules
mod vpn_kill;
mod vpn_start;

fn print_usage() {
    println!("Usage: ./vpn <command>");
    println!("Commands:");
    println!("  start    Start the VPN connection");
    println!("  stop     Stop the VPN connection");
    println!("  help     Show this help message");
}

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        print_usage();
        exit(1);
    }

    let command = &args[1];

    match command.as_str() {
        "start" => vpn_start::main_vpn(), // Updated function name based on your original code
        "stop" => vpn_kill::kill_vpn(),
        "help" => print_usage(),
        _ => {
            println!("Unknown command: {}", command);
            print_usage();
            exit(1);
        }
    }
}
