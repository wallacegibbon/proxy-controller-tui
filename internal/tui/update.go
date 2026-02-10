package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.Loading = false
		m.Err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.adjustViewport()
		return m, nil

	case proxiesLoadedMsg:
		m.Loading = false
		m.Proxies = msg.proxies
		m.Groups = msg.groups
		if m.CurrentIdx >= len(m.Groups) {
			m.CurrentIdx = 0
		}
		if len(m.Groups) > 0 && m.CurrentIdx < len(m.Groups) {
			if proxy, ok := m.Proxies[m.Groups[m.CurrentIdx]]; ok {
				cursorFound := false
				for i, p := range proxy.All {
					if p == proxy.Now {
						m.Cursor = i
						cursorFound = true
						break
					}
				}
				if !cursorFound && len(proxy.All) > 0 {
					m.Cursor = 0
				}
			}
		}
		m.ViewportOffset = 0
		m.adjustViewport()
		return m, nil

	case tea.KeyMsg:
		if m.Loading {
			return m, nil
		}

		switch msg.Type {
		case tea.KeyUp, tea.KeyCtrlK:
			if m.CurrentIdx < len(m.Groups) {
				group := m.Groups[m.CurrentIdx]
				if proxy, ok := m.Proxies[group]; ok && len(proxy.All) > 0 {
					if m.Cursor > 0 {
						m.Cursor--
						m.adjustViewport()
					}
				}
			}
			return m, nil

		case tea.KeyDown, tea.KeyCtrlJ:
			if m.CurrentIdx < len(m.Groups) {
				group := m.Groups[m.CurrentIdx]
				if proxy, ok := m.Proxies[group]; ok && len(proxy.All) > 0 {
					if m.Cursor < len(proxy.All)-1 {
						m.Cursor++
						m.adjustViewport()
					}
				}
			}
			return m, nil

		case tea.KeyLeft:
			return m.navigateGroup(-1)

		case tea.KeyRight:
			return m.navigateGroup(1)

		case tea.KeyEnter:
			if m.CurrentIdx < len(m.Groups) {
				group := m.Groups[m.CurrentIdx]
				if proxy, ok := m.Proxies[group]; ok && m.Cursor < len(proxy.All) {
					selectedProxy := proxy.All[m.Cursor]
					if err := m.Client.SelectProxy(group, selectedProxy); err != nil {
						m.Err = err
						return m, nil
					}
					return m, LoadProxiesCmd(m.Client)
				}
			}
			return m, nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}

		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "r":
			m.Loading = true
			return m, LoadProxiesCmd(m.Client)
		case "h":
			return m.navigateGroup(-1)
		case "l":
			return m.navigateGroup(1)
		case "k":
			if m.CurrentIdx < len(m.Groups) {
				group := m.Groups[m.CurrentIdx]
				if proxy, ok := m.Proxies[group]; ok && len(proxy.All) > 0 {
					if m.Cursor > 0 {
						m.Cursor--
						m.adjustViewport()
					}
				}
			}
			return m, nil
		case "j":
			if m.CurrentIdx < len(m.Groups) {
				group := m.Groups[m.CurrentIdx]
				if proxy, ok := m.Proxies[group]; ok && len(proxy.All) > 0 {
					if m.Cursor < len(proxy.All)-1 {
						m.Cursor++
						m.adjustViewport()
					}
				}
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) navigateGroup(direction int) (tea.Model, tea.Cmd) {
	newIdx := m.CurrentIdx + direction
	if newIdx >= 0 && newIdx < len(m.Groups) {
		m.CurrentIdx = newIdx
		group := m.Groups[m.CurrentIdx]
		if proxy, ok := m.Proxies[group]; ok {
			for i, p := range proxy.All {
				if p == proxy.Now {
					m.Cursor = i
					break
				}
			}
		} else {
			m.Cursor = 0
		}
		m.ViewportOffset = 0
		m.adjustViewport()
	}
	return *m, nil
}

func (m *Model) adjustViewport() {
	if len(m.Groups) == 0 {
		return
	}
	group := m.Groups[m.CurrentIdx]
	proxy, ok := m.Proxies[group]
	if !ok {
		return
	}

	// Calculate max visible proxies based on terminal height
	// Footer takes: separator (1) + scrollbar (1 if needed) + help (1)
	footerRows := minHelpRows
	if len(proxy.All) > maxVisibleProxies {
		footerRows++ // Add scrollbar row
	}
	availableRows := m.Height - len(m.Groups) - footerRows
	if availableRows < 1 {
		availableRows = 1
	}
	visibleCount := min(maxVisibleProxies, availableRows)

	if m.Cursor < m.ViewportOffset {
		m.ViewportOffset = m.Cursor
	} else if m.Cursor >= m.ViewportOffset+visibleCount {
		m.ViewportOffset = m.Cursor - visibleCount + 1
	}

	if m.ViewportOffset < 0 {
		m.ViewportOffset = 0
	}

	maxOffset := len(proxy.All) - visibleCount
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.ViewportOffset > maxOffset {
		m.ViewportOffset = maxOffset
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
