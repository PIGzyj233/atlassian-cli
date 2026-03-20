package convert

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// MarkdownToStorage converts Markdown to Confluence Storage Format (XHTML).
// Covers: paragraphs, headings, code blocks, lists, bold, italic, links.
func MarkdownToStorage(md string) string {
	var result strings.Builder
	lines := strings.Split(md, "\n")
	i := 0

	for i < len(lines) {
		line := lines[i]

		// Code block
		if strings.HasPrefix(line, "```") {
			lang := strings.TrimPrefix(line, "```")
			lang = strings.TrimSpace(lang)
			var codeLines []string
			i++
			for i < len(lines) && !strings.HasPrefix(lines[i], "```") {
				codeLines = append(codeLines, lines[i])
				i++
			}
			i++ // skip closing ```

			result.WriteString(`<ac:structured-macro ac:name="code">`)
			if lang != "" {
				result.WriteString(`<ac:parameter ac:name="language">` + lang + `</ac:parameter>`)
			}
			result.WriteString(`<ac:plain-text-body><![CDATA[` + strings.Join(codeLines, "\n") + `]]></ac:plain-text-body>`)
			result.WriteString(`</ac:structured-macro>`)
			continue
		}

		// Heading
		if m := storageHeadingRe.FindStringSubmatch(line); m != nil {
			level := len(m[1])
			text := convertInlineToHTML(m[2])
			result.WriteString(fmt.Sprintf("<h%d>%s</h%d>", level, text, level))
			i++
			continue
		}

		// Bullet list
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			result.WriteString("<ul>")
			for i < len(lines) && (strings.HasPrefix(lines[i], "- ") || strings.HasPrefix(lines[i], "* ")) {
				text := strings.TrimPrefix(lines[i], "- ")
				text = strings.TrimPrefix(text, "* ")
				result.WriteString("<li>" + convertInlineToHTML(text) + "</li>")
				i++
			}
			result.WriteString("</ul>")
			continue
		}

		// Empty line — skip
		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		// Paragraph (default)
		result.WriteString("<p>" + convertInlineToHTML(line) + "</p>")
		i++
	}

	return result.String()
}

var storageHeadingRe = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
var boldRe = regexp.MustCompile(`\*\*(.+?)\*\*`)
var italicRe = regexp.MustCompile(`\*(.+?)\*`)
var codeInlineRe = regexp.MustCompile("`([^`]+)`")
var linkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

// convertInlineToHTML converts inline Markdown to HTML.
func convertInlineToHTML(text string) string {
	// Order matters: bold before italic since ** contains *
	text = boldRe.ReplaceAllString(text, "<strong>$1</strong>")
	text = italicRe.ReplaceAllString(text, "<em>$1</em>")
	text = codeInlineRe.ReplaceAllString(text, "<code>$1</code>")
	text = linkRe.ReplaceAllString(text, `<a href="$2">$1</a>`)
	return text
}

// StorageToMarkdown converts Confluence Storage Format (XHTML) to Markdown.
// This is a simplified converter for common elements.
func StorageToMarkdown(storage string) string {
	result := storage

	// Headings
	for level := 1; level <= 6; level++ {
		prefix := strings.Repeat("#", level)
		openTag := fmt.Sprintf("<h%d>", level)
		closeTag := fmt.Sprintf("</h%d>", level)
		for strings.Contains(result, openTag) {
			start := strings.Index(result, openTag)
			end := strings.Index(result, closeTag)
			if start >= 0 && end > start {
				content := result[start+len(openTag) : end]
				replacement := prefix + " " + content + "\n\n"
				result = result[:start] + replacement + result[end+len(closeTag):]
			} else {
				break
			}
		}
	}

	// Paragraphs
	result = strings.ReplaceAll(result, "<p>", "")
	result = strings.ReplaceAll(result, "</p>", "\n\n")

	// Bold
	result = strings.ReplaceAll(result, "<strong>", "**")
	result = strings.ReplaceAll(result, "</strong>", "**")

	// Italic
	result = strings.ReplaceAll(result, "<em>", "*")
	result = strings.ReplaceAll(result, "</em>", "*")

	// Code
	result = strings.ReplaceAll(result, "<code>", "`")
	result = strings.ReplaceAll(result, "</code>", "`")

	// Lists
	result = strings.ReplaceAll(result, "<ul>", "")
	result = strings.ReplaceAll(result, "</ul>", "")
	result = strings.ReplaceAll(result, "<li>", "- ")
	result = strings.ReplaceAll(result, "</li>", "\n")

	// HTML unescape
	result = html.UnescapeString(result)

	// Clean up excessive whitespace
	result = strings.TrimSpace(result)

	return result
}
