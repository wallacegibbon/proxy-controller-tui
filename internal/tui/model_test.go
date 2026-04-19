package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/wallacegibbon/proxy-controller-tui/internal/clash"
)

func TestURLTestFixedIndicator(t *testing.T) {
	// Test that URLTest groups show [fixed] indicator when Fixed field is set
	m := Model{
		Proxies: map[string]clash.Proxy{
			"Auto": {
				Name:  "Auto",
				Type:  "URLTest",
				Now:   "Auto-2",
				Fixed: "Auto-2",
				All:   []string{"Auto-1", "Auto-2", "Auto-3", "Auto-4"},
			},
		},
		Groups:     []string{"Auto"},
		CurrentIdx: 0,
		Cursor:     0,
		Loading:    false,
		Height:     24,
	}

	v := m.View()
	out := v.Content
	t.Logf("View output:\n%s", out)

	// Should show [fixed] indicator for URLTest group with fixed proxy
	if !strings.Contains(out, "Auto (URLTest)") {
		t.Errorf("Expected group name with type 'Auto (URLTest)' in output")
	}

	// Check for fixed indicator (in styled form)
	if !strings.Contains(out, "[fixed]") {
		t.Errorf("Expected '[fixed]' indicator in output for URLTest group with fixed proxy")
	}

}

func TestURLTestWithoutFixed(t *testing.T) {
	// Test that URLTest groups without Fixed don't show [fixed] indicator
	m := Model{
		Proxies: map[string]clash.Proxy{
			"Auto": {
				Name:  "Auto",
				Type:  "URLTest",
				Now:   "Auto-2",
				Fixed: "",
				All:   []string{"Auto-1", "Auto-2", "Auto-3", "Auto-4"},
			},
		},
		Groups:     []string{"Auto"},
		CurrentIdx: 0,
		Cursor:     0,
		Loading:    false,
		Height:     24,
	}

	v := m.View()
	out := v.Content
	t.Logf("View output:\n%s", out)

	// Should NOT show [fixed] indicator
	if strings.Contains(out, "[fixed]") {
		t.Errorf("Should NOT show '[fixed]' indicator for URLTest group without fixed proxy")
	}

	// Should NOT show "↩ Auto" reset option
	if strings.Contains(out, "Auto (restore auto-selection)") {
		t.Errorf("Should NOT show reset option for URLTest group without fixed proxy")
	}
}

func TestResetFixedKeyBinding(t *testing.T) {
	// Test pressing 'a' key on URLTest group with fixed proxy
	m := Model{
		Client: clash.NewClient(""),
		Proxies: map[string]clash.Proxy{
			"Auto": {
				Name:  "Auto",
				Type:  "URLTest",
				Now:   "Auto-2",
				Fixed: "Auto-2",
				All:   []string{"Auto-1", "Auto-2", "Auto-3", "Auto-4"},
			},
		},
		Groups:     []string{"Auto"},
		CurrentIdx: 0,
		Cursor:     0,
		Loading:    false,
		Height:     24,
	}

	// Press 'a' key to reset fixed
	newModel, cmd := m.Update(tea.KeyPressMsg(tea.Key{Text: "a", Code: 'a'}))
	m2 := newModel.(Model)

	// Should be loading after pressing 'a'
	if !m2.Loading {
		t.Errorf("Expected Loading to be true after pressing 'a' on URLTest group with fixed proxy")
	}

	// Should have a command to execute
	if cmd == nil {
		t.Errorf("Expected a command to be returned after pressing 'a'")
	}
}

