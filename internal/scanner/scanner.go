package scanner

import (
	"container/heap"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jcuel/disk-tool/internal/model"
)

const maxLargestFiles = 100

type Options struct {
	Root       string
	MaxDepth   int // 0 = overview only; >0 = depth limit; -1 = full tree
	OnProgress func(model.ProgressEvent)
}

type Scanner struct{}

func New() *Scanner {
	return &Scanner{}
}

type fileHeap []model.FileEntry

func (h fileHeap) Len() int           { return len(h) }
func (h fileHeap) Less(i, j int) bool { return h[i].Size < h[j].Size }
func (h fileHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *fileHeap) Push(x any)        { *h = append(*h, x.(model.FileEntry)) }
func (h *fileHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

type walkState struct {
	dirs    int64
	files   int64
	bytes   int64
	cur     string
	largest fileHeap
	visited map[visitKey]struct{}
	mu      sync.Mutex
}

func newWalkState() *walkState {
	st := &walkState{visited: make(map[visitKey]struct{})}
	heap.Init(&st.largest)
	return st
}

func (st *walkState) emit(opts Options) {
	if opts.OnProgress == nil {
		return
	}
	opts.OnProgress(model.ProgressEvent{
		Type:         "progress",
		DirsScanned:  st.dirs,
		FilesScanned: st.files,
		BytesScanned: st.bytes,
		CurrentPath:  st.cur,
	})
}

func (st *walkState) noteFile(path, name string, size int64) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.files++
	st.bytes += size
	pushLargest(&st.largest, model.FileEntry{Path: path, Name: name, Size: size})
}

func (st *walkState) largestFiles() []model.FileEntry {
	st.mu.Lock()
	defer st.mu.Unlock()
	out := make([]model.FileEntry, st.largest.Len())
	for i := st.largest.Len() - 1; i >= 0; i-- {
		out[st.largest.Len()-1-i] = st.largest[i]
	}
	return out
}

// ScanOverview lists immediate children with accurate aggregate sizes (parallel per top-level folder).
func (sc *Scanner) ScanOverview(ctx context.Context, opts Options) (*model.ScanNode, []model.FileEntry, error) {
	root, err := prepareRoot(opts.Root)
	if err != nil {
		return nil, nil, err
	}
	st := newWalkState()
	st.dirs++
	st.cur = root
	st.emit(opts)

	node := newNode(root)
	entries, err := os.ReadDir(root)
	if err != nil {
		return node, st.largestFiles(), nil
	}

	type dirResult struct {
		child *model.ScanNode
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	results := make([]dirResult, 0)

	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			wg.Wait()
			return nil, nil, err
		}
		name := entry.Name()
		full := filepath.Join(root, name)

		if entry.Type()&os.ModeSymlink != 0 && skipSymlink(full, st.visited) {
			continue
		}

		if entry.IsDir() {
			wg.Add(1)
			go func(dirPath, dirName string) {
				defer wg.Done()
				size, files, err := sc.aggregateDir(ctx, dirPath, st, opts)
				if err != nil {
					return
				}
				child := &model.ScanNode{
					Name:       dirName,
					Path:       dirPath,
					Size:       size,
					FileCount:  files,
					IsDir:      true,
					Scanned:    false,
					Expandable: true,
				}
				mu.Lock()
				results = append(results, dirResult{child: child})
				mu.Unlock()
			}(full, name)
			continue
		}

		fi, err := entry.Info()
		if err != nil {
			continue
		}
		size := fi.Size()
		st.noteFile(full, name, size)
		node.Size += size
		node.FileCount++
	}

	wg.Wait()
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	for _, r := range results {
		node.Children = append(node.Children, r.child)
		node.Size += r.child.Size
		node.FileCount += r.child.FileCount
	}
	node.Scanned = true
	sortChildren(node)
	return node, st.largestFiles(), nil
}

