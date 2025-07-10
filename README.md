# Pulse - Docker Healthcheck Tool

A lightweight CLI tool for Docker container healthchecks with minimal dependencies. Pulse provides TCP port checking and HTTP endpoint health monitoring with configurable timeouts and verbosity levels.

## Features

- **TCP Port Check**: Simple TCP connectivity testing (like netcat)
- **UDP Port Check**: UDP connectivity testing for connectionless protocols
- **HTTP Health Check**: HTTP request testing with configurable method and path
- **Silent/Verbose Modes**: Control output verbosity for debugging
- **Multi-architecture Support**: Can be built for multiple platforms
- **Minimal Dependencies**: Uses only Go standard library
- **Docker-friendly**: Designed for container healthchecks

## Installation

### Build from source

```bash
# Build for current platform
go build -o pulse .

# Build for multiple architectures
make build-all
```

### Pre-built binaries

Download the latest release for your platform from the releases page.

## Usage

### TCP Mode (Default)

Check if a TCP port is open:

```bash
# Basic TCP check
./pulse -host localhost -port 8080

# With custom timeout
./pulse -host myservice -port 5432 -timeout 10s

# Silent mode (no output, only exit code)
./pulse -host localhost -port 80 -silent

# Verbose mode with debug information
./pulse -host localhost -port 3306 -verbose
```

### UDP Mode

Check if a UDP port is reachable:

```bash
# Basic UDP check
./pulse -mode udp -host localhost -port 53

# With custom timeout
./pulse -mode udp -host dns.example.com -port 53 -timeout 10s

# Silent mode
./pulse -mode udp -host localhost -port 161 -silent

# Verbose mode with debug information
./pulse -mode udp -host localhost -port 514 -verbose
```

### HTTP Mode

Check HTTP endpoints:

```bash
# Basic HTTP GET check
./pulse -mode http -host localhost -port 8080 -path /health

# Custom HTTP method
./pulse -mode http -host api.example.com -port 443 -path /status -method HEAD

# With timeout
./pulse -mode http -host myservice -port 3000 -path /ready -timeout 30s
```

### Docker Healthcheck Examples

```dockerfile
# TCP healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["pulse", "-host", "localhost", "-port", "8080", "-silent"]

# UDP healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["pulse", "-mode", "udp", "-host", "localhost", "-port", "53", "-silent"]

# HTTP healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["pulse", "-mode", "http", "-host", "localhost", "-port", "8080", "-path", "/health", "-silent"]
```

## Command Line Options

| Flag       | Default     | Description                   |
| ---------- | ----------- | ----------------------------- |
| `-mode`    | `tcp`       | Mode: `tcp`, `udp`, or `http` |
| `-host`    | `localhost` | Host to check                 |
| `-port`    | `80`        | Port to check                 |
| `-timeout` | `5s`        | Timeout for the check         |
| `-silent`  | `false`     | Silent mode (no output)       |
| `-verbose` | `false`     | Verbose mode (debug output)   |
| `-path`    | `/`         | HTTP path (for http mode)     |
| `-method`  | `GET`       | HTTP method (for http mode)   |

## Exit Codes

- `0`: Health check passed
- `1`: Health check failed or error occurred

## Building for Multiple Architectures

Use the provided Makefile to build for multiple platforms:

```bash
# Build for all supported architectures
make build-all

# Build for specific architecture
make build-linux-amd64
make build-linux-arm64
make build-darwin-amd64
make build-darwin-arm64
```

## Supported Architectures

- Linux: amd64, arm64, arm/v7, arm/v6
- macOS: amd64, arm64
- Windows: amd64, arm64

## License

MIT License - see LICENSE file for details.
