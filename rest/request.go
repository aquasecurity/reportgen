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

	wg.Add(1)
	go func() {
		defer wg.Done()
		defaultPageSize := 100
		vulnCount := 0
		var maxVulnerabilities int
		for page:=1; vulnCount < maxVulnerabilities || vulnCount == 0;page++ {
			urlForVulnerabilitiesD := fmt.Sprintf("%s%s?pagesize=%d&page=%d", urlBase, vulnerabiliti_url, defaultPageSize, page)
			vulnerabiliti := getData( urlForVulnerabilitiesD, user, password)
			vuln := new (data.VulnerabilitiesType)
			if err := json.Unmarshal(vulnerabiliti, &vuln); err != nil {
				fmt.Println("Can't parse response from server (vulnerabiliti):")
				fmt.Println(string(vulnerabiliti))
				os.Exit(1)
			}
			if maxVulnerabilities == 0 {
				result.Vulnerabilities = vuln
				maxVulnerabilities = vuln.Count
			} else {
				result.Vulnerabilities.Results = append(result.Vulnerabilities.Results, vuln.Results...)
			}
			vulnCount += len(vuln.Results)
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
