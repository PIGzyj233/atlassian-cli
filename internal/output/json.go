package output

import (
	"encoding/json"
	"io"
)

// JSONFormatter outputs data as indented JSON.
type JSONFormatter struct{}

func (f *JSONFormatter) Write(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}
