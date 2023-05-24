package updater

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/netip"
)

// Check out the JSON data from https://api.nordvpn.com/v1/servers?limit=10
type serverData struct {
	// Name is the server name, for example 'Poland #128'
	Name string `json:"name"`
	// Stations is, it seems, the entry IP address.
	// However it is ignored in favor of the 'ips' entry field.
	Station netip.Addr `json:"station"`
	// IPv6Station is mostly empty, so we ignore it for now.
	IPv6Station netip.Addr `json:"station_ipv6"`
	// Hostname is the server hostname, for example 'pl128.nordvpn.com'
	Hostname string
	// Status is the server status, for example 'online'
	Status string `json:"status"`
	// Locations is the list of locations for the server.
	// Only the first location is taken into account for now.
	Locations []struct {
		Country struct {
			// Name is the country name, for example 'Poland'.
			Name string `json:"name"`
			City struct {
				// Name is the city name, for example 'Warsaw'.
				Name string `json:"name"`
			} `json:"city"`
		} `json:"country"`
	} `json:"locations"`
	Technologies []struct {
		// Identifier is the technology id name, it can notably be:
		// - openvpn_udp
		// - openvpn_tcp
		// - wireguard_udp
		Identifier string `json:"identifier"`
		// Metadata is notably useful for the Wireguard public key.
		Metadata []struct {
			// Name can notably be 'public_key'.
			Name string `json:"name"`
			// Value can notably the Wireguard public key value.
			Value string `json:"value"`
		} `json:"metadata"`
	} `json:"technologies"`
	Groups []struct {
		// Title can notably be the region name, for example 'Europe',
		// if the group's type/identifier is 'regions'.
		Title string `json:"title"`
		Type  struct {
			// Identifier can be 'regions'.
			Identifier string `json:"identifier"`
		} `json:"type"`
	} `json:"groups"`
	// IPs is the list of IP addresses for the server.
	IPs []struct {
		// Type can notably be 'entry'.
		Type string `json:"type"`
		IP   struct {
			IP netip.Addr `json:"ip"`
		} `json:"ip"`
	} `json:"ips"`
}

// country returns the country name of the server.
func (s *serverData) country() (country string) {
	if len(s.Locations) == 0 {
		return ""
	}
	return s.Locations[0].Country.Name
}

// region returns the region name of the server.
func (s *serverData) region() (region string) {
	for _, group := range s.Groups {
		if group.Type.Identifier == "regions" {
			return group.Title
		}
	}
	return ""
}

// city returns the city name of the server.
func (s *serverData) city() (city string) {
	if len(s.Locations) == 0 {
		return ""
	}
	return s.Locations[0].Country.City.Name
}

// ips returns the list of IP addresses for the server.
func (s *serverData) ips() (ips []netip.Addr) {
	ips = make([]netip.Addr, 0, len(s.IPs))
	for _, ipObject := range s.IPs {
		if ipObject.Type != "entry" {
			continue
		}
		ips = append(ips, ipObject.IP.IP)
	}
	return ips
}

var (
	ErrWireguardPublicKeyMalformed = errors.New("wireguard public key is malformed")
	ErrWireguardPublicKeyNotFound  = errors.New("wireguard public key not found")
)

// wireguardPublicKey returns the Wireguard public key for the server.
func (s *serverData) wireguardPublicKey() (wgPubKey string, err error) {
	for _, technology := range s.Technologies {
		if technology.Identifier != "wireguard_udp" {
			continue
		}
		for _, metadata := range technology.Metadata {
			if metadata.Name != "public_key" {
				continue
			}
			wgPubKey = metadata.Value
			_, err = base64.StdEncoding.DecodeString(wgPubKey)
			if err != nil {
				return "", fmt.Errorf("%w: %s cannot be decoded: %s",
					ErrWireguardPublicKeyMalformed, wgPubKey, err)
			}
			return metadata.Value, nil
		}
	}
	return "", fmt.Errorf("%w", ErrWireguardPublicKeyNotFound)
}
