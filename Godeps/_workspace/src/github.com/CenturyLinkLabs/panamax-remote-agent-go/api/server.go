package api

import (
	"github.com/gorilla/mux"
)

// A Server is the HTTP server which responds to API requests.
type Server interface {
	Start(string)
	newRouter() *mux.Router
}
