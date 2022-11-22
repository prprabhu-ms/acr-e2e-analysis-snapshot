package transform

import "log"

type FlattenedPlaywrightRun struct {
	// Used to merge playwright runs across fetches.
	MetaId int64 `json:"meta_id"`

	WorkflowId         int64  `json:"workflow_id"`
	WorkflowRunNumber  int64  `json:"workflow_run_number"`
	WorkflowRunAttempt int64  `json:"workflow_run_attempt"`
	WorkflowCreatedAt  string `json:"workflow_created_at"`
	WorkflowEvent      string `json:"workflow_event"`
	WorkflowHeadBranch string `json:"workflow_head_branch"`

	CommitSha   string `json:"commit_sha"`
	BuildFlavor string `json:"build_flavor"`
	Composite   string `json:"composite"`

	PlaywrightSuiteTitle         string `json:"pw_suite_title"`
	PlaywrightSpecTitle          string `json:"pw_spec_title"`
	PlaywrightTestProjectName    string `json:"pw_test_project_name"`
	PlaywrightTestExpectedStatus string `json:"pw_test_expected_status"`
	PlaywrightTestStatus         string `json:"pw_test_status"`
	PlaywrightResultStatus       string `json:"pw_result_status"`
	PlaywrightResultDuration     int64  `json:"pw_result_duration"`
	PlaywrightResultRetry        int64  `json:"pw_result_retry"`
}

func flatten(run MatchedRun) []FlattenedPlaywrightRun {
	prs := run.Workflow.PullRequests
	sha := ""
	if len(prs) != 1 {
		log.Printf("WARNING: Found %d pull requests for workflow %d/%d, want 1", len(prs), run.Workflow.Id, run.Workflow.RunAttempt)
	} else {
		sha = prs[0].Head.Sha
	}

	resp := []FlattenedPlaywrightRun{}
	for _, sw := range run.Playwright.Stats.Suites {
		for _, suite := range sw.Suites {
			for _, spec := range suite.Specs {
				for _, test := range spec.Tests {
					for _, result := range test.Results {
						resp = append(resp, FlattenedPlaywrightRun{
							MetaId: run.Playwright.Metadata.ArtifactId,

							WorkflowId:         run.Workflow.Id,
							WorkflowRunNumber:  run.Workflow.RunNumber,
							WorkflowRunAttempt: run.Workflow.RunAttempt,
							WorkflowCreatedAt:  run.Workflow.CreatedAt,
							WorkflowEvent:      run.Workflow.Event,
							WorkflowHeadBranch: run.Workflow.HeadBranch,

							CommitSha:   sha,
							BuildFlavor: run.Playwright.Metadata.BuildFlavor,
							Composite:   run.Playwright.Metadata.Composite,

							PlaywrightSuiteTitle:         suite.Title,
							PlaywrightSpecTitle:          spec.Title,
							PlaywrightTestProjectName:    test.ProjectName,
							PlaywrightTestExpectedStatus: test.ExpectedStatus,
							PlaywrightTestStatus:         test.Status,
							PlaywrightResultStatus:       result.Status,
							PlaywrightResultDuration:     result.Duration,
							PlaywrightResultRetry:        result.Retry,
						})
					}
				}
			}
		}
	}
	return resp
}

func FlattenPlaywrightRuns(runs []MatchedRun) []FlattenedPlaywrightRun {
	resp := make([]FlattenedPlaywrightRun, 0, len(runs))
	for _, r := range runs {
		resp = append(resp, flatten(r)...)
	}
	return resp
}
