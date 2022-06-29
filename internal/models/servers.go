package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/providers"
)

type AllServers struct {
	Version           uint16 // used for migration of the top level scheme
	ProviderToServers map[string]Servers
}

var _ json.Marshaler = (*AllServers)(nil)

// MarshalJSON marshals all servers to JSON.
// Note you need to use a pointer to all servers
// for it to work with native json methods such as
// json.Marshal.
func (a *AllServers) MarshalJSON() (data []byte, err error) {
	buffer := bytes.NewBuffer(nil)

	_, err = buffer.WriteString("{")
	if err != nil {
		return nil, fmt.Errorf("cannot write opening bracket: %w", err)
	}

	versionString := fmt.Sprintf(`"version":%d`, a.Version)
	_, err = buffer.WriteString(versionString)
	if err != nil {
		return nil, fmt.Errorf("cannot write schema version string: %w", err)
	}

	sortedProviders := make(sort.StringSlice, 0, len(a.ProviderToServers))
	for provider := range a.ProviderToServers {
		sortedProviders = append(sortedProviders, provider)
	}
	sortedProviders.Sort()

	for _, provider := range sortedProviders {
		providerKey := fmt.Sprintf(`,"%s":`, provider)
		_, err = buffer.WriteString(providerKey)
		if err != nil {
			return nil, fmt.Errorf("cannot write provider key %s: %w",
				providerKey, err)
		}

		servers := a.ProviderToServers[provider]
		serversJSON, err := json.Marshal(servers)
		if err != nil {
			return nil, fmt.Errorf("failed encoding servers for provider %s: %w",
				provider, err)
		}
		_, err = buffer.Write(serversJSON)
		if err != nil {
			return nil, fmt.Errorf("cannot write JSON servers data for provider %s: %w",
				provider, err)
		}
	}

	_, err = buffer.WriteString("}")
	if err != nil {
		return nil, fmt.Errorf("cannot write closing bracket: %w", err)
	}

	return buffer.Bytes(), nil
}

var _ json.Unmarshaler = (*AllServers)(nil)

func (a *AllServers) UnmarshalJSON(data []byte) (err error) {
	keyValues := make(map[string]interface{})
	err = json.Unmarshal(data, &keyValues)
	if err != nil {
		return err
	}

	versionUnmarshaled := keyValues["version"]
	if versionUnmarshaled != nil { // defaults to 0
		version, ok := versionUnmarshaled.(float64)
		if !ok {
			return &json.UnmarshalTypeError{
				Value:  fmt.Sprintf("number %v", versionUnmarshaled),
				Type:   reflect.TypeOf(uint16(0)),
				Struct: "models.AllServers",
				Field:  "Version",
			}
		}

		if math.Round(version) != version ||
			version < 0 || version > float64(^uint16(0)) {
			return &json.UnmarshalTypeError{
				Value:  fmt.Sprintf("number %v", version),
				Type:   reflect.TypeOf(uint16(0)),
				Struct: "models.AllServers",
				Field:  "Version",
			}
		}

		a.Version = uint16(version)
		delete(keyValues, "version")
	}

	if len(keyValues) == 0 {
		return nil
	}

	a.ProviderToServers = make(map[string]Servers, len(keyValues))

	allProviders := providers.All()
	allProvidersSet := make(map[string]struct{}, len(allProviders))
	for _, provider := range allProviders {
		allProvidersSet[provider] = struct{}{}
	}

	for key, value := range keyValues {
		if _, ok := allProvidersSet[key]; !ok {
			// not a provider known by Gluetun
			// or a non-servers field.
			continue
		}

		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("cannot marshal %s servers: %w",
				key, err)
		}

		var servers Servers
		err = json.Unmarshal(jsonValue, &servers)
		if err != nil {
			return fmt.Errorf("cannot unmarshal %s servers: %w",
				key, err)
		}

		a.ProviderToServers[key] = servers
	}

	return nil
}

func (a *AllServers) Count() (count int) {
	for _, servers := range a.ProviderToServers {
		count += len(servers.Servers)
	}
	return count
}

type Servers struct {
	Version   uint16   `json:"version"`
	Timestamp int64    `json:"timestamp"`
	Servers   []Server `json:"servers,omitempty"`
}
