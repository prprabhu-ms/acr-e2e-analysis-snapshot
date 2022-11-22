package gh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"scripts/gh/internal/api"
)

func FetchCIWorkflowRuns(oFile string) {
	qArgs := url.Values{}
	qArgs.Add("status", "completed")
	data := api.Run(fmt.Sprintf("repos/%s/actions/workflows/ci.yml/runs?%s", ghRepo, qArgs.Encode()))
	writeFile(oFile, data)
}

func ParseCIWorkflowRuns(iFile string) []WorkflowRun {
	data, err := ioutil.ReadFile(iFile)
	if err != nil {
		log.Fatal(err)
	}

	response := make([]WorkflowRun, 0)
	for i, part := range SplitPages(string(data)) {
		fmt.Printf("Part %d... %d characters long\n", i, len(part))
		pageResponse := WorkflowRunAttempts{}
		if err := json.Unmarshal([]byte(part), &pageResponse); err != nil {
			log.Fatal(err)
		}
		response = append(response, pageResponse.WorkflowRuns...)
	}
	return response
}
