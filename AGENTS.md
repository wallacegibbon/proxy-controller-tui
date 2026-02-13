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
- Uses alternate screen buffer for proper display cleanup on exit
- Small terminal support with dynamic viewport calculation
- Beautiful turquoise background (color 45) for all groups, selected group in white
- Active proxy marked with `>` in orange (color 208), cursor marked with `>` in cyan (color 51), or `>>` when cursor is on active proxy
- Inline position indicator `(x/xx)` shows current cursor position
- Proper multi-byte character support for Chinese/English names
- Group type displayed in parentheses after group name (e.g., "MyGroup (Selector)")
- Help text fixed at bottom: `[←h]Prev [→l]Next  [↑k]↑ [↓j]↓  [Ent]Select  [r]Reload  [q]Quit`
- Top group always on top line, no gaps between unselected groups
- Bottom group directly above help line with proper padding distribution

## Agent Instructions
- **Read STATE.md** at the start of every conversation
- **Update STATE.md** after completing any meaningful work (features, bug fixes, etc.)
- Keep STATE.md as the single source of truth for project status
