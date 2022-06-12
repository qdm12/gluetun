package custom

import "github.com/qdm12/gluetun/internal/models"

type Extractor interface {
	Data(filepath string) (lines []string,
		connection models.Connection, err error)
}
