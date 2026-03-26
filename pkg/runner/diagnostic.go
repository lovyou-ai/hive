package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// PhaseEvent is a diagnostic event emitted by a runner phase on error or failure.
type PhaseEvent struct {
	Phase     string  `json:"phase"`
	Error     string  `json:"error,omitempty"`
	CostUSD   float64 `json:"cost_usd"`
	Timestamp string  `json:"timestamp"`
}

// appendDiagnostic appends a PhaseEvent to loop/diagnostics.jsonl.
// Silently skips if HiveDir is empty.
func (r *Runner) appendDiagnostic(e PhaseEvent) {
	if r.cfg.HiveDir == "" {
		return
	}
	e.Timestamp = time.Now().UTC().Format(time.RFC3339)
	data, err := json.Marshal(e)
	if err != nil {
		log.Printf("[runner] marshal diagnostic: %v", err)
		return
	}
	path := filepath.Join(r.cfg.HiveDir, "loop", "diagnostics.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[runner] open diagnostics.jsonl: %v", err)
		return
	}
	defer f.Close()
	_, _ = fmt.Fprintf(f, "%s\n", data)
}
