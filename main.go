package main

import (
	"./pdfrender"
	"./rest"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type strslice []string

var (
	serverUrl string
	registryName string
	imageName string
	user string
	password string
	output string
	severities strslice

	severitiesTypes = []string{
		"critical",
		"high",
		"medium",
		"low",
	}
)

const (
	cmdRegistry = "registry"
	cmdServer = "server"
	cmdImage = "image"
	cmdUser = "user"
	cmdPassword = "password"
	cmdOutput = "output"
	cmdSeverity = "severity"


)

func (str *strslice) String() string {
	return fmt.Sprintf("%s", *str)
}

func (str *strslice) Set(value string) error {
	*str = append(*str, strings.ToLower(value))
	return nil
}

func init()  {
	flag.StringVar(&serverUrl, cmdServer, "", "URL of a data server")
	flag.StringVar(&registryName, cmdRegistry, "", "name of a registry")
	flag.StringVar(&imageName, cmdImage, "", "name of an image")
	flag.StringVar(&user, cmdUser, "", "a user for the basic authentication")
	flag.StringVar(&password, cmdPassword, "", "a user's password for the basic authentication")
	flag.StringVar(&output, cmdOutput, "report.pdf", "a name of output pdf file")
	flag.Var( &severities, cmdSeverity, "to get list of vulnerabilities. critical/high/medium/low" )
}

func checkRequiredParams() bool {
	var missingRequiredFlags []string
	if registryName == "" {
		missingRequiredFlags = append(missingRequiredFlags, cmdRegistry)
	}
	if imageName == "" {
		missingRequiredFlags = append(missingRequiredFlags, cmdImage)
	}

	if serverUrl == "" {
		if serverUrl = os.Getenv("server"); serverUrl == "" {
			fmt.Println("Server isn't setup (as -server param or environment variable)")
			return false
		}
	}

	if user == ""  {
		if user = os.Getenv("user"); user == "" {
			fmt.Println("User isn't setup (as -user param or environment variable)")
			return false
		}

	}
	if password == "" {
		if password = os.Getenv("password");password == "" {
			fmt.Println("Password isn't setup (as -password param or environment variable)")
			return false
		}
	}

	if len(missingRequiredFlags) > 0 {
		message := "Param '%s' is missing or the value is empty.\n"
		for _, f := range missingRequiredFlags {
			fmt.Printf(message, f)
		}
		return false
	}

	if severities != nil {
		for _, severity := range severities {
			var count int
			for _,v := range severitiesTypes {
				if severity == v {
					break
				}
				count++
			}
			if count == len(severitiesTypes) {
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

	if ok:=checkRequiredParams(); !ok {
		fmt.Println("All params are required!")
		fmt.Println("Run with key '-h' for usage.")
		os.Exit(1)
	}

	var filename string
	if !strings.HasSuffix(output, ".pdf") {
		filename = output + ".pdf"
	} else {
		filename = output
	}
	data := rest.GetData(serverUrl, user, password, registryName, imageName, severities)
	data.ServerUrl = serverUrl
	err := pdfrender.Render(filename, data)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Report was written to", filename)
}
