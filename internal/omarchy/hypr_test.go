package omarchy

import "testing"

func TestPanelAt(t *testing.T) {
	// 2528x960 at 1.6 is 1580x600 in layout coordinates; Omarchy's bar
	// takes 26 at the top. A second monitor to its right, turned, at 1.
	mons := []hyprMonitor{
		{Name: "HDMI-A-2", Width: 2528, Height: 960, Scale: 1.6, Reserved: [4]int{0, 26, 0, 0}},
		{Name: "DP-1", X: 1580, Width: 1920, Height: 1080, Scale: 1, Transform: 1},
	}
	for _, c := range []struct {
		name       string
		monitor    string
		cx, cy     int
		x, y, room int
		ok         bool
	}{
		{"middle", "HDMI-A-2", 800, 10, 580, 32, 562, true},
		{"right edge", "HDMI-A-2", 1570, 10, 1580 - 6 - 440, 32, 562, true},
		{"left edge", "HDMI-A-2", 20, 10, 6, 32, 562, true},
		{"turned monitor", "DP-1", 1700, 1500, 6, 6, 1920 - 12, true},
		{"off every monitor", "", -5, 10, 0, 0, 0, false},
	} {
		monitor, x, y, room, ok := panelAt(440, c.cx, c.cy, mons)
		if monitor != c.monitor || x != c.x || y != c.y || room != c.room || ok != c.ok {
			t.Errorf("%s: %s %d,%d room %d %v; want %s %d,%d room %d %v", c.name, monitor, x, y, room, ok, c.monitor, c.x, c.y, c.room, c.ok)
		}
	}
}

// Monitor offsets must not change the rule position or the available height.
func TestPanelAtMonitorOrigin(t *testing.T) {
	for _, origin := range [][2]int{{0, 0}, {1920, 0}, {-1920, 0}, {0, 1080}, {0, -1080}, {-1920, -1080}} {
		m := hyprMonitor{Name: "DP-1", X: origin[0], Y: origin[1], Width: 1920, Height: 1080, Scale: 1, Reserved: [4]int{10, 34, 20, 30}}
		monitor, x, y, room, ok := panelAt(440, m.X+1538, m.Y+17, []hyprMonitor{m})
		if !ok || monitor != m.Name || x != 1318 || y != 40 || room != 1004 {
			t.Errorf("origin %v: %s %d,%d room %d %v; want DP-1 1318,40 room 1004 true", origin, monitor, x, y, room, ok)
		}
	}
}

func TestRegexpQuote(t *testing.T) {
	for in, want := range map[string]string{
		"magpie panel": "magpie panel",
		"a.b(c)":       `a\.b\(c\)`,
		"x+y*":         `x\+y\*`,
	} {
		if got := regexpQuote(in); got != want {
			t.Errorf("regexpQuote(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestOutside(t *testing.T) {
	at, size := [2]int{1134, 32}, [2]int{440, 520}
	for _, c := range []struct {
		name string
		x, y int
		want bool
	}{
		{"on the panel", 1300, 300, false},
		{"on its top-left corner", 1134, 32, false},
		{"left of it", 900, 300, true},
		{"under it", 1300, 560, true},
		{"on the bar", 1400, 10, false},
		{"in the gap under the bar", 900, 28, true},
	} {
		if got := outside(c.x, c.y, at, size); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
}
