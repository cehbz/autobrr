package autobrr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is used to interact with the Autobrr API
type Client struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// Filter represents an Autobrr filter
type Filter struct {
	ID                    int64            `json:"id,omitempty"`
	Name                  string           `json:"name"`
	Enabled               bool             `json:"enabled"`
	Priority              int              `json:"priority"`
	SmartEpisode          bool             `json:"smart_episode"`
	Shows                 []string         `json:"shows,omitempty"`
	Resolutions           []string         `json:"resolutions,omitempty"`
	Sources               []string         `json:"sources,omitempty"`
	Codecs                []string         `json:"codecs,omitempty"`
	Containers            []string         `json:"containers,omitempty"`
	MatchReleases         string           `json:"match_releases,omitempty"`
	ExceptReleases        string           `json:"except_releases,omitempty"`
	Years                 string           `json:"years,omitempty"`
	Tags                  string           `json:"tags,omitempty"`
	ExceptTags            string           `json:"except_tags,omitempty"`
	MatchReleaseGroups    string           `json:"match_release_groups,omitempty"`
	ExceptReleaseGroups   string           `json:"except_release_groups,omitempty"`
	MaxSize               string           `json:"max_size,omitempty"`
	MinSize               string           `json:"min_size,omitempty"`
	Actions               []Action         `json:"actions,omitempty"`
	External              []ExternalFilter `json:"external,omitempty"`
	IndexerIDs            []int            `json:"indexer_ids,omitempty"`
	MatchCategories       string           `json:"match_categories,omitempty"`
	ExceptCategories      string           `json:"except_categories,omitempty"`
	MatchUploaders        string           `json:"match_uploaders,omitempty"`
	ExceptUploaders       string           `json:"except_uploaders,omitempty"`
	MatchLanguage         []string         `json:"match_language,omitempty"`
	ExceptLanguage        []string         `json:"except_language,omitempty"`
	UseRegex              bool             `json:"use_regex"`
	UseRegexReleaseGroups bool             `json:"use_regex_release_groups"`
	Scene                 bool             `json:"scene"`
	Origins               []string         `json:"origins,omitempty"`
	ExceptOrigins         []string         `json:"except_origins,omitempty"`
	Bonus                 []string         `json:"bonus,omitempty"`
	Freeleech             bool             `json:"freeleech"`
	FreeleechPercent      string           `json:"freeleech_percent,omitempty"`
	Description           string           `json:"description,omitempty"`
	CreatedAt             string           `json:"created_at,omitempty"`
	UpdatedAt             string           `json:"updated_at,omitempty"`
}

// Action represents an action to be taken when a filter matches
type Action struct {
	ID                   int64    `json:"id,omitempty"`
	Name                 string   `json:"name"`
	Type                 string   `json:"type"`
	Enabled              bool     `json:"enabled"`
	ExecCmd              string   `json:"exec_cmd,omitempty"`
	ExecArgs             string   `json:"exec_args,omitempty"`
	WatchFolder          string   `json:"watch_folder,omitempty"`
	Category             string   `json:"category,omitempty"`
	Tags                 string   `json:"tags,omitempty"`
	Label                string   `json:"label,omitempty"`
	SavePath             string   `json:"save_path,omitempty"`
	Paused               bool     `json:"paused"`
	IgnoreRules          bool     `json:"ignore_rules"`
	SkipHashCheck        bool     `json:"skip_hash_check"`
	ContentLayout        string   `json:"content_layout,omitempty"`
	FirstLastPiecePrio   bool     `json:"first_last_piece_prio"`
	Priority             string   `json:"priority,omitempty"`
	LimitDownloadSpeed   int64    `json:"limit_download_speed,omitempty"`
	LimitUploadSpeed     int64    `json:"limit_upload_speed,omitempty"`
	LimitRatio           float64  `json:"limit_ratio,omitempty"`
	LimitSeedTime        int64    `json:"limit_seed_time,omitempty"`
	ReannounceInterval   int64    `json:"reannounce_interval,omitempty"`
	ClientID             int      `json:"client_id,omitempty"`
	WebhookHost          string   `json:"webhook_host,omitempty"`
	WebhookType          string   `json:"webhook_type,omitempty"`
	WebhookMethod        string   `json:"webhook_method,omitempty"`
	WebhookData          string   `json:"webhook_data,omitempty"`
	WebhookHeaders       []string `json:"webhook_headers,omitempty"`
	ExternalDownloadOnly bool     `json:"external_download_only"`
	FilterID             int64    `json:"filter_id,omitempty"`
}

// ExternalFilter represents external filter criteria
type ExternalFilter struct {
	ID               int64  `json:"id,omitempty"`
	Name             string `json:"name"`
	Index            int    `json:"index"`
	Type             string `json:"type"`
	Value            string `json:"value"`
	ExecCmd          string `json:"exec_cmd,omitempty"`
	ExecArgs         string `json:"exec_args,omitempty"`
	ExecExpectOutput string `json:"exec_expect_output,omitempty"`
	FilterID         int64  `json:"filter_id,omitempty"`
}

