package safety

import (
	"github.com/jcuel/disk-tool/internal/model"
)

// Preset describes a one-click maintenance strategy.
type Preset struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AutoSelect  bool   `json:"autoSelect"`
}

var presetDefs = []Preset{
	{
		ID:          "dev-reclaim",
		Name:        "Dev reclaim",
		Description: "Regenerable project artifacts: node_modules, target, dist, build, __pycache__, .venv",
		AutoSelect:  true,
	},
	{
		ID:          "cache-review",
		Name:        "Cache review",
		Description: "Shared package caches — opens review only, never auto-selected",
		AutoSelect:  false,
	},
	{
		ID:          "downloads-sweep",
		Name:        "Downloads sweep",
		Description: "Stale installers and large archives in Downloads (requires age filter)",
		AutoSelect:  true,
	},
	{
		ID:          "temp-cleanup",
		Name:        "Temp cleanup",
		Description: "User temp folders tagged as maintenance zone",
		AutoSelect:  true,
	},
}

// AllPresets returns maintenance preset definitions.
func AllPresets() []Preset {
	out := make([]Preset, len(presetDefs))
	copy(out, presetDefs)
	return out
}

// PresetMatch holds preset id and matching deletable paths.
type PresetMatch struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MatchCount  int      `json:"matchCount"`
	MatchBytes  int64    `json:"matchBytes"`
	Paths       []string `json:"paths"`
}

// MatchPresets filters cleanup candidates into preset buckets.
func MatchPresets(candidates []model.CleanupCandidate) []PresetMatch {
	out := make([]PresetMatch, 0, len(presetDefs))
	for _, p := range presetDefs {
		m := PresetMatch{ID: p.ID, Name: p.Name, Description: p.Description}
		for _, c := range candidates {
			if !c.Deletable {
				continue
			}
			if presetMatches(p.ID, c) {
				m.MatchCount++
				m.MatchBytes += c.Size
				m.Paths = append(m.Paths, c.Path)
			}
		}
		out = append(out, m)
	}
	return out
}

func presetMatches(id string, c model.CleanupCandidate) bool {
	switch id {
	case "dev-reclaim":
		return c.Zone == string(ZoneReview) || (c.Zone == string(ZoneNormal) && c.Risk == "review")
	case "cache-review":
		return c.Zone == string(ZoneCaution) || c.Risk == "caution"
	case "downloads-sweep":
		return c.Category == model.CategoryDownload || c.Category == model.CategoryStaleLarge
	case "temp-cleanup":
		return c.Zone == string(ZoneMaintenance)
	default:
		return false
	}
}
