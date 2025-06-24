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
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Priority  int    `json:"priority"`
	Shows     string `json:"shows,omitempty"`
	Seasons   string `json:"seasons,omitempty"`
	Episodes  string `json:"episodes,omitempty"`

	// Match
	AnnounceTypes       []string `json:"announce_types,omitempty"`
	Resolutions         []string `json:"resolutions,omitempty"`
	Sources             []string `json:"sources,omitempty"`
	Codecs              []string `json:"codecs,omitempty"`
	Containers          []string `json:"containers,omitempty"`
	MatchReleases       string   `json:"match_releases,omitempty"`
	ExceptReleases      string   `json:"except_releases,omitempty"`
	Years               string   `json:"years,omitempty"`
	Tags                string   `json:"tags,omitempty"`
	ExceptTags          string   `json:"except_tags,omitempty"`
	MatchReleaseGroups  string   `json:"match_release_groups,omitempty"`
	ExceptReleaseGroups string   `json:"except_release_groups,omitempty"`

	// Size
	MaxSize string `json:"max_size,omitempty"`
	MinSize string `json:"min_size,omitempty"`

	// Indexers
	IndexerIDs []int `json:"indexer_ids,omitempty"`

	// Categories
	MatchCategories  string `json:"match_categories,omitempty"`
	ExceptCategories string `json:"except_categories,omitempty"`

	// Uploaders
	MatchUploaders  string `json:"match_uploaders,omitempty"`
	ExceptUploaders string `json:"except_uploaders,omitempty"`

	// Language
	MatchLanguage  []string `json:"match_language,omitempty"`
	ExceptLanguage []string `json:"except_language,omitempty"`

	// Regex
	UseRegex              bool `json:"use_regex"`
	UseRegexReleaseGroups bool `json:"use_regex_release_groups"`

	// Other
	Scene         bool     `json:"scene"`
	Origins       []string `json:"origins,omitempty"`
	ExceptOrigins []string `json:"except_origins,omitempty"`
	Bonus         []string `json:"bonus,omitempty"`

	// Freeleech
	Freeleech        bool   `json:"freeleech"`
	FreeleechPercent string `json:"freeleech_percent,omitempty"`

	// Other
	Description string `json:"description,omitempty"`

	SmartEpisode bool `json:"smart_episode"`

	ExceptOther []string `json:"except_other,omitempty"`

	MaxDownloads     int    `json:"max_downloads,omitempty"`
	MaxDownloadsUnit string `json:"max_downloads_unit,omitempty"`

	ActionsCount        int `json:"actions_count,omitempty"`
	ActionsEnabledCount int `json:"actions_enabled_count,omitempty"`

	IsAutoUpdated bool `json:"is_auto_updated,omitempty"`

	ReleaseProfileDuplicate interface{} `json:"release_profile_duplicate,omitempty"`

	// Relations
	Actions  []Action   `json:"actions,omitempty"`
	External []External `json:"external,omitempty"`
	Indexers []Indexer  `json:"indexers,omitempty"`

	Downloads *Downloads `json:"downloads,omitempty"`
}

// Action represents an action to be taken when a filter matches
type Action struct {
	ID                    int64    `json:"id,omitempty"`
	Name                  string   `json:"name"`
	Type                  string   `json:"type"`
	Enabled               bool     `json:"enabled"`
	Category              string   `json:"category,omitempty"`
	ReannounceInterval    int64    `json:"reannounce_interval,omitempty"`
	ReannounceMaxAttempts int64    `json:"reannounce_max_attempts,omitempty"`
	ClientID              int      `json:"client_id,omitempty"`
	ExecCmd               string   `json:"exec_cmd,omitempty"`
	ExecArgs              string   `json:"exec_args,omitempty"`
	WatchFolder           string   `json:"watch_folder,omitempty"`
	Tags                  string   `json:"tags,omitempty"`
	Label                 string   `json:"label,omitempty"`
	SavePath              string   `json:"save_path,omitempty"`
	Paused                bool     `json:"paused"`
	IgnoreRules           bool     `json:"ignore_rules"`
	SkipHashCheck         bool     `json:"skip_hash_check"`
	ContentLayout         string   `json:"content_layout,omitempty"`
	FirstLastPiecePrio    bool     `json:"first_last_piece_prio"`
	Priority              string   `json:"priority,omitempty"`
	LimitDownloadSpeed    int64    `json:"limit_download_speed,omitempty"`
	LimitUploadSpeed      int64    `json:"limit_upload_speed,omitempty"`
	LimitRatio            float64  `json:"limit_ratio,omitempty"`
	LimitSeedTime         int64    `json:"limit_seed_time,omitempty"`
	WebhookHost           string   `json:"webhook_host,omitempty"`
	WebhookType           string   `json:"webhook_type,omitempty"`
	WebhookMethod         string   `json:"webhook_method,omitempty"`
	WebhookData           string   `json:"webhook_data,omitempty"`
	WebhookHeaders        []string `json:"webhook_headers,omitempty"`
	ExternalDownloadOnly  bool     `json:"external_download_only"`
	FilterID              int64    `json:"filter_id,omitempty"`
}

// External represents an external filter.
type External struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Index    int    `json:"index"`
	Type     string `json:"type"`
	Enabled  bool   `json:"enabled"`
	ExecCmd  string `json:"exec_cmd,omitempty"`
	ExecArgs string `json:"exec_args,omitempty"`
}

// Indexer represents an indexer.
type Indexer struct {
	ID                 int                    `json:"id"`
	Name               string                 `json:"name"`
	Identifier         string                 `json:"identifier"`
	IdentifierExternal string                 `json:"identifier_external"`
	Enabled            bool                   `json:"enabled"`
	Implementation     string                 `json:"implementation"`
	BaseURL            string                 `json:"base_url"`
	UseProxy           bool                   `json:"use_proxy"`
	Proxy              interface{}            `json:"proxy"`
	ProxyID            int                    `json:"proxy_id"`
	Settings           map[string]interface{} `json:"settings"`
}

// Downloads represents download statistics.
type Downloads struct {
	HourCount  int `json:"hour_count"`
	DayCount   int `json:"day_count"`
	WeekCount  int `json:"week_count"`
	MonthCount int `json:"month_count"`
	TotalCount int `json:"total_count"`
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

	var response []Filter
	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to decode filters response: %v", err)
	}

	return response, nil
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
