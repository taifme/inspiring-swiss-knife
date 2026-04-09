package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inspiring-group/inspiring-swiss-knife/pkgs"
)

// ─── Helpers ─────────────────────────────────────────────────────────────────

// regSet generates a PowerShell Set-ItemProperty command.
func regSet(path, name, value, regType string) string {
	return fmt.Sprintf(
		`$null = New-Item -Path '%s' -Force -ErrorAction SilentlyContinue; Set-ItemProperty -Path '%s' -Name '%s' -Value %s -Type %s -Force`,
		path, path, name, value, regType,
	)
}

// joinPS joins multiple PS commands into one semicolon-separated line.
func joinPS(cmds ...string) string { return strings.Join(cmds, "; ") }

// ─── Data types ───────────────────────────────────────────────────────────────

type TweakCategory int

const (
	TweakEssential TweakCategory = iota
	TweakAdvanced
)

type Tweak struct {
	Name   string
	Desc   string
	Script string // PowerShell
}

type Toggle struct {
	Name          string
	Desc          string
	DefaultOn     bool
	Enabled       bool
	EnableScript  string // PS when toggled ON
	DisableScript string // PS when toggled OFF
}

// ─── Essential Tweaks (WinUtil official data) ─────────────────────────────────

var EssentialTweaks = []Tweak{
	{
		Name: "Create Restore Point",
		Desc: "Creates a system restore checkpoint before applying changes.",
		Script: `Enable-ComputerRestore -Drive "$env:SystemDrive"; Checkpoint-Computer -Description "InspiringSwissKnife" -RestorePointType "MODIFY_SETTINGS"`,
	},
	{
		Name: "Delete Temporary Files",
		Desc: "Erases TEMP folders to free disk space.",
		Script: joinPS(
			`Remove-Item -Path "$Env:Temp\*" -Recurse -Force -ErrorAction SilentlyContinue`,
			`Remove-Item -Path "$Env:SystemRoot\Temp\*" -Recurse -Force -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Disable Activity History",
		Desc: "Erases recent docs, clipboard, and run history. Stops Windows from tracking activities.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\System`, "EnableActivityFeed", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\System`, "PublishUserActivities", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\System`, "UploadUserActivities", "0", "DWord"),
		),
	},
	{
		Name: "Disable Consumer Features",
		Desc: "Windows won't auto-install games or third-party apps from the Store.",
		Script: regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\CloudContent`, "DisableWindowsConsumerFeatures", "1", "DWord"),
	},
	{
		Name: "Disable Explorer Auto Folder Discovery",
		Desc: "Stops Explorer from guessing folder types and applying unwanted grouping.",
		Script: joinPS(
			`Remove-Item -Path 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\Bags' -Recurse -Force -ErrorAction SilentlyContinue`,
			`Remove-Item -Path 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\BagMRU' -Recurse -Force -ErrorAction SilentlyContinue`,
			`$null = New-Item -Path 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\Bags\AllFolders\Shell' -Force`,
			`Set-ItemProperty -Path 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\Bags\AllFolders\Shell' -Name 'FolderType' -Value 'NotSpecified' -Force`,
		),
	},
	{
		Name: "Disable Hibernation",
		Desc: "Removes hiberfil.sys and disables hibernate. Best for desktops to save disk space.",
		Script: joinPS(
			regSet(`HKLM:\System\CurrentControlSet\Control\Session Manager\Power`, "HibernateEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\FlyoutMenuSettings`, "ShowHibernateOption", "0", "DWord"),
			"powercfg.exe /hibernate off",
		),
	},
	{
		Name: "Disable Location Tracking",
		Desc: "Disables location services and GPS features. Denies app location access.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\location`, "Value", "'Deny'", "String"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Sensor\Overrides\{BFA794E4-F964-4FDB-90F6-51056BFE4B44}`, "SensorPermissionState", "0", "DWord"),
			regSet(`HKLM:\SYSTEM\Maps`, "AutoUpdateEnabled", "0", "DWord"),
			`Set-Service -Name lfsvc -StartupType Disabled -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Disable Microsoft Store Search Results",
		Desc: "Hides Store app recommendations from Start menu search results.",
		Script: `icacls "$Env:LocalAppData\Packages\Microsoft.WindowsStore_8wekyb3d8bbwe\LocalState\store.db" /deny Everyone:F 2>$null; ` +
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\Explorer`, "NoUseStoreOpenWith", "1", "DWord"),
	},
	{
		Name: "Disable PowerShell 7 Telemetry",
		Desc: "Sets POWERSHELL_TELEMETRY_OPTOUT env variable to opt out of PS7 data collection.",
		Script: `[Environment]::SetEnvironmentVariable('POWERSHELL_TELEMETRY_OPTOUT', '1', 'Machine')`,
	},
	{
		Name: "Disable Telemetry",
		Desc: "Disables Microsoft telemetry via registry, services, and advertising ID. (Source: WinUtil WPFTweaksTelemetry)",
		Script: joinPS(
			regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\AdvertisingInfo`, "Enabled", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Privacy`, "TailoredExperiencesWithDiagnosticDataEnabled", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Speech_OneCore\Settings\OnlineSpeechPrivacy`, "HasAccepted", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Input\TIPC`, "Enabled", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\InputPersonalization`, "RestrictImplicitInkCollection", "1", "DWord"),
			regSet(`HKCU:\Software\Microsoft\InputPersonalization`, "RestrictImplicitTextCollection", "1", "DWord"),
			regSet(`HKCU:\Software\Microsoft\InputPersonalization\TrainedDataStore`, "HarvestContacts", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Personalization\Settings`, "AcceptedPrivacyPolicy", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\DataCollection`, "AllowTelemetry", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Start_TrackProgs", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Siuf\Rules`, "NumberOfSIUFInPeriod", "0", "DWord"),
			`Set-MpPreference -SubmitSamplesConsent 2 -ErrorAction SilentlyContinue`,
			`Set-Service -Name diagtrack -StartupType Disabled -ErrorAction SilentlyContinue`,
			`Set-Service -Name wermgr -StartupType Disabled -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Disable Windows Platform Binary Table (WPBT)",
		Desc: "Prevents OEM/vendor software from auto-executing at boot via BIOS/UEFI.",
		Script: regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager`, "DisableWpbtExecution", "1", "DWord"),
	},
	{
		Name: "Enable End Task With Right Click",
		Desc: "Adds 'End Task' to the taskbar right-click context menu for quick process kill.",
		Script: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced\TaskbarDeveloperSettings`, "TaskbarEndTask", "1", "DWord"),
	},
	{
		Name: "Remove Widgets",
		Desc: "Removes the Widgets panel and WebExperience package from the taskbar.",
		Script: joinPS(
			`Get-Process *Widget* -ErrorAction SilentlyContinue | Stop-Process -Force`,
			`Get-AppxPackage -AllUsers Microsoft.WidgetsPlatformRuntime -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers`,
			`Get-AppxPackage -AllUsers MicrosoftWindows.Client.WebExperience -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers`,
		),
	},
	{
		Name: "Revert Start Menu Layout",
		Desc: "Restores the default Start Menu layout by re-importing the current export.",
		Script: `$f = [IO.Path]::GetTempFileName() + '.xml'; Export-StartLayout -Path $f; Import-StartLayout -LayoutPath $f -MountPath $env:SystemDrive\`,
	},
	{
		Name: "Run Disk Cleanup",
		Desc: "Runs Windows built-in Disk Cleanup to remove old update files and junk.",
		Script: `Start-Process cleanmgr -ArgumentList '/sagerun:1' -Wait`,
	},
	{
		Name: "Set Services to Manual",
		Desc: "Sets non-essential background services to Manual startup. They start on-demand.",
		Script: joinPS(
			`$svcs = @('ALG','AppMgmt','AppReadiness','CDPSvc','PcaSvc','InventorySvc','StorSvc','UsoSvc','WpnService')`,
			`foreach ($s in $svcs) { Set-Service -Name $s -StartupType Manual -ErrorAction SilentlyContinue }`,
			`Set-Service -Name DiagTrack -StartupType Disabled -ErrorAction SilentlyContinue`,
		),
	},
}

// ─── Advanced Tweaks (WinUtil official data) ──────────────────────────────────

var AdvancedTweaks = []Tweak{
	{
		Name: "Adobe Network Block",
		Desc: "Blocks Adobe telemetry and licensing servers via hosts file (CAUTION: may affect Creative Cloud).",
		Script: "Add-Content -Path \"$env:SystemRoot\\System32\\drivers\\etc\\hosts\" -Value \"`n127.0.0.1 lmlicenses.wip4.adobe.com`n127.0.0.1 lm.licenses.adobe.com`n127.0.0.1 na2.adobebettle.com`n127.0.0.1 apiu.adobe.com\"",
	},
	{
		Name: "Block Razer Software Installs",
		Desc: "Prevents Razer Synapse from auto-installing via Windows Update driver delivery.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\DeviceInstall\Restrictions`, "DenyDeviceIDs", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\DeviceInstall\Restrictions`, "DenyDeviceIDsRetroactive", "0", "DWord"),
		),
	},
	{
		Name: "Disable Background Apps",
		Desc: "Prevents Microsoft Store UWP apps from running and updating in background.",
		Script: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\BackgroundAccessApplications`, "GlobalUserDisabled", "1", "DWord"),
	},
	{
		Name: "Disable Fullscreen Optimizations",
		Desc: "Removes FSO layering in exclusive fullscreen games. Can improve frame pacing.",
		Script: regSet(`HKCU:\System\GameConfigStore`, "GameDVR_DXGIHonorFSEWindowsCompatible", "1", "DWord"),
	},
	{
		Name: "Disable IPv6",
		Desc: "Completely disables IPv6 protocol across all adapters (CAUTION: may break some VPNs).",
		Script: joinPS(
			regSet(`HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip6\Parameters`, "DisabledComponents", "255", "DWord"),
			`Disable-NetAdapterBinding -Name * -ComponentID ms_tcpip6 -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Disable Microsoft Copilot",
		Desc: "Removes Copilot AppX packages and disables the AI assistant system-wide.",
		Script: joinPS(
			`Get-AppxPackage -AllUsers *Copilot* -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers`,
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\WindowsCopilot`, "TurnOffWindowsCopilot", "1", "DWord"),
		),
	},
	{
		Name: "Disable Notification Tray/Calendar",
		Desc: "Disables all notifications including the Action Center calendar widget.",
		Script: joinPS(
			regSet(`HKCU:\Software\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter", "1", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", "0", "DWord"),
		),
	},
	{
		Name: "Disable Storage Sense",
		Desc: "Prevents Windows from automatically deleting temp files to free space.",
		Script: regSet(`HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\StorageSense\Parameters\StoragePolicy`, "01", "0", "DWord"),
	},
	{
		Name: "Disable Teredo",
		Desc: "Disables IPv6-over-IPv4 tunneling. Reduces attack surface and can lower latency.",
		Script: joinPS(
			regSet(`HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip6\Parameters`, "DisabledComponents", "1", "DWord"),
			`netsh interface teredo set state disabled`,
		),
	},
	{
		Name: "Prefer IPv4 over IPv6",
		Desc: "Sets IPv4 as preferred protocol. Benefits latency on most private networks.",
		Script: regSet(`HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip6\Parameters`, "DisabledComponents", "32", "DWord"),
	},
	{
		Name: "Remove Microsoft Edge",
		Desc: "CAUTION: Uninstalls Microsoft Edge browser completely. Some Windows features depend on it.",
		Script: joinPS(
			`Get-AppxPackage -AllUsers *MicrosoftEdge* | Remove-AppxPackage -AllUsers -ErrorAction SilentlyContinue`,
			`$edgeDir = "${env:ProgramFiles(x86)}\Microsoft\Edge\Application"; if (Test-Path $edgeDir) { Get-ChildItem $edgeDir -Filter "setup.exe" -Recurse | ForEach-Object { Start-Process $_.FullName -ArgumentList '--uninstall --system-level --verbose-logging --force-uninstall' -Wait } }`,
		),
	},
	{
		Name: "Remove OneDrive",
		Desc: "CAUTION: Uninstalls OneDrive completely. Use if you prefer Google Drive exclusively.",
		Script: joinPS(
			`Stop-Process -Name OneDrive -Force -ErrorAction SilentlyContinue`,
			`Start-Process "$env:SystemRoot\SysWOW64\OneDriveSetup.exe" -ArgumentList '/uninstall' -Wait -ErrorAction SilentlyContinue`,
			`Start-Process "$env:SystemRoot\System32\OneDriveSetup.exe" -ArgumentList '/uninstall' -Wait -ErrorAction SilentlyContinue`,
			`Remove-Item -Path "$env:UserProfile\OneDrive" -Force -Recurse -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Set Classic Right-Click Menu",
		Desc: "Restores the compact Windows 10-style context menu in File Explorer.",
		Script: joinPS(
			`$null = New-Item -Path 'HKCU:\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32' -Force`,
			`Set-ItemProperty -Path 'HKCU:\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32' -Name '(Default)' -Value '' -Force`,
		),
	},
	{
		Name: "Set Time to UTC",
		Desc: "Essential for dual-boot systems. Makes Windows use UTC hardware clock like Linux.",
		Script: regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\TimeZoneInformation`, "RealTimeIsUniversal", "1", "DWord"),
	},
}

// ─── Toggles (WinUtil official registry data) ─────────────────────────────────

var DefaultToggles = []Toggle{
	{
		Name:          "Bing Search in Start Menu",
		Desc:          "Include Bing web results in Start menu search when enabled.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Search`, "BingSearchEnabled", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Search`, "BingSearchEnabled", "0", "DWord"),
	},
	{
		Name:      "Center Taskbar Items",
		Desc:      "[Win11] Centers taskbar icons. Disable to left-align like Windows 10.",
		DefaultOn: true,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "TaskbarAl", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "TaskbarAl", "0", "DWord"),
	},
	{
		Name:      "Dark Theme for Windows",
		Desc:      "Enables dark mode for apps and system UI.",
		DefaultOn: true,
		EnableScript: joinPS(
			regSet(`HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, "AppsUseLightTheme", "0", "DWord"),
			regSet(`HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, "SystemUsesLightTheme", "0", "DWord"),
		),
		DisableScript: joinPS(
			regSet(`HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, "AppsUseLightTheme", "1", "DWord"),
			regSet(`HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, "SystemUsesLightTheme", "1", "DWord"),
		),
	},
	{
		Name:          "Detailed BSoD",
		Desc:          "Shows technical stop code details on Blue Screen of Death crashes.",
		DefaultOn:     false,
		EnableScript: joinPS(
			regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\CrashControl`, "DisplayParameters", "1", "DWord"),
			regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\CrashControl`, "DisableEmoticon", "1", "DWord"),
		),
		DisableScript: joinPS(
			regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\CrashControl`, "DisplayParameters", "0", "DWord"),
			regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\CrashControl`, "DisableEmoticon", "0", "DWord"),
		),
	},
	{
		Name:      "Disable Multiplane Overlay",
		Desc:      "Prevents GPU overlay issues. Helps with flickering or black screen in games.",
		DefaultOn: false,
		EnableScript:  regSet(`HKLM:\SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode", "5", "DWord"),
		DisableScript: `Remove-ItemProperty -Path 'HKLM:\SOFTWARE\Microsoft\Windows\Dwm' -Name 'OverlayTestMode' -Force -ErrorAction SilentlyContinue`,
	},
	{
		Name:          "Mouse Acceleration",
		Desc:          "When on, cursor speed responds to physical mouse movement speed.",
		DefaultOn:     false,
		EnableScript: joinPS(
			regSet(`HKCU:\Control Panel\Mouse`, "MouseSpeed", "1", "String"),
			regSet(`HKCU:\Control Panel\Mouse`, "MouseThreshold1", "6", "String"),
			regSet(`HKCU:\Control Panel\Mouse`, "MouseThreshold2", "10", "String"),
		),
		DisableScript: joinPS(
			regSet(`HKCU:\Control Panel\Mouse`, "MouseSpeed", "0", "String"),
			regSet(`HKCU:\Control Panel\Mouse`, "MouseThreshold1", "0", "String"),
			regSet(`HKCU:\Control Panel\Mouse`, "MouseThreshold2", "0", "String"),
		),
	},
	{
		Name:          "Num Lock on Startup",
		Desc:          "Automatically activates Num Lock key when Windows starts.",
		DefaultOn:     false,
		EnableScript: joinPS(
			regSet(`HKU:\.DEFAULT\Control Panel\Keyboard`, "InitialKeyboardIndicators", "'2'", "String"),
			regSet(`HKCU:\Control Panel\Keyboard`, "InitialKeyboardIndicators", "'2'", "String"),
		),
		DisableScript: joinPS(
			regSet(`HKU:\.DEFAULT\Control Panel\Keyboard`, "InitialKeyboardIndicators", "'0'", "String"),
			regSet(`HKCU:\Control Panel\Keyboard`, "InitialKeyboardIndicators", "'0'", "String"),
		),
	},
	{
		Name:          "Recommendations in Start Menu",
		Desc:          "Shows app recommendations and suggestions in Start menu.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Start_IrisRecommendations", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Start_IrisRecommendations", "0", "DWord"),
	},
	{
		Name:          "Search Button in Taskbar",
		Desc:          "Shows the search button/icon on the taskbar.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Search`, "SearchboxTaskbarMode", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Search`, "SearchboxTaskbarMode", "0", "DWord"),
	},
	{
		Name:      "Show File Extensions",
		Desc:      "Reveals file extensions like .txt, .jpg in File Explorer.",
		DefaultOn: true,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "HideFileExt", "0", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "HideFileExt", "1", "DWord"),
	},
	{
		Name:          "Show Hidden Files",
		Desc:          "Makes hidden system files and folders visible in File Explorer.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Hidden", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Hidden", "2", "DWord"),
	},
	{
		Name:      "Sticky Keys",
		Desc:      "Accessibility: press modifier keys (Shift/Ctrl) one at a time instead of together.",
		DefaultOn: false,
		EnableScript:  regSet(`HKCU:\Control Panel\Accessibility\StickyKeys`, "Flags", "'510'", "String"),
		DisableScript: regSet(`HKCU:\Control Panel\Accessibility\StickyKeys`, "Flags", "'506'", "String"),
	},
	{
		Name:          "Task View Button in Taskbar",
		Desc:          "Shows the Task View / virtual desktops button on the taskbar.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "ShowTaskViewButton", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "ShowTaskViewButton", "0", "DWord"),
	},
	{
		Name:          "Verbose Messages During Logon",
		Desc:          "Shows detailed status messages during Windows login for troubleshooting.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "VerboseStatus", "1", "DWord"),
		DisableScript: regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "VerboseStatus", "0", "DWord"),
	},
}

