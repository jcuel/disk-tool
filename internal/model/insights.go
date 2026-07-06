package model

type CleanupCategory string

const (
	CategoryNodeModules  CleanupCategory = "node_modules"
	CategoryPythonVenv   CleanupCategory = "python_venv"
	CategoryBuildOutput  CleanupCategory = "build_output"
	CategoryPackageCache CleanupCategory = "package_cache"
	CategoryDownload     CleanupCategory = "downloads"
	CategoryLargeFile    CleanupCategory = "large_file"
	CategoryStaleLarge   CleanupCategory = "stale-large-file"
)

type CleanupCandidate struct {
	Category  CleanupCategory `json:"category"`
	Path      string          `json:"path"`
	Size      int64           `json:"size"`
	Hint      string          `json:"hint"`
	Risk      string          `json:"risk"` // review | caution
	Zone      string          `json:"zone"`
	Deletable bool            `json:"deletable"`
}

type SafetyZoneStats struct {
	Count int   `json:"count"`
	Bytes int64 `json:"bytes"`
}

type SafetyGrid struct {
	Zones       map[string]SafetyZoneStats `json:"zones"`
	DriveRoot   bool                       `json:"driveRoot"`
	ProtectedBytes int64                      `json:"protectedBytes"`
}

type TopConsumer struct {
	Name string  `json:"name"`
	Path string  `json:"path"`
	Size int64   `json:"size"`
	Pct  float64 `json:"pct"`
}

type InsightsReport struct {
	Summary           string             `json:"summary"`
	TopConsumers      []TopConsumer      `json:"topConsumers"`
	CleanupCandidates []CleanupCandidate `json:"cleanupCandidates"`
	TotalReclaimable  int64              `json:"totalReclaimable"`
	TicketText        string             `json:"ticketText"`
	SafetyGrid        *SafetyGrid        `json:"safetyGrid,omitempty"`
}

type DuplicateGroup struct {
	Hash   string      `json:"hash"`
	Files  []FileEntry `json:"files"`
	Wasted int64       `json:"wasted"`
}
