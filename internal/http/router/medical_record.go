package router

import (
	"net/http"

	"github.com/indrasaputra/orvosi-api/internal/http/handler"
	"github.com/indrasaputra/orvosi-api/internal/http/middleware"
	"github.com/labstack/echo/v4"
)

// MedicalRecordCreator creates routes for medical record creator.
func MedicalRecordCreator(h *handler.MedicalRecordCreator) []*Route {
	var routes []*Route

	r := &Route{
		Method:      http.MethodPost,
		Path:        "/medical-records",
		Handler:     h.Create,
		Middlewares: []echo.MiddlewareFunc{middleware.WithContentType(echo.MIMEApplicationJSON)},
	}

	routes = append(routes, r)
	return routes
}

// MedicalRecordFinder creates routes for medical record finder.
func MedicalRecordFinder(h *handler.MedicalRecordFinder) []*Route {
	var routes []*Route

	fbe := &Route{
		Method:  http.MethodGet,
		Path:    "/medical-records",
		Handler: h.FindByEmail,
	}

	fbi := &Route{
		Method:  http.MethodGet,
		Path:    "/medical-records/:id",
		Handler: h.FindByID,
	}

	routes = append(routes, fbe, fbi)
	return routes
}

// MedicalRecordUpdater creates routes for medical record updater.
func MedicalRecordUpdater(h *handler.MedicalRecordUpdater) []*Route {
	var routes []*Route

	r := &Route{
		Method:      http.MethodPut,
		Path:        "/medical-records/:id",
		Handler:     h.Update,
		Middlewares: []echo.MiddlewareFunc{middleware.WithContentType(echo.MIMEApplicationJSON)},
	}

	routes = append(routes, r)
	return routes
}
