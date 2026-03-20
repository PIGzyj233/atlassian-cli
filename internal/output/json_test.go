package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFormatter_Format(t *testing.T) {
	f := &JSONFormatter{}
	data := map[string]any{
		"key":     "TEST-1",
		"summary": "Test issue",
		"status":  "Open",
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `"key": "TEST-1"`)
	assert.Contains(t, output, `"summary": "Test issue"`)
	// Verify indented output
	assert.Contains(t, output, "  ")
}

func TestJSONFormatter_FormatList(t *testing.T) {
	f := &JSONFormatter{}
	data := []map[string]string{
		{"key": "A-1"},
		{"key": "A-2"},
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), `"key": "A-1"`)
}
