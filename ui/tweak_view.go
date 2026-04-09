package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inspiring-group/inspiring-swiss-knife/pkgs"
)

func regSet(path, name, value, regType string) string {
	return fmt.Sprintf(
		`$null = New-Item -Path '%s' -Force -ErrorAction SilentlyContinue; Set-ItemProperty -Path '%s' -Name '%s' -Value %s -Type %s -Force`,
		path, path, name, value, regType,
	)
}

func regRemove(path, name string) string {
	return fmt.Sprintf(
		`Remove-ItemProperty -Path '%s' -Name '%s' -Force -ErrorAction SilentlyContinue`,
		path, name,
	)
}

func joinPS(cmds ...string) string {
	return strings.Join(cmds, "; ")
}

func enableOptionalFeatures(features []string, extras ...string) string {
	cmds := make([]string, 0, len(features)+len(extras))
	for _, feature := range features {
		cmds = append(cmds, fmt.Sprintf(
			`Enable-WindowsOptionalFeature -Online -FeatureName '%s' -All -NoRestart`,
			feature,
		))
	}
	cmds = append(cmds, extras...)
	return joinPS(cmds...)
}

const restartExplorerScript = `Stop-Process -Name explorer -Force -ErrorAction SilentlyContinue; Start-Process explorer.exe`

type Tweak struct {
	Name   string
	Desc   string
	Script string
}

type Toggle struct {
	Name          string
	Desc          string
	DefaultOn     bool
	Enabled       bool
	Dirty         bool
	EnableScript  string
	DisableScript string
}

var EssentialTweaks = []Tweak{
	{
		Name: "Create Restore Point",
		Desc: "Creates a restore checkpoint and unlocks restore-point frequency before changes are applied.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\SystemRestore`, "SystemRestorePointCreationFrequency", "0", "DWord"),
			`if (-not (Get-ComputerRestorePoint -ErrorAction SilentlyContinue)) { Enable-ComputerRestore -Drive $Env:SystemDrive }`,
			`Checkpoint-Computer -Description 'System Restore Point created by WinUtil' -RestorePointType MODIFY_SETTINGS`,
		),
	},
	{
		Name: "Delete Temporary Files",
		Desc: "Clears user and system TEMP folders to recover disk space.",
		Script: joinPS(
			`Remove-Item -Path "$Env:Temp\*" -Recurse -Force -ErrorAction SilentlyContinue`,
			`Remove-Item -Path "$Env:SystemRoot\Temp\*" -Recurse -Force -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Disable Activity History",
		Desc: "Turns off Windows activity history, upload, and publishing.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\System`, "EnableActivityFeed", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\System`, "PublishUserActivities", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\System`, "UploadUserActivities", "0", "DWord"),
		),
	},
	{
		Name:   "Disable Consumer Features",
		Desc:   "Stops Windows from silently installing promoted apps and games.",
		Script: regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\CloudContent`, "DisableWindowsConsumerFeatures", "1", "DWord"),
	},
	{
		Name: "Disable Explorer Auto Folder Discovery",
		Desc: "Resets bag views and forces generic folder templates instead of content sniffing.",
		Script: joinPS(
			`Remove-Item -Path 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\Bags' -Recurse -Force -ErrorAction SilentlyContinue`,
			`Remove-Item -Path 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\BagMRU' -Recurse -Force -ErrorAction SilentlyContinue`,
			`$null = New-Item -Path 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\Bags\AllFolders\Shell' -Force`,
			`Set-ItemProperty -Path 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\Bags\AllFolders\Shell' -Name 'FolderType' -Value 'NotSpecified' -Force`,
		),
	},
	{
		Name: "Disable Hibernation",
		Desc: "Disables hibernate and removes the hiberfil footprint from disk.",
		Script: joinPS(
			regSet(`HKLM:\System\CurrentControlSet\Control\Session Manager\Power`, "HibernateEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\FlyoutMenuSettings`, "ShowHibernateOption", "0", "DWord"),
			`powercfg.exe /hibernate off`,
		),
	},
	{
		Name: "Disable Location Tracking",
		Desc: "Disables location services, sensors, and maps auto updates.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\CapabilityAccessManager\ConsentStore\location`, "Value", "'Deny'", "String"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Sensor\Overrides\{BFA794E4-F964-4FDB-90F6-51056BFE4B44}`, "SensorPermissionState", "0", "DWord"),
			regSet(`HKLM:\SYSTEM\Maps`, "AutoUpdateEnabled", "0", "DWord"),
			`Set-Service -Name lfsvc -StartupType Disabled -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Disable Microsoft Store Search Results",
		Desc: "Removes Store suggestions from search and blocks open-with Store prompts.",
		Script: joinPS(
			`icacls "$Env:LocalAppData\Packages\Microsoft.WindowsStore_8wekyb3d8bbwe\LocalState\store.db" /deny Everyone:F 2>$null`,
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\Explorer`, "NoUseStoreOpenWith", "1", "DWord"),
		),
	},
	{
		Name:   "Disable PowerShell 7 Telemetry",
		Desc:   "Sets the machine-wide opt-out for PowerShell telemetry.",
		Script: `[Environment]::SetEnvironmentVariable('POWERSHELL_TELEMETRY_OPTOUT', '1', 'Machine')`,
	},
	{
		Name: "Disable Telemetry",
		Desc: "Disables advertising ID, tailored experiences, telemetry services, and sample submission.",
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
		Name:   "Disable Windows Platform Binary Table (WPBT)",
		Desc:   "Stops OEM firmware payloads from dropping software back into Windows at boot.",
		Script: regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager`, "DisableWpbtExecution", "1", "DWord"),
	},
	{
		Name:   "Enable End Task With Right Click",
		Desc:   "Adds the taskbar End Task entry for faster process termination.",
		Script: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced\TaskbarDeveloperSettings`, "TaskbarEndTask", "1", "DWord"),
	},
	{
		Name: "Remove Widgets",
		Desc: "Removes widget packages and restarts Explorer so the taskbar updates immediately.",
		Script: joinPS(
			`Get-Process -Name WidgetService,Widgets,msedgewebview2 -ErrorAction SilentlyContinue | Stop-Process -Force`,
			`Get-AppxPackage -AllUsers Microsoft.WidgetsPlatformRuntime -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers`,
			`Get-AppxPackage -AllUsers MicrosoftWindows.Client.WebExperience -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers`,
			restartExplorerScript,
		),
	},
	{
		Name: "Run Disk Cleanup",
		Desc: "Runs the built-in cleanup pass and then performs component store cleanup.",
		Script: joinPS(
			`cleanmgr.exe /d C: /VERYLOWDISK`,
			`Dism.exe /online /Cleanup-Image /StartComponentCleanup /ResetBase`,
		),
	},
	{
		Name: "Set Services to Manual",
		Desc: "Moves non-essential Windows services to manual start while keeping them available on demand.",
		Script: joinPS(
			`$svcs = @('ALG','AppMgmt','AppReadiness','CDPSvc','PcaSvc','InventorySvc','StorSvc','UsoSvc','WpnService')`,
			`foreach ($svc in $svcs) { Set-Service -Name $svc -StartupType Manual -ErrorAction SilentlyContinue }`,
			`Set-Service -Name DiagTrack -StartupType Disabled -ErrorAction SilentlyContinue`,
		),
	},
}

