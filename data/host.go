package data

import "time"

type HostAssuranceType struct {
	Disallowed       bool                 `json:"disallowed"`
	AssuranceResults []CheckPerformedType `json:"assurance_results"`
}

type SecurityIssuesType struct {
	CritVulns int   `json:"crit_vulns"`
	HighVulns int   `json:"high_vulns"`
	MedVulns  int   `json:"med_vulns"`
	LowVulns  int   `json:"low_vulns"`
	NegVulns  int   `json:"neg_vulns"`
	Malwares  int   `json:"malware"`
	LastScan  int64 `json:"last_vuln_scan"`
}

type HostType struct {
	Id        int    `json:"id"`
	ClusterId int    `json:"cluster_id"`
	NodeId    string `json:"node_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`

	SecurityIssues SecurityIssuesType `json:"security_issues"`

	CreatedDate string `json:"created_date"`
}

func (host *HostType) GetGeneral() *GeneralType {
	general := new(GeneralType)

	(*general).ImageName = host.Name
	(*general).Created = host.CreatedDate
	(*general).ScanDate = time.Unix(host.SecurityIssues.LastScan, 0).Format("2006-01-02T15:04:05.999999999Z07:00")
	(*general).Critical = host.SecurityIssues.CritVulns
	(*general).High = host.SecurityIssues.HighVulns
	(*general).Medium = host.SecurityIssues.MedVulns
	(*general).Low = host.SecurityIssues.LowVulns
	(*general).Negligible = host.SecurityIssues.NegVulns
	(*general).Malware = host.SecurityIssues.Malwares
	return general
}
