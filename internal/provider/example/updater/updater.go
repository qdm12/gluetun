package updater

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	// TODO: remove fields not used by the updater
	client           *http.Client
	unzipper         common.Unzipper
	parallelResolver common.ParallelResolver
	warner           common.Warner
}

func New(warner common.Warner, unzipper common.Unzipper,
	client *http.Client, parallelResolver common.ParallelResolver) *Updater {
	// TODO: remove arguments not used by the updater
	return &Updater{
		client:           client,
		unzipper:         unzipper,
		parallelResolver: parallelResolver,
		warner:           warner,
	}
}
