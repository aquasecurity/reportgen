package rest

import (
	"../data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	api_url = "api/v2/images/"
	vulnerabiliti_url = "/vulnerabilities"
	sensitive_url = "/sensitive"
	malware_url = "/malware"
	all_severities = "ALL"
	scanhistory_url = "/scan_history"
)

var (
	wg sync.WaitGroup
)

func getData(url, user, password string) []byte  {
	fmt.Println("Getting data from", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil{
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Wrong access to the URL: ", url)
		fmt.Println("Status:", resp.Status)
		os.Exit(1)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	return bodyText
}

func GetData(server, user, password, registry, image string, severities []string ) *data.Report {
	if len(severities) > 0 {
		fmt.Println("Next severities will be selected:")
		for _, s := range severities {
			fmt.Println("*", s)
		}
	} else {
		severities = []string{all_severities}
	}

	var slash string
	if strings.HasSuffix( server,"/") {
		slash = ""
	} else {
		slash = "/"
	}

	var correctImage string
	if strings.Contains(image, ":") {
		correctImage = strings.ReplaceAll(image, ":", "/")
	} else {
		correctImage = image
	}

	urlBase := server + slash+ api_url + registry + "/" + correctImage

	result := new(data.Report)
	result.General = new(data.GeneralType)
	result.Sensitive = new(data.SensitiveType)
	result.Malware = new(data.MalwareType)
	result.Vulnerabilities = new(data.VulnerabilitiesType)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defaultPageSize := 100
		result.Vulnerabilities = new(data.VulnerabilitiesType)

		for _, severity := range severities {
			vulnCount := 0
			var maxVulnerabilities int
			for page:=1; vulnCount < maxVulnerabilities || vulnCount == 0;page++ {
				urlForVulnerabilitiesD := fmt.Sprintf("%s%s?pagesize=%d&page=%d", urlBase, vulnerabiliti_url, defaultPageSize, page)
				if severity != all_severities {
					urlForVulnerabilitiesD += "&severity=" + severity
				}
				vulnerabiliti := getData( urlForVulnerabilitiesD, user, password)
				vuln := new (data.VulnerabilitiesType)
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
				result.Vulnerabilities.Results = append(result.Vulnerabilities.Results, vuln.Results...)
				vulnCount += len(vuln.Results)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pagesize := 10

		var scanMax int
		var scanCount int
		for page := 1; scanCount < scanMax || scanMax == 0; page++ {
			scanHistoryUrl := fmt.Sprintf("%s%s?order_by=-date&page=%d&page_size=%d", urlBase, scanhistory_url, page, pagesize)
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
			fmt.Println("Can't parse response from server (general):", string(general))
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sensitive := getData(urlBase+sensitive_url, user, password)
		if err := json.Unmarshal(sensitive, &result.Sensitive); err != nil {
			fmt.Println("Can't parse response from server (sensitive):", string(sensitive))
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go func() {
		wg.Done()
		malware := getData(urlBase+malware_url, user, password)
		if err := json.Unmarshal(malware, &result.Malware); err != nil {
			fmt.Println("Can't parse response from server (malware):", string(malware))
			os.Exit(1)
		}
	}()

	wg.Wait()

	return result
}
