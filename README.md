# go-port-tester

`go-port-tester` is a simple Go project designed to test the status of specified ports on multiple servers by making HTTP requests.  
This project is helpful for network administrators or developers who need to verify the accessibility of various services running on different servers and ports.

[한국어](README.ko.md)

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

### Features
- 🌐 **Multi-Server, Multi-Port Checks**: Test multiple IPs and ports simultaneously.
- ⚡ **Fast & Concurrent**: Adjustable concurrency for faster testing, optimized for large networks.
- 📊 **Easy Results**: Save results in a clean CSV format, ready for analysis.
- 🛠️ **Cross-Platform Compatibility**: Build and run on Windows, macOS, and Linux.

---

### Table of Contents
- [Requirements](#requirements)
- [Installation](#installation)
    - [Installing Go](#installing-go)
    - [Installing Make](#installing-make)
- [Build Instructions](#build-instructions)
- [Running the Application](#running-the-application)
- [Usage Examples](#usage-examples)

---

### Requirements

- **Go** (version 1.20 or higher)
- **GNU Make**

### Installation

#### Installing Go
- **Windows**: Download and install Go from the [official Go website](https://golang.org/dl/).
- **macOS**: Use Homebrew to install:
  ```bash
  brew install go
  ```
- **Linux**: Follow the instructions on the [official Go website](https://golang.org/doc/install) or use a package manager like `apt`:
  ```bash
  sudo apt update
  sudo apt install -y golang
  ```

#### Installing Make
Make is typically pre-installed on Linux and macOS. For Windows, follow these steps:

- **Windows**: Install Make via Chocolatey:
  ```powershell
  choco install make
  ```
  or install via Scoop:
  ```powershell
  scoop install make
  ```

- **macOS**: Install with Homebrew:
  ```bash
  brew install make
  ```

- **Linux**: Make is commonly pre-installed. To install, run:
  ```bash
  sudo apt install -y make
  ```

---

### Build Instructions

To build the application, use the provided `Makefile` for cross-platform compatibility.

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/go-port-tester.git
   cd go-port-tester
   ```

2. **Build for all platforms**:
   ```bash
   make all
   ```

3. **Platform-Specific Builds**:
    - **Windows**:
      ```bash
      make build-windows
      ```
    - **macOS**:
      ```bash
      make build-macos
      ```
    - **Linux**:
      ```bash
      make build-linux
      ```

The built binaries will be located in the `build` directory, organized by OS.

---

### Running the Application

#### Server (Test Subject)
The server component will host a simple HTTP service on specified ports for testing purposes.

```bash
./build/<platform>/server --config server.csv
```

#### Client (Port Tester)
The client component will send requests to specified IP and port combinations to check their status.

```bash
./build/<platform>/client --servers servers.csv --output result.csv --timeout 2 --concurrency 5
```

---

### Usage Examples

**Sample CSV for Servers** (`servers.csv`):
```csv
IP,Port
127.0.0.1,8080
192.168.1.10,8081
```
