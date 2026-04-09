$ids = @(
  "Microsoft.VisualStudioCode", "VSCodium.VSCodium", "wez.wezterm", "Yarn.Yarn",
  "Starship.Starship", "SublimeHQ.SublimeText.4", "HashiCorp.Vagrant",
  "Unity.UnityHub", "Rustlang.Rustup", "RubyInstallerTeam.Ruby",
  "Tailscale.Tailscale", "CodeSector.TeraCopy", "TeamViewer.TeamViewer",
  "TeamSpeakSystems.TeamSpeakClient", "Syncthing.Syncthing", "Eugeny.Tabby",
  "stefansundin.SuperF4", "uderzo.SpaceSniffer", "WhirlwindFX.SignalRgb",
  "Rufus.Rufus", "BurntSushi.ripgrep.MSVC", "Sandboxie.Plus",
  "RARLab.WinRAR", "AntibodySoftware.WizTree", "XnSoft.XnView",
  "ZeroTier.ZeroTier", "Microsoft.Sysinternals.ZoomIt", "ajeetdsouza.zoxide",
  "TextExpander.TextExpander", "Rakuten.Viber", "Cisco.ClamAV",
  "OBSProject.OBSStudio", "Microsoft.PowerToys", "SyncTrayzor.SyncTrayzor",
  "Spacedrive.Spacedrive", "OpenWhisperSystems.Signal", "SlackTechnologies.Slack",
  "Foxit.PhantomPDF", "Foxit.FoxitReader", "ChidiWilliams.Buzz",
  "angryziber.AngryIPScanner", "szTheory.exifcleaner", "nomic.gpt4all",
  "Hibbiki.Chromium", "Ablaze.Floorp", "LibreWolf.LibreWolf",
  "MullvadVPN.MullvadBrowser", "MoonchildProductions.PaleMoon",
  "Alex313031.Thorium.AVX2", "TorProject.TorBrowser",
  "eloston.ungoogled-chromium", "Vivaldi.Vivaldi", "Waterfox.Waterfox",
  "zen-browser.zen", "Beeper.Beeper", "Betterbird.Betterbird",
  "Ferdium.Ferdium", "HexChat.HexChat", "SFLinux.Jami",
  "BelledonneCommunications.Linphone", "Element.Element",
  "Telegram.Unigram", "Vencord.Vesktop", "Zulip.Zulip"
)
foreach ($id in $ids) {
  $r = (winget show --id $id --accept-source-agreements 2>&1 | Select-String "Found|No package found")
  Write-Host "$id`t$r"
}
