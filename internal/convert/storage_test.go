package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownToStorage_Paragraph(t *testing.T) {
	md := "Hello world"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<p>Hello world</p>")
}

func TestMarkdownToStorage_Heading(t *testing.T) {
	md := "# Title\n\n## Subtitle"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<h1>Title</h1>")
	assert.Contains(t, html, "<h2>Subtitle</h2>")
}

func TestMarkdownToStorage_CodeBlock(t *testing.T) {
	md := "```python\nprint('hello')\n```"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, `<ac:structured-macro ac:name="code">`)
	assert.Contains(t, html, "python")
	// CDATA sections preserve content without HTML-escaping
	assert.Contains(t, html, "print('hello')")
}

func TestMarkdownToStorage_BulletList(t *testing.T) {
	md := "- Item 1\n- Item 2"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<ul>")
	assert.Contains(t, html, "<li>Item 1</li>")
	assert.Contains(t, html, "<li>Item 2</li>")
}

func TestMarkdownToStorage_Bold(t *testing.T) {
	md := "Hello **bold** world"
	html := MarkdownToStorage(md)
	assert.Contains(t, html, "<strong>bold</strong>")
}

func TestStorageToMarkdown_Paragraph(t *testing.T) {
	storage := "<p>Hello world</p>"
	md := StorageToMarkdown(storage)
	assert.Contains(t, md, "Hello world")
}

func TestStorageToMarkdown_Heading(t *testing.T) {
	storage := "<h1>Title</h1>"
	md := StorageToMarkdown(storage)
	assert.Contains(t, md, "# Title")
}
