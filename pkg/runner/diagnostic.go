package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// PhaseEvent is a diagnostic event emitted by every pipeline phase.
// Not just failures — records tokens, duration, and outcome for every run.
// The Observer uses this to detect inefficiency, not just errors.
// Cost is a derivative of tokens × model rate — track the raw data.
type PhaseEvent struct {
	Phase        string  `json:"phase"`
	Outcome      string  `json:"outcome,omitempty"`       // "success", "failure", "revise", "skip"
	Error        string  `json:"error,omitempty"`
	Preview      string  `json:"preview,omitempty"`
	Model        string  `json:"model,omitempty"`          // which model was used
	InputTokens  int     `json:"input_tokens,omitempty"`
	OutputTokens int     `json:"output_tokens,omitempty"`
	DurationSecs float64 `json:"duration_secs,omitempty"`  // wall clock time
	CostUSD      float64 `json:"cost_usd,omitempty"`       // derived, kept for convenience
	Timestamp    string  `json:"timestamp"`
}

// appendDiagnostic appends a PhaseEvent as a JSON line to
// {hiveDir}/loop/diagnostics.jsonl.  It sets Timestamp if unset.
func appendDiagnostic(hiveDir string, e PhaseEvent) error {
	if e.Timestamp == "" {
		e.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshal diagnostic: %w", err)
	}
	path := filepath.Join(hiveDir, "loop", "diagnostics.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open diagnostics.jsonl: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// countDiagnostics counts newline-terminated lines in
// {hiveDir}/loop/diagnostics.jsonl. Returns 0 if the file doesn't exist.
func countDiagnostics(hiveDir string) int {
	path := filepath.Join(hiveDir, "loop", "diagnostics.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	count := 0
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}
	return count
}

// appendDiagnostic appends a PhaseEvent to loop/diagnostics.jsonl.
// Silently skips if HiveDir is empty.
func (r *Runner) appendDiagnostic(e PhaseEvent) {
	if r.cfg.HiveDir == "" {
		return
	}
	if err := appendDiagnostic(r.cfg.HiveDir, e); err != nil {
		log.Printf("[runner] appendDiagnostic: %v", err)
	}
}
