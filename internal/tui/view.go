package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/wallacegibbon/proxy-controller-tui/internal/clash"
)

var (
	fixedIndicatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	helpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

func (m Model) View() tea.View {
	if m.Loading {
		v := tea.NewView(
			separatorStyle.Render("═══════════════════════════════════════") + "\n" +
				headerStyle.Render("  Loading proxies..."),
		)
		v.AltScreen = true
		return v
	}

	if m.Err != nil {
		v := tea.NewView(
			separatorStyle.Render("═══════════════════════════════════════") + "\n" +
				headerStyle.Render("  Error") + "\n" +
				fmt.Sprintf("  %v\n", m.Err) +
				helpStyle.Render("  Press [r] retry, [q] quit"),
		)
		v.AltScreen = true
		return v
	}

	if len(m.Groups) == 0 {
		v := tea.NewView(
			separatorStyle.Render("═══════════════════════════════════════") + "\n" +
				headerStyle.Render("  No proxy groups found") + "\n" +
				helpStyle.Render("  Press [r] refresh, [q] quit"),
		)
		v.AltScreen = true
		return v
	}

	// Get selected group's proxy info
	var selectedProxy clash.Proxy
	var selectedOk bool
	if m.CurrentIdx < len(m.Groups) {
		selectedProxy, selectedOk = m.Proxies[m.Groups[m.CurrentIdx]]
	}

	var s string

	// Show only the selected group with navigation indicators
	if selectedOk {
		group := m.Groups[m.CurrentIdx]
		groupWithType := group
		if selectedProxy.Type != "" {
			groupWithType = group + " (" + selectedProxy.Type + ")"
			if selectedProxy.Type == "URLTest" && selectedProxy.Fixed != "" {
				groupWithType += " " + fixedIndicatorStyle.Render("[fixed]")
			}
		}

		// Navigation indicators
		hasLeft := m.CurrentIdx > 0
		hasRight := m.CurrentIdx < len(m.Groups)-1
		prefix := "   "
		suffix := ""
		if hasLeft {
			prefix = "<< "
		}
		if hasRight {
			suffix = " >>"
		}

		s += selectedGroupStyle.Render(prefix+groupWithType+suffix) + "\n"

		// Render proxies
		if len(selectedProxy.All) > 0 {
			maxProxyLines := m.Height - 1
			if maxProxyLines < 1 {
				maxProxyLines = 1
			}

			totalProxies := len(selectedProxy.All)
			visibleCount := maxProxyLines
			if visibleCount > totalProxies {
				visibleCount = totalProxies
			}

			startIdx := m.ViewportOffset
			if startIdx < 0 {
				startIdx = 0
			}
			if startIdx > totalProxies-visibleCount {
				startIdx = totalProxies - visibleCount
			}
			endIdx := startIdx + visibleCount
			if endIdx > totalProxies {
				endIdx = totalProxies
			}

			for j, p := range selectedProxy.All[startIdx:endIdx] {
				actualIdx := j + startIdx
				var line string
				if actualIdx == m.Cursor && p == selectedProxy.Now {
					line = cursorStyle.Render(">> ") + activeProxyStyle.Render(p)
				} else if actualIdx == m.Cursor {
					line = cursorStyle.Render(">  ") + p
				} else if p == selectedProxy.Now {
					line = " " + activeProxyMarkStyle.Render(">") + " " + activeProxyStyle.Render(p)
				} else {
					line = "   " + normalStyle.Render(p)
				}
				if actualIdx == m.Cursor && totalProxies > visibleCount {
					line += helpStyle.Render(fmt.Sprintf(" (%d/%d)", m.Cursor+1, totalProxies))
				}
				s += line + "\n"
			}
		}
	}

	v := tea.NewView(s)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
