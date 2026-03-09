// Package pipeline orchestrates the product build pipeline.
package pipeline

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/intelligence"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/roles"
	"github.com/lovyou-ai/hive/pkg/workspace"
)

// Phase represents a stage in the product pipeline.
type Phase string

const (
	PhaseResearch  Phase = "research"
	PhaseDesign    Phase = "design"
	PhaseBuild     Phase = "build"
	PhaseReview    Phase = "review"
	PhaseTest      Phase = "test"
	PhaseIntegrate Phase = "integrate"
)

// ProductInput describes how a product idea enters the hive.
type ProductInput struct {
	Name        string // Product name (used for repo and directory). If empty, CTO derives one.
	URL         string // Read from URL (Substack post, docs, etc.)
	Description string // Natural language description
	SpecFile    string // Path to a Code Graph spec file
}

// Pipeline orchestrates agents through the product build phases.
type Pipeline struct {
	store   store.Store
	ws      *workspace.Workspace
	product *workspace.Product // current product being built

	cto      *roles.Agent
	guardian *roles.Agent
	agents   map[roles.Role]*roles.Agent
}

// Config for creating a new pipeline.
type Config struct {
	Store   store.Store
	WorkDir string // Root directory for generated products
}

// New creates a pipeline and bootstraps the CTO and Guardian.
func New(ctx context.Context, cfg Config) (*Pipeline, error) {
	ws, err := workspace.New(cfg.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("workspace: %w", err)
	}

	p := &Pipeline{
		store:  cfg.Store,
		ws:     ws,
		agents: make(map[roles.Role]*roles.Agent),
	}

	// Bootstrap CTO first — architectural oversight (Opus)
	cto, err := p.ensureAgent(ctx, roles.RoleCTO, "cto")
	if err != nil {
		return nil, fmt.Errorf("bootstrap CTO: %w", err)
	}
	p.cto = cto

	// Bootstrap Guardian — independent integrity monitor (Opus)
	guardian, err := p.ensureAgent(ctx, roles.RoleGuardian, "guardian")
	if err != nil {
		return nil, fmt.Errorf("bootstrap Guardian: %w", err)
	}
	p.guardian = guardian

	return p, nil
}

// providerForRole creates an intelligence provider with the model and system prompt
// appropriate for the role. Uses Claude CLI (flat rate via Max plan).
func (p *Pipeline) providerForRole(role roles.Role) (intelligence.Provider, error) {
	model := roles.PreferredModel(role)
	return intelligence.New(intelligence.Config{
		Provider:     "claude-cli",
		Model:        model,
		SystemPrompt: roles.SystemPrompt(role),
	})
}

