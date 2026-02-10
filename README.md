# Proxy TUI Controller

A modern, compact terminal user interface (TUI) for managing Clash/Mihomo proxy services. Built with Go and charmbracelet ecosystem (bubbletea & lipgloss).

## Features

- **Modern TUI Interface**: Clean, compact design with enhanced visual styling
- **Proxy Management**: Select proxies from Selector and URLTest groups
- **Smart Navigation**: Vim-style (h/j/k/l) and arrow key support
- **Viewport Scrolling**: Handles large proxy lists efficiently (20 items visible)
- **API Authentication**: Support for Mihomo secret tokens
- **Mock Mode**: Built-in testing mode without a running proxy server
- **Cursor Alignment**: Proper cursor positioning on active proxies

## Installation

### From Source (Recommended)

```bash
go install github.com/wallacegibbon/proxy-controller-tui@latest
```

### Building from Source

```bash
# Clone the repository
git clone git@github.com:wallacegibbon/proxy-controller-tui.git
cd proxy-controller-tui

# Build binary
go build
```

## Usage

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MIHOMO_SECRET` | Mihomo API secret token | (none) |
| `MOCK_CLASH` | Enable mock mode for testing | `0` |

### Running

```bash
# With Mihomo secret
MIHOMO_SECRET=YOUR_SECRET proxy-controller-tui

# Standard Clash
proxy-controller-tui

# Mock mode for testing (no proxy server required)
MOCK_CLASH=1 proxy-controller-tui
```

## Configuration

The application connects to to Clash/Mihomo RESTful API at:

```
http://127.0.0.1:9090
```

Ensure your proxy server is running and accessible at this endpoint.

## Controls

| Key | Action |
|-----|--------|
| `←` / `h` | Previous proxy group |
| `→` / `l` | Next proxy group |
| `↑` / `k` | Previous proxy in group |
| `↓` / `j` | Next proxy in group |
| `Enter` | Select current proxy |
| `r` | Reload proxy list |
| `q` / `Ctrl+C` | Quit |

## Requirements

- Go 1.25.6 or later
- Clash or Mihomo proxy server running with RESTful API enabled

## Development

```bash
# Run tests
go test ./...

# Run with debug output
go run .

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o proxy-controller-tui-linux
GOOS=darwin GOARCH=amd64 go build -o proxy-controller-tui-mac
GOOS=windows GOARCH=amd64 go build -o proxy-controller-tui.exe
```

## Project Structure

```
proxy-controller-tui/
├── main.go                      # Application entry point
├── internal/
│   ├── clash/                   # Clash/Mihomo API client
│   └── tui/                     # TUI implementation
│       ├── model.go              # Model and initialization
│       ├── update.go             # Update logic
│       ├── view.go               # View rendering
│       └── model_test.go         # Tests
├── go.mod
├── go.sum
├── README.md
├── AGENTS.md
└── STATE.md
```

## Tech Stack

- **[bubbletea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling
- **[charmbracelet](https://github.com/charmbracelet)** ecosystem

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Links

- [GitHub Repository](https://github.com/wallacegibbon/proxy-controller-tui)
- [Gitee Repository](https://gitee.com/wallacegibbon/proxy-controller-tui)
- [Clash Documentation](https://github.com/Dreamacro/clash)
- [Mihomo Documentation](https://github.com/MetaCubeX/mihomo)