var AdvancedTweaks = []Tweak{
	{
		Name:   "Adobe Network Block",
		Desc:   "Blocks common Adobe licensing and telemetry endpoints in the hosts file.",
		Script: "Add-Content -Path \"$env:SystemRoot\\System32\\drivers\\etc\\hosts\" -Value \"`n127.0.0.1 lmlicenses.wip4.adobe.com`n127.0.0.1 lm.licenses.adobe.com`n127.0.0.1 na2.adobebettle.com`n127.0.0.1 apiu.adobe.com\"",
	},
	{
		Name: "Block Razer Software Installs",
		Desc: "Blocks Windows from auto-pulling Razer software through device installs.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\DeviceInstall\Restrictions`, "DenyDeviceIDs", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\DeviceInstall\Restrictions`, "DenyDeviceIDsRetroactive", "0", "DWord"),
		),
	},
	{
		Name: "Brave Debloat",
		Desc: "Disables Brave wallet, VPN, rewards, background apps, and noisy promotions.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Policies\BraveSoftware\Brave`, "BraveRewardsDisabled", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\BraveSoftware\Brave`, "BraveVPNDisabled", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\BraveSoftware\Brave`, "WalletDisabled", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\BraveSoftware\Brave`, "BraveNewsDisabled", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\BraveSoftware\Brave`, "ShowTalkFeature", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\BraveSoftware\Brave`, "PromotionsEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\BraveSoftware\Brave`, "BackgroundModeEnabled", "0", "DWord"),
		),
	},
	{
		Name:   "Disable Background Apps",
		Desc:   "Prevents Store apps from running and refreshing in the background.",
		Script: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\BackgroundAccessApplications`, "GlobalUserDisabled", "1", "DWord"),
	},
	{
		Name:   "Disable Fullscreen Optimizations",
		Desc:   "Disables fullscreen optimization shims for better exclusive-fullscreen behavior.",
		Script: regSet(`HKCU:\System\GameConfigStore`, "GameDVR_DXGIHonorFSEWindowsCompatible", "1", "DWord"),
	},
	{
		Name: "Disable IPv6",
		Desc: "Disables IPv6 globally and unbinds it from adapters. Use carefully on VPN-heavy setups.",
		Script: joinPS(
			regSet(`HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip6\Parameters`, "DisabledComponents", "255", "DWord"),
			`Disable-NetAdapterBinding -Name '*' -ComponentID ms_tcpip6 -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Disable Microsoft Copilot",
		Desc: "Removes Copilot packages and applies the policy and EndOfLife keys WinUtil uses.",
		Script: joinPS(
			`Get-AppxPackage -AllUsers *Copilot* -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers`,
			`Get-AppxPackage -AllUsers MicrosoftWindows.Client.WebExperience -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers`,
			`Get-AppxPackage -AllUsers Microsoft.549981C3F5F10 -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers`,
			`foreach ($user in (Get-LocalUser -ErrorAction SilentlyContinue)) { reg load "HKU\WinUtilTemp" "C:\Users\$($user.Name)\NTUSER.DAT" >$null 2>&1; $null = New-Item -Path "HKU:\WinUtilTemp\Software\Policies\Microsoft\Windows\WindowsCopilot" -Force -ErrorAction SilentlyContinue; Set-ItemProperty -Path "HKU:\WinUtilTemp\Software\Policies\Microsoft\Windows\WindowsCopilot" -Name "TurnOffWindowsCopilot" -Value 1 -Type DWord -Force -ErrorAction SilentlyContinue; reg unload "HKU\WinUtilTemp" >$null 2>&1 }`,
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\WindowsCopilot`, "TurnOffWindowsCopilot", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Communications`, "ConfigureChatAutoInstall", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Microsoft\WindowsUpdate\Orchestrator\UScheduler_Oobe`, "OutlookUpdate", "'Blocked'", "String"),
			regSet(`HKLM:\SOFTWARE\Microsoft\WindowsUpdate\Orchestrator\UScheduler_Oobe`, "MSChatUpdate", "'Blocked'", "String"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\CloudContent`, "DisableCloudOptimizedContent", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\CloudContent`, "DisableConsumerAccountStateContent", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Windows\CloudContent`, "DisableWindowsConsumerFeatures", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Appx\AppxAllUserStore\EndOfLife\{S-1-1-0}\MicrosoftWindows.Client.WebExperience_423.1000.0.0_neutral__cw5n1h2txyewy`, "", "''", "String"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Appx\AppxAllUserStore\EndOfLife\{S-1-1-0}\Microsoft.549981C3F5F10_8wekyb3d8bbwe`, "", "''", "String"),
		),
	},
	{
		Name: "Disable Notification Tray/Calendar",
		Desc: "Turns off the notification center and toast notifications.",
		Script: joinPS(
			regSet(`HKCU:\Software\Policies\Microsoft\Windows\Explorer`, "DisableNotificationCenter", "1", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\PushNotifications`, "ToastEnabled", "0", "DWord"),
		),
	},
	{
		Name:   "Disable Storage Sense",
		Desc:   "Turns off automated Storage Sense cleanup.",
		Script: regSet(`HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\StorageSense\Parameters\StoragePolicy`, "01", "0", "DWord"),
	},
	{
		Name: "Disable Teredo",
		Desc: "Disables the Teredo tunneling interface and related IPv6 preference.",
		Script: joinPS(
			regSet(`HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip6\Parameters`, "DisabledComponents", "1", "DWord"),
			`netsh interface teredo set state disabled`,
		),
	},
	{
		Name: "DNS (Cloudflare Preset)",
		Desc: "Applies Cloudflare DNS to active adapters. WinUtil exposes this through a DNS picker.",
		Script: joinPS(
			`$adapters = Get-NetAdapter | Where-Object { $_.Status -eq 'Up' -and $_.HardwareInterface }`,
			`foreach ($adapter in $adapters) { Set-DnsClientServerAddress -InterfaceIndex $adapter.InterfaceIndex -ServerAddresses @('1.1.1.1','1.0.0.1') -ErrorAction SilentlyContinue }`,
		),
	},
	{
		Name: "Edge Debloat",
		Desc: "Disables Edge promos, shopping, sidebar, wallet, telemetry hooks, and background behavior.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "EdgeEnhanceImagesEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "PersonalizationReportingEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "ShowRecommendationsEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "HideFirstRunExperience", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "UserFeedbackAllowed", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "ConfigureDoNotTrack", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "AlternateErrorPagesEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "EdgeCollectionsEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "EdgeShoppingAssistantEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "MicrosoftEdgeInsiderPromotionEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "ShowMicrosoftRewards", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "WebWidgetAllowed", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "DiagnosticData", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "WalletDonationEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "EdgeAssetDeliveryServiceEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "CryptoWalletEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "WalletEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "HubsSidebarEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "StandaloneHubsSidebarEnabled", "0", "DWord"),
			regSet(`HKLM:\SOFTWARE\Policies\Microsoft\Edge`, "BackgroundModeEnabled", "0", "DWord"),
		),
	},
	{
		Name:   "Prefer IPv4 over IPv6",
		Desc:   "Keeps IPv6 available but changes protocol preference to IPv4 first.",
		Script: regSet(`HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip6\Parameters`, "DisabledComponents", "32", "DWord"),
	},
	{
		Name: "Remove Unwanted Pre-Installed Apps",
		Desc: "Removes the standard consumer AppX set that WinUtil classifies as pre-installed bloat.",
		Script: joinPS(
			`$apps = @('*ActiproSoftwareLLC*','*AdobeSystemsIncorporated.AdobePhotoshopExpress*','*BubbleWitch3Saga*','*CandyCrush*','*DevHome*','*Disney*','*Dolby*','*Duolingo-LearnLanguagesforFree*','*EclipseManager*','*Facebook*','*Flipboard*','*gaming*','*Minecraft*','*Office*','*PandoraMediaInc*','*Royal Revolt*','*Sway*','*Twitter*','*Wunderlist*','*clipchamp*','*getstarted*','*bing*','*solit*','*MixedReality*','*speed test*','*McAfee*','*Tiktok*')`,
			`foreach ($app in $apps) { Get-AppxPackage -AllUsers $app -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers -ErrorAction SilentlyContinue }`,
			`Get-AppxPackage -AllUsers MicrosoftTeams -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers -ErrorAction SilentlyContinue`,
		),
	},
	{
		Name: "Remove Gallery from Explorer",
		Desc: "Removes the Windows 11 Gallery node from Explorer and refreshes Explorer.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Desktop\NameSpace_41040327`, "{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}", "''", "String"),
			restartExplorerScript,
		),
	},
	{
		Name: "Remove Home from Explorer",
		Desc: "Removes the Home shortcut in Explorer by blocking the namespace entry and pinning behavior.",
		Script: joinPS(
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer`, "HubMode", "1", "DWord"),
			regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Shell Extensions\Blocked`, "{e88865ea-0e1c-4e20-9aa6-edcd0212c87c}", "''", "String"),
			restartExplorerScript,
		),
	},
	{
		Name: "Remove Microsoft Edge",
		Desc: "Uninstalls Edge packages and then runs the Edge setup uninstallers if present.",
		Script: joinPS(
			`Get-AppxPackage -AllUsers *MicrosoftEdge* -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers -ErrorAction SilentlyContinue`,
			`$edgeDir = "${env:ProgramFiles(x86)}\Microsoft\Edge\Application"; if (Test-Path $edgeDir) { Get-ChildItem $edgeDir -Filter 'setup.exe' -Recurse | ForEach-Object { Start-Process $_.FullName -ArgumentList '--uninstall --system-level --verbose-logging --force-uninstall' -Wait } }`,
		),
	},
	{
		Name: "Remove OneDrive",
		Desc: "Fully removes OneDrive, blocks its folder writeback, and disables related sync services.",
		Script: joinPS(
			`Stop-Process -Name OneDrive -Force -ErrorAction SilentlyContinue`,
			`$setupPaths = @("$Env:SystemRoot\System32\OneDriveSetup.exe", "$Env:SystemRoot\SysWOW64\OneDriveSetup.exe")`,
			`foreach ($setup in $setupPaths) { if (Test-Path $setup) { Start-Process $setup -ArgumentList '/uninstall' -Wait } }`,
			`Stop-Process -Name FileCoAuth,Explorer -Force -ErrorAction SilentlyContinue`,
			`if (Test-Path $Env:OneDrive) { icacls $Env:OneDrive /deny 'Administrators:(D,DC)' | Out-Null }`,
			`Remove-Item "$Env:UserProfile\OneDrive" -Recurse -Force -ErrorAction SilentlyContinue`,
			`Remove-Item "$Env:LocalAppData\Microsoft\OneDrive" -Recurse -Force -ErrorAction SilentlyContinue`,
			`Remove-Item "$Env:ProgramData\Microsoft OneDrive" -Recurse -Force -ErrorAction SilentlyContinue`,
			`Remove-Item "$Env:SystemDrive\OneDriveTemp" -Recurse -Force -ErrorAction SilentlyContinue`,
			`Set-Service -Name OneSyncSvc -StartupType Disabled -ErrorAction SilentlyContinue`,
			`Start-Process explorer.exe`,
		),
	},
	{
		Name: "Remove Xbox & Gaming Components",
		Desc: "Removes Xbox and gaming capture packages and disables GameDVR capture hooks.",
		Script: joinPS(
			regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\GameDVR`, "AppCaptureEnabled", "0", "DWord"),
			regSet(`HKCU:\System\GameConfigStore`, "GameDVR_Enabled", "0", "DWord"),
			`$apps = @('*Xbox*','*GamingApp*','*Xbox.TCUI*','*XboxSpeechToTextOverlay*','*XboxGameOverlay*','*XboxGamingOverlay*','*XboxIdentityProvider*','*XboxGameCallableUI*','*Microsoft.GamingApp*','*Microsoft.XboxGamingOverlay*')`,
			`foreach ($app in $apps) { Get-AppxPackage -AllUsers $app -ErrorAction SilentlyContinue | Remove-AppxPackage -AllUsers -ErrorAction SilentlyContinue }`,
		),
	},
	{
		Name: "Revert the New Start Menu",
		Desc: "Rebuilds Start menu shell state and restarts Explorer to reset Win11 Start layout state.",
		Script: joinPS(
			`Remove-Item 'HKCU:\Software\Microsoft\Windows\CurrentVersion\CloudStore\Store\Cache\DefaultAccount\*windows.data.placeholdertilecollection*' -Recurse -Force -ErrorAction SilentlyContinue`,
			`Remove-Item 'HKCU:\Software\Microsoft\Windows\CurrentVersion\CloudStore\Store\Cache\DefaultAccount\*start*' -Recurse -Force -ErrorAction SilentlyContinue`,
			restartExplorerScript,
		),
	},
	{
		Name: "Run OO Shutup 10",
		Desc: "Downloads O&O ShutUp10++, applies the recommended config, and runs it silently.",
		Script: joinPS(
			`$tool = Join-Path $Env:TEMP 'OOSU10.exe'`,
			`$cfg = Join-Path $Env:TEMP 'ooshutup10.cfg'`,
			`Invoke-WebRequest -Uri 'https://dl5.oo-software.com/files/ooshutup10/OOSU10.exe' -OutFile $tool`,
			`Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/ChrisTitusTech/winutil/main/config/ooshutup10.cfg' -OutFile $cfg`,
			`Start-Process $tool -ArgumentList "$cfg /quiet" -Wait`,
		),
	},
	{
		Name: "Set Classic Right-Click Menu",
		Desc: "Restores the Windows 10 style Explorer context menu.",
		Script: joinPS(
			`$null = New-Item -Path 'HKCU:\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32' -Force`,
			`Set-ItemProperty -Path 'HKCU:\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32' -Name '(Default)' -Value '' -Force`,
			restartExplorerScript,
		),
	},
	{
		Name: "Set Display for Performance",
		Desc: "Switches Windows visual effects to the best-performance preset.",
		Script: joinPS(
			regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects`, "VisualFXSetting", "2", "DWord"),
			regSet(`HKCU:\Control Panel\Desktop\WindowMetrics`, "MinAnimate", "'0'", "String"),
			restartExplorerScript,
		),
	},
	{
		Name:   "Set Time to UTC",
		Desc:   "Sets Windows to use UTC hardware clock for dual-boot setups.",
		Script: regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\TimeZoneInformation`, "RealTimeIsUniversal", "1", "DWord"),
	},
}

