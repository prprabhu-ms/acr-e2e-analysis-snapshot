package playwright

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseAllRuns(dataDir string) []Run {
	tests, err := os.ReadDir(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	resp := make([]Run, 0, len(tests))
	for _, test := range tests {
		resp = append(resp, Run{
			Metadata: parseName(test.Name()),
			Stats:    parseStats(filepath.Join(dataDir, test.Name())),
		})
	}
	return resp
}

func parseStats(path string) Stats {
	log.Printf("Parsing %s\n", path)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	resp := Stats{}
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Fatal(err)
	}
	return resp
}

func parseName(name string) Metadata {
	name = strings.Split(name, ".")[0]
	parts := strings.Split(name, "-")
	artifactId, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatal(err)
	}
	workflowRunId, err := strconv.Atoi(parts[3])
	if err != nil {
		log.Fatal(err)
	}
	workflowAttempt, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Fatal(err)
	}
	return Metadata{
		ArtifactId:      int64(artifactId),
		WorkflowRunId:   int64(workflowRunId),
		WorkflowAttempt: int64(workflowAttempt),
		BuildFlavor:     parts[5],
		Composite:       parts[6],
	}
}
