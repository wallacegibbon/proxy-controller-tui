package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.Loading {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  Loading proxies...")
	}

	if m.Err != nil {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  Error") + "\n" +
			fmt.Sprintf("  %v\n", m.Err) +
			helpStyle.Render("  Press [r] retry, [q] quit")
	}

	if len(m.Groups) == 0 {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  No proxy groups found") + "\n" +
			helpStyle.Render("  Press [r] refresh, [q] quit")
	}

	// Calculate max group name display width for uniform padding (including type)
	maxGroupWidth := 0
	for _, group := range m.Groups {
		proxy, ok := m.Proxies[group]
		groupWithType := group
		if ok && proxy.Type != "" {
			groupWithType = group + " (" + proxy.Type + ")"
		}
		groupWidth := lipgloss.Width(groupWithType)
		if groupWidth > maxGroupWidth {
			maxGroupWidth = groupWidth
		}
	}

	var s string

	// First, count the number of proxy lines that will be shown
	proxyLines := 0
	for i, group := range m.Groups {
		if i == m.CurrentIdx {
			if proxy, ok := m.Proxies[group]; ok {
				// Calculate how many proxies we can show
				availableRows := m.Height - len(m.Groups) - minHelpRows
				if availableRows < 1 {
					availableRows = 1
				}
				visibleCount := availableRows
				if len(proxy.All) < visibleCount {
					proxyLines = len(proxy.All)
				} else {
					proxyLines = visibleCount
				}
			}
		}
	}

	// Calculate padding after selected group's proxies to push remaining groups down
	// This ensures bottom group stays near help line
	// Total content = len(m.Groups) groups + proxyLines + 1 help line
	// Available padding = m.Height - len(m.Groups) - proxyLines - 1
	paddingAfterSelected := m.Height - len(m.Groups) - proxyLines - 1
	if paddingAfterSelected < 0 {
		paddingAfterSelected = 0
	}

	for i, group := range m.Groups {
		proxy, ok := m.Proxies[group]
		if !ok {
			continue
		}

		// Pad group name to uniform display width with 3 spaces on each side
		groupWithType := group
		if proxy.Type != "" {
			groupWithType = group + " (" + proxy.Type + ")"
		}
		currentWidth := lipgloss.Width(groupWithType)
		paddedGroup := "   " + groupWithType + strings.Repeat(" ", maxGroupWidth-currentWidth) + "   "

		var groupLabel string
		if i == m.CurrentIdx {
			groupLabel = selectedGroupStyle.Render(paddedGroup)
		} else {
			groupLabel = normalGroupStyle.Render(paddedGroup)
		}
		s += groupLabel + "\n"

		if i == m.CurrentIdx {
			// Calculate how many proxies we can show
			// Footer takes: help (1 row)
			availableRows := m.Height - len(m.Groups) - minHelpRows
			if availableRows < 1 {
				availableRows = 1
			}
			visibleCount := availableRows

			visibleProxies := proxy.All
			if len(proxy.All) > visibleCount {
				startIdx := m.ViewportOffset
				if startIdx < 0 {
					startIdx = 0
				}
				endIdx := startIdx + visibleCount
				if endIdx > len(proxy.All) {
					endIdx = len(proxy.All)
				}
				visibleProxies = proxy.All[startIdx:endIdx]
			}

			for j, p := range visibleProxies {
				actualIdx := j + m.ViewportOffset
				var line string
				if actualIdx == m.Cursor && p == proxy.Now {
					line = cursorStyle.Render(">> ") + activeProxyStyle.Render(p)
				} else if actualIdx == m.Cursor {
					line = cursorStyle.Render(">  ") + p
				} else if p == proxy.Now {
					line = " " + activeProxyMarkStyle.Render(">") + " " + activeProxyStyle.Render(p)
				} else {
					line = "   " + normalStyle.Render(p)
				}
				if actualIdx == m.Cursor && len(proxy.All) > visibleCount {
					line += helpStyle.Render(fmt.Sprintf(" (%d/%d)", m.Cursor+1, len(proxy.All)))
				}
				s += line + "\n"
			}
		}

		// Add padding after selected group's proxies to push remaining groups down
		// This keeps bottom group near help line (for multiple groups when not last group)
		if i == m.CurrentIdx && i < len(m.Groups)-1 && paddingAfterSelected > 0 {
			s += strings.Repeat("\n", paddingAfterSelected)
		}
	}

	// Add padding before help line when:
	// 1. Single group (no more groups after proxies)
	// 2. Multiple groups but selected group is the last one
	needPaddingBeforeHelp := len(m.Groups) == 1 || m.CurrentIdx == len(m.Groups)-1
	if needPaddingBeforeHelp && paddingAfterSelected > 0 {
		s += strings.Repeat("\n", paddingAfterSelected)
	}

	// Add help text at bottom
	s += helpStyle.Render(" [←h]Prev [→l]Next  [↑k]↑ [↓j]↓  [Ent]Select  [r]Reload  [q]Quit")

	return s
}