// ensureAgent creates an agent of the given role if it doesn't exist yet.
// Each role gets a provider with the appropriate model and system prompt.
func (p *Pipeline) ensureAgent(ctx context.Context, role roles.Role, name string) (*roles.Agent, error) {
	if agent, ok := p.agents[role]; ok {
		return agent, nil
	}
	provider, err := p.providerForRole(role)
	if err != nil {
		return nil, fmt.Errorf("provider for %s: %w", role, err)
	}
	agent, err := roles.NewAgent(ctx, roles.AgentConfig{
		Role:     role,
		Name:     name,
		Store:    p.store,
		Provider: provider,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("  ↳ %s agent using %s\n", role, roles.PreferredModel(role))
	p.agents[role] = agent
	return agent, nil
}

// Run executes the full product pipeline for a given input.
func (p *Pipeline) Run(ctx context.Context, input ProductInput) error {
	// ── Phase 1: Research ──
	fmt.Println("═══ Phase 1: Research ═══")
	spec, err := p.research(ctx, input)
	if err != nil {
		return fmt.Errorf("research: %w", err)
	}
	p.guardianCheck(ctx, "research")

	// Derive product name if not provided
	name := input.Name
	if name == "" {
		name, err = p.deriveName(ctx, spec)
		if err != nil {
			return fmt.Errorf("derive name: %w", err)
		}
	}

	// Initialize product repo
	product, err := p.ws.InitProduct(name)
	if err != nil {
		return fmt.Errorf("init product: %w", err)
	}
	p.product = product
	fmt.Printf("Product repo: %s → %s\n", product.Dir, product.Repo)

	// ── Phase 2: Design ──
	fmt.Println("═══ Phase 2: Design ═══")
	design, err := p.design(ctx, spec)
	if err != nil {
		return fmt.Errorf("design: %w", err)
	}
	p.guardianCheck(ctx, "design")

	// ── Phase 2b: Simplify ──
	fmt.Println("═══ Phase 2b: Simplify ═══")
	design, err = p.simplify(ctx, design)
	if err != nil {
		return fmt.Errorf("simplify: %w", err)
	}

	// Save the final spec to the product repo
	if err := p.product.WriteFile("SPEC.md", design); err != nil {
		return fmt.Errorf("save spec: %w", err)
	}
	if err := p.product.Commit("docs: Code Graph specification"); err != nil {
		return fmt.Errorf("commit spec: %w", err)
	}
	fmt.Println("Spec committed to product repo.")

	// Extract language from the design
	lang := p.extractLanguage(design)
	fmt.Printf("Target language: %s\n", lang)

	// ── Phase 3: Build ──
	fmt.Println("═══ Phase 3: Build ═══")
	files, err := p.build(ctx, design, lang)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}
	p.guardianCheck(ctx, "build")

	// ── Phase 4: Review → Rebuild loop ──
	const maxReviewRounds = 3
	for round := 1; round <= maxReviewRounds; round++ {
		fmt.Printf("═══ Phase 4: Review (round %d) ═══\n", round)
		feedback, approved, err := p.review(ctx, files, design, lang)
		if err != nil {
			return fmt.Errorf("review round %d: %w", round, err)
		}
		p.guardianCheck(ctx, "review")

		if approved {
			fmt.Println("Code approved by reviewer.")
			break
		}

		if round == maxReviewRounds {
			fmt.Println("Max review rounds reached — proceeding with current code.")
			break
		}

		// Rebuild with reviewer feedback
		fmt.Printf("═══ Phase 4b: Rebuild from feedback (round %d) ═══\n", round)
		files, err = p.rebuild(ctx, files, feedback, design, lang)
		if err != nil {
			return fmt.Errorf("rebuild round %d: %w", round, err)
		}
	}

	// ── Phase 5: Test ──
	fmt.Println("═══ Phase 5: Test ═══")
	err = p.test(ctx, files, lang)
	if err != nil {
		return fmt.Errorf("test: %w", err)
	}
	p.guardianCheck(ctx, "test")

	// ── Phase 6: Integrate ──
	fmt.Println("═══ Phase 6: Integrate ═══")
	err = p.integrate(ctx)
	if err != nil {
		return fmt.Errorf("integrate: %w", err)
	}
	p.guardianCheck(ctx, "integrate")

	fmt.Println("═══ Pipeline Complete ═══")
	return nil
}

// deriveName asks the CTO to derive a short, kebab-case product name from the spec.
func (p *Pipeline) deriveName(ctx context.Context, spec string) (string, error) {
	_, name, err := p.cto.Runtime.Evaluate(ctx, "product_name",
		fmt.Sprintf(`Derive a short product name (kebab-case, 2-4 words, lowercase, no special characters) from this product idea. Reply with ONLY the name, nothing else.

Product idea:
%s`, spec))
	if err != nil {
		return "product", nil // fallback
	}
	name = strings.TrimSpace(name)
	// Sanitize: lowercase, replace spaces with hyphens, remove non-alphanumeric
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	var clean []byte
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
			clean = append(clean, c)
		}
	}
	if len(clean) == 0 {
		return "product", nil
	}
	fmt.Printf("CTO named product: %s\n", string(clean))
	return string(clean), nil
}

// research gathers information about the product idea.
func (p *Pipeline) research(ctx context.Context, input ProductInput) (string, error) {
	var spec string

	if input.SpecFile != "" {
		content, err := p.ws.ReadFile(input.SpecFile)
		if err != nil {
			return "", fmt.Errorf("read spec: %w", err)
		}
		spec = content
	} else {
		researcher, err := p.ensureAgent(ctx, roles.RoleResearcher, "researcher")
		if err != nil {
			return "", err
		}

		if input.URL != "" {
			_, evaluation, err := researcher.Runtime.Research(ctx, input.URL,
				"extract the product idea, key entities, features, and requirements. Output in Code Graph vocabulary where possible.")
			if err != nil {
				return "", fmt.Errorf("research URL: %w", err)
			}
			spec = evaluation
		} else if input.Description != "" {
			_, evaluation, err := researcher.Runtime.Evaluate(ctx, "product_idea", input.Description)
			if err != nil {
				return "", fmt.Errorf("evaluate idea: %w", err)
			}
			spec = evaluation
		}
	}

	// CTO evaluates feasibility
	_, ctoEval, err := p.cto.Runtime.Evaluate(ctx, "feasibility",
		fmt.Sprintf("Evaluate this product idea for feasibility. What agents are needed? What's the build sequence? Key risks?\n\n%s", spec))
	if err != nil {
		return "", fmt.Errorf("CTO evaluate: %w", err)
	}

	fmt.Printf("CTO Assessment:\n%s\n", ctoEval)
	return spec, nil
}

