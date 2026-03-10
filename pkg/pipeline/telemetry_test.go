package pipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPhaseTimingEntryRoundTrip(t *testing.T) {
	original := PhaseTimingEntry{
		Phase:    "build",
		Duration: 1500 * time.Millisecond,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// Verify JSON contains duration_ms as integer.
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal raw: %v", err)
	}
	if raw["duration_ms"] != float64(1500) {
		t.Errorf("duration_ms = %v, want 1500", raw["duration_ms"])
	}

	// Round-trip back.
	var decoded PhaseTimingEntry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Phase != original.Phase {
		t.Errorf("phase = %q, want %q", decoded.Phase, original.Phase)
	}
	if decoded.Duration != original.Duration {
		t.Errorf("duration = %v, want %v", decoded.Duration, original.Duration)
	}
}

func TestReadTelemetryNoDir(t *testing.T) {
	// Non-existent directory returns empty slice, no error.
	results, err := ReadTelemetry(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("got %d results, want 0", len(results))
	}
}

func TestReadTelemetryWithFiles(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".hive", "telemetry")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write two telemetry files.
	for _, name := range []string{"run1.json", "run2.json"} {
		result := PipelineResult{
			Mode:             "targeted",
			InputDescription: "test change for " + name,
			StartedAt:        time.Now(),
			EndedAt:          time.Now(),
			PhaseTimings: []PhaseTimingEntry{
				{Phase: "build", Duration: 2 * time.Second},
			},
		}
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Write a non-JSON file that should be skipped.
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	results, err := ReadTelemetry(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	// Verify phase timing round-tripped correctly.
	for _, r := range results {
		if len(r.PhaseTimings) != 1 {
			t.Errorf("phase timings: got %d, want 1", len(r.PhaseTimings))
			continue
		}
		if r.PhaseTimings[0].Duration != 2*time.Second {
			t.Errorf("duration = %v, want 2s", r.PhaseTimings[0].Duration)
		}
	}
}

func TestParseSelfImproveRecommendation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantDesc string
		wantSkip string
	}{
		{
			name: "bare JSON",
			input: `{"description":"fix guardian cost","files_to_change":["pkg/pipeline/pipeline.go"],"expected_impact":"reduce cost by 50%","priority":"high","skip_reason":""}`,
			wantDesc: "fix guardian cost",
		},
		{
			name: "markdown code block",
			input: "Here's my analysis:\n```json\n{\"description\":\"optimize reviewer\",\"files_to_change\":[],\"expected_impact\":\"save time\",\"priority\":\"medium\",\"skip_reason\":\"\"}\n```\n",
			wantDesc: "optimize reviewer",
		},
		{
			name: "skip reason",
			input: `{"description":"","files_to_change":[],"expected_impact":"","priority":"low","skip_reason":"nothing worth fixing"}`,
			wantSkip: "nothing worth fixing",
		},
		{
			name: "no JSON",
			input: "I think everything looks fine.",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec, err := parseSelfImproveRecommendation(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rec.Description != tt.wantDesc {
				t.Errorf("description = %q, want %q", rec.Description, tt.wantDesc)
			}
			if rec.SkipReason != tt.wantSkip {
				t.Errorf("skip_reason = %q, want %q", rec.SkipReason, tt.wantSkip)
			}
		})
	}
}

func TestExtractJSONBlock(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bare JSON",
			input: `{"key":"value"}`,
			want:  `{"key":"value"}`,
		},
		{
			name:  "markdown json block",
			input: "```json\n{\"key\":\"value\"}\n```",
			want:  `{"key":"value"}`,
		},
		{
			name:  "markdown plain block",
			input: "```\n{\"key\":\"value\"}\n```",
			want:  `{"key":"value"}`,
		},
		{
			name:  "surrounded by text",
			input: "Here is the result: {\"key\":\"value\"} and more text",
			want:  `{"key":"value"}`,
		},
		{
			name:  "nested braces",
			input: `{"outer":{"inner":"val"}}`,
			want:  `{"outer":{"inner":"val"}}`,
		},
		{
			name:  "no JSON",
			input: "just plain text",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractJSONBlock(tt.input)
			if got != tt.want {
				t.Errorf("extractJSONBlock = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSummarizeTelemetryEmpty(t *testing.T) {
	summary := summarizeTelemetry(nil)
	if summary != "No telemetry data available (first run)." {
		t.Errorf("unexpected summary for empty data: %q", summary)
	}
}

func TestSummarizeTelemetryWithData(t *testing.T) {
	results := []PipelineResult{
		{
			Mode:             "targeted",
			InputDescription: "fix bug",
			StartedAt:        time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
			PhaseTimings: []PhaseTimingEntry{
				{Phase: "build", Duration: 5 * time.Second},
			},
			TokenUsage: []RoleTokenUsage{
				{Role: "builder", Model: "opus", TotalTokens: 1000, CostUSD: 0.05},
			},
			GuardianAlerts: []string{"ALERT: test alert"},
			ReviewSignals:  []string{"APPROVED"},
		},
	}

	summary := summarizeTelemetry(results)
	// Verify key content is present.
	for _, want := range []string{"1 past pipeline run", "fix bug", "build", "builder", "ALERT: test alert", "APPROVED"} {
		if !contains(summary, want) {
			t.Errorf("summary missing %q", want)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
