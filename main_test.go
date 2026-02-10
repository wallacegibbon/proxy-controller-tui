package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"testing"
)

func TestCursorMovement(t *testing.T) {
	// Create a model with mock data
	m := model{
		proxies: map[string]Proxy{
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
		groups:     []string{"Proxy", "Auto"},
		currentIdx: 0,
		cursor:     0,
		loading:    false,
	}

	// Test j key moves cursor down
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m2 := newModel.(model)
	if m2.cursor != 1 {
		t.Errorf("Expected cursor to move down to 1, got %d", m2.cursor)
	}

	// Test k key moves cursor up
	m2.cursor = 2
	newModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m3 := newModel.(model)
	if m3.cursor != 1 {
		t.Errorf("Expected cursor to move up to 1, got %d", m3.cursor)
	}

	// Test cursor doesn't go below 0
	m.cursor = 0
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m4 := newModel.(model)
	if m4.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0 when at top, got %d", m4.cursor)
	}

	// Test cursor doesn't go beyond last proxy
	m.cursor = 2 // last index in Proxy group (Proxy-3)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m5 := newModel.(model)
	if m5.cursor != 2 {
		t.Errorf("Expected cursor to stay at 2 when at bottom, got %d", m5.cursor)
	}

	// Test group switching with l (right)
	m.currentIdx = 0
	m.cursor = 1
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m6 := newModel.(model)
	if m6.currentIdx != 1 {
		t.Errorf("Expected currentIdx to move to 1, got %d", m6.currentIdx)
	}
	// After switching to Auto group, cursor should point to active proxy Auto-2 (index 1)
	expectedCursor := 1 // Auto-2 is at index 1 in Auto.All
	if m6.cursor != expectedCursor {
		t.Errorf("Expected cursor to be %d (active proxy) after group switch, got %d", expectedCursor, m6.cursor)
	}
}

func TestViewCursor(t *testing.T) {
	m := model{
		proxies: map[string]Proxy{
			"Proxy": {
				Name: "Proxy",
				Type: "Selector",
				Now:  "Proxy-1",
				All:  []string{"Proxy-1", "Proxy-2", "Proxy-3"},
			},
		},
		groups:     []string{"Proxy"},
		currentIdx: 0,
		cursor:     1, // cursor on Proxy-2 (index 1)
		loading:    false,
	}
	out := m.View()
	t.Logf("View output:\n%s", out)
	// Check that cursor marker appears
	if !strings.Contains(out, ">  ") {
		t.Errorf("Expected cursor marker '>  ' in output, got:\n%s", out)
	}
	// Check that active proxy marker appears
	if !strings.Contains(out, " ● Proxy-1") {
		t.Errorf("Expected active proxy marker ' ● Proxy-1' in output, got:\n%s", out)
	}
}
func TestViewCursorOnActive(t *testing.T) {
	m := model{
		proxies: map[string]Proxy{
			"Proxy": {
				Name: "Proxy",
				Type: "Selector",
				Now:  "Proxy-1",
				All:  []string{"Proxy-1", "Proxy-2", "Proxy-3"},
			},
		},
		groups:     []string{"Proxy"},
		currentIdx: 0,
		cursor:     0, // cursor on active proxy Proxy-1 (index 0)
		loading:    false,
	}
	out := m.View()
	t.Logf("View output:\n%s", out)
	// Check that combined marker appears when cursor is on active proxy
	if !strings.Contains(out, ">● Proxy-1") {
		t.Errorf("Expected combined marker '>● Proxy-1' when cursor is on active proxy, got:\n%s", out)
	}
	// Check that active proxy marker appears
	if !strings.Contains(out, "●") {
		t.Errorf("Expected active proxy marker ● in output, got:\n%s", out)
	}
}
