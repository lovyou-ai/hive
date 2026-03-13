package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/actor"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/resources"
	"github.com/lovyou-ai/hive/pkg/roles"
	"github.com/lovyou-ai/hive/pkg/spawn"
	"github.com/lovyou-ai/hive/pkg/work"
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

func TestSelfImproveCTOModelDefault(t *testing.T) {
	// selfImproveCTOModel defaults to Sonnet — the CTO analyzes the full codebase
	// and needs stronger reasoning to find non-trivial improvements.
	p := &Pipeline{ctoModel: ""}
	model := p.selfImproveCTOModel()
	if model != "claude-sonnet-4-6" {
		t.Errorf("default self-improve CTO model = %q, want %q", model, "claude-sonnet-4-6")
	}
}

func TestSelfImproveCTOModelOverride(t *testing.T) {
	// Config.CTOModel propagates to pipeline and overrides the default.
	p := &Pipeline{ctoModel: "claude-opus-4-6"}
	model := p.selfImproveCTOModel()
	if model != "claude-opus-4-6" {
		t.Errorf("overridden self-improve CTO model = %q, want %q", model, "claude-opus-4-6")
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

// testSigner implements event.Signer for tests.
type testSigner struct{}

func (s *testSigner) Sign(data []byte) (types.Signature, error) {
	return types.NewSignature(make([]byte, 64))
}

func TestEnsureAgentNoSpawnerEmitsAuthorityEvents(t *testing.T) {
	s := store.NewInMemoryStore()
	actors := actor.NewInMemoryActorStore()

	// Register human.
	humanRawPub := spawn.DerivePublicKey("human:TestHuman")
	humanPub, err := types.NewPublicKey([]byte(humanRawPub))
	if err != nil {
		t.Fatal(err)
	}
	humanActor, err := actors.Register(humanPub, "TestHuman", event.ActorTypeHuman)
	if err != nil {
		t.Fatal(err)
	}
	humanID := humanActor.ID()

	// Bootstrap graph — ensureAgent needs a non-empty graph head.
	registry := event.DefaultRegistry()
	bsFactory := event.NewBootstrapFactory(registry)
	signer := &testSigner{}
	bootstrap, err := bsFactory.Init(humanID, signer)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Append(bootstrap); err != nil {
		t.Fatal(err)
	}

	factory := event.NewEventFactory(registry)
	convID, err := types.NewConversationID("conv_spawn_" + strings.Repeat("0", 24))
	if err != nil {
		t.Fatal(err)
	}

	p := &Pipeline{
		store:    s,
		actors:   actors,
		humanID:  humanID,
		signer:   signer,
		factory:  factory,
		convID:   convID,
		agents:   make(map[roles.Role]*roles.Agent),
		trackers: make(map[roles.Role]*resources.TrackingProvider),
		// spawner is nil — dev/bootstrap mode.
	}

	// ensureAgent emits authority events in the no-spawner branch.
	// Provider creation may succeed or fail depending on environment —
	// we only care that the authority events were emitted.
	_, _ = p.ensureAgent(context.Background(), roles.RoleBuilder, "test-builder")

	// Verify authority.requested event was emitted.
	authReqPage, err := s.ByType(event.EventTypeAuthorityRequested, 10, types.None[types.Cursor]())
	if err != nil {
		t.Fatal(err)
	}
	if len(authReqPage.Items()) == 0 {
		t.Error("expected authority.requested event")
	}

	// Verify authority.resolved event was emitted.
	authResPage, err := s.ByType(event.EventTypeAuthorityResolved, 10, types.None[types.Cursor]())
	if err != nil {
		t.Fatal(err)
	}
	if len(authResPage.Items()) == 0 {
		t.Error("expected authority.resolved event")
	}

	// Verify resolved content shows auto-approved.
	resolved := authResPage.Items()[0]
	content, ok := resolved.Content().(event.AuthorityResolvedContent)
	if !ok {
		t.Fatal("authority.resolved event has wrong content type")
	}
	if !content.Approved {
		t.Error("expected auto-approved resolution")
	}
	if content.Reason.IsSome() && content.Reason.Unwrap() != "auto-approved (no authority gate)" {
		t.Errorf("reason = %q, want %q", content.Reason.Unwrap(), "auto-approved (no authority gate)")
	}
}

func TestFindModuleDir(t *testing.T) {
	t.Run("marker in root", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module test"), 0644); err != nil {
			t.Fatal(err)
		}
		got := findModuleDir(root, "go")
		if got != root {
			t.Errorf("findModuleDir = %q, want root %q", got, root)
		}
	})

	t.Run("marker in subdir", func(t *testing.T) {
		root := t.TempDir()
		sub := filepath.Join(root, "go")
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "go.mod"), []byte("module test"), 0644); err != nil {
			t.Fatal(err)
		}
		got := findModuleDir(root, "go")
		if got != sub {
			t.Errorf("findModuleDir = %q, want subdir %q", got, sub)
		}
	})

	t.Run("no marker returns root", func(t *testing.T) {
		root := t.TempDir()
		got := findModuleDir(root, "go")
		if got != root {
			t.Errorf("findModuleDir = %q, want root %q", got, root)
		}
	})

	t.Run("multiple subdirs picks first alphabetical", func(t *testing.T) {
		root := t.TempDir()
		// Create two subdirs, only "beta" has the marker.
		for _, name := range []string{"alpha", "beta"} {
			if err := os.MkdirAll(filepath.Join(root, name), 0755); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(root, "beta", "package.json"), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}
		got := findModuleDir(root, "javascript")
		want := filepath.Join(root, "beta")
		if got != want {
			t.Errorf("findModuleDir = %q, want %q", got, want)
		}
	})

	t.Run("root preferred over subdir", func(t *testing.T) {
		root := t.TempDir()
		sub := filepath.Join(root, "sub")
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}
		// Marker in both root and subdir — root wins.
		for _, dir := range []string{root, sub} {
			if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644); err != nil {
				t.Fatal(err)
			}
		}
		got := findModuleDir(root, "go")
		if got != root {
			t.Errorf("findModuleDir = %q, want root %q", got, root)
		}
	})

	t.Run("python pyproject.toml", func(t *testing.T) {
		root := t.TempDir()
		sub := filepath.Join(root, "py")
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "pyproject.toml"), []byte("[project]"), 0644); err != nil {
			t.Fatal(err)
		}
		got := findModuleDir(root, "python")
		if got != sub {
			t.Errorf("findModuleDir = %q, want %q", got, sub)
		}
	})

	t.Run("rust Cargo.toml", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "Cargo.toml"), []byte("[package]"), 0644); err != nil {
			t.Fatal(err)
		}
		got := findModuleDir(root, "rust")
		if got != root {
			t.Errorf("findModuleDir = %q, want root %q", got, root)
		}
	})
}

