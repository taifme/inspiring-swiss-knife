package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/inspiring-group/inspiring-swiss-knife/pkgs"
)

// ProgressView shows the installation progress overlay.
type ProgressView struct {
	items   []progressItem
	done    int
	total   int
	log     []string
	width   int
	height  int
	Visible bool
}

type progressItem struct {
	app    pkgs.App
	status pkgs.InstallStatus
	msg    string
}

func NewProgressView(w, h int) ProgressView {
	return ProgressView{width: w, height: h}
}

func (p ProgressView) Init(apps []pkgs.App) ProgressView {
	items := make([]progressItem, len(apps))
	for i, a := range apps {
		items[i] = progressItem{app: a, status: pkgs.StatusPending}
	}
	p.items = items
	p.total = len(apps)
	p.done = 0
	p.Visible = true
	p.log = nil

	return p
}

func (p *ProgressView) MarkDone() {
	p.log = append(p.log, StyleSuccess.Render(fmt.Sprintf("✓ All %d installs completed.", p.total)))
}

func (p *ProgressView) applyResult(r pkgs.InstallResult) {
	for i, item := range p.items {
		if item.app.WingetID == r.App.WingetID {
			p.items[i].status = r.Status
			switch r.Status {
			case pkgs.StatusSuccess:
				p.items[i].msg = "Installed"
				p.log = append(p.log, StyleSuccess.Render("✓ "+r.App.Name))
			case pkgs.StatusSkipped:
				p.items[i].msg = "Already installed"
				p.log = append(p.log, StyleWarn.Render("~ "+r.App.Name+" (already installed)"))
			case pkgs.StatusFailed:
				p.items[i].msg = "FAILED: " + r.Error
				p.log = append(p.log, StyleFail.Render("✗ "+r.App.Name+": "+r.Error))
			}
			return
		}
	}
}

func (p ProgressView) View() string {
	if !p.Visible {
		return ""
	}

	var sb strings.Builder

	// Title
	sb.WriteString(StyleTitle.Render("Installing Applications") + "\n\n")

	// Overall progress bar
	var pct float64
	if p.total > 0 {
		pct = float64(p.done) / float64(p.total)
	}
	barWidth := p.width - 20
	if barWidth < 20 {
		barWidth = 20
	}
	sb.WriteString(fmt.Sprintf("  %s  %d / %d\n\n",
		RenderProgressBar(pct, barWidth),
		p.done, p.total))

	// App status list (show up to 20)
	shown := p.items
	if len(shown) > 20 {
		shown = p.items[:20]
	}
	cols := 2
	colW := (p.width - 6) / cols

	for rowStart := 0; rowStart < len(shown); rowStart += cols {
		var cells []string
		for c := 0; c < cols; c++ {
			idx := rowStart + c
			if idx >= len(shown) {
				cells = append(cells, strings.Repeat(" ", colW))
				continue
			}
			item := shown[idx]
			cells = append(cells, renderProgressItem(item, colW))
		}
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cells...) + "\n")
	}

	// Log tail
	sb.WriteString("\n")
	logStart := len(p.log) - 6
	if logStart < 0 {
		logStart = 0
	}
	for _, line := range p.log[logStart:] {
		sb.WriteString("  " + line + "\n")
	}

	sb.WriteString("\n" + StyleFooter.Render("  Installation running... please wait"))

	return sb.String()
}

func renderProgressItem(item progressItem, width int) string {
	var icon string
	var style lipgloss.Style

	switch item.status {
	case pkgs.StatusPending:
		icon = "○"
		style = StyleSubtitle
	case pkgs.StatusRunning:
		icon = "◉"
		style = StyleInfo
	case pkgs.StatusSuccess:
		icon = "●"
		style = StyleSuccess
	case pkgs.StatusFailed:
		icon = "✗"
		style = StyleFail
	case pkgs.StatusSkipped:
		icon = "~"
		style = StyleWarn
	}

	label := truncate(item.app.Name, width-4)
	cell := fmt.Sprintf("  %s %s", style.Render(icon), StyleAppNormal.Render(label))
	padded := cell + strings.Repeat(" ", max(0, width-lipgloss.Width(cell)))
	return padded
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
