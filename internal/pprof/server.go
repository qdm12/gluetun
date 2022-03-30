package pprof

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/qdm12/gluetun/internal/httpserver"
)

// New creates a new Pprof server and configure profiling
// with the settings given. It returns an error
// if one of the settings is not valid.
func New(settings Settings) (server *httpserver.Server, err error) {
	runtime.SetBlockProfileRate(settings.BlockProfileRate)
	runtime.SetMutexProfileFraction(settings.MutexProfileRate)

	handler := http.NewServeMux()
	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	handler.Handle("/debug/pprof/block", pprof.Handler("block"))
	handler.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	handler.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	handler.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	settings.HTTPServer.Handler = handler

	settings.SetDefaults()
	if err = settings.Validate(); err != nil {
		return nil, fmt.Errorf("pprof settings failed validation: %w", err)
	}

	return httpserver.New(settings.HTTPServer)
}