var FixTweaks = []Tweak{
	{
		Name: "Reset Network",
		Desc: "Flushes DNS, resets Winsock and TCP/IP, and renews the current lease.",
		Script: joinPS(
			`ipconfig /flushdns`,
			`netsh winsock reset`,
			`netsh int ip reset`,
			`ipconfig /release`,
			`ipconfig /renew`,
		),
	},
	{
		Name: "Reset Windows Update",
		Desc: "Stops update services, rebuilds SoftwareDistribution/catroot2, and starts services again.",
		Script: joinPS(
			`$services = 'bits','wuauserv','appidsvc','cryptsvc'`,
			`foreach ($svc in $services) { Stop-Service -Name $svc -Force -ErrorAction SilentlyContinue }`,
			`Remove-Item "$Env:ALLUSERSPROFILE\Application Data\Microsoft\Network\Downloader\qmgr*.dat" -Force -ErrorAction SilentlyContinue`,
			`Rename-Item "$Env:SystemRoot\SoftwareDistribution" 'SoftwareDistribution.bak' -ErrorAction SilentlyContinue`,
			`Rename-Item "$Env:SystemRoot\System32\catroot2" 'catroot2.bak' -ErrorAction SilentlyContinue`,
			`foreach ($svc in $services) { Start-Service -Name $svc -ErrorAction SilentlyContinue }`,
		),
	},
	{
		Name: "Set Up Auto Login",
		Desc: "Downloads Sysinternals Autologon and launches it so credentials can be configured interactively.",
		Script: joinPS(
			`$zip = Join-Path $Env:TEMP 'Autologon.zip'`,
			`$dir = Join-Path $Env:TEMP 'Autologon'`,
			`Invoke-WebRequest -Uri 'https://download.sysinternals.com/files/AutoLogon.zip' -OutFile $zip`,
			`Expand-Archive -Path $zip -DestinationPath $dir -Force`,
			`Start-Process (Join-Path $dir 'Autologon.exe')`,
		),
	},
	{
		Name: "System Corruption Scan",
		Desc: "Runs DISM health restore and an SFC pass against the live image.",
		Script: joinPS(
			`DISM /Online /Cleanup-Image /RestoreHealth`,
			`sfc /scannow`,
		),
	},
	{
		Name: "WinGet Reinstall",
		Desc: "Re-downloads the App Installer bundle from Microsoft and reinstalls it.",
		Script: joinPS(
			`$bundle = Join-Path $Env:TEMP 'Microsoft.DesktopAppInstaller.msixbundle'`,
			`Invoke-WebRequest -Uri 'https://aka.ms/getwinget' -OutFile $bundle`,
			`Add-AppxPackage -Path $bundle -ForceApplicationShutdown -ForceUpdateFromAnyVersion`,
		),
	},
}

