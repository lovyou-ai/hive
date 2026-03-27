package runner

import (
	"strings"
	"testing"
)

func TestBuildPart2Instruction(t *testing.T) {
	cases := []struct {
		name          string
		spaceSlug     string
		apiKey        string
		wantSkip      bool
		wantCurl      bool
		wantKeyInBody bool
		wantSlugInURL bool
	}{
		{
			name:      "empty apiKey returns skip text, no curl",
			spaceSlug: "hive",
			apiKey:    "",
			wantSkip:  true,
			wantCurl:  false,
		},
		{
			name:          "set apiKey returns curl with key and slug embedded",
			spaceSlug:     "hive",
			apiKey:        "lv_testkey",
			wantSkip:      false,
			wantCurl:      true,
			wantKeyInBody: true,
			wantSlugInURL: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildPart2Instruction(tc.spaceSlug, tc.apiKey)

			if tc.wantSkip && !strings.Contains(got, "Skipped") {
				t.Errorf("expected skip message, got: %q", got)
			}
			if !tc.wantSkip && strings.Contains(got, "Skipped") {
				t.Errorf("unexpected skip message, got: %q", got)
			}
			if tc.wantCurl && !strings.Contains(got, "Authorization: Bearer") {
				t.Errorf("expected curl auth command, got: %q", got)
			}
			if !tc.wantCurl && strings.Contains(got, "Authorization: Bearer") {
				t.Errorf("unexpected curl auth command, got: %q", got)
			}
			if tc.wantKeyInBody && !strings.Contains(got, tc.apiKey) {
				t.Errorf("expected API key %q in output, got: %q", tc.apiKey, got)
			}
			if tc.wantSlugInURL && !strings.Contains(got, tc.spaceSlug) {
				t.Errorf("expected slug %q in output, got: %q", tc.spaceSlug, got)
			}
		})
	}
}

func TestBuildOutputInstruction(t *testing.T) {
	cases := []struct {
		name          string
		spaceSlug     string
		apiKey        string
		wantTextFmt   bool
		wantCurl      bool
		wantKeyInBody bool
		wantSlugInURL bool
	}{
		{
			name:        "empty apiKey returns text task format, no curl",
			spaceSlug:   "hive",
			apiKey:      "",
			wantTextFmt: true,
			wantCurl:    false,
		},
		{
			name:          "set apiKey returns curl with key and slug, no text format",
			spaceSlug:     "hive",
			apiKey:        "lv_testkey",
			wantTextFmt:   false,
			wantCurl:      true,
			wantKeyInBody: true,
			wantSlugInURL: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildOutputInstruction(tc.spaceSlug, tc.apiKey)

			if tc.wantTextFmt && !strings.Contains(got, "TASK_TITLE:") {
				t.Errorf("expected text task format, got: %q", got)
			}
			if !tc.wantTextFmt && strings.Contains(got, "TASK_TITLE:") {
				t.Errorf("unexpected text task format, got: %q", got)
			}
			if tc.wantCurl && !strings.Contains(got, "Authorization: Bearer") {
				t.Errorf("expected curl auth command, got: %q", got)
			}
			if !tc.wantCurl && strings.Contains(got, "Authorization: Bearer") {
				t.Errorf("unexpected curl auth command, got: %q", got)
			}
			if tc.wantKeyInBody && !strings.Contains(got, tc.apiKey) {
				t.Errorf("expected API key %q in output, got: %q", tc.apiKey, got)
			}
			if tc.wantSlugInURL && !strings.Contains(got, tc.spaceSlug) {
				t.Errorf("expected slug %q in output, got: %q", tc.spaceSlug, got)
			}
		})
	}
}

func TestBuildObserverInstruction(t *testing.T) {
	cases := []struct {
		name      string
		repoPath  string
		spaceSlug string
		apiKey    string
		wantParts []string
		// Part 2 skip or curl
		wantSkipInPart2 bool
		wantCurlInPart2 bool
		// Output section: text format or curl
		wantTextFmtInOutput bool
		wantCurlInOutput    bool
	}{
		{
			name:      "empty apiKey: skip in part2, text format in output",
			repoPath:  "/repo",
			spaceSlug: "hive",
			apiKey:    "",
			wantParts: []string{
				"You are the Observer",
				"Part 1: Product Audit",
				"Part 2: Graph Integrity Audit",
				"/repo",
				"hive",
			},
			wantSkipInPart2:     true,
			wantCurlInPart2:     false,
			wantTextFmtInOutput: true,
			wantCurlInOutput:    false,
		},
		{
			name:      "set apiKey: curl in part2 and output with key+slug",
			repoPath:  "/repo",
			spaceSlug: "hive",
			apiKey:    "lv_testkey",
			wantParts: []string{
				"You are the Observer",
				"Part 1: Product Audit",
				"Part 2: Graph Integrity Audit",
				"/repo",
				"hive",
				"lv_testkey",
			},
			wantSkipInPart2:     false,
			wantCurlInPart2:     true,
			wantTextFmtInOutput: false,
			wantCurlInOutput:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildObserverInstruction(tc.repoPath, tc.spaceSlug, tc.apiKey)

			for _, part := range tc.wantParts {
				if !strings.Contains(got, part) {
					t.Errorf("expected %q in output, got: %q", part, got)
				}
			}

			// Part 2 integrity section appears once; check skip vs curl.
			if tc.wantSkipInPart2 && !strings.Contains(got, "Skipped") {
				t.Errorf("expected skip message in part 2, got: %q", got)
			}
			if !tc.wantSkipInPart2 && strings.Contains(got, "Skipped") {
				t.Errorf("unexpected skip message, got: %q", got)
			}

			// Count Authorization occurrences: part2 curl + output curl.
			authCount := strings.Count(got, "Authorization: Bearer")
			if tc.wantCurlInPart2 && tc.wantCurlInOutput && authCount < 2 {
				t.Errorf("expected at least 2 curl auth headers (part2 + output), got %d in: %q", authCount, got)
			}
			if !tc.wantCurlInPart2 && !tc.wantCurlInOutput && authCount > 0 {
				t.Errorf("expected no curl auth headers, got %d in: %q", authCount, got)
			}

			if tc.wantTextFmtInOutput && !strings.Contains(got, "TASK_TITLE:") {
				t.Errorf("expected TASK_TITLE text format in output section, got: %q", got)
			}
			if !tc.wantTextFmtInOutput && strings.Contains(got, "TASK_TITLE:") {
				t.Errorf("unexpected TASK_TITLE text format, got: %q", got)
			}
		})
	}
}
