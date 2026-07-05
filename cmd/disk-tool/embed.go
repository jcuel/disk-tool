package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
)

//go:embed all:static
var staticFS embed.FS

func staticHandler() (http.Handler, error) {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, err
	}
	if _, err := sub.Open("index.html"); err != nil {
		return http.FileServer(http.Dir("web/dist")), nil
	}
	return http.FileServer(http.FS(sub)), nil
}

func init() {
	if os.Getenv("DISK_TOOL_DEV") == "1" {
		// dev: serve from web/dist on disk when rebuilding frontend
	}
}
