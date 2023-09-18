package mocks

import "net/http"

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// PostDoFunc fetches the mock client's `Do` func for POST requests
	PostDoFunc func(req *http.Request) (*http.Response, error)
	// GetDoFunc fetches the mock client's `Do` func for GET requests
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

// Do is the mock client's `Do` func
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	if req.Method == "POST" {
		return PostDoFunc(req)
	}

	return GetDoFunc(req)
}