var LegacyPanels = []Tweak{
	{Name: "Computer Management", Desc: "Opens the classic Computer Management MMC.", Script: `Start-Process compmgmt.msc`},
	{Name: "Control Panel", Desc: "Opens the legacy Control Panel shell.", Script: `Start-Process control.exe`},
	{Name: "Network Connections", Desc: "Opens the classic network adapter list.", Script: `Start-Process ncpa.cpl`},
	{Name: "Power Panel", Desc: "Opens legacy power options.", Script: `Start-Process powercfg.cpl`},
	{Name: "Printer Panel", Desc: "Opens the legacy printers and devices view.", Script: `Start-Process 'shell:::{A8A91A66-3A7D-4424-8D24-04E180695C7A}'`},
	{Name: "Region", Desc: "Opens the legacy region and locale dialog.", Script: `Start-Process intl.cpl`},
	{Name: "Sound Settings", Desc: "Opens the classic sound control panel.", Script: `Start-Process mmsys.cpl`},
	{Name: "System Properties", Desc: "Opens System Properties.", Script: `Start-Process sysdm.cpl`},
	{Name: "Time and Date", Desc: "Opens the classic date and time dialog.", Script: `Start-Process timedate.cpl`},
	{Name: "Windows Restore", Desc: "Launches the System Restore wizard.", Script: `Start-Process rstrui.exe`},
}

var FeatureTweaks = []Tweak{
	{
		Name:   "Enable .NET Framework",
		Desc:   "Enables .NET Framework 3.5 and the .NET 4 advanced services feature set.",
		Script: enableOptionalFeatures([]string{"NetFx3", "NetFx4-AdvSrvs"}),
	},
	{
		Name: "Enable Hyper-V",
		Desc: "Enables the full Hyper-V stack and applies the scheduler tweak WinUtil uses.",
		Script: enableOptionalFeatures(
			[]string{"Microsoft-Hyper-V-All"},
			`bcdedit /set hypervisorschedulertype classic`,
		),
	},
	{
		Name:   "Enable Legacy Media",
		Desc:   "Enables Windows Media Player and DirectPlay compatibility features.",
		Script: enableOptionalFeatures([]string{"WindowsMediaPlayer", "DirectPlay"}),
	},
	{
		Name: "Enable NFS",
		Desc: "Installs NFS client services and starts the related services.",
		Script: enableOptionalFeatures(
			[]string{"ServicesForNFS-ClientOnly", "ClientForNFS-Infrastructure", "NFS-Administration"},
			`Start-Service NfsClnt -ErrorAction SilentlyContinue`,
			`Start-Service Rpcxdr -ErrorAction SilentlyContinue`,
			`nfsadmin client start`,
		),
	},
	{
		Name:   "Enable Windows Sandbox",
		Desc:   "Enables the Windows Sandbox optional feature.",
		Script: enableOptionalFeatures([]string{"Containers-DisposableClientVM"}),
	},
	{
		Name: "Enable Windows Subsystem for Linux",
		Desc: "Enables WSL and the virtual machine platform prerequisites.",
		Script: enableOptionalFeatures(
			[]string{"Microsoft-Windows-Subsystem-Linux", "VirtualMachinePlatform"},
			`wsl --install --no-distribution`,
		),
	},
	{
		Name:   "Enable Legacy F8 Boot Recovery",
		Desc:   "Restores the legacy F8 boot menu behavior.",
		Script: `bcdedit /set '{current}' bootmenupolicy legacy`,
	},
	{
		Name:   "Disable Legacy F8 Boot Recovery",
		Desc:   "Returns boot menu policy to the Windows standard behavior.",
		Script: `bcdedit /set '{current}' bootmenupolicy standard`,
	},
	{
		Name: "Enable Daily Registry Backup Task",
		Desc: "Re-enables Windows registry backups and creates the RegIdleBackup scheduled task.",
		Script: joinPS(
			regSet(`HKLM:\System\CurrentControlSet\Control\Session Manager\Configuration Manager`, "EnablePeriodicBackup", "1", "DWord"),
			`schtasks /Change /TN "\Microsoft\Windows\Registry\RegIdleBackup" /Enable`,
			`schtasks /Change /TN "\Microsoft\Windows\Registry\RegIdleBackup" /RI 1440`,
		),
	},
}

