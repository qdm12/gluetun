package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/qdm12/gluetun/pkg/updaters"
	"github.com/qdm12/log"
)

func main() {
	defaultOutputDir := filepath.Join("pkg", "servers")
	outputDir := flag.String("output-dir", defaultOutputDir,
		"directory to write per-provider JSON files to")
	protonEmail := flag.String("proton-email", os.Getenv("GLUETUN_PROTON_EMAIL"),
		"ProtonVPN account email")
	protonPassword := flag.String("proton-password", os.Getenv("GLUETUN_PROTON_PASSWORD"),
		"ProtonVPN account password")
	ipinfoToken := flag.String("ipinfo-token", os.Getenv("IPINFO_TOKEN"),
		"IPInfo token used for public IP lookups")
	flag.Parse()

	client := &http.Client{Timeout: time.Minute}
	settings := updaters.UpdateAllSettings{
		OutputPath:     outputDir,
		ProtonEmail:    protonEmail,
		ProtonPassword: protonPassword,
		IpinfoToken:    ipinfoToken,
	}

	logger := log.New()
	err := updaters.UpdateAll(context.Background(), client, logger, settings)
	if err != nil {
		logger.Errorf("update failed: %v", err)
		os.Exit(1)
	}

	logger.Infof("update completed successfully, data files written to %s", *outputDir)
}
