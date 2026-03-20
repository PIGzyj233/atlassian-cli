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

	// Ensure at least one content node (matches Python reference implementation)
	if len(content) == 0 {
		content = append(content, map[string]any{
			"type":    "paragraph",
			"content": []any{},
		})
	}

	doc["content"] = content
	return doc
}

var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// inlineRe matches Markdown inline formatting tokens.
// Order matters: code → bold → strikethrough → link → italic.
var inlineRe = regexp.MustCompile(
	"`([^`]+)`" + // inline code
		`|\*\*(.+?)\*\*` + // bold
		`|~~(.+?)~~` + // strikethrough
		`|\[([^\]]+)\]\(([^)]+)\)` + // link [text](url)
		`|(?:^|[^*])\*([^*]+?)\*(?:[^*]|$)`, // italic (non-greedy, avoid ** overlap)
)

// inlineToADF converts inline Markdown to ADF inline nodes with marks.
// Handles bold (**), italic (*), inline code (`), links, and strikethrough (~~).
func inlineToADF(text string) []any {
	if text == "" {
		return []any{}
	}

	var nodes []any
	pos := 0

	for _, loc := range inlineRe.FindAllStringSubmatchIndex(text, -1) {
		matchStart := loc[0]
		matchEnd := loc[1]

		// Add plain text before this match
		if matchStart > pos {
			plain := text[pos:matchStart]
			if plain != "" {
				nodes = append(nodes, map[string]any{"type": "text", "text": plain})
			}
		}

		// Determine which group matched (groups are at indices 2/3, 4/5, 6/7, 8/9+10/11, 12/13)
		switch {
		case loc[2] >= 0 && loc[3] >= 0: // inline code
			nodes = append(nodes, map[string]any{
				"type":  "text",
				"text":  text[loc[2]:loc[3]],
				"marks": []any{map[string]any{"type": "code"}},
			})
		case loc[4] >= 0 && loc[5] >= 0: // bold
			nodes = append(nodes, map[string]any{
				"type":  "text",
				"text":  text[loc[4]:loc[5]],
				"marks": []any{map[string]any{"type": "strong"}},
			})
		case loc[6] >= 0 && loc[7] >= 0: // strikethrough
			nodes = append(nodes, map[string]any{
				"type":  "text",
				"text":  text[loc[6]:loc[7]],
				"marks": []any{map[string]any{"type": "strike"}},
			})
		case loc[8] >= 0 && loc[9] >= 0: // link text
			href := text[loc[10]:loc[11]]
			nodes = append(nodes, map[string]any{
				"type": "text",
				"text": text[loc[8]:loc[9]],
				"marks": []any{map[string]any{
					"type":  "link",
					"attrs": map[string]any{"href": href},
				}},
			})
		case loc[12] >= 0 && loc[13] >= 0: // italic
			// The italic regex may capture a leading/trailing non-* char;
			// adjust matchStart/matchEnd to only consume the *...* portion.
			italicContent := text[loc[12]:loc[13]]
			// Find exact position of *content* within the overall match
			starIdx := strings.Index(text[matchStart:matchEnd], "*"+italicContent+"*")
			if starIdx >= 0 {
				actualStart := matchStart + starIdx
				actualEnd := actualStart + len(italicContent) + 2 // +2 for the two *
				// Emit any prefix char before the *
				if actualStart > pos {
					prefix := text[pos:actualStart]
					if matchStart > pos {
						// We already emitted text[pos:matchStart] above, but for italic
						// the match may include a leading char. Re-check.
					}
					if prefix != "" && actualStart > matchStart {
						// Remove already-emitted plain text, add back with correct boundary
						// Actually the plain text above used matchStart, so if actualStart > matchStart,
						// we need to emit text[matchStart:actualStart]
						if actualStart > matchStart {
							nodes = append(nodes, map[string]any{"type": "text", "text": text[matchStart:actualStart]})
						}
					}
				}
				matchEnd = actualEnd
			}
			nodes = append(nodes, map[string]any{
				"type":  "text",
				"text":  italicContent,
				"marks": []any{map[string]any{"type": "em"}},
			})
			pos = matchEnd
			continue
		}

		pos = matchEnd
	}

	// Remaining text after last match
	if pos < len(text) {
		remaining := text[pos:]
		if remaining != "" {
			nodes = append(nodes, map[string]any{"type": "text", "text": remaining})
		}
	}

	// If nothing matched, return whole text as plain
	if len(nodes) == 0 {
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
