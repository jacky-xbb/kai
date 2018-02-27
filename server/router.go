package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Router represents the routers of mux
type Router struct {
	r *mux.Router
}

// RouterArguments represents the mapping between path and method
type RouterArguments struct {
	Handler http.HandlerFunc
	Path    string
	Method  string
}

// NewRouter creates an new the router of mux
func NewRouter() *Router {
	return &Router{
		r: mux.NewRouter(),
	}
}

// Handler returns the handler of the router
func (router *Router) Handler() http.Handler {
	return router.r
}

// AddHandler adds a handler for a router
func (router *Router) AddHandler(args RouterArguments) {
	path := fmt.Sprintf("/%s", strings.Trim(args.Path, "/"))
	router.r.
		Methods(args.Method).
		Path(fmt.Sprintf("%s", path)).
		HandlerFunc(JSONHandler(args.Handler))
}