// design creates a full Code Graph spec from the product idea.
func (p *Pipeline) design(ctx context.Context, spec string) (string, error) {
	architect, err := p.ensureAgent(ctx, roles.RoleArchitect, "architect")
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`Design the full system architecture. Output a complete Code Graph spec.
Remember: derive complexity from simple compositions. Each view should have the minimal elements needed — if a view feels heavy, decompose it. Elegant, simple, beautiful.

Also specify the target language/framework at the top of your spec in a line like:
LANGUAGE: go
or
LANGUAGE: typescript

Product idea:
%s`, spec)

	_, design, err := architect.Runtime.Evaluate(ctx, "architecture", prompt)
	if err != nil {
		return "", fmt.Errorf("architect design: %w", err)
	}

	// CTO reviews the architecture
	_, review, err := p.cto.Runtime.Evaluate(ctx, "architecture_review",
		fmt.Sprintf("Review this architecture. Check: Are views minimal? Is complexity derived from composition rather than accumulated? Are there any bloated entities or views that should be decomposed? Is it elegant and simple?\n\n%s", design))
	if err != nil {
		return "", fmt.Errorf("CTO review design: %w", err)
	}

	fmt.Printf("Architecture Review:\n%s\n", review)
	return design, nil
}

// simplify reviews the Code Graph spec and reduces it to its minimal form.
func (p *Pipeline) simplify(ctx context.Context, design string) (string, error) {
	architect, err := p.ensureAgent(ctx, roles.RoleArchitect, "architect")
	if err != nil {
		return "", err
	}

	const maxRounds = 3
	current := design

	for round := 1; round <= maxRounds; round++ {
		_, analysis, err := architect.Runtime.Evaluate(ctx, "simplify",
			fmt.Sprintf(`Review this Code Graph spec for simplification opportunities.

For each View: can it be composed from fewer elements? Are any elements redundant or derivable from others?
For each Entity: is it as small as possible? Should it be split or can properties be derived?
For each State machine: are there too many states? Can transitions be reduced?
For each Layout: does it have too many children? Can sub-views be composed instead?

If you find simplifications, output the REVISED spec with the changes applied.
If the spec is already minimal, respond with exactly: MINIMAL

Current spec:
%s`, current))
		if err != nil {
			return "", fmt.Errorf("simplify round %d: %w", round, err)
		}

		upper := strings.ToUpper(strings.TrimSpace(analysis))
		if upper == "MINIMAL" || strings.HasPrefix(upper, "MINIMAL") {
			fmt.Printf("Simplification complete after %d round(s) — spec is minimal.\n", round)
			return current, nil
		}

		fmt.Printf("Simplification round %d applied.\n", round)
		current = analysis
	}

	fmt.Printf("Simplification capped at %d rounds.\n", maxRounds)
	return current, nil
}

// extractLanguage pulls the target language from the design spec.
// Looks for "LANGUAGE: xxx" in the spec. Defaults to "go".
func (p *Pipeline) extractLanguage(design string) string {
	for _, line := range strings.Split(design, "\n") {
		line = strings.TrimSpace(line)
		upper := strings.ToUpper(line)
		if strings.HasPrefix(upper, "LANGUAGE:") {
			lang := strings.TrimSpace(line[len("LANGUAGE:"):])
			lang = strings.ToLower(lang)
			if lang != "" {
				return lang
			}
		}
	}
	return "go"
}

