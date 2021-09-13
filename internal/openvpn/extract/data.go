package extract

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

var (
	ErrRead              = errors.New("cannot read file")
	ErrExtractConnection = errors.New("cannot extract connection from file")
)

// Data extracts the lines and connection from the OpenVPN configuration file.
func (e *Extractor) Data(filepath string) (lines []string,
	connection models.Connection, err error) {
	lines, err = readCustomConfigLines(filepath)
	if err != nil {
		return nil, connection, fmt.Errorf("%w: %s", ErrRead, err)
	}

	connection, err = extractDataFromLines(lines)
	if err != nil {
		return nil, connection, fmt.Errorf("%w: %s", ErrExtractConnection, err)
	}

	return lines, connection, nil
}
