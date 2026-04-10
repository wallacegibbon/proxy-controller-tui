package tui

import (
	tea "charm.land/bubbletea/v2"
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
				// Try to restore cursor position based on the proxy name we were on
				cursorFound := false
				if m.lastCursorProxy != "" {
					for i, p := range proxy.All {
						if p == m.lastCursorProxy {
							m.Cursor = i
							cursorFound = true
							break
						}
					}
				}
				// Fall back to active proxy if:
				// 1. We couldn't find the last cursor proxy by name
				// 2. OR this is the first load (lastCursorProxy is empty)
				if !cursorFound {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.Cursor = i
							break
						}
					}
				}
				// Update lastCursorProxy to match the current cursor position
				m.updateLastCursorProxy()
			}
		}
		m.adjustViewport()
		return m, nil

	case tea.KeyPressMsg:
		if m.Loading {
			return m, nil
		}

		switch key := msg.Key(); {
		case key.Code == tea.KeyUp || (key.Text == "k" && key.Mod == 0):
			if m.CurrentIdx < len(m.Groups) {
				group := m.Groups[m.CurrentIdx]
				if proxy, ok := m.Proxies[group]; ok && len(proxy.All) > 0 {
					if m.Cursor > 0 {
						m.Cursor--
						m.updateLastCursorProxy()
						m.adjustViewport()
					}
				}
			}
			return m, nil

		case key.Code == tea.KeyDown || (key.Text == "j" && key.Mod == 0):
			if m.CurrentIdx < len(m.Groups) {
				group := m.Groups[m.CurrentIdx]
				if proxy, ok := m.Proxies[group]; ok && len(proxy.All) > 0 {
					if m.Cursor < len(proxy.All)-1 {
						m.Cursor++
						m.updateLastCursorProxy()
						m.adjustViewport()
					}
				}
			}
			return m, nil

		case key.Code == tea.KeyLeft || (key.Text == "h" && key.Mod == 0):
			return m.navigateGroup(-1)

		case key.Code == tea.KeyRight || (key.Text == "l" && key.Mod == 0):
			return m.navigateGroup(1)

		case key.Code == tea.KeyEnter:
			if m.CurrentIdx < len(m.Groups) {
				group := m.Groups[m.CurrentIdx]
				if proxy, ok := m.Proxies[group]; ok && m.Cursor < len(proxy.All) {
					selectedProxy := proxy.All[m.Cursor]
					if err := m.Client.SelectProxy(group, selectedProxy); err != nil {
						m.Err = err
						return m, nil
					}
					return m, loadProxiesWithDelayCmd(m.Client)
				}
			}
			return m, nil

		case key.Text == "q" && key.Mod == 0:
			return m, tea.Quit

		case key.Text == "r" && key.Mod == 0:
			m.Loading = true
			return m, LoadProxiesCmd(m.Client)
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
					m.updateLastCursorProxy()
					break
				}
			}
		} else {
			m.Cursor = 0
			m.lastCursorProxy = ""
		}
		m.ViewportOffset = 0
		m.adjustViewport()
	}
	return *m, nil
}

func (m *Model) updateLastCursorProxy() {
	if m.CurrentIdx < len(m.Groups) {
		group := m.Groups[m.CurrentIdx]
		if proxy, ok := m.Proxies[group]; ok && m.Cursor < len(proxy.All) {
			m.lastCursorProxy = proxy.All[m.Cursor]
		}
	}
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
	// Footer takes: help (1 row)
	// We can show at most (Height - 1) lines total (leaving 1 for help)
	// Minus 1 for the selected group header
	availableHeight := m.Height - minHelpRows
	maxProxyLines := availableHeight - 1
	if maxProxyLines < 1 {
		maxProxyLines = 1
	}

	visibleCount := maxProxyLines
	if visibleCount > len(proxy.All) {
		visibleCount = len(proxy.All)
	}

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
