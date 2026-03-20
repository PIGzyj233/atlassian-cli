package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PigZyj2333/atlassian-cli/internal/config"
)

// Client is an HTTP client for Atlassian REST APIs.
type Client struct {
	HTTP       *http.Client
	BaseURL    string
	Auth       config.Authenticator
	MaxRetries int
}

// NewClient creates a new API client.
func NewClient(baseURL string, auth config.Authenticator) *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: 75 * time.Second,
		},
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Auth:       auth,
		MaxRetries: 3,
	}
}

// IsCloud returns true if the base URL is an Atlassian Cloud instance.
func (c *Client) IsCloud() bool {
	return strings.Contains(c.BaseURL, ".atlassian.net")
}

// JiraAPIPath returns the versioned Jira REST API path.
func (c *Client) JiraAPIPath(resource string) string {
	if c.IsCloud() {
		return "/rest/api/3/" + resource
	}
	return "/rest/api/2/" + resource
}

// ConfluenceAPIPath returns the Confluence REST API path.
// Note: The /wiki prefix (Cloud) is already in the BaseURL set by Factory.
func (c *Client) ConfluenceAPIPath(resource string) string {
	return "/rest/api/" + resource
}

// Get performs an authenticated GET request. If dest is non-nil, the response
// body is JSON-decoded into dest.
func (c *Client) Get(path string, dest any) (*http.Response, error) {
	return c.doWithHeaders("GET", path, nil, dest, nil)
}

// GetWithHeaders performs an authenticated GET with custom headers.
func (c *Client) GetWithHeaders(path string, dest any, headers map[string]string) (*http.Response, error) {
	return c.doWithHeaders("GET", path, nil, dest, headers)
}

// Post performs an authenticated POST request with a JSON body.
func (c *Client) Post(path string, body any, dest any) (*http.Response, error) {
	return c.doWithHeaders("POST", path, body, dest, nil)
}

// PostWithHeaders performs an authenticated POST with custom headers.
func (c *Client) PostWithHeaders(path string, body any, dest any, headers map[string]string) (*http.Response, error) {
	return c.doWithHeaders("POST", path, body, dest, headers)
}

// Put performs an authenticated PUT request with a JSON body.
func (c *Client) Put(path string, body any, dest any) (*http.Response, error) {
	return c.doWithHeaders("PUT", path, body, dest, nil)
}

// PutWithHeaders performs an authenticated PUT with custom headers.
func (c *Client) PutWithHeaders(path string, body any, dest any, headers map[string]string) (*http.Response, error) {
	return c.doWithHeaders("PUT", path, body, dest, headers)
}

// Delete performs an authenticated DELETE request.
func (c *Client) Delete(path string, dest any) (*http.Response, error) {
	return c.doWithHeaders("DELETE", path, nil, dest, nil)
}

func (c *Client) doWithHeaders(method, path string, body any, dest any, extraHeaders map[string]string) (*http.Response, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = strings.NewReader(string(data))
	}

	var lastResp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		for k, v := range extraHeaders {
			req.Header.Set(k, v)
		}

		if c.Auth != nil {
			if err := c.Auth.Apply(req); err != nil {
				return nil, fmt.Errorf("applying auth: %w", err)
			}
		}

		resp, err := c.HTTP.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			time.Sleep(backoff(attempt))
			// Reset body reader for retry
			if body != nil {
				data, _ := json.Marshal(body)
				bodyReader = strings.NewReader(string(data))
			}
			continue
		}

		lastResp = resp
		lastErr = nil

		// Retry on 429 (rate limit) and 5xx (server errors)
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			resp.Body.Close()
			if attempt < c.MaxRetries {
				time.Sleep(backoff(attempt))
				if body != nil {
					data, _ := json.Marshal(body)
					bodyReader = strings.NewReader(string(data))
				}
				continue
			}
			return resp, &APIError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("HTTP %d after %d retries", resp.StatusCode, c.MaxRetries)}
		}

		// Non-retryable error
		if resp.StatusCode >= 400 {
			defer resp.Body.Close()
			return resp, parseAPIError(resp)
		}

		// Success — always close body (prevents leak when dest == nil)
		defer resp.Body.Close()
		if dest != nil {
			if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
				return resp, fmt.Errorf("decoding response: %w", err)
			}
		}

		return resp, nil
	}

	if lastErr != nil {
		return lastResp, lastErr
	}
	return lastResp, fmt.Errorf("request failed after %d retries", c.MaxRetries)
}

func backoff(attempt int) time.Duration {
	// Exponential backoff: 100ms, 200ms, 400ms, ...
	// Capped at 5 seconds
	d := time.Duration(100<<uint(attempt)) * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
