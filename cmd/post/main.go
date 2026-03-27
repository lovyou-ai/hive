// Command post publishes a hive iteration summary to lovyou.ai.
//
// It reads loop/state.md and loop/build.md, then posts the build report
// as a feed entry in the "hive" space on lovyou.ai. If the space doesn't
// exist yet, it creates it.
//
// Configuration via environment variables:
//
//	LOVYOU_API_KEY   — required. Bearer token for lovyou.ai API.
//	LOVYOU_BASE_URL  — optional. Defaults to https://lovyou.ai.
//
// Usage:
//
//	cd /c/src/matt/lovyou3/hive
//	LOVYOU_API_KEY=lv_... go run ./cmd/post/
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	apiKey := os.Getenv("LOVYOU_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "LOVYOU_API_KEY not set, skipping post")
		return
	}

	baseURL := os.Getenv("LOVYOU_BASE_URL")
	if baseURL == "" {
		baseURL = "https://lovyou.ai"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	// Read iteration number from state.md.
	state, err := os.ReadFile("loop/state.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read state.md: %v\n", err)
		os.Exit(1)
	}
	re := regexp.MustCompile(`Iteration (\d+)`)
	m := re.FindSubmatch(state)
	if m == nil {
		fmt.Fprintln(os.Stderr, "could not find iteration number in state.md")
		os.Exit(1)
	}
	iteration := string(m[1])

	// Read build report.
	build, err := os.ReadFile("loop/build.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read build.md: %v\n", err)
		os.Exit(1)
	}

	// Ensure the "hive" space exists.
	if err := ensureSpace(apiKey, baseURL); err != nil {
		fmt.Fprintf(os.Stderr, "ensure space: %v\n", err)
		os.Exit(1)
	}

	// Extract build title from first line of build.md (strip leading "# " or "# Build: ").
	title := buildTitle(build)
	if title == "" {
		title = fmt.Sprintf("Iteration %s", iteration)
	}

	// Post build report as KindDocument. The returned node ID is used as a
	// cause for subsequent claim nodes (Invariant 2: CAUSALITY).
	buildDocID, err := post(apiKey, baseURL, title, string(build))
	if err != nil {
		fmt.Fprintf(os.Stderr, "post: %v\n", err)
		os.Exit(1)
	}

	// causeIDs links subsequent claims back to the build document that generated them.
	var causeIDs []string
	if buildDocID != "" {
		causeIDs = []string{buildDocID}
	}

	// Create task on Board and mark complete.
	if err := createTask(apiKey, baseURL, title, string(build)); err != nil {
		fmt.Fprintf(os.Stderr, "task: %v\n", err)
		// Non-fatal — document post succeeded.
	}

	// Sync loop state to the Mind.
	if err := syncMindState(apiKey, baseURL, string(state)); err != nil {
		fmt.Fprintf(os.Stderr, "sync mind state: %v\n", err)
		// Non-fatal — post succeeded.
	}

	// Sync graph claims to loop/claims.md for MCP knowledge index.
	if err := syncClaims(apiKey, baseURL, "loop/claims.md"); err != nil {
		fmt.Fprintf(os.Stderr, "sync claims: %v\n", err)
		// Non-fatal — post succeeded.
	}

	// Assert the Scout's gap as a KindClaim node so gaps are searchable
	// via knowledge_search and survive scout.md being overwritten next iteration.
	if err := assertScoutGap(apiKey, baseURL, causeIDs); err != nil {
		fmt.Fprintf(os.Stderr, "assert scout gap: %v\n", err)
		// Non-fatal — post succeeded.
	}

	// Assert the critique verdict as a KindClaim node.
	if err := assertCritique(apiKey, baseURL, causeIDs); err != nil {
		fmt.Fprintf(os.Stderr, "assert critique: %v\n", err)
		// Non-fatal — post succeeded.
	}

	// Post the latest reflection entry as a KindDocument node.
	if err := assertLatestReflection(apiKey, baseURL, causeIDs); err != nil {
		fmt.Fprintf(os.Stderr, "assert reflection: %v\n", err)
		// Non-fatal — post succeeded.
	}

	fmt.Printf("posted iteration %s to %s/app/hive/feed\n", iteration, baseURL)
}

func ensureSpace(apiKey, baseURL string) error {
	req, _ := http.NewRequest("GET", baseURL+"/app/hive", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("check space: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil // space exists
	}

	// Create the space.
	body, _ := json.Marshal(map[string]string{
		"name":        "Hive",
		"description": "The hive loop — iteration summaries, agent activity, system state.",
		"kind":        "community",
		"visibility":  "public",
	})
	req, _ = http.NewRequest("POST", baseURL+"/app/new", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("create space: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create space: HTTP %d: %s", resp.StatusCode, b)
	}

	fmt.Println("created hive space")
	return nil
}

func syncMindState(apiKey, baseURL, state string) error {
	payload, _ := json.Marshal(map[string]string{
		"key":   "loop_state",
		"value": state,
	})
	req, _ := http.NewRequest("PUT", baseURL+"/api/mind-state", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, b)
	}
	return nil
}

