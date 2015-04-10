package api

import (
	"log"
	"net/http"
	"time"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/gorilla/mux"
)

// A Server is the HTTP server which responds to API requests.
type Server interface {
	Start(string)
	newRouter() *mux.Router
}

func newRouter(dm agent.Manager, isAuthenticated func(r *http.Request) bool) *mux.Router {
	r := mux.NewRouter()

	for _, route := range routes {
		fct := route.HandlerFunc
		wrap := func(w http.ResponseWriter, r *http.Request) {

			if !isAuthenticated(r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// make it json
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			// log it
			st := time.Now()

			log.Printf(
				"%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				route.Name,
				time.Since(st),
			)

			fct(dm, w, r)
		}

		r.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(wrap)
	}

	return r
}