// ─── TweakView model ──────────────────────────────────────────────────────────

type tweakPanel int

const (
	panelEssential tweakPanel = iota
	panelAdvanced
	panelToggles
)

type TweakView struct {
	essentialChecked []bool
	advancedChecked  []bool
	toggles          []Toggle

	activePanel tweakPanel
	cursor      int
	offset      int

	width  int
	height int

	log []string
}

func NewTweakView(w, h int) TweakView {
	toggles := make([]Toggle, len(DefaultToggles))
	copy(toggles, DefaultToggles)
	for i := range toggles {
		toggles[i].Enabled = toggles[i].DefaultOn
	}
	return TweakView{
		essentialChecked: make([]bool, len(EssentialTweaks)),
		advancedChecked:  make([]bool, len(AdvancedTweaks)),
		toggles:          toggles,
		width:            w,
		height:           h,
	}
}

type TweakRunMsg struct {
	Scripts []string
	Name    string // human-readable label for logging
}
type TweakResultMsg struct {
	Output string
	Err    error
}

func (v TweakView) Update(msg tea.Msg) (TweakView, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "tab":
			v.activePanel = (v.activePanel + 1) % 3
			v.cursor = 0
			v.offset = 0
		case "shift+tab":
			v.activePanel = (v.activePanel + 2) % 3
			v.cursor = 0
			v.offset = 0
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
				if v.cursor < v.offset {
					v.offset = v.cursor
				}
			}
		case "down", "j":
			max := v.currentListLen() - 1
			if v.cursor < max {
				v.cursor++
				vis := v.visibleLines()
				if v.cursor >= v.offset+vis {
					v.offset = v.cursor - vis + 1
				}
			}
		case " ", "enter":
			v.toggleCurrent()
		case "u":
			return v, func() tea.Msg {
				out, err := pkgs.RunPowerShell(`powercfg /duplicatescheme e9a42b02-d5df-448d-aa00-03f14749eb61 2>&1; powercfg /setactive SCHEME_CURRENT`)
				return TweakResultMsg{Output: out, Err: err}
			}
		case "r":
			return v, func() tea.Msg {
				script := `$scheme = (powercfg /list | Select-String 'Ultimate Performance' | ForEach-Object { ($_ -split '\s+')[3] }) | Select-Object -First 1; if ($scheme) { powercfg /delete $scheme; Write-Host "Removed: $scheme" } else { Write-Host "Ultimate Performance plan not found" }`
				out, err := pkgs.RunPowerShell(script)
				return TweakResultMsg{Output: out, Err: err}
			}
		case "F5":
			scripts := v.collectScripts()
			if len(scripts) > 0 {
				return v, func() tea.Msg { return TweakRunMsg{Scripts: scripts, Name: "Apply Selected Tweaks"} }
			}
		}
	case TweakResultMsg:
		if m.Err != nil {
			v.log = append(v.log, StyleFail.Render("✗ "+m.Err.Error()))
		} else {
			out := strings.TrimSpace(m.Output)
			if out == "" {
				out = "Done."
			}
			v.log = append(v.log, StyleSuccess.Render("✓ "+out))
		}
	}
	return v, nil
}

