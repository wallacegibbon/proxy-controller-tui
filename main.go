package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error
type proxiesLoadedMsg struct {
	proxies map[string]Proxy
	groups  []string
}

const (
	maxVisibleProxies = 20
)

var (
	headerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("147")).Bold(true)
	selectedGroupStyle = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("231")).Bold(true)
	normalStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	activeProxyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	cursorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true)
	selectedStyle      = lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(lipgloss.Color("255"))
	helpStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	borderStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	separatorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type model struct {
	client         *ClashClient
	proxies        map[string]Proxy
	groups         []string
	currentIdx     int
	cursor         int
	loading        bool
	err            error
	viewportOffset int
}

func initialModel() model {
	client := NewClashClient("")
	return model{
		client:         client,
		proxies:        make(map[string]Proxy),
		groups:         make([]string, 0),
		currentIdx:     0,
		cursor:         0,
		loading:        true,
		err:            nil,
		viewportOffset: 0,
	}
}

func loadProxiesCmd(client *ClashClient) tea.Cmd {
	return func() tea.Msg {
		proxies, err := client.GetProxies()
		if err != nil {
			return errMsg(err)
		}

		groups := make([]string, 0)
		for name, proxy := range proxies.Proxies {
			if proxy.Type == "Selector" || proxy.Type == "URLTest" {
				groups = append(groups, name)
			}
		}

		return proxiesLoadedMsg{
			proxies: proxies.Proxies,
			groups:  groups,
		}
	}
}

func (m model) Init() tea.Cmd {
	return loadProxiesCmd(m.client)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.loading = false
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		// Adjust viewport offset if needed when window resizes
		m.adjustViewport()
		return m, nil

	case proxiesLoadedMsg:
		m.loading = false
		m.proxies = msg.proxies
		m.groups = msg.groups
		// Ensure currentIdx is valid after reload
		if m.currentIdx >= len(m.groups) {
			m.currentIdx = 0
		}
		// Set cursor to the active proxy in the current group
		if len(m.groups) > 0 && m.currentIdx < len(m.groups) {
			if proxy, ok := m.proxies[m.groups[m.currentIdx]]; ok {
				cursorFound := false
				for i, p := range proxy.All {
					if p == proxy.Now {
						m.cursor = i
						cursorFound = true
						break
					}
				}
				if !cursorFound && len(proxy.All) > 0 {
					m.cursor = 0
				}
			}
		}
		m.viewportOffset = 0
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.Type {
		case tea.KeyUp, tea.KeyCtrlK:
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && len(proxy.All) > 0 {
					if m.cursor > 0 {
						m.cursor--
						m.adjustViewport()
					}
				}
			}
			return m, nil

		case tea.KeyDown, tea.KeyCtrlJ:
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && len(proxy.All) > 0 {
					if m.cursor < len(proxy.All)-1 {
						m.cursor++
						m.adjustViewport()
					}
				}
			}
			return m, nil

		case tea.KeyLeft:
			if m.currentIdx > 0 {
				m.currentIdx--
				// Set cursor to active proxy in the new group
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.cursor = i
							break
						}
					}
					m.viewportOffset = 0
				} else {
					m.cursor = 0
				}
			}
			return m, nil

		case tea.KeyRight:
			if m.currentIdx < len(m.groups)-1 {
				m.currentIdx++
				// Set cursor to active proxy in the new group
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.cursor = i
							break
						}
					}
					m.viewportOffset = 0
				} else {
					m.cursor = 0
				}
			}
			return m, nil

		case tea.KeyEnter:
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && m.cursor < len(proxy.All) {
					selectedProxy := proxy.All[m.cursor]
					if err := m.client.SelectProxy(group, selectedProxy); err != nil {
						m.err = err
						return m, nil
					}
					// Cursor is already at the right position, just reload
					return m, loadProxiesCmd(m.client)
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
			m.loading = true
			return m, loadProxiesCmd(m.client)
		case "h":
			if m.currentIdx > 0 {
				m.currentIdx--
				// Set cursor to active proxy in the new group
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.cursor = i
							break
						}
					}
					m.viewportOffset = 0
				} else {
					m.cursor = 0
				}
			}
			return m, nil
		case "l":
			if m.currentIdx < len(m.groups)-1 {
				m.currentIdx++
				// Set cursor to active proxy in the new group
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.cursor = i
							break
						}
					}
					m.viewportOffset = 0
				} else {
					m.cursor = 0
				}
			}
			return m, nil
		case "k":
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && len(proxy.All) > 0 {
					if m.cursor > 0 {
						m.cursor--
						m.adjustViewport()
					}
				}
			}
			return m, nil
		case "j":
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && len(proxy.All) > 0 {
					if m.cursor < len(proxy.All)-1 {
						m.cursor++
						m.adjustViewport()
					}
				}
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *model) adjustViewport() {
	if len(m.groups) == 0 {
		return
	}
	group := m.groups[m.currentIdx]
	proxy, ok := m.proxies[group]
	if !ok {
		return
	}

	// Ensure cursor is visible within viewport
	if m.cursor < m.viewportOffset {
		m.viewportOffset = m.cursor
	} else if m.cursor >= m.viewportOffset+maxVisibleProxies {
		m.viewportOffset = m.cursor - maxVisibleProxies + 1
	}

	// Ensure offset doesn't go negative
	if m.viewportOffset < 0 {
		m.viewportOffset = 0
	}

	// Ensure offset doesn't exceed list length
	maxOffset := len(proxy.All) - maxVisibleProxies
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.viewportOffset > maxOffset {
		m.viewportOffset = maxOffset
	}
}

