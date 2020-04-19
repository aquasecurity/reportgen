package data

type Report struct {
	ServerUrl string
	General *GeneralType
	Sensitive *SensitiveType
	Malware *MalwareType
	Vulnerabilities *VulnerabilitiesType
}

func (report *Report) MappingImageAssuranceChecks() map[string]bool {
	result := make(map[string]bool)
	for _,v := range report.General.AssuranceResults.ChecksPerformed {
		result[v.Control] = v.Failed
	}
	return result
}