func (v *TweakView) toggleCurrent() {
	switch v.activePanel {
	case panelEssential:
		if v.cursor < len(v.essentialChecked) {
			v.essentialChecked[v.cursor] = !v.essentialChecked[v.cursor]
		}
	case panelAdvanced:
		if v.cursor < len(v.advancedChecked) {
			v.advancedChecked[v.cursor] = !v.advancedChecked[v.cursor]
		}
	case panelToggles:
		if v.cursor < len(v.toggles) {
			v.toggles[v.cursor].Enabled = !v.toggles[v.cursor].Enabled
		}
	}
}

func (v *TweakView) currentListLen() int {
	switch v.activePanel {
	case panelEssential:
		return len(EssentialTweaks)
	case panelAdvanced:
		return len(AdvancedTweaks)
	case panelToggles:
		return len(v.toggles)
	}
	return 0
}

func (v *TweakView) visibleLines() int {
	h := v.height - 6
	if h < 5 {
		h = 5
	}
	return h
}

func (v TweakView) collectScripts() []string {
	var scripts []string
	for i, c := range v.essentialChecked {
		if c {
			scripts = append(scripts, EssentialTweaks[i].Script)
		}
	}
	for i, c := range v.advancedChecked {
		if c {
			scripts = append(scripts, AdvancedTweaks[i].Script)
		}
	}
	for _, t := range v.toggles {
		if t.Enabled {
			scripts = append(scripts, t.EnableScript)
		} else {
			scripts = append(scripts, t.DisableScript)
		}
	}
	return scripts
}