func (m model) View() string {
	if m.loading {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  Loading proxies...")
	}

	if m.err != nil {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  Error") + "\n" +
			fmt.Sprintf("  %v\n", m.err) +
			helpStyle.Render("  Press [r] retry, [q] quit")
	}

	if len(m.groups) == 0 {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  No proxy groups found") + "\n" +
			helpStyle.Render("  Press [r] refresh, [q] quit")
	}

	var s string

	for i, group := range m.groups {
		proxy, ok := m.proxies[group]
		if !ok {
			continue
		}

		var groupLabel string
		if i == m.currentIdx {
			groupLabel = selectedGroupStyle.Render("● " + group)
		} else {
			groupLabel = normalStyle.Render("○ " + group)
		}
		s += groupLabel + "\n"

		if i == m.currentIdx {
			visibleProxies := proxy.All
			if len(proxy.All) > maxVisibleProxies {
				startIdx := m.viewportOffset
				if startIdx < 0 {
					startIdx = 0
				}
				endIdx := startIdx + maxVisibleProxies
				if endIdx > len(proxy.All) {
					endIdx = len(proxy.All)
				}
				visibleProxies = proxy.All[startIdx:endIdx]
			}

			for j, p := range visibleProxies {
				actualIdx := j + m.viewportOffset
				var line string
				if actualIdx == m.cursor && p == proxy.Now {
					line = cursorStyle.Render(">● ") + activeProxyStyle.Render(p)
				} else if actualIdx == m.cursor {
					line = cursorStyle.Render(">  ") + p
				} else if p == proxy.Now {
					line = " ● " + activeProxyStyle.Render(p)
				} else {
					line = "   " + normalStyle.Render(p)
				}
				s += line + "\n"
			}

			if len(proxy.All) > maxVisibleProxies {
				s += helpStyle.Render("  " + strings.Repeat("-", 20) + fmt.Sprintf(" %d/%d ", len(proxy.All), len(proxy.All))) + "\n"
			}
		}
	}

	s += separatorStyle.Render("═══════════════════════════════════════") + "\n"
	s += helpStyle.Render(" [←h]Prev [→l]Next  [↑k]↑ [↓j]↓  [Ent]Select  [r]Reload  [q]Quit")

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
