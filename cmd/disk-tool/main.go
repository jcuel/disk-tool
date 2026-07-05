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
	"runtime"
	"time"

	"github.com/jcuel/disk-tool/internal/api"
	"github.com/jcuel/disk-tool/internal/model"
	"github.com/jcuel/disk-tool/internal/scanner"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "serve":
		runServe(os.Args[2:])
	case "scan":
		runScan(os.Args[2:])
	case "version":
		fmt.Println("disk-tool 0.1.0")
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  disk-tool serve [--port 8080] [--no-open]\n  disk-tool scan <path> [--json] [--full]\n  disk-tool version\n")
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "HTTP port")
	noOpen := fs.Bool("no-open", false, "do not open browser")
	_ = fs.Parse(args)

	store := api.NewStore()
	static, err := staticHandler()
	if err != nil {
		log.Fatalf("static assets: %v", err)
	}
	srv := api.NewServer(store, static)
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	url := "http://" + addr

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("disk-tool listening on %s", url)
	if !*noOpen {
		go openBrowser(url)
	}
	if err := http.Serve(ln, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}

func runScan(args []string) {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "output JSON")
	full := fs.Bool("full", false, "scan entire tree (default: overview only)")
	_ = fs.Parse(args)
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "scan requires a path")
		os.Exit(1)
	}
	root, err := api.ValidateRoot(fs.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	sc := scanner.New()
	opts := scanner.Options{Root: root}
	var tree *model.ScanNode
	var largest []model.FileEntry
	if *full {
		tree, largest, err = sc.Scan(context.Background(), opts)
	} else {
		tree, largest, err = sc.ScanOverview(context.Background(), opts)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *asJSON {
		out := map[string]any{
			"root":         root,
			"tree":         tree,
			"largestFiles": largest,
			"mode":         map[bool]string{true: "full", false: "overview"}[*full],
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(out)
		return
	}
	fmt.Printf("Scanned %s (%s): %d bytes, %d files\n", root, map[bool]string{true: "full", false: "overview"}[*full], tree.Size, tree.FileCount)
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