func (v TweakView) SelectedCount() int {
	n := 0
	for _, c := range v.essentialChecked {
		if c {
			n++
		}
	}
	for _, c := range v.advancedChecked {
		if c {
			n++
		}
	}
	return n
}

// Footer returns the key-hint line (rendered by root model).
func (v TweakView) Footer() string {
	hint := "[Tab] Switch panel  [↑/↓] Navigate  [Space] Toggle  [F5] Apply tweaks  [u] Ultimate Performance  [r] Remove"
	if v.SelectedCount() > 0 {
		hint += StyleCount.Render(fmt.Sprintf("  %d tweaks selected", v.SelectedCount()))
	}
	return StyleFooter.Render(hint)
}

// View renders the tweak content (without footer).
func (v TweakView) View(contentH int) string {
	halfW := v.width/2 - 2
	if halfW < 30 {
		halfW = 30
	}

	leftLines := v.renderLeft(halfW)
	rightLines := v.renderRight(halfW)

	// Pad both sides to contentH
	for len(leftLines) < contentH {
		leftLines = append(leftLines, "")
	}
	for len(rightLines) < contentH {
		rightLines = append(rightLines, "")
	}
	if len(leftLines) > contentH {
		leftLines = leftLines[:contentH]
	}
	if len(rightLines) > contentH {
		rightLines = rightLines[:contentH]
	}

	var rows []string
	for i := 0; i < contentH; i++ {
		l := leftLines[i]
		r := rightLines[i]
		lw := lipgloss.Width(l)
		if lw < halfW {
			l += strings.Repeat(" ", halfW-lw)
		}
		rows = append(rows, l+"  "+r)
	}
	return strings.Join(rows, "\n")
}

