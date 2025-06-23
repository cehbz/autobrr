package autobrr

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test with default HTTP client
	client, err := NewClient("test-api-key", "localhost", "10798")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", client.apiKey)
	}

	if client.baseURL != "http://localhost:10798" {
		t.Errorf("Expected base URL 'http://localhost:10798', got '%s'", client.baseURL)
	}

	// Test with custom HTTP client
	customHTTP := &http.Client{}
	client2, err := NewClient("test-api-key", "localhost", "10798", customHTTP)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client2.client != customHTTP {
		t.Error("Expected custom HTTP client to be used")
	}
}

func TestGetFilters(t *testing.T) {
	mockFilters := []Filter{
		{
			ID:       1,
			Name:     "Test Filter 1",
			Enabled:  true,
			Priority: 100,
			Shows:    "Test Show",
		},
		{
			ID:       2,
			Name:     "Test Filter 2",
			Enabled:  false,
			Priority: 200,
			Shows:    "Another Show",
		},
	}

	responseBody, _ := json.Marshal(mockFilters)

	// Mock successful response
	endpointResponses := map[string]mockResponse{
		"/api/filters": {statusCode: http.StatusOK, responseBody: string(responseBody)},
	}
	expectedRequests := []expectedRequest{
		{method: "GET", url: "/api/filters"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	filters, err := client.GetFilters()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(filters) != 2 {
		t.Errorf("Expected 2 filters, got %d", len(filters))
	}

	if filters[0].Name != "Test Filter 1" {
		t.Errorf("Expected first filter name 'Test Filter 1', got '%s'", filters[0].Name)
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestGetFilter(t *testing.T) {
	mockFilter := Filter{
		ID:          1,
		Name:        "Test Filter",
		Enabled:     true,
		Priority:    100,
		Shows:       "Test Show",
		Resolutions: []string{"1080p"},
		Actions: []Action{
			{
				ID:       1,
				Name:     "qBittorrent",
				Type:     "QBITTORRENT",
				Enabled:  true,
				ClientID: 1,
			},
		},
	}

	responseBody, _ := json.Marshal(mockFilter)

	// Mock successful response
	endpointResponses := map[string]mockResponse{
		"/api/filters/1": {statusCode: http.StatusOK, responseBody: string(responseBody)},
	}
	expectedRequests := []expectedRequest{
		{method: "GET", url: "/api/filters/1"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	filter, err := client.GetFilter(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if filter.Name != "Test Filter" {
		t.Errorf("Expected filter name 'Test Filter', got '%s'", filter.Name)
	}

	if len(filter.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(filter.Actions))
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestCreateFilter(t *testing.T) {
	newFilter := &Filter{
		Name:        "New Filter",
		Enabled:     true,
		Priority:    100,
		Shows:       "New Show",
		Resolutions: []string{"1080p", "2160p"},
		Actions: []Action{
			{
				Name:     "qBittorrent",
				Type:     "QBITTORRENT",
				Enabled:  true,
				ClientID: 1,
			},
		},
	}

	createdFilter := *newFilter
	createdFilter.ID = 42 // Server assigns ID
	responseBody, _ := json.Marshal(createdFilter)

	// Mock successful response
	endpointResponses := map[string]mockResponse{
		"/api/filters": {statusCode: http.StatusCreated, responseBody: string(responseBody)},
	}
	expectedRequests := []expectedRequest{
		{method: "POST", url: "/api/filters"},
	}

	customHandler := map[string]func(*http.Request){
		"/api/filters": func(req *http.Request) {
			// Verify request body
			var receivedFilter Filter
			if err := json.NewDecoder(req.Body).Decode(&receivedFilter); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}
			if receivedFilter.Name != newFilter.Name {
				t.Errorf("Expected filter name '%s', got '%s'", newFilter.Name, receivedFilter.Name)
			}
		},
	}

	client, mockTransport, err := newMockClientWithHandler(endpointResponses, expectedRequests, customHandler)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result, err := client.CreateFilter(newFilter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.ID != 42 {
		t.Errorf("Expected filter ID 42, got %d", result.ID)
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestUpdateFilter(t *testing.T) {
	updateFilter := &Filter{
		ID:          1,
		Name:        "Updated Filter",
		Enabled:     false,
		Priority:    200,
		Shows:       "Updated Show",
		Resolutions: []string{"720p"},
	}

	responseBody, _ := json.Marshal(updateFilter)

	// Mock successful response
	endpointResponses := map[string]mockResponse{
		"/api/filters/1": {statusCode: http.StatusOK, responseBody: string(responseBody)},
	}
	expectedRequests := []expectedRequest{
		{method: "PUT", url: "/api/filters/1"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result, err := client.UpdateFilter(1, updateFilter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Name != "Updated Filter" {
		t.Errorf("Expected filter name 'Updated Filter', got '%s'", result.Name)
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestDeleteFilter(t *testing.T) {
	// Mock successful response
	endpointResponses := map[string]mockResponse{
		"/api/filters/1": {statusCode: http.StatusNoContent, responseBody: ""},
	}
	expectedRequests := []expectedRequest{
		{method: "DELETE", url: "/api/filters/1"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.DeleteFilter(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestToggleFilterEnabled(t *testing.T) {
	// Mock successful response
	endpointResponses := map[string]mockResponse{
		"/api/filters/1/toggle": {statusCode: http.StatusOK, responseBody: "{}"},
	}
	expectedRequests := []expectedRequest{
		{method: "PUT", url: "/api/filters/1/toggle"},
	}

	customHandler := map[string]func(*http.Request){
		"/api/filters/1/toggle": func(req *http.Request) {
			// Verify request body
			var data map[string]bool
			if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}
			if enabled, ok := data["enabled"]; !ok || !enabled {
				t.Errorf("Expected enabled: true, got %v", data)
			}
		},
	}

	client, mockTransport, err := newMockClientWithHandler(endpointResponses, expectedRequests, customHandler)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.ToggleFilterEnabled(1, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTestConnection(t *testing.T) {
	mockResp := FilterListResponse{Data: []Filter{}}
	responseBody, _ := json.Marshal(mockResp)

	// Mock successful response
	endpointResponses := map[string]mockResponse{
		"/api/filters": {statusCode: http.StatusOK, responseBody: string(responseBody)},
	}
	expectedRequests := []expectedRequest{
		{method: "GET", url: "/api/filters"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TestConnection()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestTestConnection_Error(t *testing.T) {
	// Mock error response
	endpointResponses := map[string]mockResponse{
		"/api/filters": {statusCode: http.StatusUnauthorized, responseBody: "Invalid API key"},
	}
	expectedRequests := []expectedRequest{
		{method: "GET", url: "/api/filters"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.TestConnection()
	if err == nil {
		t.Fatal("Expected error, got none")
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestGetFilters_Error(t *testing.T) {
	// Mock error response
	endpointResponses := map[string]mockResponse{
		"/api/filters": {statusCode: http.StatusInternalServerError, responseBody: "Server error"},
	}
	expectedRequests := []expectedRequest{
		{method: "GET", url: "/api/filters"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = client.GetFilters()
	if err == nil {
		t.Fatal("Expected error, got none")
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}

func TestDeleteFilter_Error(t *testing.T) {
	// Mock error response
	endpointResponses := map[string]mockResponse{
		"/api/filters/999": {statusCode: http.StatusNotFound, responseBody: "Filter not found"},
	}
	expectedRequests := []expectedRequest{
		{method: "DELETE", url: "/api/filters/999"},
	}

	client, mockTransport, err := newMockClient(endpointResponses, expectedRequests)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = client.DeleteFilter(999)
	if err == nil {
		t.Fatal("Expected error, got none")
	}

	// Check the request made
	if mockTransport.requestIndex != len(mockTransport.expectedRequests) {
		t.Errorf("Not all expected requests were made")
	}
}
