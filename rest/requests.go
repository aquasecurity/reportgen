package rest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"

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

func loginAquaSaas(user string) bool {
	_, err := mail.ParseAddress(user)
	return err == nil
}

func getData(url, user, password string) []byte {
	fmt.Println("Getting data from", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Can't create a request to %q: %v", url, err)
	}
	if loginAquaSaas(user) {
		token, err := aquaClient.NewClient(url, user, password, strings.HasPrefix(url, "https"), nil).GetUSEAuthToken()
		if err != nil {
			log.Fatalf("Can't get the Aqua SaaS token for access to %q: %v", url, err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req.SetBasicAuth(user, password)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
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
	if strings.HasSuffix(server, "/") {
		slash = ""
	} else {
		slash = "/"
	}
	return server + slash + api
}
