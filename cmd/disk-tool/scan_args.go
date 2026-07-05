package main

import "fmt"

type scanFlags struct {
	json bool
	full bool
	path string
}

// parseScanArgs extracts --json/--full from anywhere in args (Go's flag package
// stops at the first positional, so "scan C:\Users --json" would miss --json).
func parseScanArgs(args []string) (scanFlags, error) {
	var f scanFlags
	var paths []string
	for _, a := range args {
		switch a {
		case "--json":
			f.json = true
		case "--full":
			f.full = true
		default:
			if len(a) > 0 && a[0] == '-' {
				return scanFlags{}, fmt.Errorf("unknown flag: %s", a)
			}
			paths = append(paths, a)
		}
	}
	if len(paths) != 1 {
		return scanFlags{}, fmt.Errorf("scan requires exactly one path")
	}
	f.path = paths[0]
	return f, nil
}
