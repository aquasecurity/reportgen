package rest

import (
	"../data"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

func GetHostData(server, user, password, host string) *data.Report {
	report := new(data.Report)
	report.RequestType = data.HostRequest

	hostData := new(data.HostType)
	url := getUrlApi(server, api_host) + host
	hostSource := getData(url, user, password)
	fmt.Println("[DEBUG] Host data:",string(hostSource))

	if err := json.Unmarshal(hostSource, hostData); err != nil {
		fmt.Println("Can't parse response from server (host):")
		fmt.Println(string(hostSource))
		os.Exit(1)
	}

	report.General = hostData.GetGeneral()

	var wg sync.WaitGroup
	// To get the assurance checks:
	wg.Add(1)
	go func() {
		defer wg.Done()
		url = getUrlApi(server, "api/v2/status/host/")+ hostData.NodeId

		hostAssurances := new(data.HostAssuranceType)
		assuranceSource := getData(url, user, password)
		if err := json.Unmarshal( assuranceSource, hostAssurances); err != nil {
			fmt.Println("Can't parse response from server (status):", err.Error())
			os.Exit(1)
		}
		report.General.AssuranceResults.Disallowed = hostAssurances.Disallowed
		report.General.AssuranceResults.ChecksPerformed = hostAssurances.AssuranceResults
	}()

	// To get malware:
	wg.Add(1)
	go func() {
		defer wg.Done()
		url = getUrlApi(server, "api/v1/hosts/")+hostData.NodeId + "/malware"
		report.Malware = getMalwares(user, password, url)
	}()

	// Get vulnerabilities
	wg.Add(1)
	go func() {
		defer wg.Done()
		baseUrl := getUrlApi(server, "api/v1/hosts/")+hostData.NodeId
		report.Vulnerabilities = getVulnerabilities(user, password, []string{all_severities}, data.HostRequest, baseUrl)
	}()

	// get the bench results:
	wg.Add(1)
	go func() {
		defer wg.Done()
		url := getUrlApi(server, "api/v2/risks/bench/") + hostData.NodeId + "/bench_results"
		benchResultsSource := getData(url, user, password)
		report.BenchResults = new (data.BenchResultsType)
		if err := json.Unmarshal(benchResultsSource, report.BenchResults); err != nil {
			fmt.Println("Can't parse response from server (the bench results):", err.Error())
			os.Exit(1)
		}
	}()

	wg.Wait()
	return report
}
