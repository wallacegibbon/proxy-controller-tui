# Proxy TUI Controller

A terminal user interface (TUI) for managing Clash/Mihomo proxy services. Built with Go and charmbracelet ecosystem (bubbletea & lipgloss).

## Features

- Proxy Management: Select proxies from Selector and URLTest groups
- URLTest Fixed Indicator: Shows `[fixed]` when a URLTest group has been manually pinned
- Auto-Selection Reset: Press `a` to restore auto-selection for pinned URLTest groups
- Vim-style (h/j/k/l) and arrow key navigation
- API Authentication: Support for Mihomo secret tokens
- Mock Mode: Built-in testing mode without a running proxy server

## Installation

```bash
go install github.com/wallacegibbon/proxy-controller-tui@latest
```

## Usage

```bash
# With Mihomo secret
MIHOMO_SECRET=YOUR_SECRET proxy-controller-tui

# Standard Clash
proxy-controller-tui

# Mock mode for testing
MOCK_CLASH=1 proxy-controller-tui
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MIHOMO_SECRET` | Mihomo API secret token | (none) |
| `MOCK_CLASH` | Enable mock mode for testing | `0` |

The application connects to Clash/Mihomo RESTful API at `http://127.0.0.1:9090`.

## Controls

| Key | Action |
|-----|--------|
| `←` / `h` | Previous proxy group |
| `→` / `l` | Next proxy group |
| `↑` / `k` | Previous proxy in group |
| `↓` / `j` | Next proxy in group |
| `Enter` | Select current proxy |
| `a` | Reset to auto-selection (URLTest groups with `[fixed]`) |
| `r` | Reload proxy list |
| `q` / `Ctrl+C` | Quit |

## Requirements

- Go 1.25.6 or later
- Clash or Mihomo proxy server running with RESTful API enabled

## Links

- [GitHub Repository](https://github.com/wallacegibbon/proxy-controller-tui)
- [Clash Documentation](https://github.com/Dreamacro/clash)
- [Mihomo Documentation](https://github.com/MetaCubeX/mihomo)