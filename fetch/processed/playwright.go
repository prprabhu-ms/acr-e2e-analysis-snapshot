package processed

import (
	"scripts/transform"
)

func WritePlaywright(path string, runs []transform.FlattenedPlaywrightRun) {
	writePlaywright(path, runs, []transform.FlattenedPlaywrightRun{})
}

func WritePlaywrightIncr(path string, runs []transform.FlattenedPlaywrightRun) {
	writePlaywright(path, runs, readPlaywright(path))
}

func writePlaywright(path string, runs []transform.FlattenedPlaywrightRun, oldRuns []transform.FlattenedPlaywrightRun) {
	newIDs := make(map[int64]bool)
	for _, r := range runs {
		newIDs[r.MetaId] = true
	}

	merged := make([]transform.FlattenedPlaywrightRun, 0, len(runs)+len(oldRuns))
	for _, r := range oldRuns {
		if !newIDs[r.MetaId] {
			merged = append(merged, r)
		}
	}
	merged = append(merged, runs...)
	writeJSON(path, merged)
}

func readPlaywright(path string) []transform.FlattenedPlaywrightRun {
	r := []transform.FlattenedPlaywrightRun{}
	readJSON(path, &r)
	return r
}
