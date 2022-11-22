package processed

import (
	"fmt"
	"scripts/transform"
)

func WriteJobs(path string, jobs []transform.FlattenedWorkflowJob) {
	writeJobs(path, jobs, []transform.FlattenedWorkflowJob{})
}

func WriteJobsIncr(path string, jobs []transform.FlattenedWorkflowJob) {
	writeJobs(path, jobs, readJobs(path))
}

func writeJobs(path string, jobs []transform.FlattenedWorkflowJob, oldJobs []transform.FlattenedWorkflowJob) {
	newIDs := make(map[string]bool)
	for _, j := range jobs {
		newIDs[jobUID(j)] = true
	}

	merged := make([]transform.FlattenedWorkflowJob, 0, len(jobs)+len(oldJobs))
	for _, j := range oldJobs {
		if !newIDs[jobUID(j)] {
			merged = append(merged, j)
		}
	}
	merged = append(merged, jobs...)
	writeJSON(path, merged)
}

func readJobs(path string) []transform.FlattenedWorkflowJob {
	r := []transform.FlattenedWorkflowJob{}
	readJSON(path, &r)
	return r
}

func jobUID(j transform.FlattenedWorkflowJob) string {
	return fmt.Sprintf("w:%d,a:%d,j:%s,s:%s", j.WorkflowId, j.WorkflowRunAttempt, j.JobName, j.StepName)
}
