package pkgs

// Category represents an app category.
type Category string

const (
	CategoryCompany      Category = "Company"
	CategoryBrowsers     Category = "Browsers"
	CategoryComms        Category = "Comms"
	CategoryDevEditors   Category = "Dev: Editors"
	CategoryDevVPN       Category = "Dev: VPN"
	CategoryDevTools     Category = "Dev: Tools"
	CategoryDevLang      Category = "Dev: Lang"
	CategoryDevCLI       Category = "Dev: CLI"
	CategoryDevAI        Category = "Dev: AI"
	CategoryDevTerminal  Category = "Dev: Terminal"
	CategoryDocument     Category = "Document"
	CategoryProTools     Category = "Pro Tools"
	CategoryUtilities    Category = "Utilities"
)

// AllCategories defines the display order.
var AllCategories = []Category{
	CategoryCompany,
	CategoryBrowsers,
	CategoryComms,
	CategoryDevEditors,
	CategoryDevVPN,
	CategoryDevTools,
	CategoryDevLang,
	CategoryDevCLI,
	CategoryDevAI,
	CategoryDevTerminal,
	CategoryDocument,
	CategoryProTools,
	CategoryUtilities,
}

// AppSource indicates how the app will be installed.
type AppSource string

const (
	SourceWinget   AppSource = "winget"
	SourceMSStore  AppSource = "msstore"  // winget --source msstore
	SourceManual   AppSource = "manual"   // download manually; WingetID holds URL hint
	SourceLinux    AppSource = "linux"    // Linux-only; winget not applicable
	SourceNPM      AppSource = "npm"      // installed via npm -g
)

// App represents an installable application.
type App struct {
	Name     string
	WingetID string    // winget package ID, MS Store product ID, URL hint, or npm pkg name
	Category Category
	Source   AppSource // defaults to SourceWinget
	Note     string
}

