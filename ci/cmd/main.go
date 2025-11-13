package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/qdm12/gluetun/ci/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: " + os.Args[0] + " <command>")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	var err error
	switch os.Args[1] {
	case "mullvad":
		err = internal.MullvadTest(ctx)
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}
	stop()
	if err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
	fmt.Println("✅ Test completed successfully.")
}
