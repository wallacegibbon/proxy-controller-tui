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

func TestNavigationIndicators(t *testing.T) {
	// Test with 3 groups: First, Middle, Last
	m := Model{
		Proxies: map[string]clash.Proxy{
			"First": {
				Name: "First",
				Type: "Selector",
				Now:  "First-1",
				All:  []string{"First-1"},
			},
			"Middle": {
				Name: "Middle",
				Type: "Selector",
				Now:  "Middle-1",
				All:  []string{"Middle-1"},
			},
			"Last": {
				Name: "Last",
				Type: "Selector",
				Now:  "Last-1",
				All:  []string{"Last-1"},
			},
		},
		Groups:     []string{"First", "Middle", "Last"},
		CurrentIdx: 0,
		Cursor:     0,
		Loading:    false,
		Height:     24,
	}

	// Test first group - should show ">>" (has right, no left)
	out := m.View()
	if !strings.Contains(out, "First (Selector) >>") {
		t.Errorf("First group should show '>>' indicator, got:\n%s", out)
	}
	if strings.Contains(out, "<< First") {
		t.Errorf("First group should NOT show '<<' indicator, got:\n%s", out)
	}

	// Test middle group - should show both "<<" and ">>"
	m.CurrentIdx = 1
	out = m.View()
	if !strings.Contains(out, "<< Middle (Selector) >>") {
		t.Errorf("Middle group should show both '<<' and '>>' indicators, got:\n%s", out)
	}

	// Test last group - should show "<<" (has left, no right)
	m.CurrentIdx = 2
	out = m.View()
	if !strings.Contains(out, "<< Last (Selector)") {
		t.Errorf("Last group should show '<<' indicator, got:\n%s", out)
	}
	if strings.Contains(out, "Last (Selector) >>") {
		t.Errorf("Last group should NOT show '>>' indicator, got:\n%s", out)
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

	// Verify help is at the bottom of terminal
	lines := strings.Split(out, "\n")
	if len(lines) > m.Height {
		t.Errorf("Output exceeds terminal height: got %d lines, terminal height is %d", len(lines), m.Height)
	}
	lastLine := lines[len(lines)-1]
	if !strings.Contains(lastLine, "[q]Quit") {
		t.Errorf("Help message not on last line, got: %q", lastLine)
	}
}

func TestHelpAtBottomSmallTerminal(t *testing.T) {
	m := Model{
		Proxies: map[string]clash.Proxy{
			"Proxy": {
				Name: "Proxy",
				Type: "Selector",
				Now:  "Proxy-1",
				All:  []string{"Proxy-1", "Proxy-2"},
			},
		},
		Groups:     []string{"Proxy"},
		CurrentIdx: 0,
		Cursor:     0,
		Loading:    false,
		Height:     8, // Very small terminal
	}
	out := m.View()
	lines := strings.Split(out, "\n")

	// For terminal height 8, we expect:
	// - 1 group line
	// - 2 proxy lines
	// - Some padding
	// - 1 help line
	// Total should be 8
	if len(lines) != m.Height {
		t.Errorf("Expected output to be exactly %d lines (terminal height), got %d", m.Height, len(lines))
		t.Logf("Output:\n%s", out)
	}

	lastLine := lines[len(lines)-1]
	if !strings.Contains(lastLine, "[q]Quit") {
		t.Errorf("Help message not on last line, got: %q", lastLine)
	}
}

func TestLayoutWithMultipleGroups(t *testing.T) {
	m := Model{
		Proxies: map[string]clash.Proxy{
			"Proxy": {
				Name: "Proxy",
				Type: "Selector",
				Now:  "Proxy-1",
				All:  []string{"Proxy-1", "Proxy-2"},
			},
			"Auto": {
				Name: "Auto",
				Type: "URLTest",
				Now:  "Auto-1",
				All:  []string{"Auto-1", "Auto-2"},
			},
		},
		Groups:     []string{"Proxy", "Auto"},
		CurrentIdx: 0, // Proxy is selected
		Cursor:     0,
		Loading:    false,
		Height:     15,
	}

	out := m.View()
	lines := strings.Split(out, "\n")

	// Output should not exceed terminal height
	if len(lines) > m.Height {
		t.Errorf("Output exceeds terminal height: got %d lines, terminal height is %d", len(lines), m.Height)
	}

	// Help should be on last line
	lastLine := lines[len(lines)-1]
	if !strings.Contains(lastLine, "[q]Quit") {
		t.Errorf("Help message not on last line, got: %q", lastLine)
	}

	// Groups should be in order (Proxy then Auto) with types
	foundProxy := false
	foundAuto := false
	for i, line := range lines {
		if strings.Contains(line, "Proxy") && strings.Contains(line, "(Selector)") && foundProxy == false {
			foundProxy = true
			// Next line(s) should be proxies
			if i+1 < len(lines) && strings.Contains(lines[i+1], "Proxy-") {
				// Good, proxies follow the group
			}
		}
		if strings.Contains(line, "Auto") && strings.Contains(line, "(URLTest)") && foundAuto == false {
			foundAuto = true
			// Auto should appear after Proxy
			if !foundProxy {
				t.Errorf("Expected 'Proxy' to appear before 'Auto'")
			}
		}
	}

	if !foundProxy {
		t.Errorf("Expected group 'Proxy' to be in output")
	}
	if !foundAuto {
		t.Errorf("Expected group 'Auto' to be in output")
	}
}

func TestManyGroupsExceedsTerminalHeight(t *testing.T) {
	// Test case: more groups than terminal lines
	// The selected group should always be visible
	groups := []string{}
	proxies := map[string]clash.Proxy{}
	groupNames := []string{"GroupA", "GroupB", "GroupC", "GroupD", "GroupE", "GroupF", "GroupG", "GroupH", "GroupI", "GroupJ", "GroupK", "GroupL", "GroupM", "GroupN", "GroupO", "GroupP", "GroupQ", "GroupR", "GroupS", "GroupT"}
	for _, name := range groupNames {
		groups = append(groups, name)
		proxies[name] = clash.Proxy{
			Name: name,
			Type: "Selector",
			Now:  name + "-1",
			All:  []string{name + "-1", name + "-2"},
		}
	}

	// Test with the first group selected (CurrentIdx = 0)
	m := Model{
		Proxies:    proxies,
		Groups:     groups,
		CurrentIdx: 0, // First group selected
		Cursor:     0,
		Loading:    false,
		Height:     10, // Small terminal - can't fit all 20 groups
	}

	out := m.View()
	lines := strings.Split(out, "\n")

	// Output should not exceed terminal height
	if len(lines) != m.Height {
		t.Errorf("Expected output to be exactly %d lines, got %d", m.Height, len(lines))
	}

	// The selected group (GroupA) must be visible
	foundGroupA := false
	for _, line := range lines {
		if strings.Contains(line, "GroupA") {
			foundGroupA = true
			break
		}
	}
	if !foundGroupA {
		t.Errorf("Selected group 'GroupA' should be visible when it is selected, but not found in output:\n%s", out)
	}

	// Help should be on last line
	lastLine := lines[len(lines)-1]
	if !strings.Contains(lastLine, "[q]Quit") {
		t.Errorf("Help message not on last line, got: %q", lastLine)
	}

	// Test with last group selected (CurrentIdx = 19)
	m.CurrentIdx = 19 // Last group selected
	out = m.View()
	lines = strings.Split(out, "\n")

	if len(lines) != m.Height {
		t.Errorf("Expected output to be exactly %d lines, got %d", m.Height, len(lines))
	}

	// The selected group (GroupT) must be visible
	foundGroupT := false
	for _, line := range lines {
		if strings.Contains(line, "GroupT") {
			foundGroupT = true
			break
		}
	}
	if !foundGroupT {
		t.Errorf("Selected group 'GroupT' should be visible when it is selected, but not found in output:\n%s", out)
	}

	// Test with middle group selected (CurrentIdx = 10)
	m.CurrentIdx = 10 // Middle group selected
	out = m.View()
	lines = strings.Split(out, "\n")

	// The selected group (GroupK - index 10) must be visible
	foundGroupK := false
	for _, line := range lines {
		if strings.Contains(line, "GroupK") {
			foundGroupK = true
			break
		}
	}
	if !foundGroupK {
		t.Errorf("Selected group 'GroupK' should be visible when it is selected, but not found in output:\n%s", out)
	}
}