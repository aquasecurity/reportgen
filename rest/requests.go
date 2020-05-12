package rest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	api_image         = "api/v2/images/"
	api_host          = "api/v2/infrastructure/node/"

	image_vulnerabiliti_url = "/vulnerabilities"
	image_sensitive_url     = "/sensitive"
	image_malware_url       = "/malware"
	all_severities          = "ALL"
	image_scanhistory_url   = "/scan_history"
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

func getUrlApi(server, api string) string {
	var slash string
	if strings.HasSuffix( server,"/") {
		slash = ""
	} else {
		slash = "/"
	}
	return server + slash + api
}

