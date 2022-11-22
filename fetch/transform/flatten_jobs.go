package transform

import "scripts/gh"

type FlattenedWorkflowJob struct {
	WorkflowId         int64 `json:"workflow_id"`
	WorkflowRunAttempt int64 `json:"workflow_run_attempt"`

	JobStatus      string `json:"job_status"`
	JobConclusion  string `json:"job_conclusion"`
	JobStartedAt   string `json:"job_started_at"`
	JobCompletedAt string `json:"job_completed_at"`
	JobName        string `json:"job_name"`

	StepName        string `json:"step_name"`
	StepStatus      string `json:"step_status"`
	StepConclusion  string `json:"step_conclusion"`
	StepStartedAt   string `json:"step_started_at"`
	StepCompletedAt string `json:"step_completed_at"`
}

func FlattenWorkflowJobs(jobs []gh.WorkflowJob) []FlattenedWorkflowJob {
	resp := []FlattenedWorkflowJob{}
	for _, job := range jobs {
		for _, step := range job.Steps {
			resp = append(resp, FlattenedWorkflowJob{
				WorkflowId:         job.RunId,
				WorkflowRunAttempt: job.RunAttempt,

				JobStatus:      job.Status,
				JobConclusion:  job.Conclusion,
				JobStartedAt:   job.StartedAt,
				JobCompletedAt: job.CompletedAt,
				JobName:        job.Name,

				StepName:        step.Name,
				StepStatus:      step.Status,
				StepConclusion:  step.Conclusion,
				StepStartedAt:   step.StartedAt,
				StepCompletedAt: step.CompletedAt,
			})
		}
	}
	return resp
}
