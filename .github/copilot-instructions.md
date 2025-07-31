# Gluetun Development Guide

This file contains guidance for AI coding agents working on the Gluetun VPN client codebase.

## Architecture Overview

Gluetun is a Go-based containerized VPN client that supports multiple VPN providers. The architecture follows a modular design:

- **Entry point**: `cmd/gluetun/main.go` - Complex initialization, graceful shutdown orchestration, CLI commands
- **Core loops**: Each major component runs as a goroutine managed by `goshutdown` framework
- **Provider abstraction**: `internal/provider/` contains provider-specific implementations with unified interfaces  
- **VPN protocols**: OpenVPN (`internal/openvpn/`) and WireGuard (`internal/wireguard/`) support
- **Network management**: Firewall rules (`internal/firewall/`), routing (`internal/routing/`), netlink operations

## Key Components & Data Flow

1. **Configuration pipeline**: `internal/configuration/` reads from env vars, files, secrets
2. **Provider resolution**: Based on `VPN_SERVICE_PROVIDER`, creates provider-specific connection settings
3. **VPN loop**: `internal/vpn/` orchestrates connection setup/teardown via provider interfaces
4. **Network setup**: Firewall rules → TUN device → routing configuration → VPN tunnel
5. **Supporting services**: DNS-over-TLS, HTTP proxy, port forwarding, health monitoring

## Development Workflows

### Building & Testing
```bash
# Build Docker image (multi-stage, cached layers)
docker build -t gluetun-dev .

# Run tests with race detection
go test -race ./...

# Specific test patterns
go test ./internal/configuration/settings -run TestWireguard -v

# Lint (uses .golangci.yml)
golangci-lint run --timeout=10m
```

### Running Locally
```bash
# Basic VPN connection
docker run --cap-add=NET_ADMIN --device /dev/net/tun \
  -e VPN_SERVICE_PROVIDER=nordvpn \
  -e VPN_TYPE=wireguard \
  -e WIREGUARD_PRIVATE_KEY=... \
  gluetun-dev

# Custom WireGuard config file
docker run --cap-add=NET_ADMIN --device /dev/net/tun \
  -v /path/to/config.conf:/gluetun/wg.conf:ro \
  -e WIREGUARD_CUSTOM_CONFIG_FILE=/gluetun/wg.conf \
  -e VPN_TYPE=wireguard \
  -e VPN_SERVICE_PROVIDER=custom \
  gluetun-dev
```

## Code Patterns & Conventions

### Loop Pattern
Most services follow this pattern:
```go
type Loop struct {
    settings Settings
    // ... other fields
}

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
    defer close(done)
    // Main service logic with context cancellation
}
```

### Interface-Based Architecture
- Heavy use of interfaces for testability (see `internal/*/interfaces.go`)
- Provider interfaces in `internal/provider/` enable multi-provider support
- Mock generation via `go:generate` comments

### Error Handling
- Consistent error wrapping with `fmt.Errorf("%w: additional context", err)`
- Custom error types for specific conditions (e.g., `internal/mod/load.go`)

### Configuration
- Environment variables prefixed by component (e.g., `WIREGUARD_`, `OPENVPN_`)
- Settings validation in separate `.Validate()` methods
- Default values set via `.SetDefaults()`

## Integration Points

### VPN Provider Integration
Add new providers in `internal/provider/newprovider/`:
- Implement `Provider` interface from `internal/provider/`
- Add updater for server list management
- Register in `internal/provider/providers.go`

### Custom Configuration Support
Recent addition: WireGuard custom config file support
- Parser: `internal/wireguard/confparser.go`
- Integration: `internal/provider/utils/wireguard.go`
- Tests use `internal/wireguard/wg_test.conf` with sanitized values

### Container Networking
- Requires `NET_ADMIN` capability and `/dev/net/tun` device
- Firewall management via iptables (see `internal/firewall/`)
- Network interface creation and configuration via netlink

### Dependencies
- `golang.zx2c4.com/wireguard` for WireGuard protocol
- `github.com/vishvananda/netlink` for network interface management
- `github.com/qdm12/goshutdown` for graceful shutdown orchestration

## Common Debugging

- Health check endpoint: `localhost:9999/healthcheck`
- Control server: `localhost:8000` (various management endpoints)
- Logs include component prefixes for filtering
- IPv6 support detection affects configuration validation
- Module loading failures often indicate missing kernel capabilities
