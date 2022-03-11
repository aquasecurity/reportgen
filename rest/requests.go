package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

var accessToken = ""
var accessUrl = ""
var mutex sync.Mutex

func getToken(link, user, password string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if accessToken != "" {
		return accessToken, nil
	}
	token_url := ""
	prov_url := ""

	if strings.Contains(link, "dev-cloud.aquasec.com") {
		token_url = "https://stage.api.cloudsploit.com"
		prov_url = "https://prov-dev.cloud.aquasec.com"
	} else {
		token_url = "https://api.cloudsploit.com"
		prov_url = "https://prov.cloud.aquasec.com"
	}

	req, err := http.NewRequest("POST", token_url+"/v2/signin",
		strings.NewReader(`{"email":"`+user+`", "password":"`+password+`"}`))
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("request failed. status: %s, response: %s", resp.Status, body)
	}

	var raw map[string]interface{}
	_ = json.Unmarshal(body, &raw)
	data := raw["data"].(map[string]interface{})
	token := data["token"].(string)
	//get the ese_url to make the API requests.
	requestEvents, err := http.NewRequest("GET", prov_url+"/v1/envs", nil)
	if err != nil {
		return "", err
	}
	requestEvents.Header.Set("Authorization", "Bearer "+token)
	respEvents, err := http.DefaultClient.Do(requestEvents)
	if err != nil {
		return "", err
	}
	if respEvents.StatusCode == 200 {
		eventsBody, err := ioutil.ReadAll(respEvents.Body)
		if err != nil {
			return "", err
		}
		var raw map[string]interface{}
		_ = json.Unmarshal(eventsBody, &raw)
		data := raw["data"].(map[string]interface{})
		accessUrl = "https://" + data["ese_url"].(string)
	}
	return token, nil
}

func isAquaSaasFlow(link string) bool {
	return strings.Contains(link, "cloud.aquasec.com")
}

func getData(link, user, password string) []byte {
	var req *http.Request
	if isAquaSaasFlow(link) {
		token, err := getToken(link, user, password)
		if err != nil {
			log.Fatalf("Can't receive a token to Aqua SaaS: %v", err)
		}
		u, err := url.Parse(link)
		if err != nil {
			log.Fatalf("Can't parse a link %q: %v", link, err)
		}
		req, err = http.NewRequest("GET", accessUrl+u.Path, nil)
		if err != nil {
			log.Fatalf("Can't create a request to %q: %v", accessUrl+u.Path, err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		var err error
		req, err = http.NewRequest("GET", link, nil)
		if err != nil {
			log.Fatalf("Can't create a request to %q: %v", link, err)
		}
		req.SetBasicAuth(user, password)
	}
	fmt.Println("Getting data from", req.URL)
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
