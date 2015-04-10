package api

import (
	"log"
	"net/http"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/gorilla/mux"
)

type insecureServer struct {
	Manager agent.Manager
}

// MakeInsecureServer returns a new Server instance containting a manager to which it will defer work.
func MakeInsecureServer(dm agent.Manager) Server {
	return insecureServer{
		Manager: dm,
	}
}

func (s insecureServer) Start(addr string) {
	r := s.newRouter()

	log.Fatal(http.ListenAndServe(addr, r))
}

func (s insecureServer) newRouter() *mux.Router {
	return newRouter(s.Manager, s.isAuthenticated)
}

func (s insecureServer) isAuthenticated(r *http.Request) bool {
	return true
}
