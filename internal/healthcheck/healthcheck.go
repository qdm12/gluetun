package healthcheck

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/golibs/network/connectivity"
)

func HealthCheck() error {
	// DNS, HTTP and HTTPs check on github.com
	connectivty := connectivity.NewConnectivity(3 * time.Second)
	errs := connectivty.Checks("github.com")
	if len(errs) > 0 {
		var errsStr []string
		for _, err := range errs {
			errsStr = append(errsStr, err.Error())
		}
		return fmt.Errorf("Multiple errors: %s", strings.Join(errsStr, "; "))
	}
	// TODO check IP address is in the right region
	return nil
}
