package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBuilder_BuildPath(t *testing.T) {
	rb := NewRequestBuilder("/rest/api/3/search").
		Query("jql", "project = TEST").
		Query("maxResults", "50").
		Query("startAt", "0")

	path := rb.BuildPath()
	assert.Contains(t, path, "/rest/api/3/search?")
	assert.Contains(t, path, "jql=project+%3D+TEST")
	assert.Contains(t, path, "maxResults=50")
}

func TestRequestBuilder_EmptyValues(t *testing.T) {
	rb := NewRequestBuilder("/path").
		Query("key", "value").
		Query("empty", "")

	path := rb.BuildPath()
	assert.Contains(t, path, "key=value")
	assert.NotContains(t, path, "empty")
}
