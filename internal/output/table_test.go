package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableFormatter_MapSlice(t *testing.T) {
	f := &TableFormatter{}
	data := []map[string]any{
		{"key": "TEST-1", "summary": "First issue", "status": "Open"},
		{"key": "TEST-2", "summary": "Second issue", "status": "Closed"},
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "TEST-1")
	assert.Contains(t, output, "TEST-2")
	assert.Contains(t, output, "First issue")
}

func TestTableFormatter_SingleMap(t *testing.T) {
	f := &TableFormatter{}
	data := map[string]any{
		"key":     "TEST-1",
		"summary": "Single issue",
	}

	var buf bytes.Buffer
	err := f.Write(&buf, data)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "TEST-1")
}
