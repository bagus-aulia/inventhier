package httpclient_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	cli "github.com/bagus-aulia/inventhier/internal/adapters/helpers/http_client"
	clMock "github.com/bagus-aulia/inventhier/internal/adapters/helpers/http_client/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDoRequest(t *testing.T) {
	client := &clMock.HttpClient{}

	t.Run("success", func(t *testing.T) {
		client.DoFunc = func(r *http.Request) (*http.Response, error) {
			bNoCLoser := io.NopCloser(bytes.NewReader([]byte("abc")))

			return &http.Response{
				StatusCode: 200,
				Body:       bNoCLoser,
			}, nil
		}

		opt := cli.HttpOptions{
			Client:   client,
			Hostname: "http://abc.com",
			Path:     "/def/ghi",
			Ctx:      context.Background(),
		}

		statusCode, err := cli.DoRequest(&opt, nil)
		assert.NoError(t, err)

		assert.Equal(t, 200, statusCode)
	})

	t.Run("success with json", func(t *testing.T) {
		client.DoFunc = func(r *http.Request) (*http.Response, error) {
			jsonResp := `{"message":"success"}`
			bNoCLoser := io.NopCloser(bytes.NewReader([]byte(jsonResp)))

			return &http.Response{
				StatusCode: 200,
				Body:       bNoCLoser,
			}, nil
		}

		opt := cli.HttpOptions{
			Client:   client,
			Hostname: "http://abc.com",
			Path:     "/def/ghi",
			Ctx:      context.Background(),
		}

		var resp map[string]string

		statusCode, err := cli.DoRequest(&opt, &resp)
		assert.NoError(t, err)

		assert.Equal(t, 200, statusCode)
		assert.Equal(t, "success", resp["message"])
	})
}

func TestTimeoutClientReuse(t *testing.T) {
	client := &clMock.HttpClient{}
	timeOut := 2

	t.Run("first attempt", func(t *testing.T) {
		client.DoFunc = func(r *http.Request) (*http.Response, error) {
			select {
			case <-r.Context().Done():
				return nil, r.Context().Err()
			case <-time.After(10 * time.Second):
				jsonResp := `{"message":"success"}`
				bNoCLoser := io.NopCloser(bytes.NewReader([]byte(jsonResp)))

				return &http.Response{
					StatusCode: 200,
					Body:       bNoCLoser,
				}, nil
			}
		}

		opt := cli.HttpOptions{
			Client:   client,
			Hostname: "http://abc.com",
			Path:     "/def/ghi",
			Ctx:      context.Background(),
			Timeout:  timeOut,
		}

		statusCode, err := cli.DoRequest(&opt, nil)
		assert.Error(t, err)

		assert.Equal(t, 500, statusCode)
	})

	t.Run("first attempt", func(t *testing.T) {
		client.DoFunc = func(r *http.Request) (*http.Response, error) {
			select {
			case <-r.Context().Done():
				return nil, r.Context().Err()
			case <-time.After(10 * time.Second):
				jsonResp := `{"message":"success"}`
				bNoCLoser := io.NopCloser(bytes.NewReader([]byte(jsonResp)))

				return &http.Response{
					StatusCode: 200,
					Body:       bNoCLoser,
				}, nil
			}
		}

		opt := cli.HttpOptions{
			Client:   client,
			Hostname: "http://abc.com",
			Path:     "/def/ghi",
			Ctx:      context.Background(),
			Timeout:  timeOut,
		}

		statusCode, err := cli.DoRequest(&opt, nil)
		assert.Error(t, err)

		assert.Equal(t, 500, statusCode)
	})

}
