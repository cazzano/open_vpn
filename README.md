# 🔐 SVPN
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**SVPN** is a lightweight command-line tool for managing VPN connections with simplicity and efficiency.

## ✨ Features
- **Easy Configuration** - Initialize VPN settings with a single command
- **Background Operation** - Run your VPN connection in the background
- **Simple Management** - Start and stop your VPN connection effortlessly
- **Streamlined Interface** - Minimal commands for maximum productivity

## 📦 Installation
```bash
# Clone the repository
git clone https://github.com/cazzano/open_vpn.git

# Navigate to the directory
cd open_vpn/go/beta/src

# Install the tool
go build && mv src svpn && sudo mv svpn /usr/bin && echo "SVPN successfully installed!"
```

## 🚀 Usage
SVPN provides straightforward commands to manage your VPN connections:

### Initialize VPN Configuration
```bash
svpn init
```
This will set up the necessary VPN configuration files.

### Start VPN Connection
```bash
svpn start
```
Starts the VPN connection and runs it in the background.

### Stop VPN Connection
```bash
svpn stop
```
Terminates the active VPN connection running in the background.

### Display Help Information
```bash
svpn --h
```
Shows available commands and usage examples.

### Display Version Information
```bash
svpn --v
```
Shows the current version of SVPN along with author and repository information.

## 📝 Examples
### Example 1: Complete VPN Workflow
```bash
# Initialize VPN configurations
svpn init

# Start the VPN connection
svpn start

# Later, when you're done
svpn stop
```

## 🛠️ Command Reference
| Command | Description |
|---------|-------------|
| `init` | Initialize VPN configuration |
| `start` | Start the VPN connection in the background |
| `stop` | Stop the active VPN connection |
| `--h` | Display help message |
| `--v` | Display version information |

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support
For support, please open an issue in the GitHub repository or contact the maintainers.
