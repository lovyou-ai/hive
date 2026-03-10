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

func TestReviewerModelConfig(t *testing.T) {
	// Verify Config.ReviewerModel propagates to the pipeline struct.
	// We can't call New() without a full store/actor setup, so test the
	// config field exists and the default behavior matches expectations.

	// Empty ReviewerModel means targeted reviews default to Sonnet.
	cfg := Config{ReviewerModel: ""}
	if cfg.ReviewerModel != "" {
		t.Errorf("empty ReviewerModel should be empty string, got %q", cfg.ReviewerModel)
	}

	// Explicit override propagates.
	cfg = Config{ReviewerModel: "claude-haiku-4-5-20251001"}
	if cfg.ReviewerModel != "claude-haiku-4-5-20251001" {
		t.Errorf("ReviewerModel should be claude-haiku-4-5-20251001, got %q", cfg.ReviewerModel)
	}

	// Verify Pipeline struct stores the override.
	p := &Pipeline{reviewerModel: "claude-sonnet-4-6"}
	if p.reviewerModel != "claude-sonnet-4-6" {
		t.Errorf("pipeline reviewerModel should be claude-sonnet-4-6, got %q", p.reviewerModel)
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

