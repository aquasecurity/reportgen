package main

import (
	"flag"
	"fmt"
	"github.com/aquasecurity/reportgen/data"
	"github.com/aquasecurity/reportgen/pdfrender"
	"github.com/aquasecurity/reportgen/rest"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

var (
	serverUrl      string
	registryName   string
	imageName      string
	host           string
	user           string
	password       string
	output         string
	severities     []string
	severityParams string

	severitiesTypes = map[string]struct{}{
		"critical": {},
		"high":     {},
		"medium":   {},
		"low":      {},
	}
)

const (
	cmdRegistry = "registry"
	cmdServer   = "server"
	cmdImage    = "image"
	cmdHost     = "host"
	cmdUser     = "user"
	cmdPassword = "password"
	cmdOutput   = "output"
	cmdSeverity = "severity"
)

func init() {
	flag.StringVar(&serverUrl, cmdServer, "", "URL of a data server")
	flag.StringVar(&registryName, cmdRegistry, "", "name of a registry")
	flag.StringVar(&imageName, cmdImage, "", "name of an image")
	flag.StringVar(&user, cmdUser, "", "a user for the basic authentication")
	flag.StringVar(&password, cmdPassword, "", "a user's password for the basic authentication")
	flag.StringVar(&output, cmdOutput, "report.pdf", "a name of output pdf file")
	flag.StringVar(&severityParams, cmdSeverity, "", "to get list of vulnerabilities. critical,high,medium,low")
	flag.StringVar(&host, cmdHost, "", "PDF generation to a host name")
}

func checkRequiredParams() bool {
	if (host != "" && imageName != "") || (host == "" && imageName == "") {
		fmt.Println("Wrong params: you should setup either a host or an image!")
		return false
	}

	if serverUrl == "" {
		if serverUrl = os.Getenv("server"); serverUrl == "" {
			fmt.Println("Server isn't setup (as -server param or environment variable)")
			return false
		}
	}

	if user == "" {
		if user = os.Getenv("user"); user == "" {
			fmt.Println("User isn't setup (as -user param or environment variable)")
			return false
		}
	}
	if password == "" {
		if password = os.Getenv("password"); password == "" {
			fmt.Println("Password isn't setup (as -password param or environment variable)")
			return false
		}
	}

	if host != "" {
		return true
	}

	var missingRequiredFlags []string
	if registryName == "" {
		missingRequiredFlags = append(missingRequiredFlags, cmdRegistry)
	}
	if imageName == "" {
		missingRequiredFlags = append(missingRequiredFlags, cmdImage)
	}
	if len(missingRequiredFlags) > 0 {
		message := "Param '%s' is missing or the value is empty.\n"
		for _, f := range missingRequiredFlags {
			fmt.Printf(message, f)
		}
		return false
	}

	if severityParams != "" {
		severities = strings.Split(strings.ToLower(severityParams), ",")
		for _, severity := range severities {
			if _, ok := severitiesTypes[severity]; !ok {
				fmt.Println("Wrong severity type:", severity)
				return false
			}
		}
	}
	return true
}

func main() {
	flag.Parse()
	godotenv.Load()

	if ok := checkRequiredParams(); !ok {
		fmt.Println("Run with key '-h' for usage.")
		os.Exit(1)
	}

	var filename string
	if !strings.HasSuffix(output, ".pdf") {
		filename = output + ".pdf"
	} else {
		filename = output
	}

	var report *data.Report
	if imageName != "" {
		report = rest.GetImageData(serverUrl, user, password, registryName, imageName, severities)
	} else if host != "" {
		report = rest.GetHostData(serverUrl, user, password, host)
	}

	report.ServerUrl = serverUrl
	err := pdfrender.Render(filename, report)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Report was written to", filename)
}
