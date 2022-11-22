package gh

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"scripts/gh/internal/api"
	"scripts/gh/internal/workerpool"
)

func DownloadPlaywrightStats(outDir string, metadataDir string, batch int, after time.Time, before time.Time) {
	artifacts := listArtifacts(metadataDir)
	testArtifacts := filterArtifacts(artifacts, func(a artifact) bool {
		if !strings.Contains(a.Name, "e2e-results") {
			return false
		}

		created, err := time.Parse("2006-01-02T15:04:05Z", a.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		return (after.IsZero() || created.After(after)) && (before.IsZero() || created.Before(before))
	})

	log.Printf("Will download %d e2e result artifacts (out of %d total artifacts)\n", len(testArtifacts), len(artifacts))

	actions := make([]func(), 0, len(testArtifacts))
	for _, a := range testArtifacts {
		a := a
		actions = append(actions, func() {
			downloadOne(outDir, a.Id, a.Name)
		})
	}
	workerpool.RunBatched(maxParallelFetch(), batch, actions)
}

type listArtifactsResponse struct {
	TotalCount int64      `json:"total_count"`
	Artifacts  []artifact `json:"artifacts"`
}

type artifact struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	CreatedAt          string `json:"created_at"`
	SizeInBytes        int64  `json:"size_in_bytes"`
	Url                string `json:"url"`
	ArchiveDownloadUrl string `json:"archive_download_url"`
}

func listArtifacts(outDir string) []artifact {
	oFile := filepath.Join(outDir, "artifacts.json")
	data := api.Run(fmt.Sprintf("repos/%s/actions/artifacts", ghRepo))
	writeFile(oFile, data)

	artifacts := []artifact{}
	for _, part := range SplitPages(string(data)) {
		pageResponse := listArtifactsResponse{}
		if err := json.Unmarshal([]byte(part), &pageResponse); err != nil {
			log.Fatal(err)
		}
		artifacts = append(artifacts, pageResponse.Artifacts...)
	}
	return artifacts
}

func downloadOne(outDir string, id int64, name string) {
	data := api.Run(fmt.Sprintf("repos/%s/actions/artifacts/%d/zip", ghRepo, id))

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Fatal(err)
	}
	if len(r.File) == 0 {
		log.Fatalf("Found 0 files!")
	}

	for i, p := range r.File {
		if p.FileInfo().IsDir() {
			continue
		}
		f, err := p.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		writeFile(filepath.Join(outDir, fmt.Sprintf("%d-%s-%d", id, name, i)), data)
	}
}

func filterArtifacts(as []artifact, predicate func(a artifact) bool) []artifact {
	resp := make([]artifact, 0, len(as))
	for _, a := range as {
		if predicate(a) {
			resp = append(resp, a)
		}
	}
	return resp
}
