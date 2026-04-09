@(
  "RubyInstallerTeam.RubyWithDevKit", "RubyInstallerTeam.Ruby.3",
  "XnSoft.XnViewMP", "XnSoft.XnView.Classic", "Tox.qTox",
  "Oxen.Session", "revoltchat.RevoltDesktop", "ZeroTier.ZeroTierOne",
  "j178.ChatGPT", "Zed.Zed", "Zed-Industries.Zed", "Logitech.GHUB",
  "JAMSoftware.TreeSize.Free", "xanderfrangos.twinkle-tray",
  "QL-Win.QuickLook", "BitSum.ProcessLasso", "namazso.OpenHashTab",
  "CalcProgrammer1.OpenRGB", "Open-Shell.Open-Shell-Menu",
  "SoftFever.OrcaSlicer", "Oracle.VirtualBox", "ownCloud.ownCloudDesktop",
  "Parsec.Parsec", "Giorgiotani.Peazip", "RaspberryPiFoundation.RaspberryPiImager",
  "Prusa3D.PrusaSlicer", "qBittorrent.qBittorrent", "QL-Win.QuickLook",
  "Glarysoft.GlaryUtilities", "GlazeMDE.GlazeWM", "TechPowerUp.GPU-Z",
  "gerardog.gsudo", "REALiX.HWiNFO", "CPUID.HWMonitor",
  "AppWork.JDownloader", "KDE.KDEConnect", "KeePassXCTeam.KeePassXC",
  "BartoszCichecki.LenovoLegionToolkit", "HermannSchinagl.LinkShellExtension",
  "rocksdanister.LivelyWallpaper", "LocalSend.LocalSend",
  "agalwood.Motrix", "rcmaehl.MSEdgeRedirect", "Guru3D.Afterburner",
  "M2Team.NanaZip", "nepnep.neofetch-win", "Nextcloud.NextcloudDesktop",
  "nilesoft.shell", "Nushell.Nushell", "TechPowerUp.NVCleanstall",
  "xM4ddy.OFGB", "OPAutoClicker.OPAutoClicker", "Open-Shell.Open-Shell-Menu"
) | ForEach-Object {
  $r = (winget show --id $_ --accept-source-agreements 2>&1 | Select-String "Found|No package found")
  Write-Host "$_`t$r"
}
