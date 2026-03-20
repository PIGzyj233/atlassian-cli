package attachment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterAttachments_ByFilename(t *testing.T) {
	result := map[string]any{
		"results": []any{
			map[string]any{"title": "doc.pdf", "mediaType": "application/pdf"},
			map[string]any{"title": "image.png", "mediaType": "image/png"},
			map[string]any{"title": "doc.pdf", "mediaType": "application/pdf"},
		},
		"size": 3,
	}

	filtered := filterAttachments(result, "doc.pdf", "")
	results := filtered["results"].([]any)
	assert.Len(t, results, 2)
	assert.Equal(t, 2, filtered["size"])
}

func TestFilterAttachments_ByMediaType(t *testing.T) {
	result := map[string]any{
		"results": []any{
			map[string]any{"title": "doc.pdf", "mediaType": "application/pdf"},
			map[string]any{"title": "image.png", "mediaType": "image/png"},
		},
		"size": 2,
	}

	filtered := filterAttachments(result, "", "image/png")
	results := filtered["results"].([]any)
	assert.Len(t, results, 1)
	assert.Equal(t, "image.png", results[0].(map[string]any)["title"])
}

func TestFilterAttachments_BothFilters(t *testing.T) {
	result := map[string]any{
		"results": []any{
			map[string]any{"title": "a.pdf", "mediaType": "application/pdf"},
			map[string]any{"title": "b.pdf", "mediaType": "application/pdf"},
			map[string]any{"title": "a.pdf", "mediaType": "text/plain"},
		},
		"size": 3,
	}

	filtered := filterAttachments(result, "a.pdf", "application/pdf")
	results := filtered["results"].([]any)
	assert.Len(t, results, 1)
}

func TestFilterAttachments_NoMatch(t *testing.T) {
	result := map[string]any{
		"results": []any{
			map[string]any{"title": "doc.pdf", "mediaType": "application/pdf"},
		},
		"size": 1,
	}

	filtered := filterAttachments(result, "nonexistent.txt", "")
	results := filtered["results"]
	assert.Nil(t, results)
	assert.Equal(t, 0, filtered["size"])
}

func TestFilterAttachments_NoFilters(t *testing.T) {
	result := map[string]any{
		"results": []any{
			map[string]any{"title": "doc.pdf"},
		},
		"size": 1,
	}

	// Should pass through unmodified (but this function would only be called with filters)
	filtered := filterAttachments(result, "", "")
	results := filtered["results"].([]any)
	assert.Len(t, results, 1)
}
