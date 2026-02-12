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

	// Calculate max group name display width for uniform padding
	maxGroupWidth := 0
	for _, group := range m.Groups {
		groupWidth := lipgloss.Width(group)
		if groupWidth > maxGroupWidth {
			maxGroupWidth = groupWidth
		}
	}

	var s string

	// First, show the selected group and its proxies at the top
	selectedGroup := m.Groups[m.CurrentIdx]
	selectedProxy, ok := m.Proxies[selectedGroup]
	if ok {
		// Pad group name to uniform display width with 3 spaces on each side
		currentWidth := lipgloss.Width(selectedGroup)
		paddedGroup := "   " + selectedGroup + strings.Repeat(" ", maxGroupWidth-currentWidth) + "   "
		s += selectedGroupStyle.Render(paddedGroup) + "\n"

		// Calculate how many proxies we can show
		// Footer takes: help (1 row) + unselected groups (len(m.Groups)-1 rows)
		unselectedGroupCount := len(m.Groups) - 1
		availableRows := m.Height - unselectedGroupCount - minHelpRows - 1
		if availableRows < 1 {
			availableRows = 1
		}
		visibleCount := availableRows

		visibleProxies := selectedProxy.All
		if len(selectedProxy.All) > visibleCount {
			startIdx := m.ViewportOffset
			if startIdx < 0 {
				startIdx = 0
			}
			endIdx := startIdx + visibleCount
			if endIdx > len(selectedProxy.All) {
				endIdx = len(selectedProxy.All)
			}
			visibleProxies = selectedProxy.All[startIdx:endIdx]
		}

		for j, p := range visibleProxies {
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
			if actualIdx == m.Cursor && len(selectedProxy.All) > visibleCount {
				line += helpStyle.Render(fmt.Sprintf(" (%d/%d)", m.Cursor+1, len(selectedProxy.All)))
			}
			s += line + "\n"
		}
	}

	// Collect unselected groups to show at bottom
	var unselectedGroups []string
	for i, group := range m.Groups {
		if i != m.CurrentIdx {
			proxy, ok := m.Proxies[group]
			if !ok {
				continue
			}
			// Pad group name to uniform display width with 3 spaces on each side
			currentWidth := lipgloss.Width(group)
			paddedGroup := "   " + group + strings.Repeat(" ", maxGroupWidth-currentWidth) + "   "
			unselectedGroups = append(unselectedGroups, normalGroupStyle.Render(paddedGroup)+" "+helpStyle.Render("["+proxy.Now+"]"))
		}
	}

	// Count content lines (selected group + proxies)
	contentLines := strings.Count(s, "\n")
	if contentLines > 0 {
		// Add padding to push unselected groups and help to bottom
		// We need: contentLines + padding + unselectedGroups + 1 (help line) = m.Height
		// So: padding = m.Height - contentLines - len(unselectedGroups) - 1
		padding := m.Height - contentLines - len(unselectedGroups) - 1
		if padding > 0 {
			s += strings.Repeat("\n", padding)
		}
	}

	// Add unselected groups just above help
	for _, group := range unselectedGroups {
		s += group + "\n"
	}

	// Add help text at bottom
	var helpText string
	if m.Height < 15 {
		helpText = helpStyle.Render(" h/l:grp  j/k:prox  Ent:sel  r:reload  q:quit")
	} else {
		helpText = helpStyle.Render(" [←h]Prev [→l]Next  [↑k]↑ [↓j]↓  [Ent]Select  [r]Reload  [q]Quit")
	}
	s += helpText

	return s
}
