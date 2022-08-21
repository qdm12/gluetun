package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

func newPublicIPHandler(loop PublicIPLoop, w warner) http.Handler {
	return &publicIPHandler{
		loop:   loop,
		warner: w,
	}
}

type publicIPHandler struct {
	loop   PublicIPLoop
	warner warner
}

func (h *publicIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.TrimPrefix(r.RequestURI, "/publicip")
	switch r.RequestURI {
	case "/ip":
		switch r.Method {
		case http.MethodGet:
			h.getPublicIP(w)
		default:
			http.Error(w, "method "+r.Method+" not supported", http.StatusBadRequest)
		}
	default:
		http.Error(w, "route "+r.RequestURI+" not supported", http.StatusBadRequest)
	}
}

func (h *publicIPHandler) getPublicIP(w http.ResponseWriter) {
	data := h.loop.GetData()
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		h.warner.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
