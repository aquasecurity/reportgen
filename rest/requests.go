package rest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	aquaClient "github.com/aquasecurity/terraform-provider-aquasec/client"
)

const (
	api_image = "api/v2/images/"
	api_host  = "api/v2/infrastructure/node/"

	image_vulnerabiliti_url = "/vulnerabilities"
	image_sensitive_url     = "/sensitive"
	image_malware_url       = "/malware"
	all_severities          = "ALL"
	image_scanhistory_url   = "/scan_history"
)

func isAquaSaasFlow(link string) bool {
	return strings.Contains(link, "cloud.aquasec.com")
}

var accessToken = ""
var mutex sync.Mutex

func getToken(link, user, password string) string {
	mutex.Lock()
	defer mutex.Unlock()
	if accessToken != "" {
		return accessToken
	}
	u, err := url.Parse(link)
	if err != nil {
		log.Fatalf("Can't parse a link %q: %v", link, err)
	}
	accessToken, err = aquaClient.NewClient(u.Host, user, password, u.Scheme == "https", nil).GetUSEAuthToken()
	if err != nil {
		log.Fatalf("Can't get the Aqua SaaS token for access to %q: %v", link, err)
	}
	return accessToken
}

func getData(link, user, password string) []byte {
	fmt.Println("Getting data from", link)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatalf("Can't create a request to %q: %v", link, err)
	}
	if isAquaSaasFlow(link) {
		req.Header.Set("Authorization", "Bearer "+getToken(link, user, password))
	} else {
		req.SetBasicAuth(user, password)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Wrong access to the URL: %q. Status: %s", link, resp.Status)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Can't read a body from %q: %v", link, err)
	}
	return bodyText
}

func getUrlApi(server, api string) string {
	var slash string
	if strings.HasSuffix(server, "/") {
		slash = ""
	} else {
		slash = "/"
	}
	return server + slash + api
}
