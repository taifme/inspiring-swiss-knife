package pkgs

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// InstallStatus represents the state of an installation
type InstallStatus int

const (
	StatusPending  InstallStatus = iota
	StatusRunning
	StatusSuccess
	StatusFailed
	StatusSkipped
)

// InstallResult holds the outcome of installing a single app
type InstallResult struct {
	App    App
	Status InstallStatus
	Output string
	Error  string
}

// InstallMsg is sent via Bubble Tea channel when an install completes
type InstallMsg struct {
	Result InstallResult
}

// AllDoneMsg is sent when all installs are complete
type AllDoneMsg struct {
	Results []InstallResult
}

// InstallProgressMsg carries a log line during installation
type InstallProgressMsg struct {
	AppName string
	Line    string
}

// InstallApps runs winget for each selected app and sends results on the returned channel.
// The caller must drain the channel to prevent goroutine leaks.
func InstallApps(apps []App, resultCh chan<- InstallResult) {
	var wg sync.WaitGroup

	Logger.Info("Starting installation batch", "count", len(apps))

	// Run up to 3 installs concurrently to avoid hammering the network
	sem := make(chan struct{}, 3)

	for _, app := range apps {
		wg.Add(1)
		go func(a App) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			Logger.Info("Installing", "app", a.Name, "id", a.WingetID)
			result := runWinget(a)
			switch result.Status {
			case StatusSuccess:
				Logger.Info("Installed successfully", "app", a.Name)
			case StatusSkipped:
				Logger.Warn("Already installed, skipped", "app", a.Name)
			case StatusFailed:
				Logger.Error("Installation failed", "app", a.Name, "err", result.Error)
			}
			resultCh <- result
		}(app)
	}

	wg.Wait()
	Logger.Info("Installation batch complete", "count", len(apps))
	close(resultCh)
}

func runWinget(app App) InstallResult {
	args := []string{
		"install",
		"--id", app.WingetID,
		"--silent",
		"--accept-package-agreements",
		"--accept-source-agreements",
		"--disable-interactivity",
	}

	cmd := exec.Command("winget", args...)

	// Capture combined output
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	output := out.String()

	// Check if already installed
	if strings.Contains(output, "No applicable update found") ||
		strings.Contains(output, "already installed") {
		return InstallResult{
			App:    app,
			Status: StatusSkipped,
			Output: output,
		}
	}

	if err != nil {
		return InstallResult{
			App:    app,
			Status: StatusFailed,
			Output: output,
			Error:  err.Error(),
		}
	}

	return InstallResult{
		App:    app,
		Status: StatusSuccess,
		Output: output,
	}
}

// RunPowerShell executes a PowerShell command with elevated awareness.
// Returns stdout+stderr and any error.
func RunPowerShell(script string) (string, error) {
	Logger.Debug("Running PowerShell", "script", script)
	cmd := exec.Command("powershell.exe",
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy", "Bypass",
		"-Command", script,
	)

	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	output := out.String()
	if err != nil {
		Logger.Error("PowerShell error", "err", err, "output", strings.TrimSpace(output))
	} else {
		Logger.Debug("PowerShell OK", "output", strings.TrimSpace(output))
	}
	return output, err
}

// RunPowerShellStreaming runs a PowerShell script and sends output lines to lineCh.
func RunPowerShellStreaming(script string, lineCh chan<- string) error {
	cmd := exec.Command("powershell.exe",
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy", "Bypass",
		"-Command", script,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		lineCh <- scanner.Text()
	}

	return cmd.Wait()
}

// IsAdmin returns true if the current process has administrator privileges.
func IsAdmin() bool {
	_, err := exec.Command("net", "session").Output()
	return err == nil
}

// WingetAvailable checks if winget is installed and reachable.
func WingetAvailable() bool {
	cmd := exec.Command("winget", "--version")
	return cmd.Run() == nil
}