// ScanBranch scans a folder up to maxDepth levels and merges structure into the returned subtree root.
func (sc *Scanner) ScanBranch(ctx context.Context, branchPath, scanRoot string, maxDepth int, opts Options) (*model.ScanNode, []model.FileEntry, error) {
	branchPath, err := filepath.Abs(branchPath)
	if err != nil {
		return nil, nil, err
	}
	rootAbs, err := filepath.Abs(scanRoot)
	if err != nil {
		return nil, nil, err
	}
	if !isUnderRoot(rootAbs, branchPath) {
		return nil, nil, os.ErrPermission
	}
	info, err := os.Stat(branchPath)
	if err != nil || !info.IsDir() {
		return nil, nil, os.ErrInvalid
	}

	st := newWalkState()
	var scanDir func(string, int) (*model.ScanNode, error)
	scanDir = func(dir string, depth int) (*model.ScanNode, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		st.dirs++
		st.cur = dir
		st.emit(opts)

		node := newNode(dir)
		node.IsDir = true
		entries, err := os.ReadDir(dir)
		if err != nil {
			node.Scanned = true
			return node, nil
		}

		atLimit := depth >= maxDepth

		for _, entry := range entries {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			name := entry.Name()
			full := filepath.Join(dir, name)

			if entry.Type()&os.ModeSymlink != 0 && skipSymlink(full, st.visited) {
				continue
			}

			if entry.IsDir() {
				if atLimit {
					size, files, err := sc.aggregateDir(ctx, full, st, opts)
					if err != nil {
						return nil, err
					}
					node.Children = append(node.Children, &model.ScanNode{
						Name:       name,
						Path:       full,
						Size:       size,
						FileCount:  files,
						IsDir:      true,
						Scanned:    false,
						Expandable: true,
					})
					node.Size += size
					node.FileCount += files
					continue
				}
				child, err := scanDir(full, depth+1)
				if err != nil {
					return nil, err
				}
				node.Children = append(node.Children, child)
				node.Size += child.Size
				node.FileCount += child.FileCount
				continue
			}

			fi, err := entry.Info()
			if err != nil {
				continue
			}
			size := fi.Size()
			st.noteFile(full, name, size)
			node.Size += size
			node.FileCount++
		}

		node.Scanned = true
		node.Expandable = atLimit && len(node.Children) > 0
		sortChildren(node)
		return node, nil
	}

	tree, err := scanDir(branchPath, 1)
	if err != nil {
		return nil, nil, err
	}
	return tree, st.largestFiles(), nil
}

// Scan walks the full tree (CLI --full).
func (sc *Scanner) Scan(ctx context.Context, opts Options) (*model.ScanNode, []model.FileEntry, error) {
	opts.MaxDepth = -1
	root, err := prepareRoot(opts.Root)
	if err != nil {
		return nil, nil, err
	}
	st := newWalkState()

	var scanDir func(string) (*model.ScanNode, error)
	scanDir = func(dir string) (*model.ScanNode, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		st.dirs++
		st.cur = dir
		st.emit(opts)

		node := newNode(dir)
		node.IsDir = true
		entries, err := os.ReadDir(dir)
		if err != nil {
			node.Scanned = true
			return node, nil
		}

		for _, entry := range entries {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			name := entry.Name()
			full := filepath.Join(dir, name)

			if entry.Type()&os.ModeSymlink != 0 && skipSymlink(full, st.visited) {
				continue
			}

			if entry.IsDir() {
				child, err := scanDir(full)
				if err != nil {
					return nil, err
				}
				node.Children = append(node.Children, child)
				node.Size += child.Size
				node.FileCount += child.FileCount
				continue
			}

			fi, err := entry.Info()
			if err != nil {
				continue
			}
			size := fi.Size()
			st.noteFile(full, name, size)
			node.Size += size
			node.FileCount++
		}
		node.Scanned = true
		sortChildren(node)
		return node, nil
	}

	tree, err := scanDir(root)
	if err != nil {
		return nil, nil, err
	}
	return tree, st.largestFiles(), nil
}