// build generates multi-file code from the design spec.
// The builder outputs files with --- FILE: path --- markers, which are parsed
// into individual files and committed to the product repo.
func (p *Pipeline) build(ctx context.Context, design string, lang string) (map[string]string, error) {
	builder, err := p.ensureAgent(ctx, roles.RoleBuilder, "builder")
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`Generate production-quality %s code from this specification.

Output ALL files needed for a complete, runnable project. Use this format for each file:

--- FILE: path/to/file.ext ---
<file contents>

Include:
- Project config (go.mod, package.json, Cargo.toml, etc.)
- Source files (organized in packages/modules)
- Test files alongside the code they test
- A README.md with build and run instructions

Do NOT include explanation text outside of file blocks. Every line must be inside a file block.

Specification:
%s`, lang, design)

	code, err := builder.Runtime.CodeWrite(ctx, prompt, lang)
	if err != nil {
		return nil, fmt.Errorf("builder code: %w", err)
	}

	// Parse multi-file output
	files := parseFiles(code)
	if len(files) == 0 {
		// Fallback: treat entire output as a single file
		ext := langExtension(lang)
		files = map[string]string{"main" + ext: code}
	}

	// Write all files to product repo
	for path, content := range files {
		if err := p.product.WriteFile(path, content); err != nil {
			return nil, fmt.Errorf("write %s: %w", path, err)
		}
	}
	if err := p.product.Commit(fmt.Sprintf("feat: initial %s code generation from spec", lang)); err != nil {
		return nil, fmt.Errorf("commit code: %w", err)
	}

	fmt.Printf("Generated %d files, committed.\n", len(files))
	return files, nil
}

// rebuild sends reviewer feedback to the builder and generates revised code.
func (p *Pipeline) rebuild(ctx context.Context, currentFiles map[string]string, feedback string, design string, lang string) (map[string]string, error) {
	builder, err := p.ensureAgent(ctx, roles.RoleBuilder, "builder")
	if err != nil {
		return nil, err
	}

	// Build a summary of current files
	var filesSummary strings.Builder
	for path, content := range currentFiles {
		filesSummary.WriteString(fmt.Sprintf("--- FILE: %s ---\n%s\n", path, content))
	}

	prompt := fmt.Sprintf(`The reviewer provided feedback on the code. Fix the issues and output ALL files again using the same format.

Reviewer feedback:
%s

Original specification:
%s

Current code:
%s

Output the COMPLETE revised files using --- FILE: path --- markers. Include ALL files, not just changed ones.`, feedback, design, filesSummary.String())

	code, err := builder.Runtime.CodeWrite(ctx, prompt, lang)
	if err != nil {
		return nil, fmt.Errorf("rebuild: %w", err)
	}

	files := parseFiles(code)
	if len(files) == 0 {
		return currentFiles, nil // no parseable output, keep current
	}

	for path, content := range files {
		if err := p.product.WriteFile(path, content); err != nil {
			return nil, fmt.Errorf("write %s: %w", path, err)
		}
	}
	if err := p.product.Commit("fix: address reviewer feedback"); err != nil {
		return nil, fmt.Errorf("commit rebuild: %w", err)
	}

	fmt.Printf("Rebuilt %d files from feedback, committed.\n", len(files))
	return files, nil
}

// review checks code quality and spec compliance. Returns feedback and whether approved.
func (p *Pipeline) review(ctx context.Context, files map[string]string, design string, lang string) (feedback string, approved bool, err error) {
	reviewer, err := p.ensureAgent(ctx, roles.RoleReviewer, "reviewer")
	if err != nil {
		return "", false, err
	}

	// Build code summary for review
	var codeSummary strings.Builder
	for path, content := range files {
		codeSummary.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n", path, content))
	}
	allCode := codeSummary.String()

	// Code review
	_, codeReview, err := reviewer.Runtime.CodeReview(ctx, allCode, lang)
	if err != nil {
		return "", false, fmt.Errorf("code review: %w", err)
	}

	// Spec compliance
	_, specReview, err := reviewer.Runtime.Evaluate(ctx, "spec_compliance",
		fmt.Sprintf("Does this code match the design spec? Flag any deviations.\n\nDesign:\n%s\n\nCode:\n%s", design, allCode))
	if err != nil {
		return "", false, fmt.Errorf("spec review: %w", err)
	}

	// Simplicity check
	_, simplicityReview, err := reviewer.Runtime.Evaluate(ctx, "simplicity_check",
		fmt.Sprintf(`Review this code for unnecessary complexity:
- Components that could be derived from simpler compositions?
- Redundant abstractions or over-engineered patterns?
- Did the builder add extras beyond the spec?

Code:
%s`, allCode))
	if err != nil {
		return "", false, fmt.Errorf("simplicity review: %w", err)
	}

	// Final verdict
	_, verdict, err := reviewer.Runtime.Decide(ctx, "approve_or_reject",
		fmt.Sprintf(`Based on your reviews, should this code be APPROVED or does it need CHANGES?

Code Review: %s
Spec Compliance: %s
Simplicity: %s

Reply with exactly APPROVED if the code is ready, or CHANGES NEEDED followed by the specific issues to fix.`, codeReview, specReview, simplicityReview))
	if err != nil {
		return "", false, fmt.Errorf("verdict: %w", err)
	}

	fmt.Printf("Code Review:\n%s\n\nSpec Compliance:\n%s\n\nSimplicity:\n%s\n\nVerdict: %s\n",
		codeReview, specReview, simplicityReview, verdict)

	approved = strings.Contains(strings.ToUpper(verdict), "APPROVED") &&
		!strings.Contains(strings.ToUpper(verdict), "CHANGES")
	return verdict, approved, nil
}

