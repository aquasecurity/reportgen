package rest

import (
	"../data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	bodyText, err := ioutil.ReadAll(resp.Body)
	return bodyText
}

func GetData(server, user, password, registry, image string ) *data.Report {
	urlBase := server + api_url + registry + "/" + image

	result := new(data.Report)
	result.General = new(data.GeneralType)
	result.Sensitive = new(data.SensitiveType)
	result.Malware = new(data.MalwareType)
	result.Vulnerabilities = new(data.VulnerabilitiesType)

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

//	fmt.Println(string(general))
	fmt.Println("===============================================================")
	fmt.Println("sensitive: ")
//	fmt.Println(string(sensitive))
	fmt.Println("===============================================================")
	fmt.Println("malware:")
//	fmt.Println(string(malware))
	fmt.Println("===============================================================")
	/*
	fmt.Println("vulnerabiliti:")
	fmt.Println(string(getData(urlBase+vulnerabiliti, user, password)))
	fmt.Println("===============================================================")
*/
	return result

}
