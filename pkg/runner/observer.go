package runner

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/decision"
)

// runObserver looks at the product from a human's perspective.
// Uses Operate() — needs to grep code, fetch URLs, read multiple files.
func (r *Runner) runObserver(ctx context.Context) {
	// Run every 16th tick (~4 minutes) in continuous mode.
	if !r.cfg.OneShot && r.tick%16 != 0 {
		return
	}

	log.Printf("[observer] tick %d: observing product", r.tick)

	// Build the observation instruction.
	apiKey := os.Getenv("LOVYOU_API_KEY")
	if apiKey == "" {
		log.Printf("[observer] LOVYOU_API_KEY not set; graph integrity audit will be skipped")
	}
	instruction := buildObserverInstruction(r.cfg.RepoPath, r.cfg.SpaceSlug, apiKey)

	// Use Operate() — the Observer needs file access to grep patterns, read templates, etc.
	op, ok := r.cfg.Provider.(decision.IOperator)
	if !ok {
		log.Printf("[observer] provider does not support Operate, falling back to Reason")
		r.runObserverReason(ctx)
		return
	}

	result, err := op.Operate(ctx, decision.OperateTask{
		WorkDir:      r.cfg.RepoPath,
		Instruction:  instruction,
		AllowedTools: []string{"Read", "Glob", "Grep", "Bash", "mcp__knowledge__knowledge_search"},
	})
	if err != nil {
		log.Printf("[observer] Operate error: %v", err)
		return
	}

	r.cost.Record(result.Usage)
	r.dailyBudget.Record(result.Usage.CostUSD)
	log.Printf("[observer] observation done (cost=$%.4f)", result.Usage.CostUSD)

	// Observer creates tasks directly via Operate() — no text parsing needed.

	if r.cfg.OneShot {
		r.done = true
	}
}

// runObserverReason is a fallback when Operate() isn't available.
// Gathers context manually and uses Reason().
func (r *Runner) runObserverReason(ctx context.Context) {
	// Gather what we can without tool access.
	routes := r.grepRoutes()
	kinds := r.grepEntityKinds()
	health := r.checkLiveHealth()

	prompt := fmt.Sprintf(`You are the Observer. Analyze this product for gaps.

## Routes registered
%s

## Entity kinds defined
%s

## Live health check
%s

Report findings as:
TASK_TITLE: <title>
TASK_PRIORITY: <priority>
TASK_DESCRIPTION: <description>

You may report up to 2 findings. If everything looks good, say "No issues found."`, routes, kinds, health)

	resp, err := r.cfg.Provider.Reason(ctx, prompt, nil)
	if err != nil {
		log.Printf("[observer] Reason error: %v", err)
		return
	}
	r.cost.Record(resp.Usage())
	r.dailyBudget.Record(resp.Usage().CostUSD)

	tasks := parseObserverTasks(resp.Content())
	for _, t := range tasks {
		task, err := r.cfg.APIClient.CreateTask(r.cfg.SpaceSlug, t.title, t.desc, t.priority)
		if err != nil {
			continue
		}
		log.Printf("[observer] created task %s: %s", task.ID, t.title)
		if r.cfg.AgentID != "" {
			_ = r.cfg.APIClient.ClaimTask(r.cfg.SpaceSlug, task.ID)
		}
	}

	if r.cfg.OneShot {
		r.done = true
	}
}

type observerTask struct {
	title    string
	desc     string
	priority string
}

func parseObserverTasks(content string) []observerTask {
	var tasks []observerTask
	var current observerTask

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TASK_TITLE:") {
			if current.title != "" {
				tasks = append(tasks, current)
			}
			current = observerTask{title: strings.TrimSpace(strings.TrimPrefix(line, "TASK_TITLE:"))}
		} else if strings.HasPrefix(line, "TASK_PRIORITY:") {
			current.priority = strings.TrimSpace(strings.TrimPrefix(line, "TASK_PRIORITY:"))
		} else if strings.HasPrefix(line, "TASK_DESCRIPTION:") {
			current.desc = strings.TrimSpace(strings.TrimPrefix(line, "TASK_DESCRIPTION:"))
		}
	}
	if current.title != "" {
		tasks = append(tasks, current)
	}

	// Validate priorities.
	for i := range tasks {
		switch tasks[i].priority {
		case "urgent", "high", "medium", "low":
		default:
			tasks[i].priority = "medium"
		}
	}
	return tasks
}

