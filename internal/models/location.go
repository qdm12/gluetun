package models

// SurfsharkLocationData is required to keep location data on Surfshark
// servers that are not obtained through their API.
type SurfsharkLocationData struct {
	Region   string
	Country  string
	City     string
	RetroLoc string // TODO remove in v4
	Hostname string
	MultiHop bool
}
