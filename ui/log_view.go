package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/inspiring-group/inspiring-swiss-knife/pkgs"
)

// LogView displays the live operation log.
type LogView struct {
	width  int
	height int
	offset int // scroll offset (lines from top)
}

func NewLogView(w, h int) LogView {
	return LogView{width: w, height: h}
}

// View renders the log tab into exactly contentH lines.
func (v LogView) View(contentH int) string {
	lines := pkgs.GetLogLines()

	visibleH := contentH - 3 // header + border + footer hint
	if visibleH < 1 {
		visibleH = 1
	}

	// Clamp scroll offset
	maxOffset := len(lines) - visibleH
	if maxOffset < 0 {
		maxOffset = 0
	}
	off := v.offset
	if off > maxOffset {
		off = maxOffset
	}
	if off < 0 {
		off = 0
	}

	header := fmt.Sprintf("  %s  (%d entries)",
		StyleSectionHeader.Render("Operation Log"),
		len(lines))
	if off > 0 {
		header += StyleTextMuted("  ↑ scroll")
	}
	if off+visibleH < len(lines) {
		header += StyleTextMuted("  ↓ more")
	}

	divider := StyleDivider.Render(strings.Repeat("─", v.width-2))

	var visibleLines []string
	end := off + visibleH
	if end > len(lines) {
		end = len(lines)
	}
	slice := lines[off:end]

	for _, line := range slice {
		rendered := colorizeLogLine(line)
		visibleLines = append(visibleLines, "  "+rendered)
	}

	for len(visibleLines) < visibleH {
		visibleLines = append(visibleLines, "")
	}

	// Assemble
	result := []string{header, divider}
	result = append(result, visibleLines...)

	// Pad to exact height
	for len(result) < contentH {
		result = append(result, "")
	}
	if len(result) > contentH {
		result = result[:contentH]
	}
	return strings.Join(result, "\n")
}

// Footer returns the key hint for the log tab.
func (v LogView) Footer() string {
	return StyleFooter.Render("  [↑/↓] Scroll   [c] Clear log   [1-4] Switch tab")
}

// Update handles key input for the log view. Returns the updated view.
func (v LogView) Update(key string) LogView {
	lines := pkgs.GetLogLines()
	switch key {
	case "up", "k":
		if v.offset > 0 {
			v.offset--
		}
	case "down", "j":
		maxOff := len(lines) - (v.height - 10)
		if maxOff < 0 {
			maxOff = 0
		}
		if v.offset < maxOff {
			v.offset++
		}
	case "pgup":
		v.offset -= 10
		if v.offset < 0 {
			v.offset = 0
		}
	case "pgdown":
		v.offset += 10
		maxOff := len(lines) - (v.height - 10)
		if maxOff < 0 {
			maxOff = 0
		}
		if v.offset > maxOff {
			v.offset = maxOff
		}
	case "G":
		// Jump to bottom
		maxOff := len(lines) - (v.height - 10)
		if maxOff < 0 {
			maxOff = 0
		}
		v.offset = maxOff
	case "g":
		v.offset = 0
	case "c":
		pkgs.ClearLogs()
		v.offset = 0
	}
	return v
}

// colorizeLogLine applies terminal colors based on log level keywords.
func colorizeLogLine(line string) string {
	lower := strings.ToLower(line)
	switch {
	case strings.Contains(lower, "erro"):
		return lipgloss.NewStyle().Foreground(colorRed).Render(line)
	case strings.Contains(lower, "warn"):
		return lipgloss.NewStyle().Foreground(colorYellow).Render(line)
	case strings.Contains(lower, "info"):
		return lipgloss.NewStyle().Foreground(colorText).Render(line)
	case strings.Contains(lower, "debu"):
		return lipgloss.NewStyle().Foreground(colorTextMuted).Render(line)
	default:
		return lipgloss.NewStyle().Foreground(colorText).Render(line)
	}
}
