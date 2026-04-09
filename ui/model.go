package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inspiring-group/inspiring-swiss-knife/pkgs"
)

// Tab indices
const (
	TabInstall = iota
	TabTweaks
	TabLogs
	TabAbout
)

var tabNames = []string{
	"  Install Apps  ",
	"  Tweaks        ",
	"  Logs          ",
	"  About         ",
}

// AppModel is the root Bubble Tea model.
type AppModel struct {
	activeTab int
	install   InstallView
	tweaks    TweakView
	logs      LogView
	progress  ProgressView

	width  int
	height int

	// Warning/notification banner
	banner    string
	bannerErr bool

	isAdmin      bool
	wingetOK     bool
}

// New creates the root application model.
func New() AppModel {
	return AppModel{
		activeTab: TabInstall,
		install:   NewInstallView(120, 40),
		tweaks:    NewTweakView(120, 40),
		logs:      NewLogView(120, 40),
		progress:  NewProgressView(120, 40),
		isAdmin:   pkgs.IsAdmin(),
		wingetOK:  pkgs.WingetAvailable(),
	}
}

func (m AppModel) Init() tea.Cmd {
	// Kick off terminal size detection
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// ── Window resize ─────────────────────────────────────
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.install.width = msg.Width
		m.install.height = msg.Height
		m.tweaks.width = msg.Width
		m.tweaks.height = msg.Height
		m.logs.width = msg.Width
		m.logs.height = msg.Height
		m.progress.width = msg.Width
		m.progress.height = msg.Height
		return m, nil

	// ── Global key handling ───────────────────────────────
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "1":
			m.activeTab = TabInstall
			m.progress.Visible = false
			return m, nil
		case "2":
			m.activeTab = TabTweaks
			m.progress.Visible = false
			return m, nil
		case "3":
			m.activeTab = TabLogs
			m.progress.Visible = false
			return m, nil
		case "4":
			m.activeTab = TabAbout
			m.progress.Visible = false
			return m, nil

		case "esc":
			if m.progress.Visible {
				m.progress.Visible = false
			}
			return m, nil
		}

		// Route keys to active tab
		if m.progress.Visible {
			// No key handling during install
			return m, nil
		}

		switch m.activeTab {
		case TabInstall:
			updated, cmd, wantsInstall := m.install.Update(msg)
			m.install = updated
			if wantsInstall {
				return m.startInstall()
			}
			return m, cmd

		case TabTweaks:
			updated, cmd := m.tweaks.Update(msg)
			m.tweaks = updated
			if cmd != nil {
				return m, cmd
			}
			return m, nil

		case TabLogs:
			m.logs = m.logs.Update(msg.String())
			return m, nil
		}

	// ── Install progress messages ──────────────────────────
	case installTickChanMsg:
		m.progress.applyResult(msg.result)
		m.progress.done++
		return m, func() tea.Msg {
			result, ok := <-msg.ch
			if !ok {
				return installDoneMsg{}
			}
			return installTickChanMsg{ch: msg.ch, result: result}
		}

	case installDoneMsg:
		m.progress.MarkDone()
		m.banner = fmt.Sprintf("Installation complete: %d apps processed.", m.progress.total)
		m.bannerErr = false
		return m, nil

	// ── Tweak result messages ─────────────────────────────
	case TweakResultMsg:
		updated, cmd := m.tweaks.Update(msg)
		m.tweaks = updated
		return m, cmd

	// ── Tweak run batch ───────────────────────────────────
	case TweakRunMsg:
		name := msg.Name
		if name == "" {
			name = "Tweaks"
		}
		return m, func() tea.Msg {
			pkgs.Logger.Info("Applying tweaks", "batch", name, "count", len(msg.Scripts))
			combined := strings.Join(msg.Scripts, "; ")
			out, err := pkgs.RunPowerShell(combined)
			if err != nil {
				pkgs.Logger.Error("Tweak batch failed", "batch", name, "err", err)
			} else {
				pkgs.Logger.Info("Tweak batch complete", "batch", name)
			}
			return TweakResultMsg{Output: out, Err: err}
		}
	}

	return m, nil
}

