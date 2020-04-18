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
	"time"
)

const (
	api_url = "api/v2/images/"
	vulnerabiliti = "/vulnerabilities"
	sensitive = "/sensitive"
	malware = "/malware"
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

func buildReportData( general *data.General  ) *data.Report  {
	result := new(data.Report)
	result.ImageName = general.Name
	result.Registry = general.Registry
	result.Os = general.Os
	result.OsVersion = general.OsVersion
	result.ImageAllowed = general.Disallowed
	var err error
	result.Created, err = time.Parse("2006-01-02T15:04:05.000000Z", general.Created)
	if err != nil {
		fmt.Println(err.Error)
	}

	return result
}

func GetData(server, user, password, registry, image string ) *data.Report {
	urlBase := server + api_url + registry + "/" + image

	var general data.General
	generalData := getData(urlBase, user, password)

	if err := json.Unmarshal(generalData, &general); err != nil {
		fmt.Println("Can't parse response from server (general):", string(generalData))
		os.Exit(1)
	}



//	fmt.Println(string())
	fmt.Println("===============================================================")
	fmt.Println("vulnerabiliti:")
	fmt.Println(string(getData(urlBase+vulnerabiliti, user, password)))
	fmt.Println("===============================================================")
	fmt.Println("sensitive: ")
	fmt.Println(string(getData(urlBase+sensitive, user, password)))
	fmt.Println("===============================================================")
	fmt.Println("malware:")
	fmt.Println(string(getData(urlBase+malware, user, password)))
	fmt.Println("===============================================================")

	return buildReportData(&general)

}
