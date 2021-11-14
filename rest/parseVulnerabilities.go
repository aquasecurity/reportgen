package rest

import (
	"encoding/json"
	"fmt"
	"github.com/aquasecurity/reportgen/data"
	"os"
)

func getVulnerabilitiesUrl(requestType int, currentPage, pageSize int, urlBase string) string {
	var url string
	switch requestType {
	case data.ImageRequest:
		url = fmt.Sprintf("%s%s?pagesize=%d&page=%d", urlBase, image_vulnerabiliti_url, pageSize, currentPage)
		break
	case data.HostRequest:
		url = fmt.Sprintf("%s/vulnerabilities?include_vpatch_info=true&page=%d&pagesize=%d&skip_count=true&hide_base_image=false&node_id=",
			urlBase, currentPage, pageSize)
		break
	}
	return url
}

func getVulnerabilities(user, password string, severities []string, requestType int, urlBase string) *data.VulnerabilitiesType {
	result := new(data.VulnerabilitiesType)

	defaultPageSize := 100
	for _, severity := range severities {
		vulnCount := 0
		var maxVulnerabilities int
		for page := 1; vulnCount < maxVulnerabilities || vulnCount == 0; page++ {
			urlForVulnerabilitiesD := getVulnerabilitiesUrl(requestType, page, defaultPageSize, urlBase)
			if severity != all_severities {
				urlForVulnerabilitiesD += "&severity=" + severity
			}
			vulnerabiliti := getData(urlForVulnerabilitiesD, user, password)
			vuln := new(data.VulnerabilitiesType)
			if err := json.Unmarshal(vulnerabiliti, &vuln); err != nil {
				fmt.Println("Can't parse response from server (vulnerabiliti):")
				fmt.Println(string(vulnerabiliti))
				os.Exit(1)
			}
			if vuln.Count == 0 {
				break
			}
			if maxVulnerabilities == 0 {
				maxVulnerabilities = vuln.Count
			}
			result.Results = append(result.Results, vuln.Results...)
			vulnCount += len(vuln.Results)
		}
	}
	return result
}