func (m AppModel) startInstall() (AppModel, tea.Cmd) {
	apps := m.install.SelectedApps()
	if len(apps) == 0 {
		m.banner = "No apps selected. Use Space to select apps."
		m.bannerErr = true
		return m, nil
	}
	if !m.wingetOK {
		m.banner = "winget not found! Please install App Installer from the Microsoft Store."
		m.bannerErr = true
		return m, nil
	}
	m.progress = m.progress.Init(apps)
	m.activeTab = TabInstall // keep on install tab, overlay progress

	// Start the install goroutine and wire results back
	resultCh := make(chan pkgs.InstallResult, len(apps))
	go pkgs.InstallApps(apps, resultCh)
	return m, func() tea.Msg {
		result, ok := <-resultCh
		if !ok {
			return installDoneMsg{}
		}
		return installTickChanMsg{ch: resultCh, result: result}
	}
}

// installTickChanMsg carries a result and the channel to keep reading from.
type installTickChanMsg struct {
	ch     chan pkgs.InstallResult
	result pkgs.InstallResult
}

// installDoneMsg is sent when all installs are complete.
type installDoneMsg struct{}

func (m AppModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// ── Fixed sections ────────────────────────────────────
	header := m.renderHeader() // 3 lines: title+status / tabs / divider
	headerH := lipgloss.Height(header)

	var topExtras []string
	if warn := m.systemWarnings(); warn != "" {
		topExtras = append(topExtras, warn)
	}
	if m.banner != "" {
		topExtras = append(topExtras, m.renderBanner())
	}

	// ── Footer (always 1 line, pinned to bottom) ──────────
	var footer string
	switch {
	case m.progress.Visible:
		footer = StyleFooter.Render("  [Esc] Close overlay   Installation running…")
	case m.activeTab == TabInstall:
		footer = m.install.Footer()
	case m.activeTab == TabTweaks:
		footer = m.tweaks.Footer()
	case m.activeTab == TabLogs:
		footer = m.logs.Footer()
	default:
		footer = StyleFooter.Render("  [1] Install Apps  [2] Tweaks  [3] Logs  [4] About  [q] Quit")
	}

	// ── Calculate content area height ─────────────────────
	usedLines := headerH + len(topExtras) + 1 // +1 for footer
	contentH := m.height - usedLines
	if contentH < 3 {
		contentH = 3
	}

	// ── Render content ────────────────────────────────────
	var content string
	if m.progress.Visible {
		content = m.progress.View()
	} else {
		switch m.activeTab {
		case TabInstall:
			content = m.install.View(contentH)
		case TabTweaks:
			content = m.tweaks.View(contentH)
		case TabLogs:
			content = m.logs.View(contentH)
		case TabAbout:
			content = aboutView(m.width, contentH)
		}
	}

	// Constrain content to exactly contentH lines
	contentLines := strings.Split(content, "\n")
	for len(contentLines) < contentH {
		contentLines = append(contentLines, "")
	}
	if len(contentLines) > contentH {
		contentLines = contentLines[:contentH]
	}

	// ── Assemble ──────────────────────────────────────────
	parts := []string{header}
	parts = append(parts, topExtras...)
	parts = append(parts, strings.Join(contentLines, "\n"))
	parts = append(parts, footer)
	return strings.Join(parts, "\n")
}

