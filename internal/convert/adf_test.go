package convert

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownToADF_Paragraph(t *testing.T) {
	md := "Hello world"
	adf := MarkdownToADF(md)

	assert.Equal(t, "doc", adf["type"])
	assert.Equal(t, 1, adf["version"])

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	para := content[0].(map[string]any)
	assert.Equal(t, "paragraph", para["type"])
}

func TestMarkdownToADF_Headings(t *testing.T) {
	md := "# Title\n\nSome text\n\n## Subtitle"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.GreaterOrEqual(t, len(content), 3)

	h1 := content[0].(map[string]any)
	assert.Equal(t, "heading", h1["type"])
	attrs := h1["attrs"].(map[string]any)
	assert.Equal(t, 1, attrs["level"])
}

func TestMarkdownToADF_CodeBlock(t *testing.T) {
	md := "```python\nprint('hello')\n```"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	cb := content[0].(map[string]any)
	assert.Equal(t, "codeBlock", cb["type"])
	attrs := cb["attrs"].(map[string]any)
	assert.Equal(t, "python", attrs["language"])
}

func TestMarkdownToADF_BulletList(t *testing.T) {
	md := "- Item 1\n- Item 2\n- Item 3"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	list := content[0].(map[string]any)
	assert.Equal(t, "bulletList", list["type"])
}

func TestADFToMarkdown_Paragraph(t *testing.T) {
	adf := map[string]any{
		"type":    "doc",
		"version": 1,
		"content": []any{
			map[string]any{
				"type": "paragraph",
				"content": []any{
					map[string]any{"type": "text", "text": "Hello world"},
				},
			},
		},
	}

	md := ADFToMarkdown(adf)
	assert.Equal(t, "Hello world", md)
}

func TestADFToMarkdown_RoundTrip(t *testing.T) {
	// Simple paragraphs should survive a round-trip
	original := "Hello world"
	adf := MarkdownToADF(original)
	result := ADFToMarkdown(adf)
	assert.Equal(t, original, result)
}

func TestMarkdownToADF_JSON(t *testing.T) {
	md := "Hello **bold** world"
	adf := MarkdownToADF(md)

	// Should be valid JSON
	_, err := json.Marshal(adf)
	require.NoError(t, err)
}
