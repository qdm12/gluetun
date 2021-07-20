package models

type BuildInformation struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Created string `json:"created"`
}