func createTask(apiKey, baseURL, title, description string) error {
	// Create task (intend op with kind=task).
	payload, _ := json.Marshal(map[string]string{
		"op":          "intend",
		"kind":        "task",
		"title":       title,
		"description": description,
	})
	req, _ := http.NewRequest("POST", baseURL+"/app/hive/op", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("intend: HTTP %d: %s", resp.StatusCode, b)
	}

	// Parse response to get node ID.
	var result struct {
		Node struct {
			ID string `json:"id"`
		} `json:"node"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.Node.ID == "" {
		// Can't complete without ID, but task was created.
		return nil
	}

	// Complete the task.
	completePayload, _ := json.Marshal(map[string]string{
		"op":      "complete",
		"node_id": result.Node.ID,
	})
	req2, _ := http.NewRequest("POST", baseURL+"/app/hive/op", bytes.NewReader(completePayload))
	req2.Header.Set("Authorization", "Bearer "+apiKey)
	req2.Header.Set("Accept", "application/json")
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode >= 400 {
		b, _ := io.ReadAll(resp2.Body)
		return fmt.Errorf("complete: HTTP %d: %s", resp2.StatusCode, b)
	}
	return nil
}

// buildTitle extracts the title from the first line of build.md.
// It strips markdown heading markers and the "Build: " prefix.
// e.g. "# Build: Fix: foo" → "Fix: foo"
func buildTitle(build []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(build))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Strip leading '#' characters and whitespace.
		line = strings.TrimLeft(line, "#")
		line = strings.TrimSpace(line)
		// Strip optional "Build: " prefix so feed titles are clean.
		line = strings.TrimPrefix(line, "Build: ")
		return strings.TrimSpace(line)
	}
	return ""
}

// syncClaims fetches KindClaim nodes from the graph and writes them to outPath
// as a markdown file. The MCP knowledge server indexes this file so that
// knowledge_search can find claims asserted via the assert op.
func syncClaims(apiKey, baseURL, outPath string) error {
	u := baseURL + "/app/hive/knowledge?tab=claims&limit=200"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, b)
	}

	var result struct {
		Claims []struct {
			Title     string   `json:"title"`
			Body      string   `json:"body"`
			State     string   `json:"state"`
			Author    string   `json:"author"`
			CreatedAt string   `json:"created_at"`
			Causes    []string `json:"causes"`
		} `json:"claims"`
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("decode claims: %w", err)
	}

	if len(result.Claims) == 0 {
		return nil // nothing to write
	}

	var sb strings.Builder
	sb.WriteString("# Knowledge Claims\n\n")
	sb.WriteString("Asserted knowledge claims from the hive graph store.\n\n")
	for _, c := range result.Claims {
		sb.WriteString(fmt.Sprintf("## %s\n\n", c.Title))
		if c.State != "" || c.Author != "" {
			sb.WriteString("**State:** " + c.State)
			if c.Author != "" {
				sb.WriteString(" | **Author:** " + c.Author)
			}
			sb.WriteString("\n\n")
		}
		if len(c.Causes) > 0 {
			sb.WriteString("**Causes:** " + strings.Join(c.Causes, ", ") + "\n\n")
		}
		if c.Body != "" {
			sb.WriteString(c.Body + "\n\n")
		}
		sb.WriteString("---\n\n")
	}

	if err := os.WriteFile(outPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}

	fmt.Printf("synced %d claims to %s\n", len(result.Claims), outPath)
	return nil
}

// assertScoutGap reads loop/scout.md, extracts the gap title and iteration number,
// and creates a KindClaim node on the graph via op=assert. This makes Scout gaps
// searchable via knowledge_search and preserves them across iterations (scout.md
// is overwritten each iteration; graph nodes are permanent).
//
// causeIDs should contain the build document node ID so the claim is causally
// linked to the iteration that generated it (Invariant 2: CAUSALITY).
func assertScoutGap(apiKey, baseURL string, causeIDs []string) error {
	data, err := os.ReadFile("loop/scout.md")
	if err != nil {
		return fmt.Errorf("read scout.md: %w", err)
	}

	iteration := extractIterationFromScout(data)
	gapTitle := extractGapTitle(data)
	if gapTitle == "" {
		return fmt.Errorf("could not find gap title in scout.md")
	}

	body := fmt.Sprintf("Iteration %s\n\n%s", iteration, gapTitle)

	p := map[string]string{
		"op":    "assert",
		"kind":  "claim",
		"title": gapTitle,
		"body":  body,
	}
	if len(causeIDs) > 0 {
		p["causes"] = strings.Join(causeIDs, ",")
	}
	payload, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", baseURL+"/app/hive/op", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, b)
	}

	fmt.Printf("asserted scout gap as claim: %q (iteration %s)\n", gapTitle, iteration)
	return nil
}

// extractIterationFromScout parses "Iteration N" from scout.md header.
func extractIterationFromScout(data []byte) string {
	re := regexp.MustCompile(`Iteration (\d+)`)
	m := re.FindSubmatch(data)
	if m == nil {
		return "unknown"
	}
	return string(m[1])
}

// assertCritique reads loop/critique.md and creates a KindClaim node so
// critique verdicts are searchable via knowledge_search and survive
// critique.md being overwritten next iteration.
//
// causeIDs should contain the build document node ID so the claim is causally
// linked to the iteration that generated it (Invariant 2: CAUSALITY).
func assertCritique(apiKey, baseURL string, causeIDs []string) error {
	data, err := os.ReadFile("loop/critique.md")
	if err != nil {
		return fmt.Errorf("read critique.md: %w", err)
	}

	title := extractCritiqueTitle(data)
	if title == "" {
		return fmt.Errorf("could not find critique title in critique.md")
	}

	p := map[string]string{
		"op":    "assert",
		"kind":  "claim",
		"title": title,
		"body":  string(data),
	}
	if len(causeIDs) > 0 {
		p["causes"] = strings.Join(causeIDs, ",")
	}
	payload, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", baseURL+"/app/hive/op", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, b)
	}

	fmt.Printf("asserted critique as claim: %q\n", title)
	return nil
}

// extractCritiqueTitle extracts the title from the first heading in critique.md.
// e.g. "# Critique: Fix: foo" → "Critique: Fix: foo"
func extractCritiqueTitle(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			line = strings.TrimLeft(line, "#")
			return strings.TrimSpace(line)
		}
	}
	return ""
}

// assertLatestReflection reads loop/reflections.md, extracts the most recent
// entry (the first ## section), and creates a KindDocument node so reflections
// are persistent and browsable in the Documents lens.
//
// causeIDs should contain the build document node ID so the reflection is
// causally linked to the iteration that generated it (Invariant 2: CAUSALITY).
func assertLatestReflection(apiKey, baseURL string, causeIDs []string) error {
	data, err := os.ReadFile("loop/reflections.md")
	if err != nil {
		return fmt.Errorf("read reflections.md: %w", err)
	}

	title, body := extractLatestReflection(data)
	if title == "" {
		return fmt.Errorf("could not find reflection entry in reflections.md")
	}

	p := map[string]string{
		"op":          "intend",
		"kind":        "document",
		"title":       "Reflection: " + title,
		"description": body,
	}
	if len(causeIDs) > 0 {
		p["causes"] = strings.Join(causeIDs, ",")
	}
	payload, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", baseURL+"/app/hive/op", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, b)
	}

	fmt.Printf("asserted reflection as document: %q\n", title)
	return nil
}

// extractLatestReflection parses the first ## section from reflections.md
// (which is the most recent entry, since entries are prepended).
// Returns the section heading and body text.
func extractLatestReflection(data []byte) (title, body string) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var lines []string
	inEntry := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			if inEntry {
				break // hit next entry, stop
			}
			title = strings.TrimPrefix(line, "## ")
			inEntry = true
			continue
		}
		if inEntry {
			lines = append(lines, line)
		}
	}

	body = strings.TrimSpace(strings.Join(lines, "\n"))
	return title, body
}

// extractGapTitle parses the "**Gap:** ..." line from scout.md.
func extractGapTitle(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "**Gap:**") {
			gap := strings.TrimPrefix(line, "**Gap:**")
			return strings.TrimSpace(gap)
		}
	}
	return ""
}

// post creates a KindDocument node for the build report via the intend op.
// Build reports are structured documents (not casual posts) and belong in the
// Documents lens, not the Feed.
// Returns the created node ID so it can be used as a cause for subsequent claims.
func post(apiKey, baseURL, title, body string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"op":          "intend",
		"kind":        "document",
		"title":       title,
		"description": body,
	})
	req, _ := http.NewRequest("POST", baseURL+"/app/hive/op", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, b)
	}

	var result struct {
		Node struct {
			ID string `json:"id"`
		} `json:"node"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Node.ID, nil
}
