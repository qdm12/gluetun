package updater

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	client           *http.Client
	parallelResolver common.ParallelResolver
	warner           common.Warner
	storage          common.Storage
}

func New(client *http.Client, warner common.Warner,
	parallelResolver common.ParallelResolver, storage common.Storage,
) *Updater {
	return &Updater{
		client:           client,
		parallelResolver: parallelResolver,
		warner:           warner,
		storage:          storage,
	}
}