func TestCTOAnalysisSkipsUnderstandPhase(t *testing.T) {
	// Verify that ProductInput.CTOAnalysis is properly formatted from a
	// SelfImproveRecommendation and that the field flows through to targeted
	// pipeline input, which would skip the CTO Evaluate call in Phase 2.
	rec := SelfImproveRecommendation{
		Description:    "Refactor the build phase",
		FilesToChange:  []string{"pkg/pipeline/pipeline.go", "pkg/pipeline/pipeline_test.go"},
		ExpectedImpact: "Reduce build phase duration by 30%",
		Priority:       "high",
	}

	// Format the same way RunSelfImprove does (FILES_TO_CHANGE: structured section).
	analysis := fmt.Sprintf("Description: %s\nFILES_TO_CHANGE:\n%s\nExpected impact: %s",
		rec.Description, strings.Join(rec.FilesToChange, "\n"), rec.ExpectedImpact)

	input := ProductInput{
		RepoPath:    "/tmp/fake-repo",
		Description: rec.Description,
		CTOAnalysis: analysis,
	}

	// CTOAnalysis must be non-empty so the Understand phase is skipped.
	if input.CTOAnalysis == "" {
		t.Fatal("CTOAnalysis should be non-empty")
	}

	// Verify the formatted string contains all recommendation fields.
	if !strings.Contains(input.CTOAnalysis, rec.Description) {
		t.Error("CTOAnalysis missing Description")
	}
	for _, f := range rec.FilesToChange {
		if !strings.Contains(input.CTOAnalysis, f) {
			t.Errorf("CTOAnalysis missing file %q", f)
		}
	}
	if !strings.Contains(input.CTOAnalysis, rec.ExpectedImpact) {
		t.Error("CTOAnalysis missing ExpectedImpact")
	}

	// Verify that an empty CTOAnalysis would NOT skip (the default path).
	defaultInput := ProductInput{
		RepoPath:    "/tmp/fake-repo",
		Description: "some change",
	}
	if defaultInput.CTOAnalysis != "" {
		t.Error("default ProductInput should have empty CTOAnalysis")
	}
}

