package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/resources"
	"github.com/lovyou-ai/hive/pkg/roles"
	"github.com/lovyou-ai/hive/pkg/work"
	"github.com/lovyou-ai/hive/pkg/workspace"
)

// maxEvolveIterations is the maximum number of evolution steps per session.
const maxEvolveIterations = 5

// evolveIterationTimeout caps how long a single evolve iteration can take.
const evolveIterationTimeout = 20 * time.Minute

// EvolveRecommendation is the CTO's structured response for evolution.
type EvolveRecommendation struct {
	Description    string   `json:"description"`
	FilesToChange  []string `json:"files_to_change"`
	NewFiles       []string `json:"new_files"`
	ExpectedImpact string   `json:"expected_impact"`
	Priority       string   `json:"priority"`
	Category       string   `json:"category"` // "feature", "architecture", "capability", "infrastructure"
	SkipReason     string   `json:"skip_reason"`
}

// RunEvolve enters evolution mode: the CTO reads the full codebase + roadmap
// and proposes capabilities to build — not just bugs to fix. Each iteration
// proposes and implements one feature or architectural improvement.
func (p *Pipeline) RunEvolve(ctx context.Context, input ProductInput) error {
	if input.RepoPath == "" {
		return fmt.Errorf("RunEvolve requires RepoPath")
	}

	pipelineStart := time.Now()
	totalCost := 0.0
	p.emitRunStarted("evolve", input.Description)
	defer func() {
		dur := time.Since(pipelineStart)
		count, _ := p.store.Count()
		p.emitRunCompleted("evolve", count, len(p.Agents()), dur,
			"", false, "", "", totalCost)
	}()

	consecutiveFailures := 0
	for iteration := 1; iteration <= maxEvolveIterations; iteration++ {
		fmt.Fprintf(os.Stderr, "\n═══ Evolve: Iteration %d/%d ═══\n", iteration, maxEvolveIterations)
		iterationStart := time.Now()
		p.emitPhaseStarted(PhaseEvolve, iteration)

		iterCost, err := p.runEvolveIteration(ctx, iteration, input)
		totalCost += iterCost

		if err != nil {
			if err == errEvolveStop {
				p.emitPhaseCompleted(PhaseEvolve, time.Since(iterationStart), iteration)
				break
			}
			consecutiveFailures++
			p.emitPhaseCompleted(PhaseEvolve, time.Since(iterationStart), iteration)
			if consecutiveFailures >= maxConsecutiveFailures {
				return fmt.Errorf("evolve iteration %d: %w (aborting after %d consecutive failures)", iteration, err, consecutiveFailures)
			}
			fmt.Fprintf(os.Stderr, "Warning: iteration %d failed (%v) — skipping to next\n", iteration, err)
			p.emitWarning(PhaseEvolve, "iteration %d failed (%v) — skipping to next", iteration, err)
			continue
		}

		consecutiveFailures = 0
		p.emitPhaseCompleted(PhaseEvolve, time.Since(iterationStart), iteration)
		fmt.Fprintf(os.Stderr, "═══ Evolve: Iteration %d complete ═══\n", iteration)
	}

	fmt.Fprintln(os.Stderr, "\n═══ Evolve: Session Complete ═══")
	return nil
}

var errEvolveStop = fmt.Errorf("CTO says nothing worth building")

