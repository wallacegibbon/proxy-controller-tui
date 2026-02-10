package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"

	"github.com/wallacegibbon/proxy-controller-tui/internal/clash"
)

func TestCursorMovement(t *testing.T) {
	m := Model{
		Proxies: map[string]clash.Proxy{
			"Proxy": {
				Name: "Proxy",
				Type: "Selector",
				Now:  "Proxy-1",
				All:  []string{"Proxy-1", "Proxy-2", "Proxy-3"},
			},
			"Auto": {
				Name: "Auto",
				Type: "URLTest",
				Now:  "Auto-2",
				All:  []string{"Auto-1", "Auto-2", "Auto-3", "Auto-4"},
			},
		},
		Groups:     []string{"Proxy", "Auto"},
		CurrentIdx: 0,
		Cursor:     0,
		Loading:    false,
		Height:     24,
	}

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m2 := newModel.(Model)
	if m2.Cursor != 1 {
		t.Errorf("Expected cursor to move down to 1, got %d", m2.Cursor)
	}

	m2.Cursor = 2
	newModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m3 := newModel.(Model)
	if m3.Cursor != 1 {
		t.Errorf("Expected cursor to move up to 1, got %d", m3.Cursor)
	}

	m.Cursor = 0
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m4 := newModel.(Model)
	if m4.Cursor != 0 {
		t.Errorf("Expected cursor to stay at 0 when at top, got %d", m4.Cursor)
	}

	m.Cursor = 2
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m5 := newModel.(Model)
	if m5.Cursor != 2 {
		t.Errorf("Expected cursor to stay at 2 when at bottom, got %d", m5.Cursor)
	}

	m.CurrentIdx = 0
	m.Cursor = 1
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m6 := newModel.(Model)
	if m6.CurrentIdx != 1 {
		t.Errorf("Expected currentIdx to move to 1, got %d", m6.CurrentIdx)
	}
	expectedCursor := 1
	if m6.Cursor != expectedCursor {
		t.Errorf("Expected cursor to be %d (active proxy) after group switch, got %d", expectedCursor, m6.Cursor)
	}
}

func TestViewCursor(t *testing.T) {
	m := Model{
		Proxies: map[string]clash.Proxy{
			"Proxy": {
				Name: "Proxy",
				Type: "Selector",
				Now:  "Proxy-1",
				All:  []string{"Proxy-1", "Proxy-2", "Proxy-3"},
			},
		},
		Groups:     []string{"Proxy"},
		CurrentIdx: 0,
		Cursor:     1,
		Loading:    false,
		Height:     24,
	}
	out := m.View()
	t.Logf("View output:\n%s", out)
	if !strings.Contains(out, ">  ") {
		t.Errorf("Expected cursor marker '>  ' in output, got:\n%s", out)
	}
	if !strings.Contains(out, " > Proxy-1") {
		t.Errorf("Expected active proxy marker ' > Proxy-1' in output, got:\n%s", out)
	}
}

func TestViewCursorOnActive(t *testing.T) {
	m := Model{
		Proxies: map[string]clash.Proxy{
			"Proxy": {
				Name: "Proxy",
				Type: "Selector",
				Now:  "Proxy-1",
				All:  []string{"Proxy-1", "Proxy-2", "Proxy-3"},
			},
		},
		Groups:     []string{"Proxy"},
		CurrentIdx: 0,
		Cursor:     0,
		Loading:    false,
		Height:     24,
	}
	out := m.View()
	t.Logf("View output:\n%s", out)
	if !strings.Contains(out, ">> Proxy-1") {
		t.Errorf("Expected combined marker '>> Proxy-1' when cursor is on active proxy, got:\n%s", out)
	}
	if !strings.Contains(out, ">") {
		t.Errorf("Expected active proxy marker > in output, got:\n%s", out)
	}
}