func TestParseRelevantFiles(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantAny  []string // paths that must appear in result
		wantNone []string // tokens that must NOT appear
	}{
		{
			name: "self-improve bracket format",
			input: "Description: Refactor build phase\n" +
				"Files to change: [pkg/pipeline/pipeline_targeted.go pkg/pipeline/pipeline_helpers.go]\n" +
				"Expected impact: Reduce tokens by 30%",
			wantAny: []string{"pkg/pipeline/pipeline_targeted.go", "pkg/pipeline/pipeline_helpers.go"},
		},
		{
			name: "bullet point format",
			input: "- pkg/pipeline/pipeline_targeted.go: filter existingFiles in text mode\n" +
				"- pkg/pipeline/pipeline_helpers.go: add parseRelevantFiles helper\n" +
				"Key risks: may over-filter if CTO analysis is vague",
			wantAny: []string{"pkg/pipeline/pipeline_targeted.go", "pkg/pipeline/pipeline_helpers.go"},
		},
		{
			name:    "URLs not extracted",
			input:   "See https://example.com/foo.go for details\nFile: pkg/foo.go",
			wantAny: []string{"pkg/foo.go"},
			wantNone: []string{"https://example.com/foo.go"},
		},
		{
			name:    "version strings not extracted",
			input:   "Use version v1.2.3\nChange: pkg/bar.go",
			wantAny: []string{"pkg/bar.go"},
			wantNone: []string{"v1.2.3", "1.2.3"},
		},
		{
			name:    "no paths returns empty slice",
			input:   "No changes needed here at all",
			wantAny: nil,
		},
		{
			name:    "go.mod recognised",
			input:   "Update go.mod to add dependency",
			wantAny: []string{"go.mod"},
		},
		{
			name:    "deduplication",
			input:   "pkg/foo.go: change A\npkg/foo.go: change B",
			wantAny: []string{"pkg/foo.go"},
		},
		{
			name: "FILES_TO_CHANGE section — only section files extracted",
			input: "FILES_TO_CHANGE:\n" +
				"pkg/roles/roles.go — update CTO prompt\n" +
				"pkg/pipeline/pipeline_helpers.go — update parseRelevantFiles\n" +
				"\n" +
				"Key risks: parseRelevantFiles is also called from pipeline_targeted.go " +
				"and telemetry.go is mentioned here but should NOT be included.",
			wantAny:  []string{"pkg/roles/roles.go", "pkg/pipeline/pipeline_helpers.go"},
			wantNone: []string{"pipeline_targeted.go", "telemetry.go"},
		},
		{
			name: "FILES_TO_CHANGE section stops at non-path line",
			input: "FILES_TO_CHANGE:\n" +
				"pkg/foo/bar.go — change something\n" +
				"Key risks: some risk\n" +
				"pkg/should/not/be/included.go — this is after a non-path line",
			wantAny:  []string{"pkg/foo/bar.go"},
			wantNone: []string{"pkg/should/not/be/included.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRelevantFiles(tt.input)
			gotSet := make(map[string]bool, len(got))
			for _, p := range got {
				gotSet[p] = true
			}
			for _, want := range tt.wantAny {
				if !gotSet[want] {
					t.Errorf("parseRelevantFiles missing %q; got %v", want, got)
				}
			}
			for _, bad := range tt.wantNone {
				if gotSet[bad] {
					t.Errorf("parseRelevantFiles should not contain %q; got %v", bad, got)
				}
			}
			// Deduplication: each path appears at most once
			seen := make(map[string]bool)
			for _, p := range got {
				if seen[p] {
					t.Errorf("parseRelevantFiles returned duplicate path %q", p)
				}
				seen[p] = true
			}
		})
	}
}

