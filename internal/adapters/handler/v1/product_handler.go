package v1

import (
	"encoding/json"
	"net/http"

	"github.com/bagus-aulia/inventhier/internal/core/ports"
)

// HTTPHandler handles HTTP requests for products under API v1.
type HTTPHandler struct {
	svc ports.ProductService
}

// NewHTTPHandler creates a new HTTP handler instance.
func NewHTTPHandler(svc ports.ProductService) *HTTPHandler {
	return &HTTPHandler{
		svc: svc,
	}
}

// HandleGetProduct retrieves a product by its ID.
func (h *HTTPHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing product id", http.StatusBadRequest)
		return
	}

	product, err := h.svc.GetProduct(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(product)
}

// HandleListProducts lists all products.
func (h *HTTPHandler) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := h.svc.ListProducts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(products)
}
