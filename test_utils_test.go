package autobrr

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockRoundTripper is used to mock http.Client responses
type mockRoundTripper struct {
	responses        map[string]mockResponse
	expectedRequests []expectedRequest
	requestIndex     int
	t                *testing.T
	customHandler    map[string]func(*http.Request)
}

// mockResponse represents a mock HTTP response
type mockResponse struct {
	statusCode   int
	responseBody string
}

// expectedRequest represents an expected HTTP request
type expectedRequest struct {
	method string
	url    string
}

// RoundTrip implements the RoundTripper interface
func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.requestIndex >= len(m.expectedRequests) {
		m.t.Errorf("Unexpected request: %s %s", req.Method, req.URL.Path)
		return nil, fmt.Errorf("unexpected request")
	}

	expected := m.expectedRequests[m.requestIndex]
	m.requestIndex++

	// Check method and path
	if req.Method != expected.method {
		m.t.Errorf("Expected method %s, got %s", expected.method, req.Method)
	}
	if req.URL.Path != expected.url {
		m.t.Errorf("Expected URL %s, got %s", expected.url, req.URL.Path)
	}

	// Check API key header
	if apiKey := req.Header.Get("X-API-Token"); apiKey != "test-api-key" {
		m.t.Errorf("Expected API key header 'test-api-key', got '%s'", apiKey)
	}

	// Call custom handler if present
	if handler, ok := m.customHandler[req.URL.Path]; ok && handler != nil {
		handler(req)
	}

	resp := m.responses[req.URL.Path]
	return &http.Response{
		StatusCode: resp.statusCode,
		Body:       io.NopCloser(strings.NewReader(resp.responseBody)),
		Header:     make(http.Header),
	}, nil
}

// newMockClient creates a mock client with predefined responses
func newMockClient(responses map[string]mockResponse, expectedRequests []expectedRequest) (*Client, *mockRoundTripper, error) {
	return newMockClientWithHandler(responses, expectedRequests, nil)
}

// newMockClientWithHandler creates a mock client with predefined responses and custom handlers
func newMockClientWithHandler(responses map[string]mockResponse, expectedRequests []expectedRequest, customHandler map[string]func(*http.Request)) (*Client, *mockRoundTripper, error) {
	transport := &mockRoundTripper{
		responses:        responses,
		expectedRequests: expectedRequests,
		t:                &testing.T{},
		customHandler:    customHandler,
	}

	httpClient := &http.Client{Transport: transport}
	client, err := NewClient("test-api-key", "localhost", "10798", httpClient)
	return client, transport, err
}
