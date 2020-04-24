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

func GetData(server, user, password, registry, image string ) *data.Report {
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
//	result.ScanHistory = new(data.ScanHistoryType)

	general := getData(urlBase, user, password)
	if err := json.Unmarshal(general, result.General); err != nil {
		fmt.Println("Can't parse response from server (general):", string(general))
		os.Exit(1)
	}

	sensitive := getData(urlBase+sensitive_url, user, password)
	if err := json.Unmarshal(sensitive, result.Sensitive); err != nil {
		fmt.Println("Can't parse response from server (sensitive):", string(sensitive))
		os.Exit(1)
	}

	malware := getData(urlBase+malware_url, user, password)
	if err := json.Unmarshal(malware, result.Malware); err != nil {
		fmt.Println("Can't parse response from server (malware):", string(malware))
		os.Exit(1)
	}

	vulnerabiliti := getData(urlBase+vulnerabiliti_url, user, password)
	if err := json.Unmarshal(vulnerabiliti, result.Vulnerabilities); err != nil {
		fmt.Println("Can't parse response from server (vulnerabiliti):")
		os.Exit(1)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		pagesize := 10

		var scanMax int
		var scanCount int
		for page := 1; scanCount < scanMax || scanMax == 0; page++ {
			scanHistoryUrl := fmt.Sprintf("%s%s?order_by=-date&page=%d&page_size=%d", urlBase, scanhistory_url, page, pagesize)
			scanHistorySource := getData( scanHistoryUrl, user, password)
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
	wg.Wait()
	return result
}
