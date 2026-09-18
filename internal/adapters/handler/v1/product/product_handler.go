package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bagus-aulia/go-tools/tools/handler/response"
	authHelper "github.com/bagus-aulia/inventhier/internal/adapters/helpers/auth"
	productDTO "github.com/bagus-aulia/inventhier/internal/core/dto/product"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
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

// StockIn to input product stock in
func (h *HTTPHandler) StockIn(w http.ResponseWriter, r *http.Request) {
	logger := helpers.GetZerologWithContext(r.Context()).
		With().
		Str("handler", "v1.product").
		Str("function", "StockIn").
		Logger()

	resp := response.Factory[string]()

	authToken, err := authHelper.ExtractTokenFromHeader(r.Header)
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Empty or invalid auth token")

		resp.BadRequest().
			ErrorReason("invalid auth token")

		respJSON, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		w.Write(respJSON)

		return
	}

	staffUUID, err := authHelper.ExtractUserID(authToken)
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("invalid auth token")

		resp.BadRequest().
			ErrorReason("invalid auth token")

		respJSON, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		w.Write(respJSON)

		return
	}

	var req productDTO.StockInReq
	json.NewDecoder(r.Body).Decode(&req)

	if req.ProductSKU == "" {
		req.ProductSKU = r.URL.Query().Get("product_sku")
	}

	if req.ProductSKU == "" {
		logger.Warn().
			Msg("Empty product sku param")

		resp.BadRequest().
			ErrorReason("empty product sku param")

		respJSON, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		w.Write(respJSON)

		return
	}

	if req.ProductQty == 0 {
		productQty := r.URL.Query().Get("product_qty")
		qtyInt, err := strconv.Atoi(productQty)
		if err != nil {
			logger.Warn().
				Err(err).
				Msg("Product quantity param is not integer")

			resp.BadRequest().
				ErrorReason("Product quantity param is not integer")

			respJSON, _ := json.Marshal(resp)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.Status)
			w.Write(respJSON)

			return
		}

		req.ProductQty = qtyInt
	}

	if req.ProductQty == 0 {
		logger.Warn().
			Msg("Zero product quantity param")

		resp.BadRequest().
			ErrorReason("Zero product quantity param")

		respJSON, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		w.Write(respJSON)

		return
	}

	payload := productDTO.StockInPayload{
		StaffUUID:  staffUUID,
		StockInReq: req,
	}
	err = h.svc.StockIn(r.Context(), payload)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to restock product")

		resp.InternalServerError().
			ErrorReason(err.Error())

		respJSON, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Status)
		w.Write(respJSON)

		return
	}

	resp.Success().
		WithData("Product successfully restock")

	respJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	w.Write(respJSON)
}