func (m AppModel) renderHeader() string {
	// Logo / title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorCyan).
		Render("⚔  Inspiring Swiss Knife")

	subtitle := StyleSubtitle.Render("  Windows Onboarding & Optimization Tool")

	// Status pills
	var pills []string
	if m.isAdmin {
		pills = append(pills, StyleSuccess.Render("● Admin"))
	} else {
		pills = append(pills, StyleFail.Render("● No Admin"))
	}
	if m.wingetOK {
		pills = append(pills, StyleSuccess.Render("● WinGet OK"))
	} else {
		pills = append(pills, StyleFail.Render("● WinGet Missing"))
	}

	statusLine := strings.Join(pills, "  ")

	// Tab bar
	var tabs []string
	for i, name := range tabNames {
		label := fmt.Sprintf("[%d]%s", i+1, name)
		if i == m.activeTab {
			tabs = append(tabs, StyleTabActive.Render(label))
		} else {
			tabs = append(tabs, StyleTabInactive.Render(label))
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	rightPad := m.width - lipgloss.Width(title+subtitle) - lipgloss.Width(statusLine) - 4
	if rightPad < 0 {
		rightPad = 0
	}

	topLine := title + subtitle + strings.Repeat(" ", rightPad) + statusLine
	divider := StyleDivider.Render(strings.Repeat("─", m.width))

	return topLine + "\n" + tabBar + "\n" + divider
}

func (m AppModel) systemWarnings() string {
	if !m.isAdmin {
		return StyleWarn.Render("  ⚠  Running without Administrator privileges — tweaks and some installs may fail. Restart as Admin.")
	}
	return ""
}

func (m AppModel) renderBanner() string {
	style := StyleInfo
	if m.bannerErr {
		style = StyleFail
	}
	return style.Render("  " + m.banner)
}

// iskLogo is the ISK ASCII art logo.
const iskLogo = ` ___  ________  ___  __
|\  \|\   ____\|\  \|\  \
\ \  \ \  \___|\ \  \/  /|_
 \ \  \ \_____  \ \   ___  \
  \ \  \|____|\  \ \  \\ \  \
   \ \__\____\_\  \ \__\\ \__\
    \|__|\_________\|__| \|__|
        \|_________|`

// kb renders a keyboard key as a coloured badge: [key]
func kb(key string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(colorBg).
		Background(colorBlue).
		Padding(0, 1).
		Render(key)
}

// kbRow renders a row of key→action pairs separated by spaces.
func kbRow(pairs ...string) string {
	// pairs: key, action, key, action, ...
	var parts []string
	for i := 0; i+1 < len(pairs); i += 2 {
		parts = append(parts, kb(pairs[i])+" "+
			lipgloss.NewStyle().Foreground(colorText).Render(pairs[i+1]))
	}
	return strings.Join(parts, "   ")
}

func aboutView(width, contentH int) string {
	w := width - 4
	if w < 60 {
		w = 60
	}
	// -2 for the border lines themselves
	boxH := contentH - 2
	if boxH < 10 {
		boxH = 10
	}

	logoRendered := lipgloss.NewStyle().
		Foreground(colorBlue).
		Bold(true).
		Render(iskLogo)

	titleRendered := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorCyan).
		Render("INSPIRING SWISS KNIFE")

	shortcutsHeader := StyleSectionHeader.Render("Keyboard Shortcuts")
	shortcuts := kbRow("1", "Install Apps", "2", "Tweaks", "3", "Logs", "4", "About", "q", "Quit") + "\n" +
		kbRow("←/→", "Category", "↑/↓", "Navigate", "Space", "Select", "Enter", "Install") + "\n" +
		kbRow("Tab", "Switch panel", "F5", "Apply tweaks", "c", "Clear log", "a", "Select all")

	techHeader := StyleSectionHeader.Render("Technology")
	techLine := StyleTextMuted("Go + Bubble Tea · WinGet · PowerShell · charmbracelet/log")

	content := logoRendered + "\n\n" +
		titleRendered + "\n" +
		StyleSubtitle.Render("Windows Onboarding & Optimization Tool") + "\n\n" +
		StyleAppNormal.Render("A comprehensive tool for setting up new employee laptops at Inspiring Group.\n"+
			"Install apps via WinGet, apply Windows tweaks, and configure system preferences\n"+
			"— all from a single, portable terminal UI.") + "\n\n" +
		shortcutsHeader + "\n" +
		shortcuts + "\n\n" +
		techHeader + "\n" +
		techLine + "\n\n" +
		StyleSuccess.Render("github.com/taifme/inspiring-swiss-army")

	box := lipgloss.NewStyle().
		Width(w).
		Height(boxH).
		Padding(1, 3).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder)

	return box.Render(content)
}
