package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jcuel/disk-tool/internal/api"
	"github.com/jcuel/disk-tool/internal/model"
	"github.com/jcuel/disk-tool/internal/scanner"
)

const readyPrefix = "disk-tool-ready port="

func main() {
	if len(os.Args) < 2 {
		// Windows users often double-click the CLI binary from Releases; start the server
		// instead of flashing a console with usage text and exiting.
		if runtime.GOOS == "windows" {
			runServe(nil)
			return
		}
		printUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "serve":
		runServe(os.Args[2:])
	case "scan":
		runScan(os.Args[2:])
	case "version":
		fmt.Println("disk-tool 1.5.0")
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  disk-tool serve [--port 8080] [--no-open] [--ready-stdout]\n  disk-tool scan <path> [--json] [--full]\n  disk-tool version\n")
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "HTTP port (0 = OS-assigned)")
	noOpen := fs.Bool("no-open", false, "do not open browser")
	readyStdout := fs.Bool("ready-stdout", false, "emit disk-tool-ready line on stdout when listening")
	_ = fs.Parse(args)

	store := api.NewStore()
	static, err := staticHandler()
	if err != nil {
		log.Fatalf("static assets: %v", err)
	}
	srv := api.NewServer(store, static)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	actualPort := ln.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d", actualPort)

	if *readyStdout {
		fmt.Fprintf(os.Stdout, "%s%d\n", readyPrefix, actualPort)
		_ = os.Stdout.Sync()
	}

	log.Printf("disk-tool listening on %s", url)
	if !*noOpen {
		go openBrowser(url)
	}

	httpServer := &http.Server{Handler: srv.Handler()}
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.Serve(ln)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("disk-tool shutting down (%v)", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func parseReadyPort(line string) (int, bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, readyPrefix) {
		return 0, false
	}
	port, err := strconv.Atoi(strings.TrimPrefix(line, readyPrefix))
	if err != nil || port <= 0 {
		return 0, false
	}
	return port, true
}

func runScan(args []string) {
	flags, err := parseScanArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	root, err := api.ValidateRoot(flags.path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	sc := scanner.New()
	opts := scanner.Options{Root: root}
	var tree *model.ScanNode
	var largest []model.FileEntry
	if flags.full {
		tree, largest, err = sc.Scan(context.Background(), opts)
	} else {
		tree, largest, err = sc.ScanOverview(context.Background(), opts)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if flags.json {
		out := map[string]any{
			"root":         root,
			"tree":         tree,
			"largestFiles": largest,
			"mode":         map[bool]string{true: "full", false: "overview"}[flags.full],
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(out)
		return
	}
	fmt.Printf("Scanned %s (%s): %d bytes, %d files\n", root, map[bool]string{true: "full", false: "overview"}[flags.full], tree.Size, tree.FileCount)
}

func openBrowser(url string) {
	time.Sleep(300 * time.Millisecond)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
