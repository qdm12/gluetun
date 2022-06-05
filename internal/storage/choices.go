package storage

import (
	"github.com/qdm12/gluetun/internal/configuration/settings/validation"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
)

func (s *Storage) GetFilterChoices(provider string) models.FilterChoices {
	if provider == providers.Custom {
		return models.FilterChoices{}
	}

	s.mergedMutex.RLock()
	defer s.mergedMutex.RUnlock()

	serversObject := s.getMergedServersObject(provider)
	servers := serversObject.Servers
	return models.FilterChoices{
		Countries: validation.ExtractCountries(servers),
		Regions:   validation.ExtractRegions(servers),
		Cities:    validation.ExtractCities(servers),
		ISPs:      validation.ExtractISPs(servers),
		Names:     validation.ExtractServerNames(servers),
		Hostnames: validation.ExtractHostnames(servers),
	}
}