func (v TweakView) renderLeft(w int) []string {
	var lines []string

	lines = append(lines, StyleSectionHeader.Render("Essential Tweaks"))
	for i, t := range EssentialTweaks {
		lines = append(lines, v.tweakRow(i, t.Name, v.essentialChecked[i], panelEssential, w))
	}
	lines = append(lines, "")
	lines = append(lines, StyleWarn.Render("Advanced Tweaks — CAUTION"))
	for i, t := range AdvancedTweaks {
		lines = append(lines, v.tweakRow(i, t.Name, v.advancedChecked[i], panelAdvanced, w))
	}
	lines = append(lines, "")
	// Description tooltip
	switch v.activePanel {
	case panelEssential:
		if v.cursor < len(EssentialTweaks) {
			lines = append(lines, StyleAppNote.Render("  ℹ "+EssentialTweaks[v.cursor].Desc))
		}
	case panelAdvanced:
		if v.cursor < len(AdvancedTweaks) {
			lines = append(lines, StyleAppNote.Render("  ℹ "+AdvancedTweaks[v.cursor].Desc))
		}
	}
	return lines
}

func (v TweakView) tweakRow(idx int, name string, checked bool, panel tweakPanel, w int) string {
	isCursor := v.activePanel == panel && v.cursor == idx
	prefix := "  "
	if isCursor {
		prefix = StyleInfo.Render("▶ ")
	}
	cb := CheckboxStr(checked)
	label := truncate(name, w-6)
	if isCursor {
		label = StyleAppCursor.Render(label)
	} else if checked {
		label = StyleAppSelected.Render(label)
	} else {
		label = StyleAppNormal.Render(label)
	}
	return fmt.Sprintf("%s%s %s", prefix, cb, label)
}