// Catalog is the full app list.
var Catalog = []App{

	// ─── Company Apps ────────────────────────────────────────────────────────────
	{Name: "Microsoft Office", WingetID: "Microsoft.Office", Category: CategoryCompany, Note: "Word, Excel, PowerPoint, Outlook & more"},
	{Name: "Microsoft Teams", WingetID: "Microsoft.Teams", Category: CategoryCompany, Note: "Team collaboration and video calls"},
	{Name: "OneDrive", WingetID: "Microsoft.OneDrive", Category: CategoryCompany, Note: "Microsoft cloud storage client"},
	{Name: "Bitrix24", WingetID: "https://www.bitrix24.com/download/", Category: CategoryCompany, Source: SourceManual, Note: "CRM and project management (manual download)"},
	{Name: "TextExpander", WingetID: "TextExpander.TextExpander", Category: CategoryCompany, Note: "Text snippet expansion tool"},
	{Name: "WhatsApp", WingetID: "9NKSQGP7F2NH", Category: CategoryCompany, Source: SourceMSStore, Note: "WhatsApp desktop (MS Store)"},
	{Name: "Line", WingetID: "https://line.me/en/download", Category: CategoryCompany, Source: SourceManual, Note: "Line messaging app (manual download)"},

	// ─── Browsers ────────────────────────────────────────────────────────────────
	{Name: "Firefox", WingetID: "Mozilla.Firefox", Category: CategoryBrowsers, Note: "Mozilla Firefox web browser"},
	{Name: "Chromium", WingetID: "Hibbiki.Chromium", Category: CategoryBrowsers, Note: "Open-source Chromium build"},
	{Name: "Brave", WingetID: "Brave.Brave", Category: CategoryBrowsers, Note: "Privacy-focused Chromium browser"},
	{Name: "LibreWolf", WingetID: "LibreWolf.LibreWolf", Category: CategoryBrowsers, Note: "Privacy-hardened Firefox fork"},
	{Name: "Waterfox", WingetID: "Waterfox.Waterfox", Category: CategoryBrowsers, Note: "Privacy-focused Firefox-based browser"},
	{Name: "Tor Browser", WingetID: "TorProject.TorBrowser", Category: CategoryBrowsers, Note: "Anonymous browsing via the Tor network"},
	{Name: "Google Chrome", WingetID: "Google.Chrome", Category: CategoryBrowsers, Note: "Google web browser"},
	{Name: "Zen Browser", WingetID: "zen-browser.zen", Category: CategoryBrowsers, Note: "Minimalist Firefox-based browser"},
	{Name: "Helium", WingetID: "Alex313031.Helium", Category: CategoryBrowsers, Note: "Chromium build optimised for older hardware"},
	{Name: "Vivaldi", WingetID: "Vivaldi.Vivaldi", Category: CategoryBrowsers, Note: "Feature-rich customisable browser"},
	{Name: "Ungoogled Chromium", WingetID: "eloston.ungoogled-chromium", Category: CategoryBrowsers, Note: "Chromium without any Google integration"},
	{Name: "Konqueror", WingetID: "https://apps.kde.org/konqueror/", Category: CategoryBrowsers, Source: SourceLinux, Note: "KDE file manager / browser (Linux)"},
	{Name: "Falkon", WingetID: "KDE.Falkon", Category: CategoryBrowsers, Note: "Lightweight KDE QtWebEngine browser"},
	{Name: "GNOME Web", WingetID: "https://apps.gnome.org/Epiphany/", Category: CategoryBrowsers, Source: SourceLinux, Note: "GNOME's Epiphany browser (Linux)"},
	{Name: "qutebrowser", WingetID: "qutebrowser.qutebrowser", Category: CategoryBrowsers, Note: "Keyboard-driven Vim-like browser"},
	{Name: "Nyxt", WingetID: "https://nyxt.atlas.engineer/", Category: CategoryBrowsers, Source: SourceManual, Note: "Keyboard-optimised hacker browser (manual download)"},
	{Name: "Floorp", WingetID: "Ablaze.Floorp", Category: CategoryBrowsers, Note: "Privacy Firefox fork with extra features"},
	{Name: "Mullvad Browser", WingetID: "MullvadVPN.MullvadBrowser", Category: CategoryBrowsers, Note: "Privacy browser by Mullvad VPN"},
	{Name: "Microsoft Edge", WingetID: "Microsoft.Edge", Category: CategoryBrowsers, Note: "Microsoft Edge (Chromium-based)"},
	{Name: "Opera", WingetID: "Opera.Opera", Category: CategoryBrowsers, Note: "Opera browser with built-in VPN & AI"},

	// ─── Communications ──────────────────────────────────────────────────────────
	{Name: "Slack", WingetID: "SlackTechnologies.Slack", Category: CategoryComms, Note: "Team messaging platform"},
	{Name: "Signal", WingetID: "OpenWhisperSystems.Signal", Category: CategoryComms, Note: "Secure end-to-end encrypted messaging"},
	{Name: "Viber", WingetID: "Rakuten.Viber", Category: CategoryComms, Note: "Viber messaging and calls"},
	{Name: "Telegram", WingetID: "Telegram.TelegramDesktop", Category: CategoryComms, Note: "Fast cloud-based messaging app"},
	{Name: "Discord", WingetID: "Discord.Discord", Category: CategoryComms, Note: "Gaming and community voice/text chat"},
	{Name: "Zoom", WingetID: "Zoom.Zoom", Category: CategoryComms, Note: "Video conferencing platform"},
	{Name: "Element", WingetID: "Element.Element", Category: CategoryComms, Note: "Matrix decentralised end-to-end chat"},
	{Name: "Beeper", WingetID: "Beeper.Beeper", Category: CategoryComms, Note: "Universal messaging app (all in one)"},
	{Name: "Ferdium", WingetID: "Ferdium.Ferdium", Category: CategoryComms, Note: "All-in-one messaging aggregator"},
	{Name: "Vesktop", WingetID: "Vencord.Vesktop", Category: CategoryComms, Note: "Enhanced Discord client (Vencord)"},
	{Name: "Proton Mail", WingetID: "Proton.ProtonMail", Category: CategoryComms, Note: "End-to-end encrypted email client"},
	{Name: "Thunderbird", WingetID: "Mozilla.Thunderbird", Category: CategoryComms, Note: "Mozilla open-source email client"},
	{Name: "Betterbird", WingetID: "Betterbird.Betterbird", Category: CategoryComms, Note: "Enhanced Thunderbird fork for email"},
	{Name: "Skype", WingetID: "Microsoft.Skype", Category: CategoryComms, Note: "Microsoft Skype calls and messaging"},
	{Name: "HexChat", WingetID: "HexChat.HexChat", Category: CategoryComms, Note: "IRC chat client"},
	{Name: "Jami", WingetID: "SFLinux.Jami", Category: CategoryComms, Note: "Decentralised peer-to-peer communication"},
	{Name: "Zulip", WingetID: "Zulip.Zulip", Category: CategoryComms, Note: "Topic-based team messaging"},

	// ─── Dev: Editors ─────────────────────────────────────────────────────────────
	{Name: "VS Code", WingetID: "Microsoft.VisualStudioCode", Category: CategoryDevEditors, Note: "Lightweight extensible code editor"},
	{Name: "VSCodium", WingetID: "VSCodium.VSCodium", Category: CategoryDevEditors, Note: "VS Code without telemetry (open-source)"},
	{Name: "Cursor", WingetID: "Anysphere.Cursor", Category: CategoryDevEditors, Note: "AI-first code editor (GPT-4 built-in)"},
	{Name: "Neovim", WingetID: "Neovim.Neovim", Category: CategoryDevEditors, Note: "Hyperextensible Vim-based text editor"},
	{Name: "Helix", WingetID: "Helix.Helix", Category: CategoryDevEditors, Note: "Post-modern modal text editor"},
	{Name: "Sublime Text 4", WingetID: "SublimeHQ.SublimeText.4", Category: CategoryDevEditors, Note: "Sophisticated text and code editor"},
	{Name: "Zed", WingetID: "https://zed.dev/download", Category: CategoryDevEditors, Source: SourceManual, Note: "High-performance collaborative code editor (manual)"},
	{Name: "Visual Studio 2022", WingetID: "Microsoft.VisualStudio.2022.Community", Category: CategoryDevEditors, Note: "Microsoft's full-featured IDE (Community free)"},
	{Name: "JetBrains Toolbox", WingetID: "JetBrains.Toolbox", Category: CategoryDevEditors, Note: "Manage all JetBrains IDEs from one place"},
	{Name: "Emacs", WingetID: "GNU.Emacs", Category: CategoryDevEditors, Note: "Extensible, self-documenting text editor"},

	// ─── Dev: VPN & Network ───────────────────────────────────────────────────────
	{Name: "IVPN", WingetID: "IVPN.IVPN", Category: CategoryDevVPN, Note: "Privacy-first no-logs VPN service"},
	{Name: "Proton VPN", WingetID: "Proton.ProtonVPN", Category: CategoryDevVPN, Note: "Swiss-based secure VPN service"},
	{Name: "Mullvad VPN", WingetID: "MullvadVPN.MullvadVPN", Category: CategoryDevVPN, Note: "Privacy-first no-logs VPN service"},
	{Name: "Tailscale", WingetID: "Tailscale.Tailscale", Category: CategoryDevVPN, Note: "Zero-config mesh VPN (WireGuard-based)"},
	{Name: "WireGuard", WingetID: "WireGuard.WireGuard", Category: CategoryDevVPN, Note: "Modern fast VPN protocol client"},
	{Name: "OpenVPN", WingetID: "OpenVPNTechnologies.OpenVPNConnect", Category: CategoryDevVPN, Note: "Official OpenVPN client"},
	{Name: "ZeroTier", WingetID: "ZeroTier.ZeroTierOne", Category: CategoryDevVPN, Note: "Software-defined networking platform"},
	{Name: "Netbird", WingetID: "Netbird.Netbird", Category: CategoryDevVPN, Note: "WireGuard-based zero-config mesh VPN"},
	{Name: "Nmap", WingetID: "Insecure.Nmap", Category: CategoryDevVPN, Note: "Network exploration and security auditing"},
	{Name: "OpenSSH", WingetID: "Microsoft.OpenSSH.Beta", Category: CategoryDevVPN, Note: "OpenSSH client and server for Windows"},
	{Name: "PuTTY", WingetID: "PuTTY.PuTTY", Category: CategoryDevVPN, Note: "Classic SSH, Telnet and serial client"},
	{Name: "WinSCP", WingetID: "WinSCP.WinSCP", Category: CategoryDevVPN, Note: "SFTP, FTP and SCP client"},
	{Name: "Remmina", WingetID: "https://remmina.org/", Category: CategoryDevVPN, Source: SourceLinux, Note: "Remote desktop client — RDP/VNC/SSH (Linux)"},

	// ─── Dev: Tools ───────────────────────────────────────────────────────────────
	{Name: "Git", WingetID: "Git.Git", Category: CategoryDevTools, Note: "Distributed version control system"},
	{Name: "Git LFS", WingetID: "GitHub.GitLFS", Category: CategoryDevTools, Note: "Git extension for large file storage"},
	{Name: "LazyGit", WingetID: "JesseDuffield.lazygit", Category: CategoryDevTools, Note: "Simple terminal UI for git commands"},
	{Name: "GitHub CLI", WingetID: "GitHub.cli", Category: CategoryDevTools, Note: "GitHub from the command line"},
	{Name: "Docker Desktop", WingetID: "Docker.DockerDesktop", Category: CategoryDevTools, Note: "Container development platform"},
	{Name: "Podman", WingetID: "RedHat.Podman", Category: CategoryDevTools, Note: "Daemonless container engine by Red Hat"},
	{Name: "Podman Desktop", WingetID: "RedHat.Podman-Desktop", Category: CategoryDevTools, Note: "GUI for Podman containers"},
	{Name: "kubectl", WingetID: "Kubernetes.kubectl", Category: CategoryDevTools, Note: "Kubernetes command-line tool"},
	{Name: "Vagrant", WingetID: "Hashicorp.Vagrant", Category: CategoryDevTools, Note: "Virtual machine lifecycle management"},
	{Name: "VirtualBox", WingetID: "Oracle.VirtualBox", Category: CategoryDevTools, Note: "Free x86 virtualisation platform"},
	{Name: "Incus", WingetID: "https://linuxcontainers.org/incus/", Category: CategoryDevTools, Source: SourceLinux, Note: "System container and VM manager (Linux)"},
	{Name: "GNOME Boxes", WingetID: "https://apps.gnome.org/Boxes/", Category: CategoryDevTools, Source: SourceLinux, Note: "GNOME VM manager (Linux)"},
	{Name: "Virt Manager", WingetID: "https://virt-manager.org/", Category: CategoryDevTools, Source: SourceLinux, Note: "QEMU/KVM virtual machine manager (Linux)"},
	{Name: "DBeaver", WingetID: "dbeaver.dbeaver", Category: CategoryDevTools, Note: "Universal database GUI tool"},
	{Name: "Meld", WingetID: "Meld.Meld", Category: CategoryDevTools, Note: "Visual diff and merge tool"},
	{Name: "Wireshark", WingetID: "WiresharkFoundation.Wireshark", Category: CategoryDevTools, Note: "Network protocol analyser"},
	{Name: "Postman", WingetID: "Postman.Postman", Category: CategoryDevTools, Note: "API development and testing platform"},
	{Name: "Bruno", WingetID: "Bruno.Bruno", Category: CategoryDevTools, Note: "Offline-first open-source API client"},
	{Name: "Hoppscotch", WingetID: "Hoppscotch.Hoppscotch", Category: CategoryDevTools, Note: "Open-source API development ecosystem"},
	{Name: "Yaak", WingetID: "yaak-app.yaak", Category: CategoryDevTools, Note: "Offline API client (REST/GraphQL/gRPC)"},
	{Name: "ImHex", WingetID: "WerWolv.ImHex", Category: CategoryDevTools, Note: "Feature-rich hex editor for reverse engineering"},
	{Name: "CMake", WingetID: "Kitware.CMake", Category: CategoryDevTools, Note: "Cross-platform build system generator"},
	{Name: "Sublime Merge", WingetID: "SublimeHQ.SublimeMerge", Category: CategoryDevTools, Note: "Fast Git client by Sublime HQ"},
	{Name: "GitHub Desktop", WingetID: "GitHub.GitHubDesktop", Category: CategoryDevTools, Note: "GitHub GUI client"},

	// ─── Dev: Languages ───────────────────────────────────────────────────────────
	{Name: "Python 3", WingetID: "Python.Python.3.13", Category: CategoryDevLang, Note: "Python 3 programming language"},
	{Name: "Node.js LTS", WingetID: "OpenJS.NodeJS.LTS", Category: CategoryDevLang, Note: "Node.js JavaScript runtime (LTS)"},
	{Name: "Go", WingetID: "GoLang.Go", Category: CategoryDevLang, Note: "Go programming language toolchain"},
	{Name: "Rust", WingetID: "Rustlang.Rustup", Category: CategoryDevLang, Note: "Rust systems programming language via rustup"},
	{Name: "Ruby", WingetID: "RubyInstallerTeam.Ruby", Category: CategoryDevLang, Note: "Ruby programming language (RubyInstaller)"},
	{Name: "PHP", WingetID: "PHP.PHP", Category: CategoryDevLang, Note: "PHP scripting language"},
	{Name: "OpenJDK 21 (LTS)", WingetID: "EclipseAdoptium.Temurin.21.JDK", Category: CategoryDevLang, Note: "OpenJDK 21 LTS by Eclipse Adoptium"},
	{Name: "Deno", WingetID: "DenoLand.Deno", Category: CategoryDevLang, Note: "Secure modern JavaScript/TypeScript runtime"},
	{Name: "Bun", WingetID: "Oven-sh.Bun", Category: CategoryDevLang, Note: "All-in-one fast JavaScript runtime & toolkit"},
	{Name: "pnpm", WingetID: "pnpm.pnpm", Category: CategoryDevLang, Note: "Fast, disk-efficient Node.js package manager"},
	{Name: "Yarn", WingetID: "Yarn.Yarn", Category: CategoryDevLang, Note: "Fast reliable Node.js package manager"},
	{Name: "uv", WingetID: "astral-sh.uv", Category: CategoryDevLang, Note: "Extremely fast Python package manager (Rust)"},

	// ─── Dev: CLI Tools ───────────────────────────────────────────────────────────
	{Name: "btop", WingetID: "aristocratos.btop4win", Category: CategoryDevCLI, Note: "Beautiful terminal resource monitor (Windows port)"},
	{Name: "htop", WingetID: "htop.htop", Category: CategoryDevCLI, Note: "Interactive process viewer"},
	{Name: "fastfetch", WingetID: "Fastfetch-cli.Fastfetch", Category: CategoryDevCLI, Note: "Fast neofetch-like system info display"},
	{Name: "eza", WingetID: "eza-community.eza", Category: CategoryDevCLI, Note: "Modern ls replacement with colour and icons"},
	{Name: "bat", WingetID: "sharkdp.bat", Category: CategoryDevCLI, Note: "Cat clone with syntax highlighting and Git"},
	{Name: "fzf", WingetID: "junegunn.fzf", Category: CategoryDevCLI, Note: "General-purpose command-line fuzzy finder"},
	{Name: "ripgrep", WingetID: "BurntSushi.ripgrep.MSVC", Category: CategoryDevCLI, Note: "Blazingly fast line-oriented grep tool"},
	{Name: "zoxide", WingetID: "ajeetdsouza.zoxide", Category: CategoryDevCLI, Note: "Smarter cd command that learns your habits"},
	{Name: "tldr", WingetID: "dbrgn.tealdeer", Category: CategoryDevCLI, Note: "Fast tldr client with offline cache (Rust)"},
	{Name: "wget", WingetID: "JernejSimoncic.Wget", Category: CategoryDevCLI, Note: "Non-interactive network downloader"},
	{Name: "curl", WingetID: "cURL.cURL", Category: CategoryDevCLI, Note: "Command-line tool for URL data transfer"},
	{Name: "aria2", WingetID: "aria2.aria2", Category: CategoryDevCLI, Note: "Lightweight multi-protocol download utility"},
	{Name: "yazi", WingetID: "sxyazi.yazi", Category: CategoryDevCLI, Note: "Blazing fast terminal file manager (Rust)"},
	{Name: "ranger", WingetID: "https://ranger.github.io/", Category: CategoryDevCLI, Source: SourceLinux, Note: "Terminal file manager with Vim keybindings (Linux/Mac)"},
	{Name: "fd", WingetID: "sharkdp.fd", Category: CategoryDevCLI, Note: "Simple, fast alternative to find"},
	{Name: "Zellij", WingetID: "Zellij.Zellij", Category: CategoryDevCLI, Note: "Terminal multiplexer with layouts and plugins"},
	{Name: "tmux", WingetID: "https://github.com/tmux/tmux", Category: CategoryDevCLI, Source: SourceLinux, Note: "Terminal multiplexer (Linux/Mac — use Zellij on Windows)"},
	{Name: "Superfile", WingetID: "yorukot.superfile", Category: CategoryDevCLI, Note: "Pretty fancy terminal file manager"},
	{Name: "Nushell", WingetID: "Nushell.Nushell", Category: CategoryDevCLI, Note: "Modern structured-data shell"},
	{Name: "rsync", WingetID: "cwRsync.cwRsync", Category: CategoryDevCLI, Note: "Fast file sync/transfer (cwRsync Windows port)"},

	// ─── Dev: AI Tools ────────────────────────────────────────────────────────────
	{Name: "Claude Code", WingetID: "Anthropic.ClaudeCode", Category: CategoryDevAI, Note: "Anthropic Claude AI coding agent for terminal"},
	{Name: "Claude Desktop", WingetID: "Anthropic.Claude", Category: CategoryDevAI, Note: "Anthropic Claude AI desktop assistant"},
	{Name: "Ollama", WingetID: "Ollama.Ollama", Category: CategoryDevAI, Note: "Run large language models locally"},
	{Name: "Jan", WingetID: "Jan.Jan", Category: CategoryDevAI, Note: "Open-source offline ChatGPT alternative"},
	{Name: "LM Studio", WingetID: "ElementLabs.LMStudio", Category: CategoryDevAI, Note: "Discover and run local LLMs with GUI"},
	{Name: "GPT4All", WingetID: "nomic.gpt4all", Category: CategoryDevAI, Note: "Run LLMs locally — privacy-first"},
	{Name: "OpenCode", WingetID: "opencode-ai.opencode", Category: CategoryDevAI, Note: "AI coding agent for the terminal"},
	{Name: "OpenAI Codex CLI", WingetID: "@openai/codex", Category: CategoryDevAI, Source: SourceNPM, Note: "OpenAI Codex CLI (npm install -g @openai/codex)"},
	{Name: "Gemini CLI", WingetID: "@google/gemini-cli", Category: CategoryDevAI, Source: SourceNPM, Note: "Google Gemini CLI (npm install -g @google/gemini-cli)"},
	{Name: "ChatGPT", WingetID: "j178.ChatGPT", Category: CategoryDevAI, Note: "ChatGPT desktop client"},
	{Name: "Cursor", WingetID: "Anysphere.Cursor", Category: CategoryDevAI, Note: "AI-first code editor (GPT-4 built-in)"},
	{Name: "Buzz", WingetID: "ChidiWilliams.Buzz", Category: CategoryDevAI, Note: "Offline audio transcription via Whisper"},
	{Name: "DeepL", WingetID: "DeepL.DeepL", Category: CategoryDevAI, Note: "AI-powered translation tool"},

	// ─── Dev: Terminal ────────────────────────────────────────────────────────────
	{Name: "Windows Terminal", WingetID: "Microsoft.WindowsTerminal", Category: CategoryDevTerminal, Note: "Modern tabbed terminal by Microsoft"},
	{Name: "WezTerm", WingetID: "wez.wezterm", Category: CategoryDevTerminal, Note: "GPU-accelerated terminal and multiplexer"},
	{Name: "Alacritty", WingetID: "Alacritty.Alacritty", Category: CategoryDevTerminal, Note: "GPU-accelerated minimal terminal emulator"},
	{Name: "Kitty", WingetID: "kovidgoyal.kitty", Category: CategoryDevTerminal, Note: "Fast, feature-rich GPU-based terminal"},
	{Name: "Ghostty", WingetID: "Ghostty.Ghostty", Category: CategoryDevTerminal, Note: "Fast, native terminal emulator"},
	{Name: "Tabby", WingetID: "Eugeny.Tabby", Category: CategoryDevTerminal, Note: "Modern terminal with SSH and serial support"},
	{Name: "Starship", WingetID: "Starship.Starship", Category: CategoryDevTerminal, Note: "Minimal, fast cross-shell prompt"},
	{Name: "Fish Shell", WingetID: "fish-shell.fish", Category: CategoryDevTerminal, Note: "User-friendly interactive shell"},
	{Name: "Zsh", WingetID: "https://www.zsh.org/", Category: CategoryDevTerminal, Source: SourceLinux, Note: "Z shell — use via WSL on Windows"},
	{Name: "Oh My Zsh", WingetID: "https://ohmyz.sh/", Category: CategoryDevTerminal, Source: SourceLinux, Note: "Zsh framework — use via WSL on Windows"},
	{Name: "Foot", WingetID: "https://codeberg.org/dnkl/foot", Category: CategoryDevTerminal, Source: SourceLinux, Note: "Wayland native terminal emulator (Linux)"},
	{Name: "Ptyxis", WingetID: "https://gitlab.gnome.org/chergert/ptyxis", Category: CategoryDevTerminal, Source: SourceLinux, Note: "GNOME container-aware terminal (Linux)"},

	// ─── Document ─────────────────────────────────────────────────────────────────
	{Name: "Adobe Acrobat Reader", WingetID: "Adobe.Acrobat.Reader.64-bit", Category: CategoryDocument, Note: "Industry-standard PDF viewer by Adobe"},
	{Name: "Foxit PDF Editor", WingetID: "Foxit.PhantomPDF", Category: CategoryDocument, Note: "Full-featured PDF editor by Foxit"},
	{Name: "Foxit PDF Reader", WingetID: "Foxit.FoxitReader", Category: CategoryDocument, Note: "Fast, lightweight PDF reader"},
	{Name: "Sumatra PDF", WingetID: "SumatraPDF.SumatraPDF", Category: CategoryDocument, Note: "Lightweight, fast PDF/eBook reader"},
	{Name: "Obsidian", WingetID: "Obsidian.Obsidian", Category: CategoryDocument, Note: "Markdown-based personal knowledge base"},
	{Name: "Joplin", WingetID: "Joplin.Joplin", Category: CategoryDocument, Note: "Open-source note-taking and to-do app"},
	{Name: "Logseq", WingetID: "Logseq.Logseq", Category: CategoryDocument, Note: "Privacy-first knowledge base (outliner)"},
	{Name: "Notion", WingetID: "Notion.Notion", Category: CategoryDocument, Note: "All-in-one workspace for notes and docs"},
	{Name: "LibreOffice", WingetID: "TheDocumentFoundation.LibreOffice", Category: CategoryDocument, Note: "Free open-source office suite"},
	{Name: "ONLYOffice Desktop", WingetID: "ONLYOFFICE.DesktopEditors", Category: CategoryDocument, Note: "Office suite with strong MS format compatibility"},
	{Name: "Notepad++", WingetID: "Notepad++.Notepad++", Category: CategoryDocument, Note: "Popular source code and text editor"},
	{Name: "Zotero", WingetID: "Zotero.Zotero", Category: CategoryDocument, Note: "Research reference manager"},
	{Name: "Calibre", WingetID: "calibre.calibre", Category: CategoryDocument, Note: "E-book management and conversion"},
	{Name: "WinMerge", WingetID: "WinMerge.WinMerge", Category: CategoryDocument, Note: "Visual file and folder diff/merge"},
	{Name: "Xournal++", WingetID: "Xournal++.Xournal++", Category: CategoryDocument, Note: "Handwriting notes and PDF annotation"},
	{Name: "PDF24 Creator", WingetID: "geeksoftwareGmbH.PDF24Creator", Category: CategoryDocument, Note: "Free PDF creation and editing toolkit"},
	{Name: "Anki", WingetID: "Anki.Anki", Category: CategoryDocument, Note: "Spaced-repetition flashcard learning"},

	// ─── Pro Tools ───────────────────────────────────────────────────────────────
	{Name: "Advanced IP Scanner", WingetID: "Famatech.AdvancedIPScanner", Category: CategoryProTools, Note: "Scan and manage network devices"},
	{Name: "Angry IP Scanner", WingetID: "angryziber.AngryIPScanner", Category: CategoryProTools, Note: "Fast open-source network scanner"},
	{Name: "Portmaster", WingetID: "Safing.Portmaster", Category: CategoryProTools, Note: "Open-source application firewall"},
	{Name: "RustDesk", WingetID: "RustDesk.RustDesk", Category: CategoryProTools, Note: "Open-source remote desktop software"},
	{Name: "AnyDesk", WingetID: "AnyDesk.AnyDesk", Category: CategoryProTools, Note: "Fast remote desktop software"},
	{Name: "mRemoteNG", WingetID: "mRemoteNG.mRemoteNG", Category: CategoryProTools, Note: "Multi-protocol remote connections manager"},
	{Name: "HeidiSQL", WingetID: "HeidiSQL.HeidiSQL", Category: CategoryProTools, Note: "Free database management GUI"},
	{Name: "Ventoy", WingetID: "Ventoy.Ventoy", Category: CategoryProTools, Note: "Multi-ISO bootable USB creation tool"},
	{Name: "XPipe", WingetID: "xpipe-io.xpipe", Category: CategoryProTools, Note: "Connection manager and shell hub"},
	{Name: "Malwarebytes", WingetID: "Malwarebytes.Malwarebytes", Category: CategoryProTools, Note: "Anti-malware and ransomware protection"},
	{Name: "ClamAV", WingetID: "Cisco.ClamAV", Category: CategoryProTools, Note: "Open-source antivirus by Cisco"},
	{Name: "Simplewall", WingetID: "henrypp.simplewall", Category: CategoryProTools, Note: "Simple Windows Filtering Platform firewall"},
	{Name: "EFI Boot Editor", WingetID: "EFIBootEditor.EFIBootEditor", Category: CategoryProTools, Note: "Manage UEFI/EFI boot entries"},
	{Name: "Sandboxie Plus", WingetID: "Sandboxie.Plus", Category: CategoryProTools, Note: "Run apps in isolated sandbox environment"},
	{Name: "OpenRGB", WingetID: "CalcProgrammer1.OpenRGB", Category: CategoryProTools, Note: "Open-source RGB lighting control"},
	{Name: "SignalRGB", WingetID: "WhirlwindFX.SignalRgb", Category: CategoryProTools, Note: "Unified RGB lighting ecosystem"},
	{Name: "MSI Afterburner", WingetID: "Guru3D.Afterburner", Category: CategoryProTools, Note: "GPU overclocking and monitoring utility"},
	{Name: "HWiNFO", WingetID: "REALiX.HWiNFO", Category: CategoryProTools, Note: "Comprehensive hardware analysis and monitoring"},
	{Name: "GPU-Z", WingetID: "TechPowerUp.GPU-Z", Category: CategoryProTools, Note: "Detailed GPU and video card information"},
	{Name: "CPU-Z", WingetID: "CPUID.CPU-Z", Category: CategoryProTools, Note: "Detailed CPU and system information"},
	{Name: "OBS Studio", WingetID: "OBSProject.OBSStudio", Category: CategoryProTools, Note: "Free screen recording and live streaming"},
	{Name: "Raspberry Pi Imager", WingetID: "RaspberryPiFoundation.RaspberryPiImager", Category: CategoryProTools, Note: "Write OS images to SD cards and USB drives"},
	{Name: "Rufus", WingetID: "Rufus.Rufus", Category: CategoryProTools, Note: "Create bootable USB drives from ISO images"},
	{Name: "Parsec", WingetID: "Parsec.Parsec", Category: CategoryProTools, Note: "Low-latency remote gaming and desktop access"},

	// ─── Utilities ────────────────────────────────────────────────────────────────
	{Name: "7-Zip", WingetID: "7zip.7zip", Category: CategoryUtilities, Note: "Free open-source file archiver"},
	{Name: "NanaZip", WingetID: "M2Team.NanaZip", Category: CategoryUtilities, Note: "Modern Windows archive manager (7-Zip fork)"},
	{Name: "PeaZip", WingetID: "Giorgiotani.Peazip", Category: CategoryUtilities, Note: "Free open-source file archiver"},
	{Name: "WinRAR", WingetID: "RARLab.WinRAR", Category: CategoryUtilities, Note: "Popular archive manager (shareware)"},
	{Name: "1Password", WingetID: "AgileBits.1Password", Category: CategoryUtilities, Note: "Premium password manager"},
	{Name: "Bitwarden", WingetID: "Bitwarden.Bitwarden", Category: CategoryUtilities, Note: "Open-source password manager"},
	{Name: "KeePassXC", WingetID: "KeePassXCTeam.KeePassXC", Category: CategoryUtilities, Note: "Cross-platform open-source password manager"},
	{Name: "Ente Auth", WingetID: "ente-io.auth-desktop", Category: CategoryUtilities, Note: "Open-source 2FA authenticator app"},
	{Name: "Everything Search", WingetID: "voidtools.Everything", Category: CategoryUtilities, Note: "Instant filename search across all drives"},
	{Name: "Flow Launcher", WingetID: "Flow-Launcher.Flow-Launcher", Category: CategoryUtilities, Note: "Windows app and file launcher"},
	{Name: "PowerToys", WingetID: "Microsoft.PowerToys", Category: CategoryUtilities, Note: "Microsoft utilities for power users"},
	{Name: "AutoHotkey", WingetID: "AutoHotkey.AutoHotkey", Category: CategoryUtilities, Note: "Windows automation scripting language"},
	{Name: "Espanso", WingetID: "Espanso.Espanso", Category: CategoryUtilities, Note: "Cross-platform text expander"},
	{Name: "LocalSend", WingetID: "LocalSend.LocalSend", Category: CategoryUtilities, Note: "AirDrop-like LAN file sharing (open-source)"},
	{Name: "KDE Connect", WingetID: "KDE.KDEConnect", Category: CategoryUtilities, Note: "Connect Android/iOS phone to PC"},
	{Name: "Nextcloud Desktop", WingetID: "Nextcloud.NextcloudDesktop", Category: CategoryUtilities, Note: "Self-hosted cloud file sync client"},
	{Name: "Google Drive", WingetID: "Google.GoogleDrive", Category: CategoryUtilities, Note: "Google cloud storage client"},
	{Name: "Dropbox", WingetID: "Dropbox.Dropbox", Category: CategoryUtilities, Note: "Cloud file storage and sync"},
	{Name: "Syncthing", WingetID: "Syncthing.Syncthing", Category: CategoryUtilities, Note: "Peer-to-peer continuous file sync"},
	{Name: "WizTree", WingetID: "AntibodySoftware.WizTree", Category: CategoryUtilities, Note: "Fastest disk space analyser for Windows"},
	{Name: "SpaceSniffer", WingetID: "uderzo.SpaceSniffer", Category: CategoryUtilities, Note: "Visual disk space usage browser"},
	{Name: "TreeSize Free", WingetID: "JAMSoftware.TreeSize.Free", Category: CategoryUtilities, Note: "Disk space manager and visualiser"},
	{Name: "CopyQ", WingetID: "hluk.CopyQ", Category: CategoryUtilities, Note: "Clipboard manager with scripting support"},
	{Name: "Ditto", WingetID: "Ditto.Ditto", Category: CategoryUtilities, Note: "Powerful clipboard history manager"},
	{Name: "Auto Dark Mode", WingetID: "Armin2208.WindowsAutoNightMode", Category: CategoryUtilities, Note: "Automatically switch Windows light/dark theme"},
	{Name: "F.lux", WingetID: "flux.flux", Category: CategoryUtilities, Note: "Adaptive screen colour temperature by time"},
	{Name: "Twinkle Tray", WingetID: "xanderfrangos.twinkle-tray", Category: CategoryUtilities, Note: "Monitor brightness control from system tray"},
	{Name: "FanControl", WingetID: "Rem0o.FanControl", Category: CategoryUtilities, Note: "PC fan speed control with curves"},
	{Name: "QuickLook", WingetID: "QL-Win.QuickLook", Category: CategoryUtilities, Note: "Instant file preview via Spacebar (macOS-like)"},
	{Name: "Lively Wallpaper", WingetID: "rocksdanister.LivelyWallpaper", Category: CategoryUtilities, Note: "Animated and interactive live wallpapers"},
	{Name: "Open-Shell", WingetID: "Open-Shell.Open-Shell-Menu", Category: CategoryUtilities, Note: "Customise the Windows Start menu"},
	{Name: "GlazeWM", WingetID: "GlazeMDE.GlazeWM", Category: CategoryUtilities, Note: "Tiling window manager for Windows (i3-like)"},
	{Name: "Gsudo", WingetID: "gerardog.gsudo", Category: CategoryUtilities, Note: "sudo for Windows — elevate any command"},
	{Name: "Process Lasso", WingetID: "BitSum.ProcessLasso", Category: CategoryUtilities, Note: "Automated CPU process optimisation"},
	{Name: "BleachBit", WingetID: "BleachBit.BleachBit", Category: CategoryUtilities, Note: "System cleaner and privacy guard"},
	{Name: "Bulk Crap Uninstaller", WingetID: "Klocman.BulkCrapUninstaller", Category: CategoryUtilities, Note: "Bulk app uninstaller with deep clean"},
	{Name: "OBS Studio", WingetID: "OBSProject.OBSStudio", Category: CategoryUtilities, Note: "Free screen recording and live streaming"},
	{Name: "VLC", WingetID: "VideoLAN.VLC", Category: CategoryUtilities, Note: "Universal media player"},
	{Name: "Figma", WingetID: "Figma.Figma", Category: CategoryUtilities, Note: "Collaborative UI/UX design tool"},
	{Name: "Inkscape", WingetID: "Inkscape.Inkscape", Category: CategoryUtilities, Note: "Free vector graphics editor"},
	{Name: "GIMP", WingetID: "GIMP.GIMP", Category: CategoryUtilities, Note: "GNU Image Manipulation Program"},
	{Name: "qBittorrent", WingetID: "qBittorrent.qBittorrent", Category: CategoryUtilities, Note: "Free open-source BitTorrent client"},
	{Name: "JDownloader", WingetID: "AppWork.JDownloader", Category: CategoryUtilities, Note: "Download manager for files and streams"},
	{Name: "OFGB", WingetID: "xM4ddy.OFGB", Category: CategoryUtilities, Note: "Remove ads injected into Windows Settings"},
	{Name: "OpenHashTab", WingetID: "namazso.OpenHashTab", Category: CategoryUtilities, Note: "File hash verification in Explorer properties"},
	{Name: "ExifCleaner", WingetID: "szTheory.exifcleaner", Category: CategoryUtilities, Note: "Remove EXIF metadata from photos and files"},
	{Name: "DevToys", WingetID: "DevToys-app.DevToys", Category: CategoryUtilities, Note: "Swiss army knife offline developer tools"},
}

// ByCategory returns all apps belonging to the given category.
func ByCategory(cat Category) []App {
	var result []App
	for _, a := range Catalog {
		if a.Category == cat {
			result = append(result, a)
		}
	}
	return result
}
