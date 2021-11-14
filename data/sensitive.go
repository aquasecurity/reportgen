package data

type SensitiveResult struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type SensitiveType struct {
	Count   int               `json:"count"`
	Results []SensitiveResult `json:"result"`
}
