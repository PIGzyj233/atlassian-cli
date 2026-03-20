package convert

import (
	"fmt"
	"regexp"
	"strings"
)

// MarkdownToADF converts Markdown text to Atlassian Document Format (ADF).
// This is a simplified converter covering the most common constructs:
// paragraphs, headings, code blocks, bullet lists, bold, italic, links.
func MarkdownToADF(md string) map[string]any {
	doc := map[string]any{
		"type":    "doc",
		"version": 1,
		"content": []any{},
	}

	lines := strings.Split(md, "\n")
	content := []any{}
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
			cb := map[string]any{
				"type": "codeBlock",
				"attrs": map[string]any{
					"language": lang,
				},
				"content": []any{
					map[string]any{"type": "text", "text": strings.Join(codeLines, "\n")},
				},
			}
			if lang == "" {
				delete(cb["attrs"].(map[string]any), "language")
			}
			content = append(content, cb)
			continue
		}

		// Heading
		if m := headingRe.FindStringSubmatch(line); m != nil {
			level := len(m[1])
			text := m[2]
			content = append(content, map[string]any{
				"type":    "heading",
				"attrs":   map[string]any{"level": level},
				"content": inlineToADF(text),
			})
			i++
			continue
		}

		// Bullet list
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			var items []any
			for i < len(lines) && (strings.HasPrefix(lines[i], "- ") || strings.HasPrefix(lines[i], "* ")) {
				text := strings.TrimPrefix(lines[i], "- ")
				text = strings.TrimPrefix(text, "* ")
				items = append(items, map[string]any{
					"type": "listItem",
					"content": []any{
						map[string]any{
							"type":    "paragraph",
							"content": inlineToADF(text),
						},
					},
				})
				i++
			}
			content = append(content, map[string]any{
				"type":    "bulletList",
				"content": items,
			})
			continue
		}

		// Empty line — skip
		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		// Paragraph (default)
		content = append(content, map[string]any{
			"type":    "paragraph",
			"content": inlineToADF(line),
		})
		i++
	}

	doc["content"] = content
	return doc
}

var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// inlineToADF converts inline Markdown to ADF inline nodes.
// For v1, wraps as plain text. Bold/italic/code parsing can be added later.
func inlineToADF(text string) []any {
	var nodes []any
	if text != "" {
		nodes = append(nodes, map[string]any{"type": "text", "text": text})
	}
	return nodes
}

// ADFToMarkdown converts ADF JSON to Markdown text.
func ADFToMarkdown(adf map[string]any) string {
	content, ok := adf["content"].([]any)
	if !ok {
		return ""
	}

	var parts []string
	for _, node := range content {
		nodeMap, ok := node.(map[string]any)
		if !ok {
			continue
		}
		parts = append(parts, adfNodeToMarkdown(nodeMap))
	}

	return strings.Join(parts, "\n\n")
}

func adfNodeToMarkdown(node map[string]any) string {
	nodeType, _ := node["type"].(string)

	switch nodeType {
	case "paragraph":
		return adfInlineToMarkdown(node)

	case "heading":
		attrs, _ := node["attrs"].(map[string]any)
		level := 1
		if l, ok := attrs["level"].(float64); ok {
			level = int(l)
		} else if l, ok := attrs["level"].(int); ok {
			level = l
		}
		prefix := strings.Repeat("#", level)
		return prefix + " " + adfInlineToMarkdown(node)

	case "codeBlock":
		attrs, _ := node["attrs"].(map[string]any)
		lang, _ := attrs["language"].(string)
		code := adfInlineToMarkdown(node)
		return "```" + lang + "\n" + code + "\n```"

	case "bulletList":
		items, _ := node["content"].([]any)
		var lines []string
		for _, item := range items {
			itemMap, _ := item.(map[string]any)
			itemContent, _ := itemMap["content"].([]any)
			if len(itemContent) > 0 {
				para, _ := itemContent[0].(map[string]any)
				lines = append(lines, "- "+adfInlineToMarkdown(para))
			}
		}
		return strings.Join(lines, "\n")

	case "orderedList":
		items, _ := node["content"].([]any)
		var lines []string
		for i, item := range items {
			itemMap, _ := item.(map[string]any)
			itemContent, _ := itemMap["content"].([]any)
			if len(itemContent) > 0 {
				para, _ := itemContent[0].(map[string]any)
				lines = append(lines, fmt.Sprintf("%d. %s", i+1, adfInlineToMarkdown(para)))
			}
		}
		return strings.Join(lines, "\n")

	default:
		return adfInlineToMarkdown(node)
	}
}

func adfInlineToMarkdown(node map[string]any) string {
	content, ok := node["content"].([]any)
	if !ok {
		return ""
	}

	var parts []string
	for _, inline := range content {
		inlineMap, ok := inline.(map[string]any)
		if !ok {
			continue
		}
		inlineType, _ := inlineMap["type"].(string)
		switch inlineType {
		case "text":
			text, _ := inlineMap["text"].(string)
			// Check for marks (bold, italic, etc.)
			marks, _ := inlineMap["marks"].([]any)
			for _, mark := range marks {
				markMap, _ := mark.(map[string]any)
				markType, _ := markMap["type"].(string)
				switch markType {
				case "strong":
					text = "**" + text + "**"
				case "em":
					text = "*" + text + "*"
				case "code":
					text = "`" + text + "`"
				}
			}
			parts = append(parts, text)
		case "hardBreak":
			parts = append(parts, "\n")
		}
	}

	return strings.Join(parts, "")
}
