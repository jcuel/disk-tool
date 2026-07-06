package dedup

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/jcuel/disk-tool/internal/model"
	"github.com/jcuel/disk-tool/internal/safety"
)

const defaultMinSize = 4096

// FindDuplicates hashes files under paths and returns groups with identical content.
func FindDuplicates(files []model.FileEntry, minSize int64) []model.DuplicateGroup {
	if minSize <= 0 {
		minSize = defaultMinSize
	}
	byHash := make(map[string][]model.FileEntry)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 4)

	for _, f := range files {
		if f.Size < minSize {
			continue
		}
		if safety.ClassifyPath(f.Path) != safety.ZoneNormal &&
			safety.ClassifyPath(f.Path) != safety.ZoneReview &&
			safety.ClassifyPath(f.Path) != safety.ZoneCaution &&
			safety.ClassifyPath(f.Path) != safety.ZoneMaintenance {
			continue
		}
		wg.Add(1)
		go func(entry model.FileEntry) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			h, err := hashFile(entry.Path)
			if err != nil {
				return
			}
			mu.Lock()
			byHash[h] = append(byHash[h], entry)
			mu.Unlock()
		}(f)
	}
	wg.Wait()

	groups := make([]model.DuplicateGroup, 0)
	for hash, entries := range byHash {
		if len(entries) < 2 {
			continue
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
		var wasted int64
		for i := 1; i < len(entries); i++ {
			wasted += entries[i].Size
		}
		groups = append(groups, model.DuplicateGroup{
			Hash:   hash,
			Files:  entries,
			Wasted: wasted,
		})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Wasted > groups[j].Wasted })
	return groups
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil))[:16], nil
}

// CollectFilesFromTree gathers file paths from scanned tree nodes (best-effort from largest list).
func CollectFilesFromTree(job *model.ScanJob) []model.FileEntry {
	if job == nil {
		return nil
	}
	if len(job.LargestFiles) > 0 {
		return job.LargestFiles
	}
	return nil
}

// WalkTreeFiles walks tree for regular files when largest list is incomplete.
func WalkTreeFiles(root *model.ScanNode) []model.FileEntry {
	if root == nil {
		return nil
	}
	var out []model.FileEntry
	var walk func(*model.ScanNode)
	walk = func(n *model.ScanNode) {
		if n == nil {
			return
		}
		if !n.IsDir && n.Size > 0 {
			out = append(out, model.FileEntry{Path: n.Path, Name: n.Name, Size: n.Size})
		}
		for _, ch := range n.Children {
			walk(ch)
		}
	}
	walk(root)
	return out
}
