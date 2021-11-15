package data

type VulnerabilitiesResourceType struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Path    string `json:"path"`
}

type ReferencedVulnerabilitiesType struct {
	Name         string `json:"name"`
	AquaSeverity string `json:"aqua_severity"`
}

type VulnerabilitiesResultType struct {
	Name                      string                          `json:"name"`
	AquaSeverity              string                          `json:"aqua_severity"`
	AquaScore                 float64                         `json:"aqua_score"`
	Resource                  VulnerabilitiesResourceType     `json:"resource"`
	FixVersion                string                          `json:"fix_version"`
	Description               string                          `json:"description"`
	Solution                  string                          `json:"solution"`
	ReferencedVulnerabilities []ReferencedVulnerabilitiesType `json:"referenced_vulnerabilities"`
}

type VulnerabilitiesType struct {
	Count   int                         `json:"count"`
	Results []VulnerabilitiesResultType `json:"result"`
}
