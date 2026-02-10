package tui

import (
	"fmt"
	"strings"
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

	// Calculate max group name length for uniform padding
	maxGroupLen := 0
	for _, group := range m.Groups {
		if len(group) > maxGroupLen {
			maxGroupLen = len(group)
		}
	}

	var s string

	for i, group := range m.Groups {
		proxy, ok := m.Proxies[group]
		if !ok {
			continue
		}

		// Pad group name to uniform length
		paddedGroup := group + strings.Repeat(" ", maxGroupLen-len(group))

		var groupLabel string
		if i == m.CurrentIdx {
			groupLabel = selectedGroupStyle.Render("● " + paddedGroup)
		} else {
			groupLabel = normalGroupStyle.Render("○ " + paddedGroup)
		}
		s += groupLabel + "\n"

		if i == m.CurrentIdx {
			// Calculate how many proxies we can show
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
	}

	s += separatorStyle.Render("═══════════════════════════════════════") + "\n"
	if m.Height < 15 {
		s += helpStyle.Render(" h/l:grp  j/k:prox  Ent:sel  r:reload  q:quit")
	} else {
		s += helpStyle.Render(" [←h]Prev [→l]Next  [↑k]↑ [↓j]↓  [Ent]Select  [r]Reload  [q]Quit")
	}

	return s
}
