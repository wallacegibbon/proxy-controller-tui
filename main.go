package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/wallacegibbon/proxy-controller-tui/internal/tui"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	p := tea.NewProgram(
		tui.InitialModel(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
