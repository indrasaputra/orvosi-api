package router

import (
	"net/http"

	"github.com/indrasaputra/orvosi-api/internal/http/handler"
)

// Signer creates routes for sign in.
func Signer(h *handler.Signer) []*Route {
	var routes []*Route

	r := &Route{
		Method:  http.MethodPost,
		Path:    "/sign-in",
		Handler: h.SignIn,
	}

	routes = append(routes, r)
	return routes
}