func (v TweakView) renderRight(w int) []string {
	var lines []string

	lines = append(lines, StyleSectionHeader.Render("Customize Preferences"))
	for i, t := range v.toggles {
		isCursor := v.activePanel == panelToggles && v.cursor == i
		prefix := "  "
		if isCursor {
			prefix = StyleInfo.Render("▶ ")
		}
		tog := ToggleStr(t.Enabled)
		label := truncate(t.Name, w-12)
		if isCursor {
			label = StyleAppCursor.Render(label)
		} else {
			label = StyleToggleLabel.Render(label)
		}
		lines = append(lines, fmt.Sprintf("%s%s %s", prefix, tog, label))
	}
	lines = append(lines, "")
	lines = append(lines, StyleSectionHeader.Render("Performance Plans"))
	lines = append(lines, "  "+StyleButton.Render(" Add and Activate Ultimate Performance Profile ")+"  [u]")
	lines = append(lines, "  "+StyleButtonDanger.Render(" Remove Ultimate Performance Profile ")+"  [r]")

	// Log tail
	if len(v.log) > 0 {
		lines = append(lines, "")
		lines = append(lines, StyleSectionHeader.Render("Output"))
		start := len(v.log) - 4
		if start < 0 {
			start = 0
		}
		for _, l := range v.log[start:] {
			lines = append(lines, "  "+l)
		}
	}

	// Toggle description
	if v.activePanel == panelToggles && v.cursor < len(v.toggles) {
		lines = append(lines, "")
		lines = append(lines, StyleAppNote.Render("  ℹ "+v.toggles[v.cursor].Desc))
	}

	return lines
}
