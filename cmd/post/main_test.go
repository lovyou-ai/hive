package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"standard format", "# Build: Fix: foo bar\n\nmore content", "Fix: foo bar"},
		{"heading only", "# Some Title\nbody", "Some Title"},
		{"leading blank lines", "\n\n# Build: Hello\n", "Hello"},
		{"empty input", "", ""},
		{"whitespace only", "   \n  \n", ""},
		{"multi-hash", "## Build: Nested", "Nested"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildTitle([]byte(tt.input))
			if got != tt.want {
				t.Errorf("buildTitle(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestPostCreatesNode verifies that the post() function sends op=express with
// kind=post, title, and body to /app/hive/op.
func TestPostCreatesNode(t *testing.T) {
	var received map[string]string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/app/hive/op" {
			http.NotFound(w, r)
			return
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"node":{"id":"test-id"}}`))
	}))
	defer srv.Close()

	err := post("lv_testkey", srv.URL, "Fix: some bug", "## What Was Built\nFixed the bug.")
	if err != nil {
		t.Fatalf("post() error: %v", err)
	}

	if received["op"] != "express" {
		t.Errorf("op = %q, want %q", received["op"], "express")
	}
	if received["kind"] != "post" {
		t.Errorf("kind = %q, want %q", received["kind"], "post")
	}
	if received["title"] != "Fix: some bug" {
		t.Errorf("title = %q, want %q", received["title"], "Fix: some bug")
	}
	if received["body"] == "" {
		t.Error("body is empty, want non-empty build summary")
	}
}

// TestSyncClaimsWritesFile verifies that syncClaims fetches claims from the API
// and writes them as markdown to the given output path.
func TestSyncClaimsWritesFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/app/hive/knowledge" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"claims": []map[string]any{
				{
					"title":      "Absence is invisible to traversal",
					"body":       "The Scout traverses what exists. Tests don't exist, so the Scout never encounters them.",
					"state":      "claimed",
					"author":     "Reflector",
					"created_at": "2026-03-01T00:00:00Z",
				},
				{
					"title":      "Ship what you build",
					"body":       "Every build iteration should deploy.",
					"state":      "verified",
					"author":     "Builder",
					"created_at": "2026-03-02T00:00:00Z",
				},
			},
		})
	}))
	defer srv.Close()

	outPath := filepath.Join(t.TempDir(), "claims.md")
	if err := syncClaims("lv_testkey", srv.URL, outPath); err != nil {
		t.Fatalf("syncClaims() error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("claims.md not written: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "# Knowledge Claims") {
		t.Error("missing heading")
	}
	if !strings.Contains(content, "Absence is invisible to traversal") {
		t.Error("missing first claim title")
	}
	if !strings.Contains(content, "Ship what you build") {
		t.Error("missing second claim title")
	}
	if !strings.Contains(content, "Every build iteration should deploy") {
		t.Error("missing second claim body")
	}
	if !strings.Contains(content, "**State:** verified") {
		t.Error("missing state for verified claim")
	}
}

// TestSyncClaimsEmptyDoesNotWrite verifies that syncClaims does not write a
// file when the API returns zero claims.
func TestSyncClaimsEmptyDoesNotWrite(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"claims": []any{}})
	}))
	defer srv.Close()

	outPath := filepath.Join(t.TempDir(), "claims.md")
	if err := syncClaims("lv_testkey", srv.URL, outPath); err != nil {
		t.Fatalf("syncClaims() error: %v", err)
	}

	if _, err := os.Stat(outPath); err == nil {
		t.Error("claims.md should not be written when there are no claims")
	}
}

// TestBuildTitleExtractedOnPost verifies that buildTitle + post produces a
// feed node whose title comes from build.md (not just "Iteration N").
func TestBuildTitleExtractedOnPost(t *testing.T) {
	buildMD := []byte("# Build: Fix: Observer AllowedTools missing knowledge.search\n\n## What Was Built\nFixed it.")

	var receivedTitle string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &payload)
		if payload["op"] == "express" {
			receivedTitle = payload["title"]
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"node":{"id":"test-id"}}`))
	}))
	defer srv.Close()

	title := buildTitle(buildMD)
	if title == "" {
		t.Fatal("buildTitle returned empty for valid build.md")
	}

	if err := post("lv_testkey", srv.URL, title, string(buildMD)); err != nil {
		t.Fatalf("post() error: %v", err)
	}

	want := "Fix: Observer AllowedTools missing knowledge.search"
	if receivedTitle != want {
		t.Errorf("post title = %q, want %q", receivedTitle, want)
	}
}
