# Project State - Last updated: 2026-02-11

## Status
**Complete and production-ready.**

All features implemented and tested. Application successfully deployed to GitHub and Gitee.

## Project Structure
- Module: `github.com/wallacegibbon/proxy-controller-tui`
- Binary: `proxy-controller-tui`
- Project layout:
  - `main.go` - Application entry point
  - `internal/clash/` - Clash/Mihomo API client
  - `internal/tui/` - TUI implementation (model, update, view)

## Installation
```bash
go install github.com/wallacegibbon/proxy-controller-tui@latest
```

## Features
- Uses alternate screen buffer for proper terminal cleanup on exit
- Compact, modern TUI interface for Clash/Mihomo proxy management
- Small terminal support with dynamic viewport and compact help text
- Beautiful group styling with dark blue background color
- Proper multi-byte character support (Chinese/English mixed names)
- Active proxy marked with `>` in orange (cursor: `>` in cyan)
- Inline position indicator (x/xx) following cursor
- Uniform group name padding using lipgloss.Width() for proper display
- Viewport automatically adjusts based on terminal height
- Mihomo API authentication via `MIHOMO_SECRET`
- Vim-style navigation (h/j/k/l) and arrow keys
- Mock mode for testing (`MOCK_CLASH=1`)

## Tech Stack
- bubbletea - TUI framework
- lipgloss - Styling
- charmbracelet ecosystem

## Controls
- `←/h` / `→/l`: Group navigation
- `↑/k` / `↓/j`: Proxy navigation
- `Enter`: Select proxy
- `r`: Refresh, `q`: Quit

## UI Design
- **Groups**: Dark blue background (color 60), selected group in white, 3-space padding
- **Proxies**: 
  - Active proxy: `>` marker in orange (color 208)
  - Cursor: `>` marker in cyan (color 51), or `>>` when on active proxy
  - Normal: No marker
- **Position indicator**: `(x/xx)` shown next to cursor when scrolling needed
- **Help**: Compact format on terminals < 15 rows, full format otherwise
