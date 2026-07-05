package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/jcuel/disk-tool/internal/insights"
	"github.com/jcuel/disk-tool/internal/model"
	"github.com/jcuel/disk-tool/internal/scanner"
)

type subscriber chan model.ProgressEvent

type Store struct {
	mu          sync.RWMutex
	jobs        map[string]*model.ScanJob
	cancels     map[string]context.CancelFunc
	expandMu    map[string]*sync.Mutex
	subscribers map[string]map[subscriber]struct{}
	scanner     *scanner.Scanner
}

func NewStore() *Store {
	return &Store{
		jobs:        make(map[string]*model.ScanJob),
		cancels:     make(map[string]context.CancelFunc),
		expandMu:    make(map[string]*sync.Mutex),
		subscribers: make(map[string]map[subscriber]struct{}),
		scanner:     scanner.New(),
	}
}

func (s *Store) expandLock(id string) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.expandMu[id] == nil {
		s.expandMu[id] = &sync.Mutex{}
	}
	return s.expandMu[id]
}

func (s *Store) Start(root string) (*model.ScanJob, error) {
	validRoot, err := ValidateRoot(root)
	if err != nil {
		return nil, err
	}

	id := newScanID()
	job := &model.ScanJob{
		ID:        id,
		Root:      validRoot,
		Status:    model.ScanStatusRunning,
		StartedAt: time.Now(),
	}

	s.mu.Lock()
	s.jobs[id] = job
	s.subscribers[id] = make(map[subscriber]struct{})
	ctx, cancel := context.WithCancel(context.Background())
	s.cancels[id] = cancel
	s.mu.Unlock()

	go s.runOverview(ctx, id, validRoot)
	return job, nil
}

func (s *Store) runOverview(ctx context.Context, id, root string) {
	s.broadcast(id, model.ProgressEvent{Type: "started", ScanID: id, Status: string(model.ScanStatusRunning), Root: root})

	tree, largest, err := s.scanner.ScanOverview(ctx, scanner.Options{
		Root: root,
		OnProgress: func(ev model.ProgressEvent) {
			ev.Type = "progress"
			ev.ScanID = id
			s.updateProgress(id, ev)
			s.broadcast(id, ev)
		},
	})

	s.finishScan(ctx, id, tree, largest, err)
}

func (s *Store) Expand(id, path string, depth int) error {
	job, ok := s.Get(id)
	if !ok {
		return errNotFound
	}
	if job.Status != model.ScanStatusCompleted {
		return errNotReady
	}
	target, err := PathWithinRoot(job.Root, path)
	if err != nil {
		return err
	}
	if depth <= 0 {
		depth = model.DefaultDrillDepth
	}

	mu := s.expandLock(id)
	if !mu.TryLock() {
		return errExpandBusy
	}

	s.broadcast(id, model.ProgressEvent{
		Type:       "expand-started",
		ScanID:     id,
		TargetPath: target,
	})

	go func() {
		defer mu.Unlock()
		ctx := context.Background()
		branch, largest, err := s.scanner.ScanBranch(ctx, target, job.Root, depth, scanner.Options{
			Root: job.Root,
			OnProgress: func(ev model.ProgressEvent) {
				ev.Type = "expand-progress"
				ev.ScanID = id
				ev.TargetPath = target
				s.updateProgress(id, ev)
				s.broadcast(id, ev)
			},
		})

		s.mu.Lock()
		defer s.mu.Unlock()
		j, ok := s.jobs[id]
		if !ok || j.Tree == nil {
			return
		}
		if err != nil {
			s.broadcastLocked(id, model.ProgressEvent{
				Type:       "expand-error",
				ScanID:     id,
				TargetPath: target,
				Error:      err.Error(),
			})
			return
		}
		if target == j.Root {
			j.Tree = branch
		} else {
			scanner.MergeBranch(j.Tree, target, branch)
		}
		j.LargestFiles = scanner.MergeLargest(j.LargestFiles, largest)
		j.Insights = insights.Analyze(j)
		s.broadcastLocked(id, model.ProgressEvent{
			Type:       "expand-completed",
			ScanID:     id,
			TargetPath: target,
			Status:     string(model.ScanStatusCompleted),
		})
	}()
	return nil
}

func (s *Store) finishScan(ctx context.Context, id string, tree *model.ScanNode, largest []model.FileEntry, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[id]
	if !ok {
		return
	}
	now := time.Now()
	job.CompletedAt = &now

	if err != nil {
		if ctx.Err() != nil {
			job.Status = model.ScanStatusCancelled
			s.broadcastLocked(id, model.ProgressEvent{Type: "cancelled", ScanID: id, Status: string(model.ScanStatusCancelled)})
		} else {
			job.Status = model.ScanStatusFailed
			job.Error = err.Error()
			s.broadcastLocked(id, model.ProgressEvent{Type: "error", ScanID: id, Status: string(model.ScanStatusFailed), Error: err.Error()})
		}
		delete(s.cancels, id)
		return
	}

	job.Status = model.ScanStatusCompleted
	job.Tree = tree
	job.LargestFiles = largest
	job.Insights = insights.Analyze(job)
	s.broadcastLocked(id, model.ProgressEvent{Type: "completed", ScanID: id, Status: string(model.ScanStatusCompleted)})
	delete(s.cancels, id)
}

func (s *Store) updateProgress(id string, ev model.ProgressEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[id]
	if !ok {
		return
	}
	job.DirsScanned = ev.DirsScanned
	job.FilesScanned = ev.FilesScanned
	job.BytesScanned = ev.BytesScanned
	job.CurrentPath = ev.CurrentPath
}

func (s *Store) Get(id string) (*model.ScanJob, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	if !ok {
		return nil, false
	}
	copy := *job
	return &copy, true
}

func (s *Store) GetMutable(id string) (*model.ScanJob, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	return job, ok
}

func (s *Store) Cancel(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cancel, ok := s.cancels[id]
	if !ok {
		return false
	}
	cancel()
	return true
}

func (s *Store) Subscribe(id string) (subscriber, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[id]; !ok {
		return nil, false
	}
	ch := make(subscriber, 32)
	s.subscribers[id][ch] = struct{}{}
	return ch, true
}

func (s *Store) Unsubscribe(id string, ch subscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if subs, ok := s.subscribers[id]; ok {
		delete(subs, ch)
	}
}

func (s *Store) broadcast(id string, ev model.ProgressEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.broadcastLocked(id, ev)
}

func newScanID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Store) broadcastLocked(id string, ev model.ProgressEvent) {
	subs := s.subscribers[id]
	for ch := range subs {
		select {
		case ch <- ev:
		default:
		}
	}
}
