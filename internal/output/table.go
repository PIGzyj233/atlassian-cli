package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
)

// TableFormatter outputs data as an ASCII table.
type TableFormatter struct{}

func (f *TableFormatter) Write(w io.Writer, data any) error {
	switch v := data.(type) {
	case []map[string]any:
		return f.writeMapSlice(w, v)
	case map[string]any:
		return f.writeSingleMap(w, v)
	default:
		// Fallback to JSON for unsupported types
		return (&JSONFormatter{}).Write(w, data)
	}
}

func (f *TableFormatter) writeMapSlice(w io.Writer, items []map[string]any) error {
	if len(items) == 0 {
		fmt.Fprintln(w, "No results")
		return nil
	}

	// Extract headers from first item, sorted
	var headers []string
	for k := range items[0] {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	table := tablewriter.NewWriter(w)
	table.Header(toAnySlice(headers)...)

	for _, item := range items {
		var row []any
		for _, h := range headers {
			row = append(row, fmt.Sprintf("%v", item[h]))
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

func (f *TableFormatter) writeSingleMap(w io.Writer, item map[string]any) error {
	// Key-value format for single items
	var keys []string
	for k := range item {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	table := tablewriter.NewWriter(w)

	for _, k := range keys {
		table.Append(k, fmt.Sprintf("%v", item[k]))
	}

	table.Render()
	return nil
}

func toAnySlice(s []string) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}
