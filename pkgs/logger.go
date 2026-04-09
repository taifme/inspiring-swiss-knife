package pkgs

import (
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

// logEntry holds a single log line with level metadata.
type logEntry struct {
	Level   log.Level
	Message string
}

var (
	logMu      sync.Mutex
	logEntries []logEntry
	// Logger is the global charmbracelet/log logger wired to our buffer.
	Logger *log.Logger
)

// bufWriter is an io.Writer that appends each Write to logEntries.
type bufWriter struct{}

func (bufWriter) Write(p []byte) (n int, err error) {
	line := strings.TrimRight(string(p), "\n")
	if line == "" {
		return len(p), nil
	}
	logMu.Lock()
	logEntries = append(logEntries, logEntry{Message: line})
	logMu.Unlock()
	return len(p), nil
}

func init() {
	Logger = log.New(bufWriter{})
	Logger.SetLevel(log.DebugLevel)
	Logger.SetReportTimestamp(true)
	Logger.SetReportCaller(false)
}

// GetLogLines returns a snapshot of all log lines.
func GetLogLines() []string {
	logMu.Lock()
	defer logMu.Unlock()
	out := make([]string, len(logEntries))
	for i, e := range logEntries {
		out[i] = e.Message
	}
	return out
}

// ClearLogs removes all buffered log entries.
func ClearLogs() {
	logMu.Lock()
	defer logMu.Unlock()
	logEntries = logEntries[:0]
}
