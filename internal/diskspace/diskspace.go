package diskspace

import (
	"path/filepath"
)

type Info struct {
	Path  string `json:"path"`
	Total int64  `json:"total"`
	Free  int64  `json:"free"`
	Used  int64  `json:"used"`
}

func VolumeRoot(path string) string {
	if vol := filepath.VolumeName(path); vol != "" {
		return vol + string(filepath.Separator)
	}
	return string(filepath.Separator)
}

func ForPath(path string) (Info, error) {
	root := VolumeRoot(path)
	total, free, err := platformUsage(root)
	if err != nil {
		return Info{}, err
	}
	t := int64(total)
	f := int64(free)
	return Info{
		Path:  root,
		Total: t,
		Free:  f,
		Used:  t - f,
	}, nil
}
