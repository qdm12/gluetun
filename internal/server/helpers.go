package server

import (
	"net/http"
)

func errMethodNotSupported(w http.ResponseWriter, method string) {
	http.Error(w, "method "+method+" not supported", http.StatusBadRequest)
}

func errRouteNotSupported(w http.ResponseWriter, route string) {
	http.Error(w, "route "+route+" not supported", http.StatusBadRequest)
}
