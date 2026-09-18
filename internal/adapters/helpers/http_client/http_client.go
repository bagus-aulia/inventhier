package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bagus-aulia/inventhier/internal/core/constants"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
)

// HTTPClient representation of http client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HttpOptions object needed for request
type HttpOptions struct {
	Client     HTTPClient
	Ctx        context.Context
	Hostname   string
	Path       string
	Timeout    int
	Headers    map[string]string
	URLQueries map[string]string
	Body       []byte
	Method     string
}

// HttpLog object to log error http connection
type HttpLog struct {
	Method     string
	URL        string
	StatusCode int
	Headers    map[string]string
	URLQueries map[string]string
	BodyParam  interface{}
	Response   interface{}
}

// setTracingHeaders extracts the request id from the
// request context and sets them as outbound headers. Centralized here so any
// future change to x-request-id propagation only needs to be
// made in one place.
func setTracingHeaders(req *http.Request, ctx context.Context) {
	// extract x-request-id from context
	if xRequestID, ok := ctx.Value(constants.XRequestIDKey).(string); ok && len(xRequestID) > 0 {
		// Fallback to XRequestIDStandardKey
		req.Header.Set(constants.XRequestIDHeader, xRequestID)
	}
}

// DoRequest func for executing http call
func DoRequest(opt *HttpOptions, rs interface{}) (int, error) {
	logger := helpers.GetZerologWithContext(opt.Ctx).With().
		Str("client", "http").
		Str("function", "DoRequest").
		Str("method", opt.Method).
		Str("hostname", opt.Hostname).
		Str("path", opt.Path).
		Logger()

	statusCode := http.StatusInternalServerError

	if opt.Timeout != 0 {
		ctx, cancel := context.WithTimeout(opt.Ctx, time.Duration(opt.Timeout)*time.Second)
		defer cancel()

		opt.Ctx = ctx
	}

	body := bytes.NewBuffer(opt.Body)
	defer body.Reset()

	u, err := url.JoinPath(opt.Hostname, opt.Path)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to join path url")

		return statusCode, err
	}
	req, err := http.NewRequestWithContext(opt.Ctx, opt.Method, u, body)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to create new request")

		return statusCode, err
	}
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}
	setTracingHeaders(req, opt.Ctx)

	queryValues := req.URL.Query()
	for key, val := range opt.URLQueries {
		queryValues.Set(key, val)
	}
	req.URL.RawQuery = queryValues.Encode()

	resp, err := opt.Client.Do(req)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to execute HTTP client request")

		return statusCode, err
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if rs == nil {
		return statusCode, nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to read response body")

		return http.StatusInternalServerError, err
	}

	if statusCode != http.StatusOK && statusCode != http.StatusAccepted {
		logData := HttpLog{
			Method:     opt.Method,
			URL:        u,
			StatusCode: statusCode,
			Headers:    opt.Headers,
			URLQueries: opt.URLQueries,
			BodyParam:  string(opt.Body),
			Response:   string(respBody),
		}

		if statusCode != http.StatusNotFound {
			logger.Error().
				Int("status_code", statusCode).
				Interface("log_data", logData).
				Msg("Received non-200 HTTP response")
		}
	}

	err = json.Unmarshal(respBody, rs)
	if err != nil {
		logger.Error().
			Err(err).
			Str("response_body", string(respBody)).
			Msg("Failed to unmarshal JSON response")

		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}
