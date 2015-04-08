package api

import (
	"log"
	"net/http"
	"time"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/gorilla/mux"
)

type secureServer struct {
	Manager  agent.Manager
	username string
	password string
	certFile string
	keyFile  string
}

// MakeServer returns a new Server instance containting a manager to which it will defer work.
func MakeServer(dm agent.Manager, un string, pw string, cf string, kf string) Server {
	return secureServer{
		Manager:  dm,
		username: un,
		password: pw,
		certFile: cf,
		keyFile:  kf,
	}
}

func (s secureServer) Start(addr string) {
	r := s.newRouter()

	log.Fatal(http.ListenAndServeTLS(addr, s.certFile, s.keyFile, r))
}

func (s secureServer) newRouter() *mux.Router {
	r := mux.NewRouter()

	dm := s.Manager

	for _, route := range routes {
		fct := route.HandlerFunc
		wrap := func(w http.ResponseWriter, r *http.Request) {

			if !s.isAuthenticated(r) {
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

func (s secureServer) isAuthenticated(r *http.Request) bool {
	un, pw, ok := r.BasicAuth()

	if ok && (un == s.username) && (pw == s.password) {
		return true
	}

	return false
}
