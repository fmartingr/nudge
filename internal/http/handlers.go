package http

import (
	"net/http"

	"github.com/fmartingr/nudge/internal/pinger"
)

func healthHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(""))
}

func statusHandler(ping *pinger.Pinger) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		if ping.Ok() {
			rw.WriteHeader(200)
		} else {
			rw.WriteHeader(204)
		}
	}
}
