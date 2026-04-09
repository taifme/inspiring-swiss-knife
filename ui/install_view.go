package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inspiring-group/inspiring-swiss-knife/pkgs"
)

// InstallView manages the app installation tab.
type InstallView struct {
	// Category navigation
	categories []pkgs.Category
	activeCat  int

	// Per-category app lists
	catalog  map[pkgs.Category][]pkgs.App
	selected map[string]bool // keyed by WingetID

	// Cursor within current category
	cursor int
	offset int // scroll offset

	// Viewport dimensions (set by parent)
	width  int
	height int

	// Tooltip shown at bottom
	hoveredNote string
}

func NewInstallView(w, h int) InstallView {
	catalog := make(map[pkgs.Category][]pkgs.App)
	for _, cat := range pkgs.AllCategories {
		catalog[cat] = pkgs.ByCategory(cat)
	}

	return InstallView{
		categories: pkgs.AllCategories,
		activeCat:  0,
		catalog:    catalog,
		selected:   make(map[string]bool),
		width:      w,
		height:     h,
	}
}

// currentApps returns the visible app list for the active category.
func (v *InstallView) currentApps() []pkgs.App {
	return v.catalog[v.categories[v.activeCat]]
}

// SelectedApps returns all apps that have been checked.
func (v *InstallView) SelectedApps() []pkgs.App {
	var result []pkgs.App
	for _, apps := range v.catalog {
		for _, a := range apps {
			if v.selected[a.WingetID] {
				result = append(result, a)
			}
		}
	}
	return result
}

// SelectedCount returns the total number of selected apps.
func (v *InstallView) SelectedCount() int {
	n := 0
	for _, apps := range v.catalog {
		for _, a := range apps {
			if v.selected[a.WingetID] {
				n++
			}
		}
	}
	return n
}

// Update handles key presses for the install view.
// Returns (updated view, cmd, wantsInstall).
func (v InstallView) Update(msg tea.Msg) (InstallView, tea.Cmd, bool) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		apps := v.currentApps()
		maxCursor := len(apps) - 1

		switch m.String() {
		// ── Category navigation ─────────────────────────────
		case "left", "h":
			if v.activeCat > 0 {
				v.activeCat--
				v.cursor = 0
				v.offset = 0
			}
		case "right", "l":
			if v.activeCat < len(v.categories)-1 {
				v.activeCat++
				v.cursor = 0
				v.offset = 0
			}

		// ── App list navigation ─────────────────────────────
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
				if v.cursor < v.offset {
					v.offset = v.cursor
				}
			}
		case "down", "j":
			if v.cursor < maxCursor {
				v.cursor++
				visibleLines := v.approxListHeight()
				if v.cursor >= v.offset+visibleLines {
					v.offset = v.cursor - visibleLines + 1
				}
			}
		case "pgup":
			v.cursor -= 10
			if v.cursor < 0 {
				v.cursor = 0
			}
			v.offset = v.cursor
		case "pgdown":
			v.cursor += 10
			if v.cursor > maxCursor {
				v.cursor = maxCursor
			}
			visibleLines := v.approxListHeight()
			if v.cursor >= v.offset+visibleLines {
				v.offset = v.cursor - visibleLines + 1
			}

		// ── Selection ────────────────────────────────────────
		case " ":
			if len(apps) > 0 && v.cursor <= maxCursor {
				id := apps[v.cursor].WingetID
				v.selected[id] = !v.selected[id]
			}
		case "a":
			// Select all in current category
			allSelected := true
			for _, a := range apps {
				if !v.selected[a.WingetID] {
					allSelected = false
					break
				}
			}
			for _, a := range apps {
				v.selected[a.WingetID] = !allSelected
			}

		// ── Install trigger ──────────────────────────────────
		case "enter":
			if v.SelectedCount() > 0 {
				return v, nil, true
			}
		}

		// Update hovered note
		apps = v.currentApps()
		if v.cursor < len(apps) {
			v.hoveredNote = apps[v.cursor].Note
		}
	}

	return v, nil, false
}


