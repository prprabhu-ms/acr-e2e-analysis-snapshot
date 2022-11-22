package playwright

type Run struct {
	Metadata Metadata
	Stats    Stats
}

type Metadata struct {
	ArtifactId      int64
	WorkflowRunId   int64
	WorkflowAttempt int64
	BuildFlavor     string
	Composite       string
}

type Stats struct {
	Suites []SuiteWrapper `json:"suites"`
}

type SuiteWrapper struct {
	Suites []Suite `json:"suites"`
}

type Suite struct {
	Title string `json:"title"`
	Specs []Spec `json:"specs"`
}

type Spec struct {
	Title string `json:"title"`
	Tests []Test `json:"tests"`
}

type Test struct {
	ExpectedStatus string   `json:"expectedStatus"`
	ProjectName    string   `json:"projectName"`
	Results        []Result `json:"results"`
	Status         string   `json:"status"`
}

type Result struct {
	Status   string `json:"status"`
	Duration int64  `json:"duration"`
	Retry    int64  `json:"retry"`
}