// FilterListResponse represents the response from listing filters
type FilterListResponse struct {
	Data []Filter `json:"data"`
}

// NewClient initializes a new Autobrr client.
// If httpClient is nil, http.DefaultClient is used.
func NewClient(apiKey, addr, port string, httpClient ...*http.Client) (*Client, error) {
	// Use the provided http.Client if given, otherwise use http.DefaultClient
	client := http.DefaultClient
	if len(httpClient) > 0 && httpClient[0] != nil {
		client = httpClient[0]
	}

	// Create and return the Client instance
	abClient := &Client{
		client:  client,
		baseURL: fmt.Sprintf("http://%s:%s", addr, port),
		apiKey:  apiKey,
	}

	return abClient, nil
}

// GetFilters retrieves all filters
func (c *Client) GetFilters() ([]Filter, error) {
	respData, err := c.doGet("/api/filters")
	if err != nil {
		return nil, fmt.Errorf("get filters error: %v", err)
	}

	var response FilterListResponse
	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to decode filters response: %v", err)
	}

	return response.Data, nil
}

// GetFilter retrieves a specific filter by ID
func (c *Client) GetFilter(id int64) (*Filter, error) {
	endpoint := fmt.Sprintf("/api/filters/%d", id)
	respData, err := c.doGet(endpoint)
	if err != nil {
		return nil, fmt.Errorf("get filter error: %v", err)
	}

	var filter Filter
	if err := json.Unmarshal(respData, &filter); err != nil {
		return nil, fmt.Errorf("failed to decode filter response: %v", err)
	}

	return &filter, nil
}

// CreateFilter creates a new filter
func (c *Client) CreateFilter(filter *Filter) (*Filter, error) {
	jsonData, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filter: %v", err)
	}

	respData, err := c.doPost("/api/filters", bytes.NewReader(jsonData), "application/json")
	if err != nil {
		return nil, fmt.Errorf("create filter error: %v", err)
	}

	var createdFilter Filter
	if err := json.Unmarshal(respData, &createdFilter); err != nil {
		return nil, fmt.Errorf("failed to decode created filter: %v", err)
	}

	return &createdFilter, nil
}

// UpdateFilter updates an existing filter
func (c *Client) UpdateFilter(id int64, filter *Filter) (*Filter, error) {
	jsonData, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filter: %v", err)
	}

	endpoint := fmt.Sprintf("/api/filters/%d", id)
	respData, err := c.doPut(endpoint, bytes.NewReader(jsonData), "application/json")
	if err != nil {
		return nil, fmt.Errorf("update filter error: %v", err)
	}

	var updatedFilter Filter
	if err := json.Unmarshal(respData, &updatedFilter); err != nil {
		return nil, fmt.Errorf("failed to decode updated filter: %v", err)
	}

	return &updatedFilter, nil
}

// DeleteFilter deletes a filter by ID
func (c *Client) DeleteFilter(id int64) error {
	endpoint := fmt.Sprintf("/api/filters/%d", id)
	_, err := c.doDelete(endpoint)
	if err != nil {
		return fmt.Errorf("delete filter error: %v", err)
	}

	return nil
}

// ToggleFilterEnabled enables or disables a filter
func (c *Client) ToggleFilterEnabled(id int64, enabled bool) error {
	endpoint := fmt.Sprintf("/api/filters/%d/toggle", id)
	data := map[string]bool{"enabled": enabled}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal toggle data: %v", err)
	}

	_, err = c.doPut(endpoint, bytes.NewReader(jsonData), "application/json")
	if err != nil {
		return fmt.Errorf("toggle filter error: %v", err)
	}

	return nil
}

// TestConnection verifies the connection to Autobrr
func (c *Client) TestConnection() error {
	_, err := c.doGet("/api/filters")
	if err != nil {
		return fmt.Errorf("connection test failed: %v", err)
	}

	return nil
}

// doGet is a helper method for making GET requests to the Autobrr API
func (c *Client) doGet(endpoint string) ([]byte, error) {
	return c.doRequest("GET", endpoint, nil, "")
}

// doPost is a helper method for making POST requests to the Autobrr API
func (c *Client) doPost(endpoint string, body io.Reader, contentType string) ([]byte, error) {
	return c.doRequest("POST", endpoint, body, contentType)
}

// doPut is a helper method for making PUT requests to the Autobrr API
func (c *Client) doPut(endpoint string, body io.Reader, contentType string) ([]byte, error) {
	return c.doRequest("PUT", endpoint, body, contentType)
}

// doDelete is a helper method for making DELETE requests to the Autobrr API
func (c *Client) doDelete(endpoint string) ([]byte, error) {
	return c.doRequest("DELETE", endpoint, nil, "")
}

// doRequest is a helper function to handle HTTP requests
func (c *Client) doRequest(method, endpoint string, body io.Reader, contentType string) ([]byte, error) {
	apiURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	apiURL.Path = endpoint

	req, err := http.NewRequest(method, apiURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set API key header
	req.Header.Set("X-API-Token", c.apiKey)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check for success status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected response code: %d, response: %s", resp.StatusCode, string(responseData))
	}

	return responseData, nil
}
