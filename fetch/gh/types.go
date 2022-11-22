package gh

type Head struct {
	Sha string `json:"sha"`
}
type PullRequest struct {
	Head Head `json:"head"`
}
type WorkflowRun struct {
	Id           int64         `json:"id"`
	RunNumber    int64         `json:"run_number"`
	RunAttempt   int64         `json:"run_attempt"`
	CreatedAt    string        `json:"created_at"`
	PullRequests []PullRequest `json:"pull_requests"`
	ArtifactsUrl string        `json:"artifacts_url"`
	Event        string        `json:"event"`
	HeadBranch   string        `json:"head_branch"`
}
type WorkflowRunAttempts struct {
	TotalCount   int64         `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

// Not meant to parse JSON data
type WorkflowUid struct {
	Id         int64
	RunAttempt int64
}

type WorkflowJobResponse struct {
	TotalCount int64         `json:"total_count"`
	Jobs       []WorkflowJob `json:"jobs"`
}

type WorkflowJob struct {
	RunId       int64             `json:"run_id"`
	RunAttempt  int64             `json:"run_attempt"`
	Status      string            `json:"status"`
	Conclusion  string            `json:"conclusion"`
	StartedAt   string            `json:"started_at"`
	CompletedAt string            `json:"completed_at"`
	Name        string            `json:"name"`
	Steps       []WorkflowJobStep `json:"steps"`
}

type WorkflowJobStep struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}
