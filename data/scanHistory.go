package data

type ScanHistoryResult struct {
	Date string `json:"date"`
	ImageId string `json:"docker_id"`

	ImageCreationDate string `json:"image_creation_date"`
	
	SecurityStatus bool `json:"disallowed"`

	CriticalCount int `json:"crit_vulns"`
	HighCount int `json:"high_vulns"`
	MediumCount int `json:"med_vulns"`
	LowCount int `json:"low_vulns"`
	NegCount int `json:"neg_vulns"`

	MalwareCount int `json:"malware"`
}

type ScanHistoryType struct {
	Count int `json:"count"`
	Results []ScanHistoryResult `json:"result"`
}
