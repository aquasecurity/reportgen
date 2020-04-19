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

func (report *Report) GetImageAssurancePolicies() map[string]bool {
	result := make(map[string] bool)

	for _, policy := range report.General.AssuranceResults.ChecksPerformed {
		if _, ok := result[policy.PolicyName]; !ok {
			result[policy.PolicyName] = false
		}
		if policy.Failed {
			result[policy.PolicyName] = true
		}
	}
	return result
}