func (sc *Scanner) aggregateDir(ctx context.Context, dir string, st *walkState, opts Options) (size, files int64, err error) {
	var walk func(string) error
	walk = func(path string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		st.dirs++
		st.cur = path
		st.emit(opts)

		entries, err := os.ReadDir(path)
		if err != nil {
			return nil
		}
		for _, entry := range entries {
			if err := ctx.Err(); err != nil {
				return err
			}
			name := entry.Name()
			full := filepath.Join(path, name)

			if entry.Type()&os.ModeSymlink != 0 && skipSymlink(full, st.visited) {
				continue
			}

			if entry.IsDir() {
				if err := walk(full); err != nil {
					return err
				}
				continue
			}

			fi, err := entry.Info()
			if err != nil {
				continue
			}
			sz := fi.Size()
			size += sz
			files++
			st.noteFile(full, name, sz)
		}
		return nil
	}
	if err := walk(dir); err != nil {
		return 0, 0, err
	}
	return size, files, nil
}

func prepareRoot(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", os.ErrInvalid
	}
	return abs, nil
}

func newNode(dir string) *model.ScanNode {
	name := filepath.Base(dir)
	if name == "" || name == "." {
		name = dir
	}
	return &model.ScanNode{Name: name, Path: dir, IsDir: true}
}

func isUnderRoot(root, target string) bool {
	rel, err := filepath.Rel(filepath.Clean(root), filepath.Clean(target))
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func pushLargest(h *fileHeap, f model.FileEntry) {
	if h.Len() < maxLargestFiles {
		heap.Push(h, f)
		return
	}
	if f.Size > (*h)[0].Size {
		heap.Pop(h)
		heap.Push(h, f)
	}
}

func sortChildren(node *model.ScanNode) {
	for i := 0; i < len(node.Children); i++ {
		for j := i + 1; j < len(node.Children); j++ {
			if node.Children[j].Size > node.Children[i].Size {
				node.Children[i], node.Children[j] = node.Children[j], node.Children[i]
			}
		}
	}
}

func FindNode(root *model.ScanNode, path string) *model.ScanNode {
	if root == nil {
		return nil
	}
	target, err := filepath.Abs(path)
	if err != nil {
		return nil
	}
	rootAbs, err := filepath.Abs(root.Path)
	if err != nil {
		return nil
	}
	if target == rootAbs {
		return root
	}
	return findNodeRecursive(root, target)
}

func findNodeRecursive(node *model.ScanNode, target string) *model.ScanNode {
	for _, c := range node.Children {
		if c.Path == target {
			return c
		}
		if found := findNodeRecursive(c, target); found != nil {
			return found
		}
	}
	return nil
}

// MergeBranch replaces the node at targetPath with branch (preserves path identity).
func MergeBranch(root *model.ScanNode, targetPath string, branch *model.ScanNode) bool {
	node := FindNode(root, targetPath)
	if node == nil || branch == nil {
		return false
	}
	node.Size = branch.Size
	node.FileCount = branch.FileCount
	node.Children = branch.Children
	node.Scanned = branch.Scanned
	node.Expandable = branch.Expandable
	node.IsDir = true
	sortChildren(node)
	recomputeAncestors(root, targetPath)
	return true
}

func recomputeAncestors(root *model.ScanNode, fromPath string) {
	// Sizes from overview aggregates remain valid; only refresh parent chain file counts if needed.
	_ = fromPath
	_ = root
}

func skipSymlink(path string, visited map[visitKey]struct{}) bool {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return true
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return true
	}
	key := visitKeyFromInfo(resolved, info)
	if _, ok := visited[key]; ok {
		return true
	}
	visited[key] = struct{}{}
	return false
}

func MergeLargest(existing []model.FileEntry, added []model.FileEntry) []model.FileEntry {
	h := fileHeap(existing)
	heap.Init(&h)
	for _, f := range added {
		pushLargest(&h, f)
	}
	out := make([]model.FileEntry, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		out[h.Len()-1-i] = h[i]
	}
	return out
}
