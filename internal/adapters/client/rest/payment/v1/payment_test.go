package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/bagus-aulia/go-tools/tools/handler/response"
	payment_client "github.com/bagus-aulia/inventhier/internal/adapters/client/rest/payment/v1"
	"github.com/bagus-aulia/inventhier/internal/core/constants"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/payment"
	"github.com/bxcodec/faker"
)

func (s *suiteAPI) TestCheckout() {
	var mockPayload dto.CheckoutPayload
	faker.FakeData(&mockPayload)

	s.T().Run("success", func(t *testing.T) {
		s.client.DoFunc = func(r *http.Request) (*http.Response, error) {
			var mockResp response.Response[*dto.CheckoutResp]
			faker.FakeData(&mockResp)
			mockResp.Status = 200

			mockRespJSON, _ := json.Marshal(mockResp)

			bNoCLoser := io.NopCloser(bytes.NewReader(mockRespJSON))

			return &http.Response{
				StatusCode: 200,
				Body:       bNoCLoser,
			}, nil
		}

		cl := payment_client.NewRESTPaymentClient(s.client, s.cfg)
		resp, err := cl.Checkout(context.TODO(), mockPayload)
		s.NoError(err)
		s.Equal(200, resp.Status)
	})

	s.T().Run("error resp", func(t *testing.T) {
		s.client.DoFunc = func(r *http.Request) (*http.Response, error) {
			var mockResp response.Response[*dto.CheckoutResp]
			faker.FakeData(&mockResp)
			mockResp.Status = 500
			mockResp.Error = &response.ErrorMessage{
				Message: constants.ErrInternalServer.Error(),
			}

			mockRespJSON, _ := json.Marshal(mockResp)

			bNoCLoser := io.NopCloser(bytes.NewReader(mockRespJSON))

			return &http.Response{
				StatusCode: 200,
				Body:       bNoCLoser,
			}, nil
		}

		cl := payment_client.NewRESTPaymentClient(s.client, s.cfg)
		resp, err := cl.Checkout(context.TODO(), mockPayload)
		s.NoError(err)
		s.Equal(500, resp.Status)
	})

	s.T().Run("error", func(t *testing.T) {
		s.client.DoFunc = func(r *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("error")
		}

		cl := payment_client.NewRESTPaymentClient(s.client, s.cfg)
		resp, err := cl.Checkout(context.TODO(), mockPayload)
		s.Error(err)
		s.Equal(500, resp.Status)
	})
}