// View renders the install tab content into exactly contentH lines (no footer).
func (v InstallView) View(contentH int) string {
	apps := v.currentApps()

	// ── 3-line fixed header ───────────────────────────────────────────────────
	tabBar := v.renderCategoryTabs()
	divider := StyleDivider.Render(strings.Repeat("─", v.width))

	// Active category label shown BELOW the divider (breadcrumb)
	catLabel := "  " + StyleCatActive.Render("▶  "+string(v.categories[v.activeCat])) +
		"  " + StyleTextMuted(fmt.Sprintf("(%d apps)", len(apps)))
	if v.offset > 0 {
		catLabel += StyleTextMuted("   ↑ scroll")
	}
	listH := contentH - 3 // tabBar + divider + catLabel
	if listH < 1 {
		listH = 1
	}
	end := v.offset + listH
	if end > len(apps) {
		end = len(apps)
	}
	if end < len(apps) {
		catLabel += StyleTextMuted("   ↓ more")
	}

	fixed := []string{tabBar, divider, catLabel}

	// ── App list ──────────────────────────────────────────────────────────────
	appLines := strings.Split(v.renderAppList(listH), "\n")

	lines := append(fixed, appLines...)

	// Pad or trim to exact height
	for len(lines) < contentH {
		lines = append(lines, "")
	}
	if len(lines) > contentH {
		lines = lines[:contentH]
	}
	return strings.Join(lines, "\n")
}

// Footer returns the key-hint line rendered by the root model.
func (v InstallView) Footer() string {
	count := v.SelectedCount()
	hint := "[←/→] Category  [↑/↓] Navigate  [Space] Select  [a] All  [Enter] Install"
	var extras []string
	if count > 0 {
		extras = append(extras, StyleCount.Render(fmt.Sprintf("  %d selected", count)))
	}
	if v.hoveredNote != "" {
		extras = append(extras, StyleTextMuted("  ℹ "+v.hoveredNote))
	}
	return StyleFooter.Render(hint+strings.Join(extras, ""))
}

