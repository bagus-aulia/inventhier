package v1

import (
	"net/http"

	productHandler "github.com/bagus-aulia/inventhier/internal/adapters/handler/v1/product"
	"github.com/bagus-aulia/inventhier/internal/adapters/middleware"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
	"github.com/gorilla/mux"
)

// Router maps v1 endpoints.
type Router struct {
	router *mux.Router
	svc    ports.ProductService
}

// NewRouter initializes and registers v1 routes.
func NewRouter(svc ports.ProductService) *Router {
	r := &Router{
		router: mux.NewRouter(),
		svc:    svc,
	}
	r.registerRoutes()
	return r
}

func (r *Router) registerRoutes() {
	ph := productHandler.NewHTTPHandler(r.svc)

	// V1
	apiV1 := r.router.PathPrefix("/api/v1").Subrouter()

	// Product routes
	product := apiV1.PathPrefix("/product").Subrouter()
	product.HandleFunc("/stock-in", ph.StockIn).Methods("POST")
}

// GetHandler returns the router with middleware applied
func (r *Router) GetHandler() http.Handler {
	handler := http.Handler(r.router)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.RequestIDMiddleware(handler)
	return handler
}