var RemoteAccessTweaks = []Tweak{
	{
		Name: "Enable OpenSSH Server",
		Desc: "Installs the OpenSSH Server capability, starts it, and opens the firewall rule.",
		Script: joinPS(
			`Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0`,
			`Set-Service -Name sshd -StartupType Automatic -ErrorAction SilentlyContinue`,
			`Start-Service sshd -ErrorAction SilentlyContinue`,
			`if (-not (Get-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -ErrorAction SilentlyContinue)) { New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22 }`,
		),
	},
}

var PowerShellProfileTweaks = []Tweak{
	{
		Name:   "Install CTT PowerShell Profile",
		Desc:   "Bootstraps the Chris Titus Tech PowerShell profile in PowerShell 7+.",
		Script: `if (Get-Command pwsh -ErrorAction SilentlyContinue) { pwsh -NoProfile -ExecutionPolicy Bypass -Command "iex ((irm 'https://christitus.com/win'))" } else { Write-Error 'PowerShell 7 is required for the CTT profile install.' }`,
	},
	{
		Name:   "Uninstall CTT PowerShell Profile",
		Desc:   "Removes the profile files and stub entries created by the CTT profile installer.",
		Script: `if (Get-Command pwsh -ErrorAction SilentlyContinue) { pwsh -NoProfile -ExecutionPolicy Bypass -Command "$profilePaths = @($PROFILE.CurrentUserAllHosts,$PROFILE.CurrentUserCurrentHost); foreach ($path in $profilePaths) { Remove-Item $path -Force -ErrorAction SilentlyContinue }; Remove-Item (Join-Path $HOME '.config\powershell') -Recurse -Force -ErrorAction SilentlyContinue" } else { Write-Error 'PowerShell 7 is required for the CTT profile uninstall.' }`,
	},
}

