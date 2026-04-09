package ui

import "github.com/charmbracelet/lipgloss"

// Color palette — mirrors the WinUtil dark terminal aesthetic
var (
	colorBg         = lipgloss.Color("#0D1117")
	colorPanelBg    = lipgloss.Color("#161B22")
	colorBorder     = lipgloss.Color("#30363D")
	colorBorderHi   = lipgloss.Color("#58A6FF")

	colorText       = lipgloss.Color("#E6EDF3")
	colorTextMuted  = lipgloss.Color("#8B949E")
	colorTextDim    = lipgloss.Color("#484F58")

	colorCyan       = lipgloss.Color("#39D3F5") // category headers
	colorGreen      = lipgloss.Color("#3FB950") // selected / success
	colorBlue       = lipgloss.Color("#58A6FF") // active tab / buttons
	colorYellow     = lipgloss.Color("#D29922") // warning
	colorRed        = lipgloss.Color("#FF7B72") // error / danger
	colorPurple     = lipgloss.Color("#BC8CFF") // accent

	colorToggleOn   = lipgloss.Color("#238636")
	colorToggleOnFg = lipgloss.Color("#3FB950")
)

// ── Tab bar ──────────────────────────────────────────────────────────────────

var (
	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBlue).
			Background(colorPanelBg).
			Padding(0, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(colorBlue)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(colorTextMuted).
				Padding(0, 2)

	StyleTabBar = lipgloss.NewStyle().
			Background(colorBg).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(colorBorder)
)

// ── Panels ────────────────────────────────────────────────────────────────────

var (
	StylePanel = lipgloss.NewStyle().
			Background(colorPanelBg).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	StylePanelFocused = StylePanel.
				BorderForeground(colorBorderHi)

	StyleSectionHeader = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorCyan).
				MarginBottom(1)
)

// ── Checkboxes ───────────────────────────────────────────────────────────────

var (
	StyleCheckOn = lipgloss.NewStyle().
			Foreground(colorGreen).Bold(true)

	StyleCheckOff = lipgloss.NewStyle().
			Foreground(colorTextMuted)

	StyleAppSelected = lipgloss.NewStyle().
				Foreground(colorGreen)

	StyleAppNormal = lipgloss.NewStyle().
			Foreground(colorText)

	StyleAppCursor = lipgloss.NewStyle().
			Foreground(colorBlue).Bold(true)

	StyleAppNote = lipgloss.NewStyle().
			Foreground(colorTextDim).Italic(true)
)

// ── Toggles ──────────────────────────────────────────────────────────────────

var (
	StyleToggleOn = lipgloss.NewStyle().
			Foreground(colorToggleOnFg).
			Background(colorToggleOn).
			Padding(0, 1).
			Bold(true)

	StyleToggleOff = lipgloss.NewStyle().
			Foreground(colorTextMuted).
			Background(colorBorder).
			Padding(0, 1)

	StyleToggleLabel = lipgloss.NewStyle().Foreground(colorText)
)

// ── Buttons ───────────────────────────────────────────────────────────────────

var (
	StyleButton = lipgloss.NewStyle().
			Foreground(colorText).
			Background(colorBlue).
			Padding(0, 2).
			Bold(true)

	StyleButtonDanger = lipgloss.NewStyle().
				Foreground(colorText).
				Background(colorRed).
				Padding(0, 2).
				Bold(true)

	StyleButtonSuccess = lipgloss.NewStyle().
				Foreground(colorText).
				Background(lipgloss.Color("#238636")).
				Padding(0, 2).
				Bold(true)
)

// ── Status icons ─────────────────────────────────────────────────────────────

var (
	StyleSuccess = lipgloss.NewStyle().Foreground(colorGreen)
	StyleFail    = lipgloss.NewStyle().Foreground(colorRed)
	StyleWarn    = lipgloss.NewStyle().Foreground(colorYellow)
	StyleInfo    = lipgloss.NewStyle().Foreground(colorBlue)
)

// ── Category tabs (install view) ──────────────────────────────────────────────

var (
	StyleCatActive = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true).
			Padding(0, 1).
			Background(colorPanelBg)

	StyleCatInactive = lipgloss.NewStyle().
				Foreground(colorTextMuted).
				Padding(0, 1)
)

// ── Misc ──────────────────────────────────────────────────────────────────────

var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCyan).
			MarginBottom(1)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(colorTextMuted).
			Italic(true)

	StyleFooter = lipgloss.NewStyle().
			Foreground(colorTextDim).
			MarginTop(1)

	StyleCount = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	StyleDivider = lipgloss.NewStyle().
			Foreground(colorBorder)

	StyleProgressBar = lipgloss.NewStyle().
				Foreground(colorGreen)
)

// Helpers

func CheckboxStr(on bool) string {
	if on {
		return StyleCheckOn.Render("[✓]")
	}
	return StyleCheckOff.Render("[ ]")
}

func ToggleStr(on bool) string {
	if on {
		return StyleToggleOn.Render(" ON ")
	}
	return StyleToggleOff.Render(" OFF")
}

func RenderProgressBar(percent float64, width int) string {
	filled := int(float64(width) * percent)
	if filled > width {
		filled = width
	}
	bar := lipgloss.NewStyle().Foreground(colorGreen).Render(repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(colorBorder).Render(repeat("░", width-filled))
	return bar
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	result := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		result = append(result, s...)
	}
	return string(result)
}
