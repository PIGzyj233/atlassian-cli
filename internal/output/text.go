package output

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// TextFormatter outputs data as human-readable plain text.
type TextFormatter struct{}

func (f *TextFormatter) Write(w io.Writer, data any) error {
	switch v := data.(type) {
	case map[string]any:
		return f.writeSingleMap(w, v)
	case []map[string]any:
		return f.writeMapSlice(w, v)
	default:
		return (&JSONFormatter{}).Write(w, data)
	}
}

func (f *TextFormatter) writeSingleMap(w io.Writer, item map[string]any) error {
	var keys []string
	for k := range item {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Find longest key for alignment
	maxLen := 0
	for _, k := range keys {
		if len(k) > maxLen {
			maxLen = len(k)
		}
	}

	for _, k := range keys {
		val := fmt.Sprintf("%v", item[k])
		// Truncate long values
		if len(val) > 200 {
			val = val[:200] + "..."
		}
		padding := strings.Repeat(" ", maxLen-len(k))
		fmt.Fprintf(w, "%s:%s %s\n", k, padding, val)
	}
	return nil
}

func (f *TextFormatter) writeMapSlice(w io.Writer, items []map[string]any) error {
	if len(items) == 0 {
		fmt.Fprintln(w, "No results")
		return nil
	}
	for i, item := range items {
		if i > 0 {
			fmt.Fprintln(w, "---")
		}
		if err := f.writeSingleMap(w, item); err != nil {
			return err
		}
	}
	return nil
}
