package extract

import (
	"github.com/qdm12/gluetun/internal/models"
)

var _ Interface = (*Extractor)(nil)

type Interface interface {
	Data(filepath string) (lines []string,
		connection models.Connection, err error)
}

type Extractor struct{}

func New() *Extractor {
	return new(Extractor)
}
