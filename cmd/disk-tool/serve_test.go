package main

import (
	"bufio"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func TestParseReadyPort(t *testing.T) {
	tests := []struct {
		line string
		ok   bool
		port int
	}{
		{"disk-tool-ready port=8080", true, 8080},
		{"  disk-tool-ready port=54321  ", true, 54321},
		{"disk-tool-ready port=0", false, 0},
		{"disk-tool-ready port=abc", false, 0},
		{"listening on 8080", false, 0},
	}
	for _, tc := range tests {
		port, ok := parseReadyPort(tc.line)
		if ok != tc.ok || port != tc.port {
			t.Errorf("parseReadyPort(%q) = (%d, %v), want (%d, %v)", tc.line, port, ok, tc.port, tc.ok)
		}
	}
}

func TestServeReadyStdoutAndDynamicPort(t *testing.T) {
	if os.Getenv("DISK_TOOL_SERVE_HELPER") == "1" {
		runServe([]string{"--no-open", "--port", "0", "--ready-stdout"})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestServeReadyStdoutAndDynamicPort$")
	cmd.Env = append(os.Environ(), "DISK_TOOL_SERVE_HELPER=1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatal("no stdout from serve")
	}
	port, ok := parseReadyPort(scanner.Text())
	if !ok {
		t.Fatalf("expected ready line, got %q", scanner.Text())
	}

	deadline := time.Now().Add(3 * time.Second)
	var resp *http.Response
	for time.Now().Before(deadline) {
		resp, err = http.Get("http://127.0.0.1:" + strconv.Itoa(port) + "/api/roots")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("GET /api/roots: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
}