func (v InstallView) renderCategoryTabs() string {
	// Pre-render all tab labels and measure widths.
	type tabEntry struct {
		rendered string
		w        int
	}
	entries := make([]tabEntry, len(v.categories))
	for i, cat := range v.categories {
		label := string(cat)
		var r string
		if i == v.activeCat {
			r = StyleCatActive.Render("▶ " + label)
		} else {
			r = StyleCatInactive.Render("  " + label)
		}
		entries[i] = tabEntry{rendered: r, w: lipgloss.Width(r)}
	}

	// Find a contiguous window of tabs that contains activeCat and fits
	// within the available terminal width.
	maxW := v.width - 6 // leave room for ◀/▶ arrows
	if maxW < 20 {
		maxW = 20
	}
	start, end := v.activeCat, v.activeCat+1
	used := entries[v.activeCat].w
	for {
		grew := false
		if start > 0 && used+entries[start-1].w <= maxW {
			start--
			used += entries[start].w
			grew = true
		}
		if end < len(entries) && used+entries[end].w <= maxW {
			used += entries[end].w
			end++
			grew = true
		}
		if !grew {
			break
		}
	}

	var parts []string
	if start > 0 {
		parts = append(parts, StyleTextMuted("◀ "))
	}
	for i := start; i < end; i++ {
		parts = append(parts, entries[i].rendered)
	}
	if end < len(entries) {
		parts = append(parts, StyleTextMuted(" ▶"))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (v InstallView) renderAppList(maxLines int) string {
	apps := v.currentApps()
	if len(apps) == 0 {
		return StyleTextMuted("  No apps in this category.")
	}

	end := v.offset + maxLines
	if end > len(apps) {
		end = len(apps)
	}
	visible := apps[v.offset:end]

	var rows []string
	for i, app := range visible {
		rows = append(rows, v.renderApp(app, v.offset+i))
	}
	return strings.Join(rows, "\n")
}

// renderApp renders one app as a full-width vertical list row.
func (v InstallView) renderApp(app pkgs.App, idx int) string {
	isSelected := v.selected[app.WingetID]
	isCursor := idx == v.cursor

	rowW := v.width
	if rowW < 40 {
		rowW = 40
	}

	// ── Cursor & checkbox ─────────────────────────────────────────────────────
	prefix := "  "
	if isCursor {
		prefix = StyleInfo.Render("▶ ")
	}
	checkbox := CheckboxStr(isSelected)

	// ── Name (fixed width) ────────────────────────────────────────────────────
	nameW := 28
	if rowW > 120 {
		nameW = 34
	}
	var nameStyle lipgloss.Style
	switch {
	case isCursor:
		nameStyle = StyleAppCursor
	case isSelected:
		nameStyle = StyleAppSelected
	default:
		nameStyle = StyleAppNormal
	}
	nameTxt := nameStyle.Render(truncate(app.Name, nameW))
	// Pad name to fixed column width (visual chars)
	nameVisW := lipgloss.Width(nameTxt)
	if nameVisW < nameW {
		nameTxt += strings.Repeat(" ", nameW-nameVisW)
	}

	// ── Source badge ──────────────────────────────────────────────────────────
	var badge string
	switch app.Source {
	case pkgs.SourceManual:
		badge = lipgloss.NewStyle().Foreground(colorYellow).Render(" manual ")
	case pkgs.SourceMSStore:
		badge = lipgloss.NewStyle().Foreground(colorPurple).Render(" store  ")
	case pkgs.SourceLinux:
		badge = lipgloss.NewStyle().Foreground(colorTextDim).Render(" linux  ")
	case pkgs.SourceNPM:
		badge = lipgloss.NewStyle().Foreground(colorGreen).Render(" npm    ")
	default:
		badge = lipgloss.NewStyle().Foreground(colorTextDim).Render(" winget ")
	}
	badgeW := lipgloss.Width(badge)

	// ── Note (fills remaining width) ─────────────────────────────────────────
	// prefix(2) + checkbox(3) + space(1) + name(nameW) + badge(badgeW) + padding(2)
	usedW := 2 + 3 + 1 + nameW + badgeW + 2
	noteMaxW := rowW - usedW
	var note string
	if noteMaxW > 8 && app.Note != "" {
		note = "  " + StyleTextMuted(truncate(app.Note, noteMaxW))
		// Pad note to fill exact remaining space
		noteVisW := lipgloss.Width(note)
		padNeeded := noteMaxW - (noteVisW - 2) // subtract the "  " we added
		if padNeeded > 0 {
			note += strings.Repeat(" ", padNeeded)
		}
	}

	// ── Highlight cursor row with a subtle background ─────────────────────────
	row := fmt.Sprintf("%s%s %s%s%s", prefix, checkbox, nameTxt, note, badge)
	if isCursor {
		row = lipgloss.NewStyle().
			Background(colorPanelBg).
			Render(row)
	}

	// Ensure full width
	rowVisW := lipgloss.Width(row)
	if rowVisW < rowW {
		row += strings.Repeat(" ", rowW-rowVisW)
	}
	return row
}


// approxListHeight gives a rough estimate used during keyboard navigation.
// Accounts for: app-header(~5) + tabBar(1) + divider(1) + catLabel(1) + footer(1)
func (v *InstallView) approxListHeight() int {
	h := v.height - 9
	if h < 5 {
		h = 5
	}
	return h
}

// Helpers

func StyleTextMuted(s string) string {
	return lipgloss.NewStyle().Foreground(colorTextMuted).Render(s)
}

func truncate(s string, max int) string {
	if max <= 3 {
		return s
	}
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
