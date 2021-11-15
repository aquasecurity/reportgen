package data

type TestResultBenchType struct {
	TestNumber     string   `json:"test_number"`
	Status         string   `json:"status"`
	TestDesc       string   `json:"test_desc"`
	TestInfo       []string `json:"test_info"`
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
	Openshift BaseBenchType `json:"openshift"`
}

type HostCommonType struct {
	HostOs  string `json:"host_os"`
	Address string `json:"address"`
}
