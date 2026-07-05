package api

import "errors"

var (
	errNotFound   = errors.New("scan not found")
	errNotReady   = errors.New("overview scan not complete")
	errExpandBusy = errors.New("expand already in progress")
)
