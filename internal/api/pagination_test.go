package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationInfo_IsLast(t *testing.T) {
	// Offset-based: startAt=0, maxResults=50, total=30
	p := &PaginationInfo{StartAt: 0, MaxResults: 50, Total: 30, IsLast: true}
	assert.True(t, p.IsLast)

	// Not last
	p2 := &PaginationInfo{StartAt: 0, MaxResults: 50, Total: 100, IsLast: false}
	assert.False(t, p2.IsLast)
}

func TestPaginationInfo_JSON(t *testing.T) {
	p := &PaginationInfo{StartAt: 10, MaxResults: 50, Total: 200, IsLast: false}
	data := p.ToMap()
	assert.Equal(t, 10, data["startAt"])
	assert.Equal(t, 50, data["maxResults"])
	assert.Equal(t, 200, data["total"])
	assert.Equal(t, false, data["isLast"])
}