// test runs tests in the product directory and has the tester analyze gaps.
func (p *Pipeline) test(ctx context.Context, files map[string]string, lang string) error {
	tester, err := p.ensureAgent(ctx, roles.RoleTester, "tester")
	if err != nil {
		return err
	}

	// Actually run tests in the product directory
	testCmd, testArgs := langTestCommand(lang)
	fmt.Printf("Running: %s %s\n", testCmd, strings.Join(testArgs, " "))

	cmd := exec.Command(testCmd, testArgs...)
	cmd.Dir = p.product.Dir
	testOutput, testErr := cmd.CombinedOutput()

	testResult := string(testOutput)
	if testErr != nil {
		fmt.Printf("Tests failed:\n%s\n", testResult)
	} else {
		fmt.Printf("Tests passed:\n%s\n", testResult)
	}

	// Have the tester analyze results and coverage gaps
	var codeSummary strings.Builder
	for path, content := range files {
		codeSummary.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n", path, content))
	}

	_, testEval, err := tester.Runtime.Evaluate(ctx, "test_analysis",
		fmt.Sprintf(`Analyze the test results and code. Are there coverage gaps? What additional tests are needed?

Test output:
%s

Code:
%s`, testResult, codeSummary.String()))
	if err != nil {
		return fmt.Errorf("test analysis: %w", err)
	}

	fmt.Printf("Test Analysis:\n%s\n", testEval)

	// If tests failed, have the builder fix them
	if testErr != nil {
		fmt.Println("Attempting to fix failing tests...")
		builder, err := p.ensureAgent(ctx, roles.RoleBuilder, "builder")
		if err != nil {
			return err
		}

		fixPrompt := fmt.Sprintf(`The tests are failing. Fix the code so tests pass.

Test output:
%s

Current code:
%s

Output ALL files using --- FILE: path --- markers.`, testResult, codeSummary.String())

		fixedCode, err := builder.Runtime.CodeWrite(ctx, fixPrompt, lang)
		if err != nil {
			return fmt.Errorf("fix tests: %w", err)
		}

		fixedFiles := parseFiles(fixedCode)
		for path, content := range fixedFiles {
			if err := p.product.WriteFile(path, content); err != nil {
				return fmt.Errorf("write fix %s: %w", path, err)
			}
		}
		if len(fixedFiles) > 0 {
			_ = p.product.Commit("fix: address failing tests")
		}

		// Re-run tests
		cmd2 := exec.Command(testCmd, testArgs...)
		cmd2.Dir = p.product.Dir
		retryOutput, retryErr := cmd2.CombinedOutput()
		if retryErr != nil {
			fmt.Printf("Tests still failing after fix attempt:\n%s\n", string(retryOutput))
		} else {
			fmt.Printf("Tests now passing:\n%s\n", string(retryOutput))
		}
	}

	return nil
}

