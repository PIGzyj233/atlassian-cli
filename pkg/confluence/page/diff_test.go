package page

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnifiedDiff_IdenticalContent(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	diff := unifiedDiff(lines, lines, "v1", "v2")
	assert.Empty(t, diff, "identical content should produce empty diff")
}

func TestUnifiedDiff_Insertion(t *testing.T) {
	from := []string{"a", "b"}
	to := []string{"a", "inserted", "b"}
	diff := unifiedDiff(from, to, "v1", "v2")

	assert.Contains(t, diff, "--- v1")
	assert.Contains(t, diff, "+++ v2")
	assert.Contains(t, diff, "@@")
	assert.Contains(t, diff, "+inserted")
}

func TestUnifiedDiff_Deletion(t *testing.T) {
	from := []string{"a", "removed", "b"}
	to := []string{"a", "b"}
	diff := unifiedDiff(from, to, "v1", "v2")

	assert.Contains(t, diff, "-removed")
	assert.NotContains(t, diff, "+removed")
}

func TestUnifiedDiff_Modification(t *testing.T) {
	from := []string{"hello world"}
	to := []string{"hello go"}
	diff := unifiedDiff(from, to, "v1", "v2")

	assert.Contains(t, diff, "-hello world")
	assert.Contains(t, diff, "+hello go")
}

func TestUnifiedDiff_HunkHeaders(t *testing.T) {
	// Generate enough context to get proper hunk headers
	from := make([]string, 10)
	to := make([]string, 10)
	for i := range from {
		from[i] = "line" + strings.Repeat(" ", i)
		to[i] = from[i]
	}
	to[5] = "CHANGED"

	diff := unifiedDiff(from, to, "old", "new")
	assert.Contains(t, diff, "@@", "should contain hunk header")
	assert.Contains(t, diff, "-line     ")
	assert.Contains(t, diff, "+CHANGED")
}
