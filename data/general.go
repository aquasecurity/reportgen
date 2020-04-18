package data

type General struct {
	Name string `json:"name"`
	Registry string `json:"registry"`
	Malware int `json:"malware"`
	Disallowed bool `json:"disallowed"`
	Os string `json:"os"`
	OsVersion string `json:"os_version"`
	Created string `json:"created"`

	/*
	"created":"2019-06-12T14:22:51.717668Z",
	"os":"alpine",
	"os_version":"3.3.3",
	
	"registry":"Docker Hub",
	"name":"alpine:3.3",
	"malware":0,

	 */
}
