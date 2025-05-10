# SimplePortScanner

SimplePortScanner is a cross-platform, high-performance port scanner tool that can identify open ports on a target IP or domain name. It provides both a command-line interface and an optional web-based interface.

## Features

- Fast, concurrent port scanning
- Support for scanning individual ports, port ranges, and port lists
- Banner grabbing to identify services running on open ports
- Color-coded command-line output
- Web-based interface for easy use
- Cross-platform compatibility (Windows, macOS, Linux)

## Installation

### Prerequisites

- Go 1.16 or newer

### Building from source

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/portscanner.git
   cd portscanner
   ```

2. Build the CLI tool:
   ```
   go build -o portscanner ./cmd/cli
   ```

3. (Optional) Build the web server:
   ```
   go build -o portscanner-web ./cmd/web
   ```

## Usage

### Command Line Interface

```
./portscanner --target example.com --ports 80,443,8000-8100
```

Available options:

- `--target`: The target IP address or domain name (required)
- `--ports`: Ports to scan (default: "1-1000")
  - Single port: `80`
  - Port range: `1-1000`
  - Port list: `22,80,443`
  - Combinations: `22,80,1000-2000`
- `--threads`: Number of concurrent scanning threads (default: 100)
- `--timeout`: Timeout for connection attempts (default: 2s)

### Web Interface

1. Start the web server:
   ```
   ./portscanner-web --port 8080
   ```

2. Open your browser and navigate to `http://localhost:8080`

3. Enter the target and port specifications, then click "Start Scan"

## Example Output

### CLI Output

```
SimplePortScanner v1.0
----------------------------
Target: example.com
Scanning 1000 ports...
Started at: Mon, 01 Jan 2023 12:00:00 EST

Port    80: open    HTTP/1.1 200 OK
Port   443: open    
Port  8080: open    HTTP/1.1 302 Found

Scan completed in 2.5s
Results: 1000 ports scanned, 3 open ports
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
