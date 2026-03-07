package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/qdm12/gluetun/ci/internal"
	"github.com/qdm12/log"
)

func main() {
	logger := log.New()
	if len(os.Args) < 2 {
		logger.Error("Usage: " + os.Args[0] + " <command>")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	var err error
	switch os.Args[1] {
	case "mullvad":
		err = internal.MullvadTest(ctx, logger)
	case "protonvpn":
		err = internal.ProtonVPNTest(ctx, logger)
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}
	stop()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("test completed successfully")
}
