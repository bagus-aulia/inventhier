package v1

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/bagus-aulia/go-tools/tools/handler/response"
	"github.com/bagus-aulia/inventhier/config"
	http "github.com/bagus-aulia/inventhier/internal/adapters/helpers/http_client"
	"github.com/bagus-aulia/inventhier/internal/core/constants"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/payment"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
)

type restPaymentClient struct {
	httpClient http.HTTPClient
	cfg        *config.Config
}

// NewRESTPaymentClient creates a new payment REST client that implements the PaymentClient port.
func NewRESTPaymentClient(
	httpClient http.HTTPClient,
	cfg *config.Config,
) ports.PaymentClient {
	return &restPaymentClient{
		httpClient: httpClient,
		cfg:        cfg,
	}
}

// Checkout processes a checkout via REST API.
func (c *restPaymentClient) Checkout(ctx context.Context, payload dto.CheckoutPayload) (response.Response[*dto.CheckoutResp], error) {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("client", "rest.payment.v1").
		Str("function", "Checkout").
		Interface("payload", payload).
		Logger()

	bodyJSON, _ := json.Marshal(payload)

	var resp response.Response[*dto.CheckoutResp]
	statusCode, err := http.DoRequest(&http.HttpOptions{
		Client:   c.httpClient,
		Ctx:      ctx,
		Hostname: c.cfg.PaymentServiceURL,
		Path:     "/api/v1/checkout",
		Timeout:  c.cfg.ClientTimeout,
		Method:   "POST",
		Body:     bodyJSON,
		Headers: map[string]string{
			constants.XAPIKey: c.cfg.PaymentServiceAPIKey,
		},
	}, &resp)
	if err != nil {
		resp.Status = 500
		resp.Error = &response.ErrorMessage{
			Message: "Cannot reach checkout endpoint",
			Reason:  err.Error(),
		}

		logger.Error().
			Err(err).
			Msg("Failed to checkout")

		return resp, err
	}

	if statusCode != 200 {
		errReason := resp.Error.Reason
		err = errors.New(errReason)

		logger.Error().
			Err(err).
			Msg("There is an error on checkout service")

		return resp, err
	}

	return resp, nil
}
