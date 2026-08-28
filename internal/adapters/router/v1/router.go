package v1

import (
	"net/http"

	v1Handler "github.com/bagus-aulia/inventhier/internal/adapters/handler/v1"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
)

// Router maps v1 endpoints.
type Router struct {
	mux *http.ServeMux
	svc ports.ProductService
}

// NewRouter initializes and registers v1 routes.
func NewRouter(mux *http.ServeMux, svc ports.ProductService) *Router {
	r := &Router{
		mux: mux,
		svc: svc,
	}
	r.registerRoutes()
	return r
}

func (r *Router) registerRoutes() {
	h := v1Handler.NewHTTPHandler(r.svc)

	r.mux.HandleFunc("/api/v1/product", h.HandleGetProduct)
	r.mux.HandleFunc("/api/v1/products", h.HandleListProducts)
}