func TestCursorMovement(t *testing.T) {
	m := Model{
		Client: clash.NewClient(""),
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

	// Test 'j' key (move down)
	newModel, _ := m.Update(tea.KeyPressMsg(tea.Key{Text: "j", Code: 'j'}))
	m2 := newModel.(Model)
	if m2.Cursor != 1 {
		t.Errorf("Expected cursor to move down to 1, got %d", m2.Cursor)
	}

	// Test 'k' key (move up) from position 2
	m2.Cursor = 2
	newModel, _ = m2.Update(tea.KeyPressMsg(tea.Key{Text: "k", Code: 'k'}))
	m3 := newModel.(Model)
	if m3.Cursor != 1 {
		t.Errorf("Expected cursor to move up to 1, got %d", m3.Cursor)
	}

	// Test 'k' at top (should stay)
	m.Cursor = 0
	newModel, _ = m.Update(tea.KeyPressMsg(tea.Key{Text: "k", Code: 'k'}))
	m4 := newModel.(Model)
	if m4.Cursor != 0 {
		t.Errorf("Expected cursor to stay at 0 when at top, got %d", m4.Cursor)
	}

	// Test 'j' at bottom (should stay)
	m.Cursor = 2
	newModel, _ = m.Update(tea.KeyPressMsg(tea.Key{Text: "j", Code: 'j'}))
	m5 := newModel.(Model)
	if m5.Cursor != 2 {
		t.Errorf("Expected cursor to stay at 2 when at bottom, got %d", m5.Cursor)
	}

	// Test 'l' key (navigate to next group)
	m.CurrentIdx = 0
	m.Cursor = 1
	newModel, _ = m.Update(tea.KeyPressMsg(tea.Key{Text: "l", Code: 'l'}))
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
	v := m.View()
	out := v.Content
	if !strings.Contains(out, "First (Selector) >>") {
		t.Errorf("First group should show '>>' indicator, got:\n%s", out)
	}
	if strings.Contains(out, "<< First") {
		t.Errorf("First group should NOT show '<<' indicator, got:\n%s", out)
	}

	// Test middle group - should show both "<<" and ">>"
	m.CurrentIdx = 1
	v = m.View()
	out = v.Content
	if !strings.Contains(out, "<< Middle (Selector) >>") {
		t.Errorf("Middle group should show both '<<' and '>>' indicators, got:\n%s", out)
	}

	// Test last group - should show "<<" (has left, no right)
	m.CurrentIdx = 2
	v = m.View()
	out = v.Content
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
	v := m.View()
	out := v.Content
	t.Logf("View output:\n%s", out)
	if !strings.Contains(out, ">  ") {
		t.Errorf("Expected cursor marker '>  ' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Proxy-1") {
		t.Errorf("Expected active proxy 'Proxy-1' in output, got:\n%s", out)
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
	v := m.View()
	out := v.Content
	t.Logf("View output:\n%s", out)
	if !strings.Contains(out, ">> ") {
		t.Errorf("Expected cursor marker '>>' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Proxy-1") {
		t.Errorf("Expected active proxy 'Proxy-1' in output, got:\n%s", out)
	}

}

func TestManyGroupsExceedsTerminalHeight(t *testing.T) {
	// Test case: many groups, but only the selected group is shown
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
		CurrentIdx: 0,
		Cursor:     0,
		Loading:    false,
		Height:     10,
	}

	v := m.View()
	out := v.Content
	lines := strings.Split(out, "\n")

	// Output should not exceed terminal height (accounting for trailing newline)
	if len(lines) > m.Height+1 {
		t.Errorf("Output exceeds terminal height: got %d lines, terminal height is %d", len(lines)-1, m.Height)
	}

	// The selected group (GroupA) must be visible
	if !strings.Contains(out, "GroupA") {
		t.Errorf("Selected group 'GroupA' should be visible:\n%s", out)
	}

	// Other groups should NOT be visible
	if strings.Contains(out, "GroupB") {
		t.Errorf("Other groups should NOT be visible when GroupA is selected:\n%s", out)
	}

	// Test with last group selected (CurrentIdx = 19)
	m.CurrentIdx = 19
	v = m.View()
	out = v.Content

	if !strings.Contains(out, "GroupT") {
		t.Errorf("Selected group 'GroupT' should be visible:\n%s", out)
	}

	// Test with middle group selected (CurrentIdx = 10)
	m.CurrentIdx = 10
	v = m.View()
	out = v.Content

	if !strings.Contains(out, "GroupK") {
		t.Errorf("Selected group 'GroupK' should be visible:\n%s", out)
	}
}
