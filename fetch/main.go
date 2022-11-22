package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"scripts/gh"
	"scripts/playwright"
	"scripts/processed"
	"scripts/transform"
)

func main() {
	args := parseArgs()

	if args.workDir == "" {
		dir, err := os.MkdirTemp("", "acranalysis")
		if err != nil {
			log.Fatal(err)
		}
		args.workDir = dir
		if !args.leak {
			defer os.RemoveAll(dir)
		}
	}
	log.Printf("Using working directory %s\n", args.workDir)

	plMetaDir := filepath.Join(args.workDir, "raw-playwright")
	playwrightDir := filepath.Join(plMetaDir, "test-results")
	os.MkdirAll(playwrightDir, 0755)
	if args.fetch {
		gh.DownloadPlaywrightStats(playwrightDir, plMetaDir, args.batch, args.after, args.before)
	}

	ciRunsFile := filepath.Join(args.workDir, "ci.json")
	if args.fetch {
		gh.FetchCIWorkflowRuns(ciRunsFile)
	}

	wRuns := gh.ParseCIWorkflowRuns(ciRunsFile)
	pRuns := playwright.ParseAllRuns(playwrightDir)
	matched := transform.MatchCIRunsToPlaywrightStats(pRuns, wRuns)
	log.Printf("Parsed %d playwright runs (%d matched to %d workflow runs)", len(pRuns), len(matched), len(wRuns))

	jobsDir := filepath.Join(args.workDir, "workflows", "jobs")
	os.MkdirAll(jobsDir, 0755)
	wIds := workflowUids(matched)
	if args.fetch {
		gh.FetchAllWorkflowJobs(jobsDir, args.batch, wIds)
	}
	wJobs := gh.ParseAllWorkflowJobs(jobsDir)
	log.Printf("Parsed %d jobs\n", len(wJobs))

	writeOutput(args, transform.FlattenPlaywrightRuns(matched), transform.FlattenWorkflowJobs(wJobs))
}

type Args struct {
	after   time.Time
	batch   int
	before  time.Time
	fetch   bool
	leak    bool
	oDir    string
	purge   bool
	workDir string
}

func parseArgs() Args {
	batch := flag.Int("batch", 0, "Number of artifacts to fetch together. If greater than 0, data is fetched in batches with a gap to avoid GitHub API throttling.")
	fetch := flag.Bool("fetch", true, "Whether to fetch remote data. If set to false, `workdir` must be set and contain already fetch data.")
	leak := flag.Bool("leak", false, "Whether to leak working directory. Only useful when `workdir` is unset.")
	oDir := flag.String("outdir", "", "Directory to write final output to. Required.")
	purge := flag.Bool("purge", false, "If set, purge existing results in `outdir`. By default, results are merged.")
	rAfter := flag.String("after", "", "Only fetch artifacts created after provided date, formatted as 2006-01-02.")
	rBefore := flag.String("before", "", "Only fetch artifacts created before provided date, formatted as 2006-01-02.")
	workDir := flag.String("workdir", "", "Working directory to use. Default is to use a temporary directory.")

	flag.Parse()

	if *oDir == "" {
		log.Fatalf("Argument `outdir` is required\n")
	}

	return Args{
		after:   parseTimeIfSet(rAfter),
		batch:   *batch,
		before:  parseTimeIfSet(rBefore),
		fetch:   *fetch,
		leak:    *leak,
		oDir:    *oDir,
		purge:   *purge,
		workDir: *workDir,
	}
}

func parseTimeIfSet(v *string) time.Time {
	if *v == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", *v)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func workflowUids(matched []transform.MatchedRun) []gh.WorkflowUid {
	resp := make([]gh.WorkflowUid, 0, len(matched))
	for _, r := range matched {
		resp = append(resp, gh.WorkflowUid{
			Id:         r.Workflow.Id,
			RunAttempt: r.Workflow.RunAttempt,
		})
	}
	return resp
}

func writeOutput(args Args, pRuns []transform.FlattenedPlaywrightRun, jobs []transform.FlattenedWorkflowJob) {
	os.MkdirAll(args.oDir, 0755)
	pPath := filepath.Join(args.oDir, "playwright.json")
	jPath := filepath.Join(args.oDir, "jobs.json")
	if args.purge {
		processed.WritePlaywright(pPath, pRuns)
		processed.WriteJobs(jPath, jobs)
	} else {
		processed.WritePlaywrightIncr(pPath, pRuns)
		processed.WriteJobsIncr(jPath, jobs)
	}
}
