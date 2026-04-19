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
- `a`: Reset to auto-selection (for URLTest groups with `[fixed]`)
- `r`: Refresh, `q`: Quit

## UI Features
- Uses alternate screen buffer for proper display cleanup on exit
- Single group view: shows only the selected group with its proxies
- Navigation indicators `<<`/`>>` show if there are groups to the left/right
- Beautiful turquoise background (color 45), selected group in white
- Active proxy marked with `>` in orange (color 208), cursor in cyan (color 51)
- Position indicator `(x/xx)` when proxy list exceeds screen height
- Group type displayed after name (e.g., "MyGroup (Selector)")
- `[fixed]` indicator in red when URLTest group is pinned

## Agent Instructions
- **Read STATE.md** at the start of every conversation
- **Update STATE.md** after completing any meaningful work (features, bug fixes, etc.)
- Keep STATE.md as the single source of truth for project status
