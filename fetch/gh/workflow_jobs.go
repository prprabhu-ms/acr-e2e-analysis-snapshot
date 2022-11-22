package gh

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"scripts/gh/internal/api"
	"scripts/gh/internal/workerpool"
)

func FetchWorkflowJobs(oDir string, w WorkflowUid) {
	data := api.Run(fmt.Sprintf("repos/%s/actions/runs/%d/attempts/%d/jobs", ghRepo, w.Id, w.RunAttempt))
	writeFile(filepath.Join(oDir, fmt.Sprintf("job_%d_%d", w.Id, w.RunAttempt)), data)
}

type item struct {
	WorkflowUid WorkflowUid
	Done        bool
}

func FetchAllWorkflowJobs(oDir string, batch int, ws []WorkflowUid) {
	actions := make([]func(), 0, len(ws))
	for _, w := range ws {
		w := w
		actions = append(actions, func() {
			FetchWorkflowJobs(oDir, w)
		})
	}
	workerpool.RunBatched(maxParallelFetch(), batch, actions)
}

func ParseAllWorkflowJobs(dataDir string) []WorkflowJob {
	paths, err := os.ReadDir(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	response := []WorkflowJob{}
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(dataDir, path.Name()))
		if err != nil {
			log.Fatal(err)
		}

		for _, part := range SplitPages(string(data)) {
			pageResponse := WorkflowJobResponse{}
			if err := json.Unmarshal([]byte(part), &pageResponse); err != nil {
				log.Fatal(err)
			}
			response = append(response, pageResponse.Jobs...)
		}
	}
	return response
}
