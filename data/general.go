package data

type CheckPerformedType struct {
	PolicyName string `json:"policy_name"`
	Failed bool `json:"failed"`
	Control string `json:"control"`
}

type AssuranceResultsType struct {
	Disallowed bool `json:"disallowed"`
	ChecksPerformed []CheckPerformedType `json:"checks_performed"`
}

type GeneralType struct {
	ImageName string `json:"name"`
	Registry string `json:"registry"`
	Malware int `json:"malware"`
	Disallowed bool `json:"disallowed"`
	Os string `json:"os"`
	OsVersion string `json:"os_version"`
	Created string `json:"created"`

	Critical int `json:"crit_vulns"`
	High int `json:"high_vulns"`
	Medium int `json:"med_vulns"`
	Low int `json:"low_vulns"`
	Negligible int `json:"neg_vulns"`
	
	AssuranceResults AssuranceResultsType `json:"assurance_results"`
}
