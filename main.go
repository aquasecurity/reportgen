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

var (
	serverUrl string
	registryName string
	imageName string
	user string
	password string
	output string
)

const (
	cmdRegistry = "registry"
	cmdServer = "server"
	cmdImage = "image"
	cmdUser = "user"
	cmdPassword = "password"
	cmdOutput = "output"
)

func init()  {
	flag.StringVar(&serverUrl, cmdServer, "", "URL of a data server")
	flag.StringVar(&registryName, cmdRegistry, "", "name of a registry")
	flag.StringVar(&imageName, cmdImage, "", "name of an image")
	flag.StringVar(&user, cmdUser, "", "a user for the basic authentication")
	flag.StringVar(&password, cmdPassword, "", "a user's password for the basic authentication")
	flag.StringVar(&output, cmdOutput, "report.pdf", "a name of output pdf file")
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

	fmt.Println("Debug info:")
	fmt.Println("Server:", serverUrl)
	fmt.Println("Registry:", registryName)
	fmt.Println("Image:", imageName)
	fmt.Println("User:", user)
	fmt.Println("Password:", password)

	var filename string
	if !strings.HasSuffix(output, ".pdf") {
		filename = output + ".pdf"
	} else {
		filename = output
	}
	data := rest.GetData(serverUrl, user, password, registryName, imageName)
	data.ServerUrl = serverUrl
	pdfrender.Render(filename, data)
	fmt.Println("Report was written to", filename)
}