var DefaultToggles = []Toggle{
	{
		Name:          "Bing Search in Start Menu",
		Desc:          "Shows Bing web results in the Start menu search experience.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Search`, "BingSearchEnabled", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Search`, "BingSearchEnabled", "0", "DWord"),
	},
	{
		Name:          "Center Taskbar Items",
		Desc:          "Centers taskbar icons on Windows 11. Turn it off to left align them.",
		DefaultOn:     true,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "TaskbarAl", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "TaskbarAl", "0", "DWord"),
	},
	{
		Name:          "Cross-Device Resume",
		Desc:          "Controls the cross-device resume experience introduced in Windows 11.",
		DefaultOn:     true,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\CrossDeviceResume`, "CrossDeviceResumeEnabled", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\CrossDeviceResume`, "CrossDeviceResumeEnabled", "0", "DWord"),
	},
	{
		Name:      "Dark Theme for Windows",
		Desc:      "Enables dark mode for apps and system surfaces.",
		DefaultOn: false,
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
		Name:      "Detailed BSoD",
		Desc:      "Shows the technical crash details on blue screens.",
		DefaultOn: false,
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
		Name:          "Disable Multiplane Overlay",
		Desc:          "Applies the MPO workaround commonly used for flicker and black-screen issues.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKLM:\SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode", "5", "DWord"),
		DisableScript: regRemove(`HKLM:\SOFTWARE\Microsoft\Windows\Dwm`, "OverlayTestMode"),
	},
	{
		Name:      "Modern Standby Fix",
		Desc:      "Disables connected-standby networking to reduce overnight drain on Modern Standby systems.",
		DefaultOn: false,
		EnableScript: joinPS(
			`powercfg /SETACVALUEINDEX SCHEME_CURRENT SUB_NONE CONNECTIVITYINSTANDBY 0`,
			`powercfg /SETDCVALUEINDEX SCHEME_CURRENT SUB_NONE CONNECTIVITYINSTANDBY 0`,
			`powercfg /SETACTIVE SCHEME_CURRENT`,
		),
		DisableScript: joinPS(
			`powercfg /SETACVALUEINDEX SCHEME_CURRENT SUB_NONE CONNECTIVITYINSTANDBY 1`,
			`powercfg /SETDCVALUEINDEX SCHEME_CURRENT SUB_NONE CONNECTIVITYINSTANDBY 1`,
			`powercfg /SETACTIVE SCHEME_CURRENT`,
		),
	},
	{
		Name:      "Mouse Acceleration",
		Desc:      "Turns enhanced pointer precision on or off.",
		DefaultOn: true,
		EnableScript: joinPS(
			regSet(`HKCU:\Control Panel\Mouse`, "MouseSpeed", "1", "DWord"),
			regSet(`HKCU:\Control Panel\Mouse`, "MouseThreshold1", "6", "DWord"),
			regSet(`HKCU:\Control Panel\Mouse`, "MouseThreshold2", "10", "DWord"),
		),
		DisableScript: joinPS(
			regSet(`HKCU:\Control Panel\Mouse`, "MouseSpeed", "0", "DWord"),
			regSet(`HKCU:\Control Panel\Mouse`, "MouseThreshold1", "0", "DWord"),
			regSet(`HKCU:\Control Panel\Mouse`, "MouseThreshold2", "0", "DWord"),
		),
	},
	{
		Name:      "New Outlook",
		Desc:      "Controls the new Outlook migration switch and hides or shows the toggle.",
		DefaultOn: true,
		EnableScript: joinPS(
			regSet(`HKCU:\Software\Microsoft\Office\16.0\Outlook\Options\General`, "UseNewOutlook", "1", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Office\16.0\Outlook\Options\General`, "HideNewOutlookToggle", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Office\16.0\Outlook\Preferences`, "NewOutlookMigrationUserSetting", "0", "DWord"),
		),
		DisableScript: joinPS(
			regSet(`HKCU:\Software\Microsoft\Office\16.0\Outlook\Options\General`, "UseNewOutlook", "0", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Office\16.0\Outlook\Options\General`, "HideNewOutlookToggle", "1", "DWord"),
			regSet(`HKCU:\Software\Microsoft\Office\16.0\Outlook\Options\General`, "DoNewOutlookAutoMigration", "0", "DWord"),
			regRemove(`HKCU:\Software\Microsoft\Office\16.0\Outlook\Preferences`, "NewOutlookMigrationUserSetting"),
		),
	},
	{
		Name:      "Num Lock on Startup",
		Desc:      "Turns Num Lock on or off at logon.",
		DefaultOn: false,
		EnableScript: joinPS(
			regSet(`HKU:\.DEFAULT\Control Panel\Keyboard`, "InitialKeyboardIndicators", "2", "DWord"),
			regSet(`HKCU:\Control Panel\Keyboard`, "InitialKeyboardIndicators", "2", "DWord"),
		),
		DisableScript: joinPS(
			regSet(`HKU:\.DEFAULT\Control Panel\Keyboard`, "InitialKeyboardIndicators", "0", "DWord"),
			regSet(`HKCU:\Control Panel\Keyboard`, "InitialKeyboardIndicators", "0", "DWord"),
		),
	},
	{
		Name:          "Recommendations in Start Menu",
		Desc:          "Shows or hides recommendation content in Start.",
		DefaultOn:     true,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Start_IrisRecommendations", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Start_IrisRecommendations", "0", "DWord"),
	},
	{
		Name:          "Remove Settings Home Page",
		Desc:          "Shows the Settings app without the Windows 11 home feed.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer\SettingsPageVisibility`, "SettingsPageVisibility", "'hide:home'", "String"),
		DisableScript: regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer\SettingsPageVisibility`, "SettingsPageVisibility", "'show:home'", "String"),
	},
	{
		Name:          "S3 Sleep",
		Desc:          "Forces classic S3 sleep by disabling Modern Standby through PlatformAoAcOverride.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKLM:\SYSTEM\CurrentControlSet\Control\Power`, "PlatformAoAcOverride", "0", "DWord"),
		DisableScript: regRemove(`HKLM:\SYSTEM\CurrentControlSet\Control\Power`, "PlatformAoAcOverride"),
	},
	{
		Name:          "Search Button in Taskbar",
		Desc:          "Shows or hides the taskbar search button.",
		DefaultOn:     true,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Search`, "SearchboxTaskbarMode", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Search`, "SearchboxTaskbarMode", "0", "DWord"),
	},
	{
		Name:          "Show File Extensions",
		Desc:          "Shows file extensions in Explorer and refreshes Explorer when changed.",
		DefaultOn:     false,
		EnableScript:  joinPS(regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "HideFileExt", "0", "DWord"), restartExplorerScript),
		DisableScript: joinPS(regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "HideFileExt", "1", "DWord"), restartExplorerScript),
	},
	{
		Name:          "Show Hidden Files",
		Desc:          "Shows hidden files and folders in Explorer.",
		DefaultOn:     false,
		EnableScript:  joinPS(regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Hidden", "1", "DWord"), restartExplorerScript),
		DisableScript: joinPS(regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "Hidden", "2", "DWord"), restartExplorerScript),
	},
	{
		Name:          "Sticky Keys",
		Desc:          "Turns Sticky Keys on or off.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKCU:\Control Panel\Accessibility\StickyKeys`, "Flags", "510", "DWord"),
		DisableScript: regSet(`HKCU:\Control Panel\Accessibility\StickyKeys`, "Flags", "506", "DWord"),
	},
	{
		Name:          "Task View Button in Taskbar",
		Desc:          "Shows or hides the task view button on the taskbar.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "ShowTaskViewButton", "1", "DWord"),
		DisableScript: regSet(`HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`, "ShowTaskViewButton", "0", "DWord"),
	},
	{
		Name:          "Verbose Messages During Logon",
		Desc:          "Shows detailed Windows status messages during logon and shutdown.",
		DefaultOn:     false,
		EnableScript:  regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "VerboseStatus", "1", "DWord"),
		DisableScript: regSet(`HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, "VerboseStatus", "0", "DWord"),
	},
}

type tweakPanel int

const (
	panelEssential tweakPanel = iota
	panelAdvanced
	panelFixes
	panelLegacy
	panelFeatures
	panelRemoteAccess
	panelPowerShellProfile
	panelToggles
)

var panelOrder = []tweakPanel{
	panelEssential,
	panelAdvanced,
	panelFixes,
	panelLegacy,
	panelFeatures,
	panelRemoteAccess,
	panelPowerShellProfile,
	panelToggles,
}

var leftPanels = []tweakPanel{
	panelEssential,
	panelAdvanced,
	panelFixes,
	panelLegacy,
}

var rightPanels = []tweakPanel{
	panelFeatures,
	panelRemoteAccess,
	panelPowerShellProfile,
	panelToggles,
}

type panelState struct {
	cursor int
	offset int
}

type TweakView struct {
	essentialChecked []bool
	advancedChecked  []bool
	fixChecked       []bool
	legacyChecked    []bool
	featureChecked   []bool
	remoteChecked    []bool
	profileChecked   []bool
	toggles          []Toggle

	activePanel tweakPanel
	leftPanel   tweakPanel
	rightPanel  tweakPanel
	states      map[tweakPanel]panelState

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

	states := make(map[tweakPanel]panelState, len(panelOrder))
	for _, panel := range panelOrder {
		states[panel] = panelState{}
	}

	return TweakView{
		essentialChecked: make([]bool, len(EssentialTweaks)),
		advancedChecked:  make([]bool, len(AdvancedTweaks)),
		fixChecked:       make([]bool, len(FixTweaks)),
		legacyChecked:    make([]bool, len(LegacyPanels)),
		featureChecked:   make([]bool, len(FeatureTweaks)),
		remoteChecked:    make([]bool, len(RemoteAccessTweaks)),
		profileChecked:   make([]bool, len(PowerShellProfileTweaks)),
		toggles:          toggles,
		activePanel:      panelEssential,
		leftPanel:        panelEssential,
		rightPanel:       panelFeatures,
		states:           states,
		width:            w,
		height:           h,
	}
}

type TweakRunMsg struct {
	Scripts []string
	Name    string
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
			v.switchPanel(1)
		case "shift+tab":
			v.switchPanel(-1)
		case "up", "k":
			v.moveCursor(-1)
		case "down", "j":
			v.moveCursor(1)
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
				return v, func() tea.Msg {
					return TweakRunMsg{Scripts: scripts, Name: "Apply Selected Tweaks"}
				}
			}
		}
	case TweakResultMsg:
		if m.Err != nil {
			v.log = append(v.log, StyleFail.Render("! "+m.Err.Error()))
		} else {
			out := strings.TrimSpace(m.Output)
			if out == "" {
				out = "Done."
			}
			v.log = append(v.log, StyleSuccess.Render("> "+out))
		}
	}

	return v, nil
}

func (v *TweakView) switchPanel(delta int) {
	idx := 0
	for i, panel := range panelOrder {
		if panel == v.activePanel {
			idx = i
			break
		}
	}

	next := (idx + delta + len(panelOrder)) % len(panelOrder)
	v.activePanel = panelOrder[next]
	if v.isRightPanel(v.activePanel) {
		v.rightPanel = v.activePanel
	} else {
		v.leftPanel = v.activePanel
	}
}

func (v *TweakView) moveCursor(delta int) {
	length := v.currentListLen()
	if length == 0 {
		return
	}

	state := v.states[v.activePanel]
	state.cursor += delta
	if state.cursor < 0 {
		state.cursor = 0
	}
	if state.cursor >= length {
		state.cursor = length - 1
	}

	visible := v.visibleLines()
	if state.cursor < state.offset {
		state.offset = state.cursor
	}
	if state.cursor >= state.offset+visible {
		state.offset = state.cursor - visible + 1
	}
	if state.offset < 0 {
		state.offset = 0
	}
	v.states[v.activePanel] = state
}

func (v *TweakView) toggleCurrent() {
	state := v.states[v.activePanel]
	idx := state.cursor

	switch v.activePanel {
	case panelEssential:
		if idx < len(v.essentialChecked) {
			v.essentialChecked[idx] = !v.essentialChecked[idx]
		}
	case panelAdvanced:
		if idx < len(v.advancedChecked) {
			v.advancedChecked[idx] = !v.advancedChecked[idx]
		}
	case panelFixes:
		if idx < len(v.fixChecked) {
			v.fixChecked[idx] = !v.fixChecked[idx]
		}
	case panelLegacy:
		if idx < len(v.legacyChecked) {
			v.legacyChecked[idx] = !v.legacyChecked[idx]
		}
	case panelFeatures:
		if idx < len(v.featureChecked) {
			v.featureChecked[idx] = !v.featureChecked[idx]
		}
	case panelRemoteAccess:
		if idx < len(v.remoteChecked) {
			v.remoteChecked[idx] = !v.remoteChecked[idx]
		}
	case panelPowerShellProfile:
		if idx < len(v.profileChecked) {
			v.profileChecked[idx] = !v.profileChecked[idx]
		}
	case panelToggles:
		if idx < len(v.toggles) {
			v.toggles[idx].Enabled = !v.toggles[idx].Enabled
			v.toggles[idx].Dirty = v.toggles[idx].Enabled != v.toggles[idx].DefaultOn
		}
	}
}

func (v TweakView) currentListLen() int {
	return v.panelListLen(v.activePanel)
}

func (v TweakView) panelListLen(panel tweakPanel) int {
	switch panel {
	case panelEssential:
		return len(EssentialTweaks)
	case panelAdvanced:
		return len(AdvancedTweaks)
	case panelFixes:
		return len(FixTweaks)
	case panelLegacy:
		return len(LegacyPanels)
	case panelFeatures:
		return len(FeatureTweaks)
	case panelRemoteAccess:
		return len(RemoteAccessTweaks)
	case panelPowerShellProfile:
		return len(PowerShellProfileTweaks)
	case panelToggles:
		return len(v.toggles)
	default:
		return 0
	}
}

func (v TweakView) visibleLines() int {
	lines := v.height - 18
	if lines < 5 {
		lines = 5
	}
	return lines
}

func (v TweakView) SelectedCount() int {
	count := 0
	for _, selected := range v.essentialChecked {
		if selected {
			count++
		}
	}
	for _, selected := range v.advancedChecked {
		if selected {
			count++
		}
	}
	for _, selected := range v.fixChecked {
		if selected {
			count++
		}
	}
	for _, selected := range v.legacyChecked {
		if selected {
			count++
		}
	}
	for _, selected := range v.featureChecked {
		if selected {
			count++
		}
	}
	for _, selected := range v.remoteChecked {
		if selected {
			count++
		}
	}
	for _, selected := range v.profileChecked {
		if selected {
			count++
		}
	}
	for _, toggle := range v.toggles {
		if toggle.Dirty {
			count++
		}
	}
	return count
}

func (v TweakView) collectScripts() []string {
	var scripts []string
	scripts = appendCheckedScripts(scripts, EssentialTweaks, v.essentialChecked)
	scripts = appendCheckedScripts(scripts, AdvancedTweaks, v.advancedChecked)
	scripts = appendCheckedScripts(scripts, FixTweaks, v.fixChecked)
	scripts = appendCheckedScripts(scripts, LegacyPanels, v.legacyChecked)
	scripts = appendCheckedScripts(scripts, FeatureTweaks, v.featureChecked)
	scripts = appendCheckedScripts(scripts, RemoteAccessTweaks, v.remoteChecked)
	scripts = appendCheckedScripts(scripts, PowerShellProfileTweaks, v.profileChecked)

	for _, toggle := range v.toggles {
		if !toggle.Dirty {
			continue
		}
		if toggle.Enabled {
			scripts = append(scripts, toggle.EnableScript)
		} else {
			scripts = append(scripts, toggle.DisableScript)
		}
	}

	return scripts
}

func appendCheckedScripts(dst []string, items []Tweak, checked []bool) []string {
	for i, selected := range checked {
		if selected {
			dst = append(dst, items[i].Script)
		}
	}
	return dst
}

func (v TweakView) Footer() string {
	hint := "[Tab/Shift+Tab] Switch section  [Up/Down] Navigate  [Space] Toggle/select  [F5] Apply  [u] Add Ultimate Performance  [r] Remove"
	if selected := v.SelectedCount(); selected > 0 {
		hint += StyleCount.Render(fmt.Sprintf("  %d actions queued", selected))
	}
	return StyleFooter.Render(hint)
}

func (v TweakView) View(contentH int) string {
	halfW := v.width/2 - 2
	if halfW < 36 {
		halfW = 36
	}

	leftLines := v.renderActionColumn(halfW, contentH)
	rightLines := v.renderPreferenceColumn(halfW, contentH)

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

	rows := make([]string, 0, contentH)
	for i := 0; i < contentH; i++ {
		left := leftLines[i]
		right := rightLines[i]
		leftWidth := lipgloss.Width(left)
		if leftWidth < halfW {
			left += strings.Repeat(" ", halfW-leftWidth)
		}
		rows = append(rows, left+"  "+right)
	}
	return strings.Join(rows, "\n")
}

func (v TweakView) renderActionColumn(w, contentH int) []string {
	lines := []string{StyleSectionHeader.Render("Action Groups")}
	lines = append(lines, v.renderPanelMenu(leftPanels, v.leftPanel, w)...)
	lines = append(lines, "")
	lines = append(lines, v.renderPanelBody(v.leftPanel, w, contentH, 6)...)
	return lines
}

func (v TweakView) renderPreferenceColumn(w, contentH int) []string {
	lines := []string{StyleSectionHeader.Render("Preferences")}
	lines = append(lines, v.renderPanelMenu(rightPanels, v.rightPanel, w)...)
	lines = append(lines, "")
	lines = append(lines, v.renderPanelBody(v.rightPanel, w, contentH, 12)...)
	lines = append(lines, "")
	lines = append(lines, StyleSectionHeader.Render("Performance Plans"))
	lines = append(lines, "  "+StyleButton.Render(" Add and Activate Ultimate Performance ")+"  [u]")
	lines = append(lines, "  "+StyleButtonDanger.Render(" Remove Ultimate Performance ")+"  [r]")
	if len(v.log) > 0 {
		lines = append(lines, "")
		lines = append(lines, StyleSectionHeader.Render("Output"))
		start := len(v.log) - 4
		if start < 0 {
			start = 0
		}
		for _, line := range v.log[start:] {
			lines = append(lines, "  "+truncate(line, w-2))
		}
	}
	return lines
}
func (v TweakView) renderPanelMenu(panels []tweakPanel, current tweakPanel, w int) []string {
	lines := make([]string, 0, len(panels))
	for _, panel := range panels {
		lines = append(lines, v.renderPanelMenuRow(panel, current, w))
	}
	return lines
}

func (v TweakView) renderPanelMenuRow(panel, current tweakPanel, w int) string {
	prefix := "  "
	labelStyle := StyleAppNormal

	if panel == current {
		prefix = StyleInfo.Render("> ")
		labelStyle = StyleCatActive
	}
	if panel == v.activePanel {
		prefix = StyleInfo.Render("> ")
		labelStyle = StyleAppCursor
	}
	if panel == panelAdvanced {
		labelStyle = StyleWarn
	}

	count := v.panelSelectionCount(panel)
	title := panelTitle(panel)
	if count > 0 {
		title = fmt.Sprintf("%s (%d)", title, count)
	}
	return prefix + labelStyle.Render(truncate(title, w-2))
}

func (v TweakView) renderPanelBody(panel tweakPanel, w, contentH, reservedBottom int) []string {
	lines := []string{StyleSectionHeader.Render(panelTitle(panel))}
	if panel == panelAdvanced {
		lines[0] = StyleWarn.Render(panelTitle(panel))
	}

	visible := contentH - len(lines) - reservedBottom
	if visible < 5 {
		visible = 5
	}

	state := v.states[panel]
	length := v.panelListLen(panel)
	if state.cursor >= length && length > 0 {
		state.cursor = length - 1
	}
	if state.offset > state.cursor {
		state.offset = state.cursor
	}
	maxOffset := length - visible
	if maxOffset < 0 {
		maxOffset = 0
	}
	if state.offset > maxOffset {
		state.offset = maxOffset
	}

	for i := 0; i < visible && state.offset+i < length; i++ {
		idx := state.offset + i
		lines = append(lines, v.renderPanelRow(panel, idx, w))
	}

	if length == 0 {
		lines = append(lines, StyleTextMuted("  No items in this section."))
	}

	lines = append(lines, "")
	desc := v.panelItemDesc(panel, state.cursor)
	if desc != "" {
		lines = append(lines, StyleAppNote.Render("  "+truncate(desc, w-2)))
	}
	return lines
}

func (v TweakView) renderPanelRow(panel tweakPanel, idx, w int) string {
	isCursor := panel == v.activePanel && v.states[panel].cursor == idx
	prefix := "  "
	if isCursor {
		prefix = StyleInfo.Render("> ")
	}

	switch panel {
	case panelToggles:
		item := v.toggles[idx]
		label := truncate(item.Name, w-12)
		if item.Dirty {
			label += StyleTextMuted(" *")
		}
		if isCursor {
			label = StyleAppCursor.Render(label)
		} else {
			label = StyleToggleLabel.Render(label)
		}
		return fmt.Sprintf("%s%s %s", prefix, ToggleStr(item.Enabled), label)
	default:
		checked := v.panelChecked(panel, idx)
		label := truncate(v.panelItemName(panel, idx), w-6)
		switch {
		case isCursor:
			label = StyleAppCursor.Render(label)
		case checked:
			label = StyleAppSelected.Render(label)
		default:
			label = StyleAppNormal.Render(label)
		}
		return fmt.Sprintf("%s%s %s", prefix, CheckboxStr(checked), label)
	}
}

func (v TweakView) panelSelectionCount(panel tweakPanel) int {
	switch panel {
	case panelEssential:
		return countChecked(v.essentialChecked)
	case panelAdvanced:
		return countChecked(v.advancedChecked)
	case panelFixes:
		return countChecked(v.fixChecked)
	case panelLegacy:
		return countChecked(v.legacyChecked)
	case panelFeatures:
		return countChecked(v.featureChecked)
	case panelRemoteAccess:
		return countChecked(v.remoteChecked)
	case panelPowerShellProfile:
		return countChecked(v.profileChecked)
	case panelToggles:
		count := 0
		for _, toggle := range v.toggles {
			if toggle.Dirty {
				count++
			}
		}
		return count
	default:
		return 0
	}
}

func countChecked(values []bool) int {
	count := 0
	for _, value := range values {
		if value {
			count++
		}
	}
	return count
}

func (v TweakView) panelChecked(panel tweakPanel, idx int) bool {
	switch panel {
	case panelEssential:
		return v.essentialChecked[idx]
	case panelAdvanced:
		return v.advancedChecked[idx]
	case panelFixes:
		return v.fixChecked[idx]
	case panelLegacy:
		return v.legacyChecked[idx]
	case panelFeatures:
		return v.featureChecked[idx]
	case panelRemoteAccess:
		return v.remoteChecked[idx]
	case panelPowerShellProfile:
		return v.profileChecked[idx]
	default:
		return false
	}
}

func (v TweakView) panelItemName(panel tweakPanel, idx int) string {
	switch panel {
	case panelEssential:
		return EssentialTweaks[idx].Name
	case panelAdvanced:
		return AdvancedTweaks[idx].Name
	case panelFixes:
		return FixTweaks[idx].Name
	case panelLegacy:
		return LegacyPanels[idx].Name
	case panelFeatures:
		return FeatureTweaks[idx].Name
	case panelRemoteAccess:
		return RemoteAccessTweaks[idx].Name
	case panelPowerShellProfile:
		return PowerShellProfileTweaks[idx].Name
	case panelToggles:
		return v.toggles[idx].Name
	default:
		return ""
	}
}

func (v TweakView) panelItemDesc(panel tweakPanel, idx int) string {
	if idx < 0 {
		return ""
	}
	if panel == panelToggles {
		if idx < len(v.toggles) {
			return v.toggles[idx].Desc
		}
		return ""
	}

	items := v.panelTweaks(panel)
	if idx >= len(items) {
		return ""
	}
	return items[idx].Desc
}

func (v TweakView) panelTweaks(panel tweakPanel) []Tweak {
	switch panel {
	case panelEssential:
		return EssentialTweaks
	case panelAdvanced:
		return AdvancedTweaks
	case panelFixes:
		return FixTweaks
	case panelLegacy:
		return LegacyPanels
	case panelFeatures:
		return FeatureTweaks
	case panelRemoteAccess:
		return RemoteAccessTweaks
	case panelPowerShellProfile:
		return PowerShellProfileTweaks
	default:
		return nil
	}
}

func panelTitle(panel tweakPanel) string {
	switch panel {
	case panelEssential:
		return "Essential Tweaks"
	case panelAdvanced:
		return "Advanced Tweaks - Caution"
	case panelFixes:
		return "Fixes"
	case panelLegacy:
		return "Legacy Windows Panels"
	case panelFeatures:
		return "Features"
	case panelRemoteAccess:
		return "Remote Access"
	case panelPowerShellProfile:
		return "PowerShell Profile"
	case panelToggles:
		return "Customize Preferences"
	default:
		return "Tweaks"
	}
}

func (v TweakView) isRightPanel(panel tweakPanel) bool {
	switch panel {
	case panelFeatures, panelRemoteAccess, panelPowerShellProfile, panelToggles:
		return true
	default:
		return false
	}
}
