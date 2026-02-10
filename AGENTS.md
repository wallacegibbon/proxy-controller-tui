# Proxy TUI Controller

Go TUI application for managing Clash/Mihomo proxy services.

## Project
- Module: `github.com/wallacegibbon/proxy-controller-tui`
- Binary: `proxy-controller-tui`
- Connects to Clash/Mihomo RESTful API (`http://127.0.0.1:9090`)
- Supports proxy group selection (Selector and URLTest types)
- Built with bubbletea and lipgloss (charmbracelet)

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

## Controls
- `←/h` / `→/l`: Previous/Next group
- `↑/k` / `↓/j`: Previous/Next proxy
- `Enter`: Select proxy
- `r`: Refresh, `q`: Quit

## UI Features
- Small terminal support with dynamic viewport calculation
- Beautiful dark blue background for all groups
- Active proxy marked with `>` (orange), cursor marked with `>` (cyan)
- Inline position indicator `(x/xx)` shows current cursor position
- Proper multi-byte character support for Chinese/English names
- Compact help text on terminals < 15 rows

## Agent Instructions
- **Read STATE.md** at the start of every conversation
- **Update STATE.md** after completing any meaningful work (features, bug fixes, etc.)
- Keep STATE.md as the single source of truth for project status
