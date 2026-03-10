package pipeline

import (
	"strings"
	"testing"
)

func TestContainsAlert(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Everything looks fine", false},
		{"ALERT: trust anomaly in builder agent", true},
		{"VIOLATION: soul values breached", true},
		{"QUARANTINE agent builder_01", true},
		{"Line one\nALERT: something wrong\nLine three", true},

		// Negative — keywords embedded in prose, not line-start directives.
		{"Found a VIOLATION of soul values", false},
		{"Minor alert about formatting", false},
		{"No VIOLATIONS DETECTED", false},
		{"The code is clean", false},
		{"", false},
		{"halt operations immediately", false}, // HALT handled separately
	}

	for _, tt := range tests {
		got := containsAlert(tt.input)
		if got != tt.want {
			t.Errorf("containsAlert(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseFiles(t *testing.T) {
	input := `--- FILE: main.go ---
package main

func main() {}
--- FILE: lib/util.go ---
package lib

func Helper() string {
	return "hello"
}
--- FILE: main_test.go ---
package main

import "testing"

func TestMain(t *testing.T) {}
`
	files := parseFiles(input)

	if len(files) != 3 {
		t.Fatalf("parseFiles returned %d files, want 3", len(files))
	}

	if _, ok := files["main.go"]; !ok {
		t.Error("missing main.go")
	}
	if _, ok := files["lib/util.go"]; !ok {
		t.Error("missing lib/util.go")
	}
	if _, ok := files["main_test.go"]; !ok {
		t.Error("missing main_test.go")
	}

	if !strings.Contains(files["main.go"], "package main") {
		t.Error("main.go missing package declaration")
	}
	if !strings.Contains(files["lib/util.go"], "func Helper()") {
		t.Error("util.go missing Helper function")
	}
}

func TestParseFilesEmpty(t *testing.T) {
	files := parseFiles("just some text without markers")
	if len(files) != 0 {
		t.Errorf("parseFiles with no markers returned %d files, want 0", len(files))
	}
}

func TestParseFilesSingleFile(t *testing.T) {
	input := `--- FILE: app.py ---
def main():
    print("hello")

if __name__ == "__main__":
    main()
`
	files := parseFiles(input)
	if len(files) != 1 {
		t.Fatalf("parseFiles returned %d files, want 1", len(files))
	}
	if !strings.Contains(files["app.py"], "def main():") {
		t.Error("app.py missing main function")
	}
}

func TestLangExtension(t *testing.T) {
	tests := []struct {
		lang string
		want string
	}{
		{"go", ".go"},
		{"typescript", ".ts"},
		{"python", ".py"},
		{"rust", ".rs"},
		{"csharp", ".cs"},
		{"unknown", ".go"},
	}
	for _, tt := range tests {
		got := langExtension(tt.lang)
		if got != tt.want {
			t.Errorf("langExtension(%q) = %q, want %q", tt.lang, got, tt.want)
		}
	}
}

func TestLangTestCommand(t *testing.T) {
	cmd, args := langTestCommand("go")
	if cmd != "go" || args[0] != "test" {
		t.Errorf("go test command = %s %v", cmd, args)
	}

	cmd, args = langTestCommand("python")
	if cmd != "python" || args[1] != "pytest" {
		t.Errorf("python test command = %s %v", cmd, args)
	}

	cmd, args = langTestCommand("rust")
	if cmd != "cargo" || args[0] != "test" {
		t.Errorf("rust test command = %s %v", cmd, args)
	}
}

func TestReviewerModelDefault(t *testing.T) {
	// reviewTargeted uses Sonnet by default for targeted reviews (not the
	// reviewer role's default Opus). Verify via reviewerModel selection logic:
	// empty reviewerModel → "claude-sonnet-4-6".
	p := &Pipeline{reviewerModel: ""}
	model := p.targetedReviewModel()
	if model != "claude-sonnet-4-6" {
		t.Errorf("default targeted review model = %q, want %q", model, "claude-sonnet-4-6")
	}
}

func TestReviewerModelOverride(t *testing.T) {
	// Config.ReviewerModel propagates to pipeline and overrides the default.
	p := &Pipeline{reviewerModel: "claude-haiku-4-5-20251001"}
	model := p.targetedReviewModel()
	if model != "claude-haiku-4-5-20251001" {
		t.Errorf("overridden targeted review model = %q, want %q", model, "claude-haiku-4-5-20251001")
	}
}

func TestExtractLanguage(t *testing.T) {
	p := &Pipeline{}

	tests := []struct {
		design string
		want   string
	}{
		{"LANGUAGE: go\n\nEntity(Task)...", "go"},
		{"LANGUAGE: typescript\nSome spec", "typescript"},
		{"  LANGUAGE:  python \nstuff", "python"},
		{"No language specified here", "go"},
		{"language: rust\nspec", "rust"},
	}
	for _, tt := range tests {
		got := p.extractLanguage(tt.design)
		if got != tt.want {
			t.Errorf("extractLanguage(%q) = %q, want %q", tt.design[:20], got, tt.want)
		}
	}
}

func TestSanitizeBranchName(t *testing.T) {
	tests := []struct {
		desc string
		want string
	}{
		// Short input — no truncation
		{"add login page", "add-login-page"},

		// Exactly 40 chars — no truncation needed
		{"a234567890123456789012345678901234567890", "a234567890123456789012345678901234567890"},

		// Over 40 chars — truncate at last word boundary before 40
		{"build a task management app with kanban boards and dashboards", "build-a-task-management-app-with-kanban"},

		// Char 40 is mid-word — truncate at prior word boundary
		{"add comprehensive authentication support for enterprise users", "add-comprehensive-authentication"},

		// First word alone exceeds 40 chars — hard truncate fallback
		{"abcdefghijklmnopqrstuvwxyz1234567890abcdefghij", "abcdefghijklmnopqrstuvwxyz1234567890abcd"},

		// Empty input
		{"", "change"},

		// Non-alphanumeric only
		{"!@#$%", "change"},

		// Underscores and slashes become hyphens
		{"my_feature/branch name", "my-feature-branch-name"},
	}
	for _, tt := range tests {
		got := sanitizeBranchName(tt.desc)
		if got != tt.want {
			t.Errorf("sanitizeBranchName(%q) = %q, want %q", tt.desc, got, tt.want)
		}
	}
}