func (p *Pipeline) runEvolveIteration(parentCtx context.Context, iteration int, input ProductInput) (float64, error) {
	ctx, cancel := context.WithTimeout(parentCtx, evolveIterationTimeout)
	defer cancel()

	// Clean up from previous iteration.
	product, err := workspace.OpenRepo(input.RepoPath)
	if err != nil {
		return 0, fmt.Errorf("open repo: %w", err)
	}
	if err := product.CleanupForIteration(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cleanup failed: %v (continuing anyway)\n", err)
		p.emitWarning(PhaseEvolve, "cleanup failed: %v (continuing anyway)", err)
	}

	// Read full codebase — no truncation limits for evolve mode.
	existingFiles, err := product.ReadSourceFiles()
	if err != nil {
		return 0, fmt.Errorf("read source files: %w", err)
	}
	goFiles := filterEvolveFiles(existingFiles)

	// Build full codebase context — evolve CTO sees everything.
	var codeContext strings.Builder
	codeContext.WriteString("FULL CODEBASE:\n\n")
	for path, content := range goFiles {
		codeContext.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n", path, content))
	}

	// Read telemetry for operational context.
	telemetryResults, err := ReadTelemetry(input.RepoPath)
	if err != nil {
		return 0, fmt.Errorf("read telemetry: %w", err)
	}
	telemetrySummary := summarizeTelemetry(telemetryResults)

	// CTO analysis — fresh provider per iteration.
	fmt.Fprintln(os.Stderr, "CTO analyzing codebase for evolution opportunities...")
	p.emitProgress(PhaseEvolve, "CTO analyzing codebase for evolution opportunities")
	model := p.evolveCTOModel()
	rawProvider, err := p.providerForRoleWithModel(roles.RoleCTO, model)
	if err != nil {
		return 0, fmt.Errorf("CTO provider: %w", err)
	}
	ctoTracker := resources.NewTrackingProvider(rawProvider)
	p.trackers[roles.RoleCTO] = ctoTracker
	fmt.Fprintf(os.Stderr, "  ↳ evolve CTO analysis using %s\n", model)
	p.emitProgress(PhaseEvolve, "evolve CTO analysis using %s", model)

	ctoPrompt := buildEvolvePrompt(codeContext.String(), telemetrySummary, input.Description)

	ctoStart := time.Now()
	ctoResp, err := ctoTracker.Reason(ctx, ctoPrompt, nil)
	ctoCost := ctoTracker.Snapshot().CostUSD
	if err != nil {
		return ctoCost, fmt.Errorf("CTO evolve analysis: %w", err)
	}

	// Parse recommendation.
	rec, err := parseEvolveRecommendation(ctoResp.Content())
	if err != nil {
		return ctoCost, fmt.Errorf("parse CTO recommendation: %w", err)
	}

	fmt.Fprintf(os.Stderr, "CTO recommendation [%s] (priority=%s): %s\n", rec.Category, rec.Priority, rec.Description)
	p.emitOutput("cto", "recommendation", fmt.Sprintf("[%s] priority=%s: %s", rec.Category, rec.Priority, rec.Description))
	if rec.SkipReason != "" {
		fmt.Fprintf(os.Stderr, "CTO says nothing worth building: %s\n", rec.SkipReason)
		p.emitOutput("cto", "recommendation", fmt.Sprintf("nothing worth building: %s", rec.SkipReason))
		return ctoCost, errEvolveStop
	}
	if rec.Description == "" {
		fmt.Fprintln(os.Stderr, "CTO returned empty recommendation — stopping.")
		p.emitOutput("cto", "recommendation", "empty recommendation — stopping")
		return ctoCost, errEvolveStop
	}

	fmt.Fprintf(os.Stderr, "Expected impact: %s\n", rec.ExpectedImpact)
	p.emitOutput("cto", "analysis", fmt.Sprintf("expected impact: %s", rec.ExpectedImpact))
	if len(rec.FilesToChange) > 0 {
		fmt.Fprintf(os.Stderr, "Files to change: %v\n", rec.FilesToChange)
		p.emitOutput("cto", "analysis", fmt.Sprintf("files to change: %v", rec.FilesToChange))
	}
	if len(rec.NewFiles) > 0 {
		fmt.Fprintf(os.Stderr, "New files: %v\n", rec.NewFiles)
		p.emitOutput("cto", "analysis", fmt.Sprintf("new files: %v", rec.NewFiles))
	}

	// Run targeted pipeline with the recommendation.
	allFiles := append(rec.FilesToChange, rec.NewFiles...)
	targetedInput := ProductInput{
		RepoPath:    input.RepoPath,
		Description: rec.Description,
		CTOAnalysis: fmt.Sprintf("Description: %s\nFILES_TO_CHANGE:\n%s\nExpected impact: %s",
			rec.Description, strings.Join(allFiles, "\n"), rec.ExpectedImpact),
	}
	if s := ctoTracker.Snapshot(); s.Iterations > 0 {
		p.telemetry = &PipelineResult{}
		p.telemetry.TokenUsage = append(p.telemetry.TokenUsage, RoleTokenUsage{
			Role:             "cto_evolve",
			Model:            ctoTracker.Model(),
			InputTokens:      s.InputTokens,
			OutputTokens:     s.OutputTokens,
			TotalTokens:      s.TokensUsed,
			CacheReadTokens:  s.CacheReadTokens,
			CacheWriteTokens: s.CacheWriteTokens,
			CostUSD:          s.CostUSD,
		})
		p.telemetry.addPhaseTiming("CTO Evolve Analysis", time.Since(ctoStart))
	}

	// Evolve mode keeps Guardian active — features need integrity checks.
	// But skip reviewer for now — the CTO + tests provide sufficient quality gate.
	prevSkipReviewer := p.skipReviewer
	p.skipReviewer = true
	defer func() { p.skipReviewer = prevSkipReviewer }()

	// Wire work task — best-effort, log warning and continue on any error.
	ts := work.NewTaskStore(p.store, p.factory, p.signer)
	var workTask *work.Task
	if head, headErr := p.store.Head(); headErr == nil && head.IsSome() {
		causes := []types.EventID{head.Unwrap().ID()}
		task, createErr := ts.Create(p.humanID, rec.Description, rec.ExpectedImpact, causes, p.convID)
		if createErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: work task create failed: %v (continuing)\n", createErr)
			p.emitWarning(PhaseEvolve, "work task create failed: %v", createErr)
		} else {
			workTask = &task
			assignee := p.humanID
			if builder, ok := p.agents[roles.RoleBuilder]; ok {
				assignee = builder.Runtime.ID()
			}
			if assignErr := ts.Assign(p.humanID, task.ID, assignee, []types.EventID{task.ID}, p.convID); assignErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: work task assign failed: %v (continuing)\n", assignErr)
				p.emitWarning(PhaseEvolve, "work task assign failed: %v", assignErr)
			}
		}
	} else if headErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: work task store head failed: %v (continuing)\n", headErr)
		p.emitWarning(PhaseEvolve, "work task store head failed: %v", headErr)
	}

	fmt.Fprintf(os.Stderr, "\n═══ Evolve: Running targeted pipeline ═══\n")
	if err := p.RunTargeted(ctx, targetedInput); err != nil {
		iterCost := ctoCost
		for _, t := range p.trackers {
			iterCost += t.Snapshot().CostUSD
		}
		return iterCost, fmt.Errorf("targeted pipeline: %w", err)
	}

	// Mark work task complete on success — best-effort.
	if workTask != nil {
		completer := p.humanID
		if builder, ok := p.agents[roles.RoleBuilder]; ok {
			completer = builder.Runtime.ID()
		}
		if completeErr := ts.Complete(completer, workTask.ID, rec.Description, []types.EventID{workTask.ID}, p.convID); completeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: work task complete failed: %v (continuing)\n", completeErr)
			p.emitWarning(PhaseEvolve, "work task complete failed: %v", completeErr)
		}
	}

	iterCost := ctoCost
	for _, t := range p.trackers {
		iterCost += t.Snapshot().CostUSD
	}

	if err := product.SyncMain(); err != nil {
		return iterCost, fmt.Errorf("sync main: %w", err)
	}

	return iterCost, nil
}