func buildObserverInstruction(repoPath, spaceSlug, apiKey string) string {
	part2 := buildPart2Instruction(spaceSlug, apiKey)
	return fmt.Sprintf(`You are the Observer. Audit both the product AND the hive's own graph for integrity.

## Your repo: %s

## Part 1: Product Audit

1. **Consistency:** Grep entity kind constants. For each, verify handler, route, sidebar nav, create form exist.
2. **Route health:** curl key pages — https://lovyou.ai/, /discover, /hive — check status codes.
3. **User flow:** Is there a clear path from landing → sign in → create space → use product?

%s

Also use mcp__knowledge__knowledge_search to check:
- Are lessons (claims) being asserted? Or just documents?
- Are reflections searchable? Or only in flat files?
- Do agents use the right entity kind for each artifact?

## Output

%s

If everything looks good, say "No issues found."`, repoPath, part2, buildOutputInstruction(spaceSlug, apiKey))
}

func buildPart2Instruction(spaceSlug, apiKey string) string {
	if apiKey == "" {
		return `## Part 2: Graph Integrity Audit

(Skipped — LOVYOU_API_KEY not set. Authenticated requests require an API key.)`
	}
	return fmt.Sprintf(`## Part 2: Graph Integrity Audit

Check the hive's own data on the board for structural issues:

curl -s -H "Authorization: Bearer %s" -H "Accept: application/json" "https://lovyou.ai/app/%s/board"

Look for:
1. **Unlinked nodes** — critiques/claims that don't reference the build they review
2. **Stale data** — open tasks that have been active for days with no progress
3. **Title compounding** — tasks with "Fix: Fix: Fix:..." prefixes (should be stripped)
4. **Schema violations** — any field using display names where IDs should be used (Invariant 11)
5. **Orphaned milestones** — PM milestones still open after their subtasks completed`, apiKey, spaceSlug)
}

func buildOutputInstruction(spaceSlug, apiKey string) string {
	if apiKey == "" {
		return `Report the most important findings (max 2) as:
TASK_TITLE: <title>
TASK_PRIORITY: <priority>
TASK_DESCRIPTION: <description>`
	}
	return fmt.Sprintf(`Create tasks directly for the most important findings (max 2):

curl -s -X POST -H "Authorization: Bearer %s" -H "Content-Type: application/json" -H "Accept: application/json" "https://lovyou.ai/app/%s/op" -d '{"op":"intend","kind":"task","title":"<TITLE>","description":"<DESCRIPTION>","priority":"<PRIORITY>"}'`, apiKey, spaceSlug)
}

// Helper: grep registered routes from handlers.go.
func (r *Runner) grepRoutes() string {
	cmd := exec.Command("grep", "-n", "mux.Handle", "graph/handlers.go")
	cmd.Dir = r.cfg.RepoPath
	out, err := cmd.Output()
	if err != nil {
		return "(could not grep routes)"
	}
	return string(out)
}

// Helper: grep entity kind constants from store.go.
func (r *Runner) grepEntityKinds() string {
	cmd := exec.Command("grep", "-n", "Kind.*=", "graph/store.go")
	cmd.Dir = r.cfg.RepoPath
	out, err := cmd.Output()
	if err != nil {
		return "(could not grep kinds)"
	}
	return string(out)
}

// Helper: check if the live site is up.
func (r *Runner) checkLiveHealth() string {
	pages := []string{
		"https://lovyou.ai/",
		"https://lovyou.ai/discover",
		"https://lovyou.ai/search",
	}
	var results []string
	for _, url := range pages {
		cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "--max-time", "5", url)
		out, err := cmd.Output()
		if err != nil {
			results = append(results, fmt.Sprintf("%s → error", url))
		} else {
			results = append(results, fmt.Sprintf("%s → %s", url, string(out)))
		}
	}
	return strings.Join(results, "\n")
}
