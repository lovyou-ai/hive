package pipeline

import (
	"strings"
	"testing"
)

// ════════════════════════════════════════════════════════════════════════
// truncateLines
// ════════════════════════════════════════════════════════════════════════

func TestTruncateLines_NoTruncation(t *testing.T) {
	input := "line1\nline2\nline3"
	got := truncateLines(input, 10)
	if got != input {
		t.Errorf("truncateLines with sufficient limit = %q, want %q", got, input)
	}
}

func TestTruncateLines_ExactLimit(t *testing.T) {
	input := "a\nb\nc"
	got := truncateLines(input, 3)
	if got != input {
		t.Errorf("truncateLines at exact limit = %q, want %q", got, input)
	}
}

func TestTruncateLines_Truncation(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5"
	got := truncateLines(input, 3)
	if !strings.HasPrefix(got, "line1\nline2\nline3") {
		t.Errorf("truncateLines should keep first 3 lines, got %q", got)
	}
	if !strings.Contains(got, "[truncated: 2 lines omitted]") {
		t.Errorf("truncateLines should append truncation notice, got %q", got)
	}
}

func TestTruncateLines_NoticeFormat(t *testing.T) {
	// Verify the exact notice format and content boundary.
	lines := make([]string, 10)
	for i := range lines {
		lines[i] = "line"
	}
	input := strings.Join(lines, "\n")
	got := truncateLines(input, 5)
	want := "line\nline\nline\nline\nline\n... [truncated: 5 lines omitted]"
	if got != want {
		t.Errorf("truncateLines notice = %q, want %q", got, want)
	}
}

// ════════════════════════════════════════════════════════════════════════
// filterSelfImproveFiles
// ════════════════════════════════════════════════════════════════════════

func TestFilterSelfImproveFiles_IncludesGoFiles(t *testing.T) {
	files := map[string]string{
		"pkg/foo/foo.go": "package foo\n",
		"pkg/bar/bar.go": "package bar\n",
	}
	got := filterSelfImproveFiles(files)
	if _, ok := got["pkg/foo/foo.go"]; !ok {
		t.Error("filterSelfImproveFiles should include .go files")
	}
	if _, ok := got["pkg/bar/bar.go"]; !ok {
		t.Error("filterSelfImproveFiles should include .go files")
	}
}

func TestFilterSelfImproveFiles_ExcludesTestFiles(t *testing.T) {
	files := map[string]string{
		"pkg/foo/foo.go":      "package foo\n",
		"pkg/foo/foo_test.go": "package foo\n",
	}
	got := filterSelfImproveFiles(files)
	if _, ok := got["pkg/foo/foo_test.go"]; ok {
		t.Error("filterSelfImproveFiles should exclude _test.go files")
	}
	if _, ok := got["pkg/foo/foo.go"]; !ok {
		t.Error("filterSelfImproveFiles should include non-test .go files")
	}
}

func TestFilterSelfImproveFiles_ExcludesNonGoNonConfig(t *testing.T) {
	files := map[string]string{
		"pkg/foo/foo.go": "package foo\n",
		"Makefile":       "all:\n\tgo build\n",
		"README.md":      "# readme\n",
		"go.sum":         "hash\n",
	}
	got := filterSelfImproveFiles(files)
	if _, ok := got["Makefile"]; ok {
		t.Error("filterSelfImproveFiles should exclude Makefile")
	}
	if _, ok := got["README.md"]; ok {
		t.Error("filterSelfImproveFiles should exclude README.md")
	}
	if _, ok := got["go.sum"]; ok {
		t.Error("filterSelfImproveFiles should exclude go.sum")
	}
}

func TestFilterSelfImproveFiles_IncludesConfigFiles(t *testing.T) {
	files := map[string]string{
		"CLAUDE.md": "# hive\n",
		"go.mod":    "module test\n",
	}
	got := filterSelfImproveFiles(files)
	if _, ok := got["CLAUDE.md"]; !ok {
		t.Error("filterSelfImproveFiles should include CLAUDE.md")
	}
	if _, ok := got["go.mod"]; !ok {
		t.Error("filterSelfImproveFiles should include go.mod")
	}
}

func TestFilterSelfImproveFiles_TruncatesAtMaxLines(t *testing.T) {
	// Build content with more than maxSelfImproveFileLines lines.
	total := maxSelfImproveFileLines + 50
	src := make([]string, total)
	for i := range src {
		src[i] = "line"
	}
	content := strings.Join(src, "\n")

	files := map[string]string{
		"pkg/foo/foo.go": content,
	}
	got := filterSelfImproveFiles(files)
	result := got["pkg/foo/foo.go"]
	if !strings.Contains(result, "[truncated: 50 lines omitted]") {
		t.Errorf("filterSelfImproveFiles should truncate at %d lines; result does not contain truncation notice",
			maxSelfImproveFileLines)
	}
}

func TestFilterSelfImproveFiles_ShortFileNotTruncated(t *testing.T) {
	content := "package foo\n\nfunc Foo() {}\n"
	files := map[string]string{
		"pkg/foo/foo.go": content,
	}
	got := filterSelfImproveFiles(files)
	if got["pkg/foo/foo.go"] != content {
		t.Errorf("filterSelfImproveFiles modified a short file: %q", got["pkg/foo/foo.go"])
	}
}
