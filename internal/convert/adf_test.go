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

// --- Inline mark tests ---

func TestMarkdownToADF_Bold(t *testing.T) {
	md := "Hello **bold** world"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	para := content[0].(map[string]any)
	inlines := para["content"].([]any)

	// Should have: "Hello " (plain), "bold" (strong), " world" (plain)
	foundStrong := false
	for _, n := range inlines {
		node := n.(map[string]any)
		if marks, ok := node["marks"].([]any); ok {
			for _, m := range marks {
				mark := m.(map[string]any)
				if mark["type"] == "strong" {
					assert.Equal(t, "bold", node["text"])
					foundStrong = true
				}
			}
		}
	}
	assert.True(t, foundStrong, "expected a node with 'strong' mark")
}

func TestMarkdownToADF_Italic(t *testing.T) {
	md := "Hello *italic* world"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	para := content[0].(map[string]any)
	inlines := para["content"].([]any)

	foundEm := false
	for _, n := range inlines {
		node := n.(map[string]any)
		if marks, ok := node["marks"].([]any); ok {
			for _, m := range marks {
				mark := m.(map[string]any)
				if mark["type"] == "em" {
					assert.Equal(t, "italic", node["text"])
					foundEm = true
				}
			}
		}
	}
	assert.True(t, foundEm, "expected a node with 'em' mark")
}

func TestMarkdownToADF_InlineCode(t *testing.T) {
	md := "Hello `code` world"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	para := content[0].(map[string]any)
	inlines := para["content"].([]any)

	foundCode := false
	for _, n := range inlines {
		node := n.(map[string]any)
		if marks, ok := node["marks"].([]any); ok {
			for _, m := range marks {
				mark := m.(map[string]any)
				if mark["type"] == "code" {
					assert.Equal(t, "code", node["text"])
					foundCode = true
				}
			}
		}
	}
	assert.True(t, foundCode, "expected a node with 'code' mark")
}

func TestMarkdownToADF_Link(t *testing.T) {
	md := "Click [here](https://example.com) now"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	para := content[0].(map[string]any)
	inlines := para["content"].([]any)

	foundLink := false
	for _, n := range inlines {
		node := n.(map[string]any)
		if marks, ok := node["marks"].([]any); ok {
			for _, m := range marks {
				mark := m.(map[string]any)
				if mark["type"] == "link" {
					assert.Equal(t, "here", node["text"])
					attrs := mark["attrs"].(map[string]any)
					assert.Equal(t, "https://example.com", attrs["href"])
					foundLink = true
				}
			}
		}
	}
	assert.True(t, foundLink, "expected a node with 'link' mark")
}

func TestMarkdownToADF_Strikethrough(t *testing.T) {
	md := "Hello ~~gone~~ world"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	para := content[0].(map[string]any)
	inlines := para["content"].([]any)

	foundStrike := false
	for _, n := range inlines {
		node := n.(map[string]any)
		if marks, ok := node["marks"].([]any); ok {
			for _, m := range marks {
				mark := m.(map[string]any)
				if mark["type"] == "strike" {
					assert.Equal(t, "gone", node["text"])
					foundStrike = true
				}
			}
		}
	}
	assert.True(t, foundStrike, "expected a node with 'strike' mark")
}

func TestMarkdownToADF_MixedInline(t *testing.T) {
	md := "**bold** and *italic*"
	adf := MarkdownToADF(md)

	content := adf["content"].([]any)
	require.Len(t, content, 1)

	para := content[0].(map[string]any)
	inlines := para["content"].([]any)

	foundStrong := false
	foundEm := false
	for _, n := range inlines {
		node := n.(map[string]any)
		if marks, ok := node["marks"].([]any); ok {
			for _, m := range marks {
				mark := m.(map[string]any)
				if mark["type"] == "strong" {
					foundStrong = true
				}
				if mark["type"] == "em" {
					foundEm = true
				}
			}
		}
	}
	assert.True(t, foundStrong, "expected a 'strong' mark")
	assert.True(t, foundEm, "expected an 'em' mark")
}

func TestMarkdownToADF_EmptyInput(t *testing.T) {
	adf := MarkdownToADF("")

	assert.Equal(t, "doc", adf["type"])
	assert.Equal(t, 1, adf["version"])

	content := adf["content"].([]any)
	require.Len(t, content, 1, "empty input should produce one empty paragraph")

	para := content[0].(map[string]any)
	assert.Equal(t, "paragraph", para["type"])
}

func TestMarkdownToADF_BlankLines(t *testing.T) {
	adf := MarkdownToADF("\n\n\n")

	content := adf["content"].([]any)
	require.Len(t, content, 1, "blank-only input should produce one empty paragraph")

	para := content[0].(map[string]any)
	assert.Equal(t, "paragraph", para["type"])
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

func TestADFToMarkdown_RoundTrip_Bold(t *testing.T) {
	original := "Hello **bold** world"
	adf := MarkdownToADF(original)
	result := ADFToMarkdown(adf)
	assert.Equal(t, original, result)
}

func TestADFToMarkdown_RoundTrip_Italic(t *testing.T) {
	original := "Hello *italic* world"
	adf := MarkdownToADF(original)
	result := ADFToMarkdown(adf)
	assert.Equal(t, original, result)
}

func TestMarkdownToADF_JSON(t *testing.T) {
	md := "Hello **bold** world"
	adf := MarkdownToADF(md)

	// Should be valid JSON
	data, err := json.Marshal(adf)
	require.NoError(t, err)

	// Verify bold mark is present in serialized JSON
	assert.Contains(t, string(data), `"strong"`)
}
