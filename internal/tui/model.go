package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wallacegibbon/proxy-tui-controller/internal/clash"
)

type errMsg error

func (m Model) Init() tea.Cmd {
	return LoadProxiesCmd(m.Client)
}
type proxiesLoadedMsg struct {
	proxies map[string]clash.Proxy
	groups  []string
}

const (
	maxVisibleProxies = 20
	minHelpRows      = 2 // footer + separator
)

var (
	headerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("147")).Bold(true)
	selectedGroupStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	normalStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	activeProxyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	activeProxyMarkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
	cursorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true)
	selectedStyle      = lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(lipgloss.Color("255"))
	helpStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	borderStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	separatorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type Model struct {
	Client         *clash.Client
	Proxies        map[string]clash.Proxy
	Groups         []string
	CurrentIdx     int
	Cursor         int
	Loading        bool
	Err            error
	ViewportOffset int
	Height         int // Terminal height
}

func InitialModel() Model {
	client := clash.NewClient("")
	return Model{
		Client:         client,
		Proxies:        make(map[string]clash.Proxy),
		Groups:         make([]string, 0),
		CurrentIdx:     0,
		Cursor:         0,
		Loading:        true,
		Err:            nil,
		ViewportOffset: 0,
		Height:         24,
	}
}

func LoadProxiesCmd(client *clash.Client) tea.Cmd {
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