func TestLooksLikeFilePath(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"pkg/pipeline/pipeline_targeted.go", true},
		{"go.mod", true},
		{"go.sum", true},
		{"CLAUDE.md", true},
		{"cmd/hive/main.go", true},
		{"", false},
		{"https://example.com/foo.go", false},
		{"//some/comment", false},
		{"v1.2.3", false},
		{"1.5", false},
		{"noextension", false},
		{"Makefile", false}, // no extension
	}
	for _, tt := range tests {
		got := looksLikeFilePath(tt.s)
		if got != tt.want {
			t.Errorf("looksLikeFilePath(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestLangMarkerFile(t *testing.T) {
	tests := []struct {
		lang string
		want string
	}{
		{"go", "go.mod"},
		{"typescript", "package.json"},
		{"javascript", "package.json"},
		{"python", "pyproject.toml"},
		{"rust", "Cargo.toml"},
		{"csharp", "*.csproj"},
		{"unknown", "go.mod"},
	}
	for _, tt := range tests {
		got := langMarkerFile(tt.lang)
		if got != tt.want {
			t.Errorf("langMarkerFile(%q) = %q, want %q", tt.lang, got, tt.want)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
// buildFileListing
// ════════════════════════════════════════════════════════════════════════

func TestBuildFileListing_Header(t *testing.T) {
	got := buildFileListing(map[string]string{})
	if !strings.HasPrefix(got, "Files:\n") {
		t.Errorf("buildFileListing should start with 'Files:\\n', got %q", got)
	}
}

func TestBuildFileListing_ContainsFilesAndLineCounts(t *testing.T) {
	files := map[string]string{
		"main.go":     "package main\n\nfunc main() {}\n",
		"lib/util.go": "package lib\n",
	}
	got := buildFileListing(files)
	if !strings.Contains(got, "main.go") {
		t.Error("buildFileListing should contain main.go")
	}
	if !strings.Contains(got, "lib/util.go") {
		t.Error("buildFileListing should contain lib/util.go")
	}
	// "package main\n\nfunc main() {}\n" has 3 newlines → 4 lines.
	if !strings.Contains(got, "(4 lines)") {
		t.Errorf("buildFileListing line count wrong for main.go content, got %q", got)
	}
	// "package lib\n" has 1 newline → 2 lines.
	if !strings.Contains(got, "(2 lines)") {
		t.Errorf("buildFileListing line count wrong for lib/util.go content, got %q", got)
	}
}

// ════════════════════════════════════════════════════════════════════════
// extractKeyFiles
// ════════════════════════════════════════════════════════════════════════

func TestExtractKeyFiles_IncludesKnownFiles(t *testing.T) {
	files := map[string]string{
		"CLAUDE.md":       "# Hive\nproject docs",
		"README.md":       "# readme",
		"go.mod":          "module test",
		"main.go":         "package main",
	}
	got := extractKeyFiles(files)
	if !strings.Contains(got, "--- CLAUDE.md ---") {
		t.Error("extractKeyFiles should include CLAUDE.md")
	}
	if !strings.Contains(got, "--- README.md ---") {
		t.Error("extractKeyFiles should include README.md")
	}
	// go.mod and main.go are not key files.
	if strings.Contains(got, "--- go.mod ---") {
		t.Error("extractKeyFiles should not include go.mod")
	}
	if strings.Contains(got, "--- main.go ---") {
		t.Error("extractKeyFiles should not include main.go")
	}
}

func TestExtractKeyFiles_TruncatesLongFiles(t *testing.T) {
	// 100 lines — exceeds the 60-line cap.
	src := make([]string, 100)
	for i := range src {
		src[i] = "line"
	}
	files := map[string]string{
		"CLAUDE.md": strings.Join(src, "\n"),
	}
	got := extractKeyFiles(files)
	if !strings.Contains(got, "[truncated: 40 lines omitted]") {
		t.Errorf("extractKeyFiles should truncate long files at 60 lines, got %q", got)
	}
}

func TestExtractKeyFiles_EmptyWhenNoKeyFiles(t *testing.T) {
	files := map[string]string{
		"main.go": "package main",
		"go.mod":  "module test",
	}
	got := extractKeyFiles(files)
	if got != "" {
		t.Errorf("extractKeyFiles with no key files = %q, want empty", got)
	}
}

// ════════════════════════════════════════════════════════════════════════
// buildRelevantFileContext
// ════════════════════════════════════════════════════════════════════════

func TestBuildRelevantFileContext_IncludesRequestedFiles(t *testing.T) {
	files := map[string]string{
		"pkg/foo.go": "package foo\nfunc Foo() {}",
		"pkg/bar.go": "package bar\nfunc Bar() {}",
		"pkg/baz.go": "package baz",
	}
	paths := []string{"pkg/foo.go", "pkg/bar.go"}
	got := buildRelevantFileContext(files, paths)
	if !strings.Contains(got, "--- FILE: pkg/foo.go ---") {
		t.Error("buildRelevantFileContext should include foo.go header")
	}
	if !strings.Contains(got, "func Foo()") {
		t.Error("buildRelevantFileContext should include foo.go content")
	}
	if !strings.Contains(got, "--- FILE: pkg/bar.go ---") {
		t.Error("buildRelevantFileContext should include bar.go header")
	}
	if strings.Contains(got, "pkg/baz.go") {
		t.Error("buildRelevantFileContext should not include baz.go (not in relevantPaths)")
	}
}

func TestBuildRelevantFileContext_MissingPathSkipped(t *testing.T) {
	files := map[string]string{
		"pkg/foo.go": "package foo",
	}
	paths := []string{"pkg/foo.go", "pkg/missing.go"}
	got := buildRelevantFileContext(files, paths)
	if !strings.Contains(got, "pkg/foo.go") {
		t.Error("buildRelevantFileContext should include existing file")
	}
	if strings.Contains(got, "pkg/missing.go") {
		t.Error("buildRelevantFileContext should skip missing file")
	}
}

func TestBuildRelevantFileContext_EmptyPaths(t *testing.T) {
	files := map[string]string{"pkg/foo.go": "package foo"}
	got := buildRelevantFileContext(files, nil)
	if got != "" {
		t.Errorf("buildRelevantFileContext with no paths = %q, want empty", got)
	}
}

// ════════════════════════════════════════════════════════════════════════
// formatTaskList
// ════════════════════════════════════════════════════════════════════════

func TestFormatTaskList_Empty(t *testing.T) {
	if got := formatTaskList(nil); got != "" {
		t.Errorf("formatTaskList(nil) = %q, want empty", got)
	}
	if got := formatTaskList([]work.Task{}); got != "" {
		t.Errorf("formatTaskList([]) = %q, want empty", got)
	}
}

func TestFormatTaskList_TitleOnly(t *testing.T) {
	tasks := []work.Task{{Title: "Fix build phase"}}
	got := formatTaskList(tasks)
	if !strings.Contains(got, "- Fix build phase") {
		t.Errorf("formatTaskList missing task title, got %q", got)
	}
	// No description — colon separator must not appear.
	if strings.Contains(got, ": ") {
		t.Errorf("formatTaskList should not add colon when description is empty, got %q", got)
	}
}

func TestFormatTaskList_WithDescription(t *testing.T) {
	tasks := []work.Task{
		{Title: "Add retry logic", Description: "retry on transient errors"},
	}
	got := formatTaskList(tasks)
	if !strings.Contains(got, "- Add retry logic: retry on transient errors") {
		t.Errorf("formatTaskList missing description, got %q", got)
	}
}

func TestFormatTaskList_MultipleTasks(t *testing.T) {
	tasks := []work.Task{
		{Title: "Task A"},
		{Title: "Task B", Description: "desc B"},
	}
	got := formatTaskList(tasks)
	if !strings.Contains(got, "- Task A") {
		t.Error("formatTaskList missing Task A")
	}
	if !strings.Contains(got, "- Task B: desc B") {
		t.Error("formatTaskList missing Task B with description")
	}
}

// ════════════════════════════════════════════════════════════════════════
// detectApproval
// ════════════════════════════════════════════════════════════════════════

func TestDetectApproval(t *testing.T) {
	tests := []struct {
		review string
		want   bool
	}{
		{"Looks good.\n\nAPPROVED", true},
		{"Some issues.\n\nCHANGES NEEDED", false},
		{"Missing error handling.\n\nCHANGES REQUIRED", false},
		{"This is wrong.\n\nREJECT", false},
		// Both present — changes take precedence.
		{"APPROVED but also CHANGES NEEDED", false},
		{"No verdict at all", false},
		// Case-insensitive via ToUpper in implementation.
		{"approved", true},
		{"", false},
	}
	for _, tt := range tests {
		got := detectApproval(tt.review)
		if got != tt.want {
			t.Errorf("detectApproval(%q) = %v, want %v", tt.review, got, tt.want)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
// detectLanguage
// ════════════════════════════════════════════════════════════════════════

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		files map[string]string
		want  string
	}{
		// Marker-file priority.
		{map[string]string{"go.mod": "module test"}, "go"},
		{map[string]string{"package.json": "{}"}, "typescript"},
		{map[string]string{"Cargo.toml": "[package]"}, "rust"},
		{map[string]string{"requirements.txt": "flask"}, "python"},
		{map[string]string{"setup.py": "from setuptools import setup"}, "python"},
		{map[string]string{"pyproject.toml": "[project]"}, "python"},
		// Extension fallback.
		{map[string]string{"src/main.go": "package main"}, "go"},
		{map[string]string{"src/app.ts": "export {}"}, "typescript"},
		{map[string]string{"src/app.rs": "fn main() {}"}, "rust"},
		{map[string]string{"src/app.py": "print()"}, "python"},
		{map[string]string{"src/app.cs": "class Foo {}"}, "csharp"},
		// Empty map defaults to go.
		{map[string]string{}, "go"},
	}
	for _, tt := range tests {
		got := detectLanguage(tt.files)
		if got != tt.want {
			t.Errorf("detectLanguage(%v) = %q, want %q", tt.files, got, tt.want)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
// sanitizeGoMod
// ════════════════════════════════════════════════════════════════════════

func TestSanitizeGoMod_NoGoMod(t *testing.T) {
	files := map[string]string{"main.go": "package main"}
	sanitizeGoMod(files) // must not panic
}

func TestSanitizeGoMod_CleanFile(t *testing.T) {
	files := map[string]string{
		"go.mod": "module github.com/foo/bar\n\ngo 1.21\n",
	}
	before := files["go.mod"]
	sanitizeGoMod(files)
	if files["go.mod"] != before {
		t.Errorf("sanitizeGoMod modified a clean go.mod: %q", files["go.mod"])
	}
}

func TestSanitizeGoMod_SplitModulePath(t *testing.T) {
	// Simulates LLM corruption: module path split across lines with embedded quotes.
	files := map[string]string{
		"go.mod": "module \"github.com/\nfoo/bar\"\n\ngo 1.21\n",
	}
	sanitizeGoMod(files)
	result := files["go.mod"]

	// The rejoined module directive must be on a single line without quotes.
	found := false
	for _, l := range strings.Split(result, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "module ") {
			if strings.Contains(l, "\"") {
				t.Errorf("module directive still has quotes: %q", l)
			}
			found = true
		}
	}
	if !found {
		t.Errorf("sanitizeGoMod result has no module directive: %q", result)
	}
	// The original corrupted split must be gone.
	if strings.Contains(result, "module \"github.com/\nfoo") {
		t.Errorf("module path still split across lines: %q", result)
	}
}

// ════════════════════════════════════════════════════════════════════════
// stripMarkdownFences
// ════════════════════════════════════════════════════════════════════════

func TestStripMarkdownFences(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no fences",
			input: "package main\n\nfunc main() {}",
			want:  "package main\n\nfunc main() {}",
		},
		{
			name:  "go fence",
			input: "```go\npackage main\n\nfunc main() {}\n```",
			want:  "package main\n\nfunc main() {}",
		},
		{
			name:  "plain fence",
			input: "```\npackage main\n```",
			want:  "package main",
		},
		{
			name:  "fence with trailing prose",
			input: "```go\npackage main\n```\nThis is the main package.",
			want:  "package main",
		},
		{
			name:  "empty content",
			input: "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripMarkdownFences(tt.input)
			if got != tt.want {
				t.Errorf("stripMarkdownFences() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ════════════════════════════════════════════════════════════════════════
// extractName
// ════════════════════════════════════════════════════════════════════════

func TestExtractName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// NAME: kebab-case parsed correctly.
		{"FEASIBLE: yes\nNAME: task-manager\nOther stuff", "task-manager"},
		// Case-insensitive prefix match.
		{"name: kanban-board\n", "kanban-board"},
		// Leading/trailing whitespace trimmed.
		{"NAME:  social-graph  \n", "social-graph"},
		// Spaces in name become hyphens.
		{"NAME: social graph\n", "social-graph"},
		// Digits allowed.
		{"NAME: product2-0\n", "product2-0"},
		// Non-alphanumeric/non-hyphen chars stripped.
		{"NAME: task_manager!\n", "taskmanager"},
		// No NAME line → fallback to "product".
		{"FEASIBLE: yes\nNo name here", "product"},
		// Empty input → fallback.
		{"", "product"},
		// NAME with no valid characters after cleaning → fallback.
		{"NAME: !!!\n", "product"},
	}
	for _, tt := range tests {
		got := extractName(tt.input)
		if got != tt.want {
			t.Errorf("extractName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
// lastNonEmptyLine
// ════════════════════════════════════════════════════════════════════════

func TestLastNonEmptyLine(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Normal multi-line string.
		{"line1\nline2\nline3", "line3"},
		// Trailing blank lines skipped.
		{"line1\nline2\n\n", "line2"},
		// Warning lines before URL (typical gh pr create output).
		{"Warning: something\nhttps://github.com/foo/bar/pull/42\n", "https://github.com/foo/bar/pull/42"},
		// Single non-empty line.
		{"only-line", "only-line"},
		// Leading/trailing whitespace on last line trimmed.
		{"  line with spaces  \n", "line with spaces"},
		// Empty input returns empty string.
		{"", ""},
	}
	for _, tt := range tests {
		got := lastNonEmptyLine(tt.input)
		if got != tt.want {
			t.Errorf("lastNonEmptyLine(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
// isTransientGHError
// ════════════════════════════════════════════════════════════════════════

func TestIsTransientGHError(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		// Transient: 502 substring.
		{"error: 502 response from GitHub", true},
		// Transient: 504 substring.
		{"received 504 status", true},
		// Transient: "Gateway Timeout" phrase.
		{"Gateway Timeout occurred", true},
		// Transient: "Bad Gateway" phrase.
		{"Bad Gateway response", true},
		// Transient: ETIMEDOUT.
		{"ETIMEDOUT: connection timed out", true},
		// Transient: ECONNRESET.
		{"ECONNRESET while pushing", true},
		// Non-transient errors.
		{"404 Not Found", false},
		{"401 Unauthorized", false},
		{"Permission denied", false},
		{"", false},
	}
	for _, tt := range tests {
		got := isTransientGHError(tt.input)
		if got != tt.want {
			t.Errorf("isTransientGHError(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
// truncate
// ════════════════════════════════════════════════════════════════════════

func TestTruncate(t *testing.T) {
	tests := []struct {
		input string
		n     int
		want  string
	}{
		// Short string — no truncation.
		{"hello", 10, "hello"},
		// Exactly at limit — no truncation.
		{"hello", 5, "hello"},
		// Truncate: last 3 chars replaced by ellipsis.
		{"hello world", 8, "hello..."},
		// Empty string.
		{"", 5, ""},
	}
	for _, tt := range tests {
		got := truncate(tt.input, tt.n)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.n, got, tt.want)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════
// writeCodeAction
// ════════════════════════════════════════════════════════════════════════

func TestWriteCodeAction(t *testing.T) {
	tests := []struct {
		lang string
		want string
	}{
		{"go", ActionWriteCode + ":go"},
		{"typescript", ActionWriteCode + ":typescript"},
		{"python", ActionWriteCode + ":python"},
		{"rust", ActionWriteCode + ":rust"},
	}
	for _, tt := range tests {
		got := writeCodeAction(tt.lang)
		if got != tt.want {
			t.Errorf("writeCodeAction(%q) = %q, want %q", tt.lang, got, tt.want)
		}
	}
}
