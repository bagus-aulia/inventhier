package mocks

import (
	"net/http"
)

type HttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (h *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return h.DoFunc(req)
}
