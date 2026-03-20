package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextFormatter_SingleMap(t *testing.T) {
	f := &TextFormatter{}
	data := map[string]any{
		"key":     "TEST-1",
		"summary": "My issue title",
		"status":  "Open",
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "key:")
	assert.Contains(t, output, "TEST-1")
	assert.Contains(t, output, "summary:")
}

func TestTextFormatter_MapSlice(t *testing.T) {
	f := &TextFormatter{}
	data := []map[string]any{
		{"key": "A-1", "summary": "First"},
		{"key": "A-2", "summary": "Second"},
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "A-1")
	assert.Contains(t, output, "A-2")
}
