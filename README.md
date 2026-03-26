<p align="center">
  <h1 align="center">🔍 HostPeek</h1>
  <p align="center">
    <strong>Tiny, cross-platform, self-hosted node agent that exposes host metrics over HTTP as JSON.</strong>
  </p>
  <p align="center">
    <a href="https://github.com/Appouse/HostPeek/releases"><img src="https://img.shields.io/github/v/release/Appouse/HostPeek?style=flat-square" alt="Release"></a>
    <a href="https://github.com/Appouse/HostPeek/blob/main/LICENSE"><img src="https://img.shields.io/github/license/Appouse/HostPeek?style=flat-square" alt="License"></a>
    <a href="https://github.com/Appouse/HostPeek/actions"><img src="https://img.shields.io/github/actions/workflow/status/Appouse/HostPeek/release.yml?style=flat-square" alt="Build"></a>
  </p>
</p>

---

## Features

- **Single binary** — no runtime dependencies, no frameworks
- **Cross-platform** — works on Linux and Windows
- **JSON over HTTP** — easy to integrate with monitoring tools (Uptime Kuma, Grafana, etc.)
- **Configurable** — YAML config for port, auth, and collector toggles
- **Packaged** — `.deb`, `.rpm`, and `.exe` available in releases
- **Lightweight** — ~10 MB binary, minimal CPU/memory footprint

## Quick Start

### Linux (Debian/Ubuntu)

```bash
# Download the latest .deb from releases
wget https://github.com/Appouse/HostPeek/releases/latest/download/hostpeek_<version>_linux_amd64.deb
sudo dpkg -i hostpeek_*.deb

# The service starts automatically on port 8080
curl http://localhost:8080/metrics
```

### Linux (RHEL/Fedora/CentOS)

```bash
# Download the latest .rpm from releases
wget https://github.com/Appouse/HostPeek/releases/latest/download/hostpeek_<version>_linux_amd64.rpm
sudo rpm -i hostpeek_*.rpm

# The service starts automatically on port 8080
curl http://localhost:8080/metrics
```

### Windows

1. Download `hostpeek_<version>_windows_amd64.zip` from [Releases](https://github.com/Appouse/HostPeek/releases)
2. Extract and run:

```powershell
.\hostpeek.exe
# Or with a custom config:
.\hostpeek.exe -config .\hostpeek.yaml
```

### From Source

```bash
git clone https://github.com/Appouse/HostPeek.git
cd HostPeek
go build -o hostpeek ./cmd/hostpeek
./hostpeek
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check & version |
| `GET` | `/metrics` | All metrics combined |
| `GET` | `/metrics/os` | OS information |
| `GET` | `/metrics/cpu` | CPU metrics |
| `GET` | `/metrics/memory` | Memory metrics |
| `GET` | `/metrics/disk` | Disk partitions |
| `GET` | `/metrics/network` | Network interfaces |
| `GET` | `/metrics/uptime` | Uptime & boot time |

## Example Responses

### `GET /health`

```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

### `GET /metrics`

```json
{
  "hostname": "srv-app-01",
  "os": {
    "hostname": "srv-app-01",
    "platform": "linux",
    "distribution": "Ubuntu 24.04",
    "kernel": "6.8.0"
  },
  "cpu": {
    "cores": 8,
    "model_name": "Intel(R) Core(TM) i7-10700K",
    "usage_percent": 13.4,
    "load_avg": [0.72, 0.65, 0.61]
  },
  "memory": {
    "total_mb": 16384,
    "used_mb": 6210,
    "free_mb": 10174,
    "usage_percent": 37.9
  },
  "disk": [
    {
      "mount": "/",
      "fs": "ext4",
      "total_gb": 256.0,
      "used_gb": 108.4,
      "free_gb": 147.6,
      "usage_percent": 42.3
    }
  ],
  "network": {
    "interfaces": [
      {
        "name": "eth0",
        "ipv4": ["10.0.0.15"],
        "mac": "00:11:22:33:44:55",
        "rx_bytes": 123456789,
        "tx_bytes": 987654321
      }
    ]
  },
  "uptime": {
    "uptime_seconds": 86400,
    "boot_time": "2026-03-25T08:15:00Z"
  },
  "agent": {
    "version": "1.0.0",
    "status": "ok"
  }
}
```

### `GET /metrics/cpu`

```json
{
  "cores": 8,
  "model_name": "Intel(R) Core(TM) i7-10700K",
  "usage_percent": 13.4,
  "load_avg": [0.72, 0.65, 0.61]
}
```

### `GET /metrics/memory`

```json
{
  "total_mb": 16384,
  "used_mb": 6210,
  "free_mb": 10174,
  "usage_percent": 37.9
}
```

## Configuration

HostPeek looks for `hostpeek.yaml` in the current directory by default. Use `-config` flag to specify a custom path.

```yaml
# Server settings
server:
  listen: ":8080"         # Bind address
  read_timeout: 5s
  write_timeout: 10s

# Authentication (optional)
auth:
  enabled: false
  api_key: "your-secret-key"  # Clients must send X-API-Key header

# Toggle individual collectors
collectors:
  cpu: true
  memory: true
  disk: true
  network: true
  os: true
  uptime: true
```

### Authentication

When `auth.enabled` is `true`, all `/metrics/*` endpoints require the `X-API-Key` header:

```bash
curl -H "X-API-Key: your-secret-key" http://localhost:8080/metrics
```

The `/health` endpoint is always accessible without authentication.

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `hostpeek.yaml` | Path to YAML config file |
| `-version` | — | Print version and exit |

## Integration Examples

### Uptime Kuma

1. Add a new monitor with type **HTTP(s) - Json Query**
2. URL: `http://your-host:8080/health`
3. JSON Query: `$.status == "ok"`

For resource monitoring:
- URL: `http://your-host:8080/metrics`
- JSON Query: `$.cpu.usage_percent < 90 && $.memory.usage_percent < 85`

### curl / Scripts

```bash
# Quick health check
curl -s http://localhost:8080/health | jq .

# Get only CPU metrics
curl -s http://localhost:8080/metrics/cpu | jq .

# With API key
curl -s -H "X-API-Key: secret" http://localhost:8080/metrics | jq .
```

## Service Management (Linux)

```bash
# After installing via .deb or .rpm:
sudo systemctl status hostpeek
sudo systemctl restart hostpeek
sudo systemctl stop hostpeek

# View logs
sudo journalctl -u hostpeek -f

# Edit config
sudo nano /etc/hostpeek/hostpeek.yaml
sudo systemctl restart hostpeek
```

## Building

### Requirements

- Go 1.22+
- [GoReleaser](https://goreleaser.com) (for packaging)

### Local Build

```bash
# Build for current platform
go build -o hostpeek ./cmd/hostpeek

# Build with version info
go build -ldflags "-X main.version=1.0.0" -o hostpeek ./cmd/hostpeek

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o hostpeek-linux ./cmd/hostpeek
```

### Create Release Packages

```bash
# Snapshot (local, no publish)
goreleaser release --snapshot --clean
# Packages will be in ./dist/
```

## Release Process

Releases are fully automated via GitHub Actions:

```bash
git tag v1.0.0
git push origin v1.0.0
```

This triggers the CI/CD pipeline that:
1. Builds binaries for Linux (amd64/arm64) and Windows (amd64)
2. Creates `.deb` and `.rpm` packages
3. Publishes everything as a GitHub Release

## License

[MIT](LICENSE)
