package runner

import "testing"

func TestParseVerdict(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"pass", "Looks good.\n\nVERDICT: PASS", "PASS"},
		{"revise", "Missing allowlist.\nVERDICT: REVISE", "REVISE"},
		{"default", "No verdict line", "PASS"},
		{"whitespace", "  VERDICT:  PASS  ", "PASS"},
		{"middle", "Line 1\nVERDICT: REVISE\nLine 3", "REVISE"},
		{"invalid", "VERDICT: INVALID", "PASS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseVerdict(tt.input)
			if got != tt.expect {
				t.Errorf("parseVerdict(%q) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}

func TestExtractIssues(t *testing.T) {
	content := "Issue 1: missing allowlist entry\nIssue 2: no tests\n\nVERDICT: REVISE"
	got := extractIssues(content)
	if got != "Issue 1: missing allowlist entry\nIssue 2: no tests" {
		t.Errorf("extractIssues returned: %q", got)
	}
}

func TestBuildReviewPrompt(t *testing.T) {
	c := commit{hash: "abc123def456", subject: "[hive:builder] Add Policy"}
	diff := "+KindPolicy = \"policy\""

	prompt := buildReviewPrompt(c, diff)

	// Should contain the commit info.
	if !contains(prompt, "abc123def456") {
		t.Error("prompt missing commit hash")
	}
	if !contains(prompt, "[hive:builder] Add Policy") {
		t.Error("prompt missing commit subject")
	}
	if !contains(prompt, "+KindPolicy") {
		t.Error("prompt missing diff content")
	}
	// Should contain the checklist.
	if !contains(prompt, "Completeness") {
		t.Error("prompt missing checklist")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && searchString(s, sub)
}

func searchString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
