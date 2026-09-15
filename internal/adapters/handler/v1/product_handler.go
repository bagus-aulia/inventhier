package v1

import (
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
	// ctx := r.Context()

	// Log the incoming request with request ID
	// helpers.LogInfo(ctx, "Retrieving product")

	// id := r.URL.Query().Get("id")
	// if id == "" {
	// // 	// helpers.LogWarn(ctx, "Product ID parameter is missing")
	// 	// http.Error(w, "missing product id", http.StatusBadRequest)
	// 	return
	// }

	// // Log request details
	// helpers.LogDebug(ctx, "Looking up product with ID: "+id)

	// product, err := h.svc.GetProduct(ctx, id)
	// if err != nil {
	// 	helpers.LogError(ctx, "Failed to retrieve product", err)
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }

	// Log success
	// helpers.LogInfoWithFields(ctx, "Product retrieved successfully", map[string]interface{}{
	// 	"product_id": product.ID,
	// 	"name": product.Name,
	// })

	// w.Header().Set("Content-Type", "application/json")
	// _ = json.NewEncoder(w).Encode(product)
}

// HandleListProducts lists all products.
func (h *HTTPHandler) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// // Log the incoming request with request ID
	// helpers.LogInfo(ctx, "Listing all products")

	// products, err := h.svc.ListProducts(ctx)
	// if err != nil {
	// 	helpers.LogError(ctx, "Failed to list products", err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Log success
	// helpers.LogInfoWithFields(ctx, "Products listed successfully", map[string]interface{}{
	// 	"total": len(products),
	// })

	// w.Header().Set("Content-Type", "application/json")
	// _ = json.NewEncoder(w).Encode(products)
}
