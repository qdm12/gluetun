package updater

import (
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	ipFetcher        common.IPFetcher
	unzipper         common.Unzipper
	parallelResolver common.ParallelResolver
	warner           common.Warner
}

func New(ipFetcher common.IPFetcher, unzipper common.Unzipper,
	warner common.Warner, parallelResolver common.ParallelResolver) *Updater {
	return &Updater{
		ipFetcher:        ipFetcher,
		unzipper:         unzipper,
		parallelResolver: parallelResolver,
		warner:           warner,
	}
}
