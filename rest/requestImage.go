package rest

import (
	"encoding/json"
	"fmt"
	"github.com/aquasecurity/reportgen/data"
	"os"
	"strings"
	"sync"
)

func GetImageData(server, user, password, registry, image string, severities []string) *data.Report {
	if len(severities) > 0 {
		fmt.Println("Next severities will be selected:")
		for _, s := range severities {
			fmt.Println("*", s)
		}
	} else {
		severities = []string{all_severities}
	}

	var correctImage string
	if strings.Contains(image, ":") {
		correctImage = strings.ReplaceAll(image, ":", "/")
	} else {
		correctImage = image
	}

	urlBase := getUrlApi(server, api_image) + registry + "/" + correctImage

	result := new(data.Report)
	result.RequestType = data.ImageRequest
	result.General = new(data.GeneralType)
	result.Sensitive = new(data.SensitiveType)
	result.Malware = new(data.MalwareType)
	result.Vulnerabilities = new(data.VulnerabilitiesType)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		result.Vulnerabilities = getVulnerabilities(user, password, severities, data.ImageRequest, urlBase)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pagesize := 10

		var scanMax int
		var scanCount int
		for page := 1; scanCount < scanMax || scanMax == 0; page++ {
			scanHistoryUrl := fmt.Sprintf("%s%s?order_by=-date&page=%d&page_size=%d", urlBase, image_scanhistory_url, page, pagesize)
			scanHistorySource := getData(scanHistoryUrl, user, password)
			scanHistory := new(data.ScanHistoryType)
			if err := json.Unmarshal(scanHistorySource, scanHistory); err != nil {
				fmt.Println("Can't parse response from server (scanHistory):")
				os.Exit(1)
			}
			if result.ScanHistory == nil {
				result.ScanHistory = scanHistory
				scanMax = scanHistory.Count
			} else {
				result.ScanHistory.Results = append(result.ScanHistory.Results, scanHistory.Results...)
			}
			scanCount += len(scanHistory.Results)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		general := getData(urlBase, user, password)
		if err := json.Unmarshal(general, &result.General); err != nil {
			fmt.Println("Can't parse response from server (general):", err.Error())
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sensitive := getData(urlBase+image_sensitive_url, user, password)
		if err := json.Unmarshal(sensitive, &result.Sensitive); err != nil {
			fmt.Println("Can't parse response from server (sensitive):", err.Error())
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result.Malware = getMalwares(user, password, urlBase+image_malware_url)
	}()
	wg.Wait()
	return result
}