// integrate assembles and prepares for deployment.
func (p *Pipeline) integrate(ctx context.Context) error {
	integrator, err := p.ensureAgent(ctx, roles.RoleIntegrator, "integrator")
	if err != nil {
		return err
	}

	_, err = integrator.Runtime.Act(ctx, "integrate", "staging")
	if err != nil {
		return fmt.Errorf("integration: %w", err)
	}

	// Push to GitHub
	if err := p.product.Push(); err != nil {
		fmt.Printf("Push failed (may need manual push): %v\n", err)
	} else {
		fmt.Printf("Pushed to https://github.com/%s\n", p.product.Repo)
	}

	// Escalate to human for production approval
	humanID := types.MustActorID("actor_human_matt")
	_, err = integrator.Runtime.Escalate(ctx, humanID, "Product ready for human review before production deploy")
	if err != nil {
		return fmt.Errorf("escalate: %w", err)
	}

	fmt.Println("Product assembled and ready for human review.")
	return nil
}

// guardianCheck runs the Guardian's integrity check after a phase.
func (p *Pipeline) guardianCheck(ctx context.Context, phase string) {
	events, err := p.guardian.Runtime.Memory(20)
	if err != nil || len(events) == 0 {
		return
	}

	var summary strings.Builder
	for _, ev := range events {
		summary.WriteString(fmt.Sprintf("[%s] %s: %s\n", ev.Type().Value(), ev.Source().Value(), ev.ID().Value()))
	}

	_, eval, err := p.guardian.Runtime.Evaluate(ctx, "integrity_check_"+phase,
		fmt.Sprintf("Review these recent events (after %s phase) for policy violations, trust anomalies, or authority overreach:\n\n%s",
			phase, summary.String()))
	if err != nil {
		fmt.Printf("Guardian check failed: %v\n", err)
		return
	}

	if containsAlert(eval) {
		fmt.Printf("⚠ Guardian Alert (after %s):\n%s\n", phase, eval)
		_, _ = p.guardian.Runtime.Emit(event.AgentEscalatedContent{
			AgentID:   p.guardian.Runtime.ID(),
			Authority: types.MustActorID("actor_human_matt"),
			Reason:    fmt.Sprintf("[%s phase] %s", phase, eval),
		})
	}
}

// ════════════════════════════════════════════════════════════════════════
// File parsing utilities
// ════════════════════════════════════════════════════════════════════════

// parseFiles extracts files from builder output using --- FILE: path --- markers.
func parseFiles(output string) map[string]string {
	files := make(map[string]string)
	lines := strings.Split(output, "\n")

	var currentPath string
	var currentContent strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--- FILE:") && strings.HasSuffix(trimmed, "---") {
			// Save previous file if any
			if currentPath != "" {
				files[currentPath] = strings.TrimRight(currentContent.String(), "\n")
			}
			// Extract new path
			path := strings.TrimSpace(trimmed[len("--- FILE:") : len(trimmed)-len("---")])
			currentPath = path
			currentContent.Reset()
		} else if currentPath != "" {
			currentContent.WriteString(line)
			currentContent.WriteString("\n")
		}
	}
	// Save last file
	if currentPath != "" {
		files[currentPath] = strings.TrimRight(currentContent.String(), "\n")
	}

	return files
}

// langExtension returns the default file extension for a language.
func langExtension(lang string) string {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return ".go"
	case "typescript", "ts":
		return ".ts"
	case "javascript", "js":
		return ".js"
	case "python", "py":
		return ".py"
	case "rust", "rs":
		return ".rs"
	case "csharp", "c#", "cs":
		return ".cs"
	default:
		return ".go"
	}
}

// langTestCommand returns the test command and args for a language.
func langTestCommand(lang string) (string, []string) {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return "go", []string{"test", "./..."}
	case "typescript", "ts":
		return "npx", []string{"vitest", "run"}
	case "javascript", "js":
		return "npm", []string{"test"}
	case "python", "py":
		return "python", []string{"-m", "pytest"}
	case "rust", "rs":
		return "cargo", []string{"test"}
	case "csharp", "c#", "cs":
		return "dotnet", []string{"test"}
	default:
		return "go", []string{"test", "./..."}
	}
}

// containsAlert checks if the Guardian's evaluation contains an alert keyword.
func containsAlert(eval string) bool {
	upper := strings.ToUpper(eval)
	for _, keyword := range []string{"HALT", "ALERT", "VIOLATION", "QUARANTINE"} {
		if strings.Contains(upper, keyword) {
			return true
		}
	}
	return false
}

// Store returns the shared event graph.
func (p *Pipeline) Store() store.Store {
	return p.store
}

// Agents returns all active agents.
func (p *Pipeline) Agents() map[roles.Role]*roles.Agent {
	return p.agents
}
