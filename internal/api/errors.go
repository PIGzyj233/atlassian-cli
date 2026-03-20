package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError represents an error response from the Atlassian API.
type APIError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"error"`
	Hint       string `json:"hint,omitempty"`
}

func (e *APIError) Error() string {
	s := fmt.Sprintf("API error (HTTP %d): %s", e.StatusCode, e.Message)
	if e.Hint != "" {
		s += fmt.Sprintf(" — hint: %s", e.Hint)
	}
	return s
}

// ExitCode returns the appropriate process exit code for this error.
func (e *APIError) ExitCode() int {
	if e.StatusCode >= 400 && e.StatusCode < 500 {
		return 2 // API error
	}
	if e.StatusCode >= 500 {
		return 2 // API error
	}
	return 1 // input error
}

// parseAPIError reads the response body and builds an APIError.
func parseAPIError(resp *http.Response) *APIError {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTP %d (could not read body)", resp.StatusCode),
		}
	}

	apiErr := &APIError{StatusCode: resp.StatusCode}

	// Try Atlassian error format: {"errorMessages": [...], "errors": {...}}
	var atlassianErr struct {
		ErrorMessages []string          `json:"errorMessages"`
		Errors        map[string]string `json:"errors"`
		Message       string            `json:"message"`
	}
	if json.Unmarshal(body, &atlassianErr) == nil {
		var parts []string
		parts = append(parts, atlassianErr.ErrorMessages...)
		if atlassianErr.Message != "" {
			parts = append(parts, atlassianErr.Message)
		}
		for field, msg := range atlassianErr.Errors {
			parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
		}
		if len(parts) > 0 {
			apiErr.Message = strings.Join(parts, "; ")
		}
	}

	if apiErr.Message == "" {
		apiErr.Message = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// Add hints for common errors
	apiErr.Hint = hintForStatus(resp.StatusCode)

	return apiErr
}

func hintForStatus(status int) string {
	switch status {
	case 401:
		return "Run `auth login` to authenticate"
	case 403:
		return "Check that your account has the required permissions"
	case 404:
		return "Check that the resource key/ID is correct"
	case 429:
		return "Rate limited — wait and retry"
	default:
		return ""
	}
}
