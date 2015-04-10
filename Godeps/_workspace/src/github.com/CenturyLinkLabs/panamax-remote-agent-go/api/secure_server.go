package api

import (
	"log"
	"net/http"

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
	return newRouter(s.Manager, s.isAuthenticated)
}

func (s secureServer) isAuthenticated(r *http.Request) bool {
	un, pw, ok := r.BasicAuth()

	if ok && (un == s.username) && (pw == s.password) {
		return true
	}

	return false
}
