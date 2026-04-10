package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wallacegibbon/proxy-controller-tui/internal/clash"
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

	availableHeight := m.Height - minHelpRows

	// Get selected group's proxy info
	var selectedProxy clash.Proxy
	var selectedOk bool
	if m.CurrentIdx < len(m.Groups) {
		selectedProxy, selectedOk = m.Proxies[m.Groups[m.CurrentIdx]]
	}

	// Calculate max proxy lines we can show (leave at least 1 line for the selected group)
	maxProxyLines := availableHeight - 1
	if maxProxyLines < 1 {
		maxProxyLines = 1
	}

	// Calculate actual proxy lines to show for selected group
	proxyLines := 0
	totalProxies := 0
	if selectedOk {
		totalProxies = len(selectedProxy.All)
		if totalProxies <= maxProxyLines {
			proxyLines = totalProxies
		} else {
			proxyLines = maxProxyLines
		}
	}

	// Calculate how many groups we can show
	groupsToShow := availableHeight - proxyLines
	if groupsToShow > len(m.Groups) {
		groupsToShow = len(m.Groups)
	}
	if groupsToShow < 1 {
		groupsToShow = 1
	}

	// Calculate start index for groups - try to center CurrentIdx in the visible range
	startGroupIdx := m.CurrentIdx - groupsToShow/2
	if startGroupIdx < 0 {
		startGroupIdx = 0
	}
	endGroupIdx := startGroupIdx + groupsToShow
	if endGroupIdx > len(m.Groups) {
		endGroupIdx = len(m.Groups)
		startGroupIdx = endGroupIdx - groupsToShow
		if startGroupIdx < 0 {
			startGroupIdx = 0
		}
	}

	var s string
	linesUsed := 0

	for i := startGroupIdx; i < endGroupIdx && linesUsed < availableHeight; i++ {
		group := m.Groups[i]
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

		var paddedGroup string
		if i == m.CurrentIdx {
			// Show navigation indicators for selected group
			hasLeft := m.CurrentIdx > 0
			hasRight := m.CurrentIdx < len(m.Groups)-1

			prefix := "   " // 3 spaces when no left indicator
			suffix := ""    // empty when no right indicator
			if hasLeft {
				prefix = "<< "
			}
			if hasRight {
				suffix = " >>"
			}
			paddedGroup = prefix + groupWithType + suffix
		} else {
			// Other groups: align with 3 spaces each side + padding to maxGroupWidth
			paddedGroup = "   " + groupWithType + strings.Repeat(" ", maxGroupWidth-currentWidth) + "   "
		}

		var groupLabel string
		if i == m.CurrentIdx {
			groupLabel = selectedGroupStyle.Render(paddedGroup)
		} else {
			groupLabel = normalGroupStyle.Render(paddedGroup)
		}
		s += groupLabel + "\n"
		linesUsed++

		// Render proxies for selected group
		if i == m.CurrentIdx && selectedOk {
			// Calculate visible proxies based on ViewportOffset
			var visibleProxies []string
			if totalProxies <= proxyLines {
				visibleProxies = selectedProxy.All
			} else {
				startIdx := m.ViewportOffset
				if startIdx < 0 {
					startIdx = 0
				}
				if startIdx > totalProxies-proxyLines {
					startIdx = totalProxies - proxyLines
				}
				endIdx := startIdx + proxyLines
				if endIdx > totalProxies {
					endIdx = totalProxies
				}
				visibleProxies = selectedProxy.All[startIdx:endIdx]
			}

			for j, p := range visibleProxies {
				if linesUsed >= availableHeight {
					break
				}
				actualIdx := j + m.ViewportOffset
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
				if actualIdx == m.Cursor && totalProxies > proxyLines {
					line += helpStyle.Render(fmt.Sprintf(" (%d/%d)", m.Cursor+1, totalProxies))
				}
				s += line + "\n"
				linesUsed++
			}
		}
	}

	// Fill remaining lines with blank lines to keep help at bottom
	for linesUsed < availableHeight {
		s += "\n"
		linesUsed++
	}

	// Add help text at bottom
	s += helpStyle.Render(" [←h]Prev [→l]Next  [↑k]↑ [↓j]↓  [Ent]Select  [r]Reload  [q]Quit")

	return s
}
