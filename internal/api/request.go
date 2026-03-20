package api

import (
	"fmt"
	"net/url"
)

// RequestBuilder constructs URL paths with query parameters.
type RequestBuilder struct {
	basePath string
	params   url.Values
}

// NewRequestBuilder creates a new RequestBuilder for the given path.
func NewRequestBuilder(basePath string) *RequestBuilder {
	return &RequestBuilder{
		basePath: basePath,
		params:   url.Values{},
	}
}

// Query adds a query parameter. Empty values are skipped.
func (rb *RequestBuilder) Query(key, value string) *RequestBuilder {
	if value != "" {
		rb.params.Set(key, value)
	}
	return rb
}

// QueryInt adds an integer query parameter. Zero values are skipped.
func (rb *RequestBuilder) QueryInt(key string, value int) *RequestBuilder {
	if value != 0 {
		rb.params.Set(key, fmt.Sprintf("%d", value))
	}
	return rb
}

// BuildPath returns the full path with query string.
func (rb *RequestBuilder) BuildPath() string {
	if len(rb.params) == 0 {
		return rb.basePath
	}
	return rb.basePath + "?" + rb.params.Encode()
}
