package transform

import (
	"log"
	"scripts/gh"
	"scripts/playwright"
)

type MatchedRun struct {
	Playwright playwright.Run
	Workflow   gh.WorkflowRun
}

func MatchCIRunsToPlaywrightStats(pRuns []playwright.Run, wRuns []gh.WorkflowRun) []MatchedRun {
	mWRuns := make(map[int64]map[int64]gh.WorkflowRun)
	for _, a := range wRuns {
		r, ok := mWRuns[a.Id]
		if !ok {
			r = make(map[int64]gh.WorkflowRun)
			mWRuns[a.Id] = r
		}
		r[a.RunAttempt] = a
	}

	resp := make([]MatchedRun, 0, len(pRuns))
	for _, pRun := range pRuns {
		if wRun, ok := mWRuns[pRun.Metadata.WorkflowRunId][pRun.Metadata.WorkflowAttempt]; ok {
			resp = append(resp, MatchedRun{
				Playwright: pRun,
				Workflow:   wRun,
			})
		} else {
			log.Printf("WARNING: Failed to find a matched workflow attempt for %#v\n", pRun.Metadata)
		}
	}
	return resp
}
