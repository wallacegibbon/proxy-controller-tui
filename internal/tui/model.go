package tui

import (
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/randomradio/proxy-controller-tui/internal/clash"
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
	minHelpRows = 1 // help text only
)

var (
	headerStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("147")).Bold(true)
	selectedGroupStyle   = lipgloss.NewStyle().Background(lipgloss.Color("45")).Foreground(lipgloss.Color("231")).Bold(true)
	normalGroupStyle     = lipgloss.NewStyle().Background(lipgloss.Color("45")).Foreground(lipgloss.Color("245"))
	normalStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	activeProxyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	activeProxyMarkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
	cursorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true)
	helpStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	separatorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type Model struct {
	Client          *clash.Client
	Proxies         map[string]clash.Proxy
	Groups          []string
	CurrentIdx      int
	Cursor          int
	Loading         bool
	Err             error
	ViewportOffset  int
	Height          int    // Terminal height
	lastCursorProxy string // Track proxy name at cursor to restore position after reload
}

func InitialModel(baseURL, secret string) Model {
	client := clash.NewClient(baseURL, secret)
	return Model{
		Client:          client,
		Proxies:         make(map[string]clash.Proxy),
		Groups:          make([]string, 0),
		CurrentIdx:      0,
		Cursor:          0,
		Loading:         true,
		Err:             nil,
		ViewportOffset:  0,
		Height:          24,
		lastCursorProxy: "",
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

		// Sort groups alphabetically for consistent ordering
		sort.Strings(groups)

		return proxiesLoadedMsg{
			proxies: proxies.Proxies,
			groups:  groups,
		}
	}
}

func loadProxiesWithDelayCmd(client *clash.Client) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(200 * time.Millisecond)

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

		// Sort groups alphabetically for consistent ordering
		sort.Strings(groups)

		return proxiesLoadedMsg{
			proxies: proxies.Proxies,
			groups:  groups,
		}
	}
}
