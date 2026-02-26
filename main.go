package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/randomradio/proxy-controller-tui/internal/tui"
)

func main() {
	url := flag.String("url", "", "Clash/Mihomo REST API base URL (e.g. http://192.168.1.10:9090)")
	secret := flag.String("secret", "", "Clash/Mihomo API secret")
	flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	p := tea.NewProgram(
		tui.InitialModel(*url, *secret),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
