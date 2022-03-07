package rest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/url"
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

func getData(link, user, password string) []byte {
	fmt.Println("Getting data from", link)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatalf("Can't create a request to %q: %v", link, err)
	}
	if loginAquaSaas(user) {
		u, err := url.Parse(link)
		if err != nil {
			log.Fatalf("Can't parse a link %q: %v", link, err)
		}
		token, err := aquaClient.NewClient(u.Host, user, password, strings.HasPrefix(link, "https"), nil).GetUSEAuthToken()
		if err != nil {
			log.Fatalf("Can't get the Aqua SaaS token for access to %q: %v", link, err)
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
		fmt.Println("Wrong access to the URL: ", link)
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
