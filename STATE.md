# Project State - Last updated: 2026-02-13

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
- Small terminal support with dynamic viewport
- Beautiful group styling with turquoise background color
- Proper multi-byte character support (Chinese/English mixed names)
- Active proxy marked with `>` in orange (cursor: `>` in cyan)
- Inline position indicator (x/xx) following cursor
- Uniform group name padding using lipgloss.Width() for proper display
- Viewport automatically adjusts based on terminal height
- Mihomo API authentication via `MIHOMO_SECRET`
- Vim-style navigation (h/j/k/l) and arrow keys
- Mock mode for testing (`MOCK_CLASH=1`) with proper state persistence
- **Consistent group ordering**: Groups are sorted alphabetically regardless of API response order
- **Smart cursor positioning**:
  - On startup and group switches, cursor goes to the currently active proxy
  - After manual navigation, cursor stays on the proxy you navigated to
  - Preserved across reloads by tracking the proxy name at cursor
- Only resets cursor to active proxy if the proxy you were on no longer exists
- 200ms delay after PUT request to allow server to process selection before refreshing data

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
- **Layout**: Top group always on top line, no gaps between unselected groups, bottom group directly above help line (no padding)
- **Groups**: Turquoise background (color 45), selected group in white, 3-space padding
  - Groups displayed in original order, navigating up/down moves through all groups
  - Selected group shows its proxies below it
- **Proxies**:
  - Active proxy: `>` marker in orange (color 208)
  - Cursor: `>` marker in cyan (color 51), or `>>` when on active proxy
  - Normal: No marker
- **Position indicator**: `(x/xx)` shown next to cursor when scrolling needed
- **Help**: Fixed at bottom of terminal with format `[←h]Prev [→l]Next  [↑k]↑ [↓j]↓  [Ent]Select  [r]Reload  [q]Quit`
- **Padding**: 
  - No gaps between unselected groups
  - Padding added after selected group's proxies (if not last group) to push remaining groups down
  - Padding added before help line (when selected group is last or single group) to fill terminal height
