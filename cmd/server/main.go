// Package main is a minimal HTTP server that serves the ISK install script.
// Deploy this to Railway (or any PaaS) so users can run:
//
//	irm "https://your-domain.com/win" | iex
//
// The server proxies install.ps1 from the GitHub repository so there is a
// single source of truth. Responses are cached in memory for 5 minutes.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	repoOwner = "taifme"
	repoName  = "inspiring-swiss-knife"
	repoURL   = "https://github.com/" + repoOwner + "/" + repoName
	rawScript = "https://raw.githubusercontent.com/" + repoOwner + "/" + repoName + "/main/install.ps1"
)

// scriptCache holds a short-lived cached copy of the install script.
var cache struct {
	mu      sync.RWMutex
	content []byte
	fetched time.Time
}

func fetchScript() ([]byte, error) {
	cache.mu.RLock()
	if time.Since(cache.fetched) < 5*time.Minute && len(cache.content) > 0 {
		defer cache.mu.RUnlock()
		return cache.content, nil
	}
	cache.mu.RUnlock()

	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Double-check after acquiring write lock
	if time.Since(cache.fetched) < 5*time.Minute && len(cache.content) > 0 {
		return cache.content, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(rawScript)
	if err != nil {
		return nil, fmt.Errorf("fetch script: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	cache.content = body
	cache.fetched = time.Now()
	return body, nil
}

func serveScript(w http.ResponseWriter, r *http.Request) {
	script, err := fetchScript()
	if err != nil {
		log.Printf("ERROR serving script: %v", err)
		http.Error(w, "Failed to fetch install script. Try the direct URL:\n"+rawScript, http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Source", rawScript)
	_, _ = w.Write(script)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Primary install endpoint — matches the WinUtil irm | iex pattern
	mux.HandleFunc("/win", serveScript)
	mux.HandleFunc("/install", serveScript)
	mux.HandleFunc("/install.ps1", serveScript)

	// Health check for Railway
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "ok")
	})

	// Root → redirect to GitHub
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, repoURL, http.StatusFound)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("ISK server starting on %s", addr)
	log.Printf("Script source: %s", rawScript)
	log.Printf("Install command: irm http://localhost%s/win | iex", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
