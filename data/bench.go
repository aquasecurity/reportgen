package data

type TestResultBenchType struct {
	Status         string   `json:"status"`
	TestDesc       string   `json:"test_desc"`
	TestInfo       []string `json:"test_info"`
	TestNumber     string   `json:"test_number"`
	ActualValue    string   `json:"actual_value"`
	ExpectedResult string   `json:"expected_result"`
}

type TestBenchType struct {
	/*desc", "info", "pass", "warn" and "fail" */
	Desc    string                `json:"desc"`
	Info    int                   `json:"info"`
	Pass    int                   `json:"pass"`
	Warn    int                   `json:"warn"`
	Fail    int                   `json:"fail"`
	Results []TestResultBenchType `json:"results"`
}

type BaseResultsBenchType struct {
	Tests []TestBenchType `json:"tests"`
}

type BaseBenchType struct {
	Result BaseResultsBenchType `json:"result"`
}

type BenchResultsType struct {
	Cis       BaseBenchType `json:"cis"`
	KubeBench BaseBenchType `json:"kube_bench"`
	Linux     BaseBenchType `json:"linux"`
	OpenShift BaseBenchType `json:"openshift"`
}