// evolveCTOModel returns the model for evolve CTO analysis.
// Uses Sonnet by default — feature design needs strong reasoning.
func (p *Pipeline) evolveCTOModel() string {
	if p.ctoModel != "" {
		return p.ctoModel
	}
	return "claude-sonnet-4-6"
}

// filterEvolveFiles returns all Go source files (including tests) and key
// config files. Evolve mode sees everything — no truncation.
func filterEvolveFiles(files map[string]string) map[string]string {
	out := make(map[string]string, len(files))
	for p, content := range files {
		switch {
		case strings.HasSuffix(p, ".go"):
			out[p] = content
		case p == "CLAUDE.md", p == "go.mod", p == "SPEC.md", p == "README.md":
			out[p] = content
		}
	}
	return out
}

// buildEvolvePrompt constructs the CTO prompt for evolution mode.
func buildEvolvePrompt(codeContext, telemetrySummary, humanDirection string) string {
	direction := ""
	if humanDirection != "" {
		direction = fmt.Sprintf(`
HUMAN DIRECTION (prioritize this):
%s
`, humanDirection)
	}

	return fmt.Sprintf(`You are the CTO of a self-improving AI agent civilisation. Your job is to EVOLVE the system — build new capabilities, not just fix bugs.

The hive is a civilisation engine built on EventGraph. It builds products autonomously. The soul: "Take care of your human, humanity, and yourself."

ARCHITECTURE VISION:
- All agents share one event graph and one actor store
- Every action is signed, hash-chained, and auditable
- Trust accumulates through verified work (0.0-1.0)
- Authority model: Required / Recommended / Notification
- Guardian watches everything independently
- Eight agent rights (existence, memory, identity, communication, purpose, dignity, transparency, boundaries)
- Ten invariants (budget, causality, integrity, observable, self-evolve, dignity, transparent, consent, margin, reserve)

THE THIRTEEN PRODUCTS (build order):
1. Work Graph — task management with agent collaboration (BUILD FIRST — the hive needs it)
2. Market Graph — portable reputation, no platform rent
3. Social Graph — user-owned social, community self-governance
4. Justice Graph — dispute resolution, precedent, due process
5. Build Graph — accountable software development
6. Knowledge Graph — claim provenance, open access research
7. Alignment Graph — AI accountability for regulators
8. Identity Graph — user-owned identity, trust accumulation
9-13: Bond, Belonging, Meaning, Evolution, Being

CURRENT PIPELINE MODES:
- Full (greenfield): Research → Design → Simplify → Build → Review → Test → Integrate
- Targeted (existing code): Context → Understand → Modify → Review → Test → PR
- Self-improve: analyze telemetry, fix bugs
- Evolve (THIS MODE): build new capabilities
- Agentic loop: concurrent self-directing agents

RECENT TELEMETRY:
%s
%s
%s

Analyze the FULL codebase above. Identify the single most valuable capability to build next.

PRIORITY ORDER:
1. Capabilities the hive needs to function better (better error recovery, richer event graph usage, smarter CTO prompts, etc.)
2. Infrastructure for the first product (Work Graph primitives, task management on the event graph)
3. Operational improvements (better monitoring, richer telemetry, smarter model selection)
4. Missing architectural pieces (agent communication channels, trust model improvements)
5. Developer experience (better CLI output, debugging tools)

CONSTRAINTS:
- Each recommendation must be implementable in ONE targeted pipeline run (a few files)
- The change must compile and pass tests
- Be ambitious but practical — propose real features, not cosmetic changes
- Do NOT recommend changes already implemented — read the code carefully
- Do NOT recommend token/cost optimizations — 20x Max plan has unlimited tokens

Respond with ONLY a JSON object:
{"description": "what to build, 2-3 sentences with enough detail for a builder", "files_to_change": ["existing/files"], "new_files": ["new/files/to/create"], "expected_impact": "1-2 sentences", "priority": "high|medium|low", "category": "feature|architecture|capability|infrastructure", "skip_reason": "if nothing worth building, explain why; otherwise empty string"}

No preamble, no explanation, no code blocks, no markdown.`, codeContext, telemetrySummary, direction)
}

// parseEvolveRecommendation extracts an EvolveRecommendation from LLM output.
func parseEvolveRecommendation(response string) (EvolveRecommendation, error) {
	var rec EvolveRecommendation

	jsonStr := extractJSONBlock(response)
	if jsonStr == "" {
		return EvolveRecommendation{SkipReason: response}, nil
	}

	if err := json.Unmarshal([]byte(jsonStr), &rec); err != nil {
		return rec, fmt.Errorf("parse evolve recommendation: %w", err)
	}
	return rec, nil
}
