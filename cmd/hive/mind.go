package main

import (
	"bufio"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store/pgstore"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/mind"
	"github.com/lovyou-ai/hive/pkg/pipeline"
)

// runMind starts an interactive chat session with the hive's mind.
func runMind(ctx context.Context, dsn, repoPath, model string) error {
	if dsn == "" {
		return fmt.Errorf("mind requires --store or DATABASE_URL (needs persistent state)")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer pool.Close()

	s, err := pgstore.NewPostgresStoreFromPool(ctx, pool)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	defer s.Close()

	// Register event types so the store can deserialize mind events.
	mind.RegisterEventTypes()
	pipeline.RegisterEventTypes()

	factory := event.NewEventFactory(event.DefaultRegistry())

	// Derive signer from "mind" identity.
	// TODO: replace with Google auth when available.
	h := sha256.Sum256([]byte("signer:mind"))
	priv := ed25519.NewKeyFromSeed(h[:])
	signer := &mindSigner{key: priv}

	mindStore := mind.NewMindStore(s, factory, signer)

	// Resolve human actor ID for MCP config.
	// TODO: this should come from Google auth, not CLI derivation.
	humanID, err := resolveHumanID("Matt")
	if err != nil {
		return fmt.Errorf("resolve human ID: %w", err)
	}

	// Generate MCP config so the mind can query the event graph.
	mcpConfigPath, cleanup, err := writeMCPConfig(dsn, humanID, repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: MCP config failed: %v (mind will run without event graph access)\n", err)
		mcpConfigPath = ""
	} else {
		defer cleanup()
	}

	// Telemetry summary.
	telemetrySummary := loadTelemetrySummary(repoPath)

	// Doc paths the mind should know about.
	absRepo, _ := filepath.Abs(repoPath)
	docPaths := discoverDocPaths(absRepo)

	provider, err := mind.CreateProvider(model, mcpConfigPath)
	if err != nil {
		return fmt.Errorf("mind provider: %w", err)
	}

	m := mind.New(mind.Config{
		Provider:         provider,
		Store:            mindStore,
		RepoPath:         absRepo,
		TelemetrySummary: telemetrySummary,
		DocPaths:         docPaths,
	})

	fmt.Fprintf(os.Stderr, "Mind loaded: %d lines of context\n", m.ContextLines())
	if mcpConfigPath != "" {
		fmt.Fprintf(os.Stderr, "MCP: event graph connected\n")
	}
	fmt.Fprintf(os.Stderr, "Docs: %d sources available\n", len(docPaths))
	fmt.Fprintf(os.Stderr, "Type your message. Empty line or Ctrl+C to exit.\n\n")

	scanner := bufio.NewScanner(os.Stdin)
	// Increase scanner buffer for long inputs.
	scanner.Buffer(make([]byte, 0, 64*1024), 64*1024)

	for {
		fmt.Fprint(os.Stderr, "Matt> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			break
		}

		fmt.Fprintf(os.Stderr, "  ⏳ thinking...\n")
		resp, err := m.Chat(ctx, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "\n%s\n\n", resp)
	}
	return nil
}

// writeMCPConfig creates a temporary MCP config file that points at the
// hive MCP server with the right store DSN and actor IDs.
func writeMCPConfig(dsn, humanID, repoPath string) (string, func(), error) {
	// Find the MCP server binary — build it if needed.
	mcpBin, err := buildMCPServer(repoPath)
	if err != nil {
		return "", nil, fmt.Errorf("build mcp-server: %w", err)
	}

	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"hive": map[string]interface{}{
				"command": mcpBin,
				"args": []string{
					"--store", dsn,
					"--agent-id", humanID, // mind acts as the human
					"--human-id", humanID,
				},
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", nil, err
	}

	tmpFile, err := os.CreateTemp("", "hive-mind-mcp-*.json")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", nil, err
	}
	tmpFile.Close()

	cleanup := func() { os.Remove(tmpFile.Name()) }
	return tmpFile.Name(), cleanup, nil
}

// buildMCPServer builds the MCP server binary and returns its path.
func buildMCPServer(repoPath string) (string, error) {
	outPath := filepath.Join(os.TempDir(), "hive-mcp-server.exe")

	// Skip build if binary exists and is recent.
	if info, err := os.Stat(outPath); err == nil && info.Size() > 0 {
		return outPath, nil
	}

	absRepo, _ := filepath.Abs(repoPath)
	cmd := exec.Command("go", "build", "-o", outPath, "./cmd/mcp-server/")
	cmd.Dir = absRepo
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outPath, nil
}

// resolveHumanID derives the human actor ID using the same method as
// registerHuman in main.go + pgactor's deriveActorID.
// TODO: replace with Google auth.
func resolveHumanID(name string) (string, error) {
	h := sha256.Sum256([]byte("human:" + name))
	priv := ed25519.NewKeyFromSeed(h[:])
	pub := priv.Public().(ed25519.PublicKey)
	// Same derivation as pgactor.deriveActorID: SHA256 of public key, hex first 16 bytes.
	pkHash := sha256.Sum256([]byte(pub))
	return fmt.Sprintf("actor_%s", hex.EncodeToString(pkHash[:16])), nil
}

// discoverDocPaths finds documentation sources the mind should know about.
func discoverDocPaths(repoPath string) []string {
	var paths []string

	// Hive project docs.
	candidates := []string{
		filepath.Join(repoPath, "CLAUDE.md"),
		filepath.Join(repoPath, "docs"),
	}

	// EventGraph docs (sibling repo).
	egPath := filepath.Join(filepath.Dir(repoPath), "eventgraph")
	if info, err := os.Stat(egPath); err == nil && info.IsDir() {
		candidates = append(candidates,
			filepath.Join(egPath, "docs"),
			filepath.Join(egPath, "README.md"),
		)
	}

	// Hive telemetry.
	telemetryDir := filepath.Join(repoPath, ".hive", "telemetry")
	if info, err := os.Stat(telemetryDir); err == nil && info.IsDir() {
		candidates = append(candidates, telemetryDir)
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			paths = append(paths, p)
		}
	}
	return paths
}

// loadTelemetrySummary reads telemetry and produces a summary string.
func loadTelemetrySummary(repoPath string) string {
	results, err := pipeline.ReadTelemetry(repoPath)
	if err != nil || len(results) == 0 {
		return ""
	}

	var sb strings.Builder
	n := len(results)
	start := 0
	if n > 10 {
		start = n - 10
	}
	for _, r := range results[start:] {
		mode := r.Mode
		if mode == "" {
			mode = "unknown"
		}
		status := "completed"
		if r.FailedPhase != "" {
			status = fmt.Sprintf("failed at %s", r.FailedPhase)
		}
		merged := ""
		if r.PRURL != "" {
			merged = fmt.Sprintf(" PR: %s", r.PRURL)
			if r.Merged {
				merged += " (merged)"
			}
		}
		var cost float64
		for _, u := range r.TokenUsage {
			cost += u.CostUSD
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s — %s, $%.2f%s\n",
			r.StartedAt.Format("2006-01-02 15:04"), mode, status, cost, merged))
	}
	sb.WriteString(fmt.Sprintf("\nTotal runs: %d\n", n))
	return sb.String()
}

// mindSigner signs mind events with a deterministic key.
type mindSigner struct {
	key ed25519.PrivateKey
}

func (s *mindSigner) Sign(data []byte) (types.Signature, error) {
	sig := ed25519.Sign(s.key, data)
	return types.NewSignature(sig)
}
