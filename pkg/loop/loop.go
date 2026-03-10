// Package loop implements the agentic loop — sustained autonomy for agents.
//
// An agent doesn't just respond to a prompt — it observes the world,
// decides what to do, acts, and observes again. The loop runs until:
//   - Quiescence — no new events, nothing to do
//   - Escalation — the agent needs human approval
//   - HALT — the Guardian stopped the agent
//   - Budget — token/cost/iteration/time limit reached
//
// The loop transforms the hive from a pipeline into a society.
package loop

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/bus"
	"github.com/lovyou-ai/eventgraph/go/pkg/decision"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/resources"
	"github.com/lovyou-ai/hive/pkg/roles"
)

// StopReason describes why a loop stopped.
type StopReason string

const (
	StopQuiescence  StopReason = "quiescence"
	StopEscalation  StopReason = "escalation"
	StopHalt        StopReason = "halt"
	StopBudget      StopReason = "budget"
	StopError       StopReason = "error"
	StopCancelled   StopReason = "cancelled"
	StopTaskDone    StopReason = "task_done"
)

// Result is the outcome of a loop run.
type Result struct {
	Reason     StopReason
	Iterations int
	Budget     resources.BudgetSnapshot
	Detail     string // human-readable explanation
}

// Config configures an agentic loop.
type Config struct {
	// Agent is the agent to run. Required.
	Agent *roles.Agent

	// HumanID is the human operator's ID (for escalation attribution).
	HumanID types.ActorID

	// Budget limits for this loop run. Required.
	Budget resources.BudgetConfig

	// ObservationWindow is how many recent events to include in context.
	// Defaults to 20.
	ObservationWindow int

	// Task is the initial task description that seeds the loop.
	// If empty, the agent observes the graph and self-directs.
	Task string

	// Bus is an optional event bus for real-time notifications.
	// When set, the loop subscribes to relevant events and wakes
	// the agent when new work arrives.
	Bus bus.IBus

	// QuiescenceDelay is how long to wait for new events before
	// declaring quiescence. Defaults to 5 seconds.
	QuiescenceDelay time.Duration

	// OnIteration is called after each loop iteration (for monitoring).
	// Optional.
	OnIteration func(iteration int, response string)
}

// Loop runs an agent's observe-reason-act-reflect cycle.
type Loop struct {
	agent   *roles.Agent
	humanID types.ActorID
	budget  *resources.Budget
	config  Config

	// mu protects pendingEvents.
	mu            sync.Mutex
	pendingEvents []event.Event
	wake          chan struct{} // signaled when new events arrive via bus
}

// New creates a new agentic loop.
func New(cfg Config) (*Loop, error) {
	if cfg.Agent == nil {
		return nil, fmt.Errorf("agent is required")
	}
	if cfg.ObservationWindow <= 0 {
		cfg.ObservationWindow = 20
	}
	if cfg.QuiescenceDelay <= 0 {
		cfg.QuiescenceDelay = 5 * time.Second
	}

	return &Loop{
		agent:   cfg.Agent,
		humanID: cfg.HumanID,
		budget:  resources.NewBudget(cfg.Budget),
		config:  cfg,
		wake:    make(chan struct{}, 1),
	}, nil
}

// Run executes the agentic loop until a stopping condition is met.
func (l *Loop) Run(ctx context.Context) Result {
	// Subscribe to bus events if available.
	var subID bus.SubscriptionID
	if l.config.Bus != nil {
		// Subscribe to all events — the agent sees everything on the graph.
		pattern := types.MustSubscriptionPattern("*")
		subID = l.config.Bus.Subscribe(pattern, l.onEvent)
		defer l.config.Bus.Unsubscribe(subID)
	}

	iteration := 0
	consecutiveEmpty := 0

	for {
		// Check context cancellation.
		if ctx.Err() != nil {
			return l.result(StopCancelled, iteration, "context cancelled")
		}

		// Check budget before each iteration.
		if err := l.budget.Check(); err != nil {
			return l.result(StopBudget, iteration, err.Error())
		}

		iteration++

		// 1. OBSERVE — gather context from the graph.
		observation, err := l.observe(ctx)
		if err != nil {
			return l.result(StopError, iteration, fmt.Sprintf("observe: %v", err))
		}

		// 2. REASON + ACT — kick the agent with context and task.
		prompt := l.buildPrompt(observation, iteration)
		response, usage, err := l.reason(ctx, prompt)
		if err != nil {
			return l.result(StopError, iteration, fmt.Sprintf("reason: %v", err))
		}

		// Record resource consumption.
		l.budget.RecordUsage(usage)

		if l.config.OnIteration != nil {
			l.config.OnIteration(iteration, response)
		}

		// 3. CHECK stopping conditions in the response.
		if stop := l.checkResponse(ctx, response, iteration); stop != nil {
			return *stop
		}

		// 4. Check for quiescence — agent said nothing useful, no new events.
		if l.isQuiescent(response) {
			consecutiveEmpty++
			if consecutiveEmpty >= 2 {
				// Wait for new events from bus or timeout.
				if l.config.Bus != nil {
					if l.waitForEvents(ctx) {
						consecutiveEmpty = 0
						continue
					}
				}
				return l.result(StopQuiescence, iteration, "no new work after multiple iterations")
			}
		} else {
			consecutiveEmpty = 0
		}
	}
}

// observe gathers recent events from the graph as context.
func (l *Loop) observe(ctx context.Context) (string, error) {
	rt := l.agent.Runtime

	// Record observation.
	_, err := rt.Observe(ctx, l.config.ObservationWindow)
	if err != nil {
		return "", err
	}

	// Get recent events for context.
	events, err := rt.Memory(l.config.ObservationWindow)
	if err != nil {
		return "", err
	}

	// Also include any pending bus events.
	l.mu.Lock()
	pending := l.pendingEvents
	l.pendingEvents = nil
	l.mu.Unlock()

	var sb strings.Builder
	sb.WriteString("## Recent Events\n")
	for _, ev := range events {
		sb.WriteString(fmt.Sprintf("- [%s] %s by %s\n",
			ev.Type().Value(), ev.ID().Value(), ev.Source().Value()))
	}

	if len(pending) > 0 {
		sb.WriteString("\n## New Events (since last check)\n")
		for _, ev := range pending {
			sb.WriteString(fmt.Sprintf("- [%s] %s by %s\n",
				ev.Type().Value(), ev.ID().Value(), ev.Source().Value()))
		}
	}

	return sb.String(), nil
}

// buildPrompt constructs the reasoning prompt for this iteration.
func (l *Loop) buildPrompt(observation string, iteration int) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("You are %s (%s), iteration %d of your agentic loop.\n\n",
		l.agent.Name, l.agent.Role, iteration))

	if l.config.Task != "" && iteration == 1 {
		sb.WriteString(fmt.Sprintf("## Your Task\n%s\n\n", l.config.Task))
	}

	sb.WriteString(observation)

	sb.WriteString(`

## Instructions
Based on the events above, decide what to do next.

End your response with exactly one signal on its own line, in this exact JSON format:
/signal {"signal": "IDLE"}

Valid signals:
- IDLE       — nothing to do right now, waiting for events
- TASK_DONE  — all work is complete
- ESCALATE   — need human approval or decision (include "reason")
- HALT       — policy violation or integrity issue (include "reason")

Examples:
  I reviewed the latest build events and found no issues.
  /signal {"signal": "IDLE"}

  The code violates security policy.
  /signal {"signal": "HALT", "reason": "SQL injection in user input handler"}

Respond concisely. Focus on actions, not explanations.
Every response MUST end with exactly one /signal line.
`)

	return sb.String()
}

// reason calls the agent's LLM and returns the response text and token usage.
func (l *Loop) reason(ctx context.Context, prompt string) (string, decision.TokenUsage, error) {
	rt := l.agent.Runtime
	memory, _ := rt.Memory(10)
	resp, err := rt.Provider().Reason(ctx, prompt, memory)
	if err != nil {
		return "", decision.TokenUsage{}, err
	}
	return resp.Content(), resp.Usage(), nil
}

// Signal is the structured JSON signal emitted by agents at the end of each response.
type Signal struct {
	Signal string `json:"signal"` // IDLE, TASK_DONE, ESCALATE, HALT
	Reason string `json:"reason,omitempty"`
}

// SignalIDLE is the signal value for "nothing to do".
const SignalIDLE = "IDLE"

// SignalTaskDone is the signal value for "work complete".
const SignalTaskDone = "TASK_DONE"

// SignalEscalate is the signal value for "needs human".
const SignalEscalate = "ESCALATE"

// SignalHalt is the signal value for "policy violation".
const SignalHalt = "HALT"

// parseSignal extracts the /signal JSON line from a response.
// Returns nil if no valid signal is found.
func parseSignal(response string) *Signal {
	lines := strings.Split(response, "\n")
	// Search from the bottom — signal should be the last line.
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "/signal ") {
			jsonStr := strings.TrimPrefix(trimmed, "/signal ")
			var sig Signal
			if err := json.Unmarshal([]byte(jsonStr), &sig); err == nil && sig.Signal != "" {
				sig.Signal = strings.ToUpper(sig.Signal)
				return &sig
			}
		}
	}
	return nil
}

// ContainsSignal checks whether a signal keyword appears as a directive in the
// response. A directive must appear at the start of a line (possibly after
// whitespace). This is the text-based fallback when the LLM doesn't emit
// a /signal JSON line.
//
// Valid examples: "HALT: violation", "  HALT", "\nHALT\n", "TASK_DONE"
// Invalid: "No HALT required", "we should not ESCALATE"
func ContainsSignal(response, signal string) bool {
	signal = strings.ToUpper(signal)
	for _, line := range strings.Split(response, "\n") {
		trimmed := strings.ToUpper(strings.TrimSpace(line))
		if trimmed == signal || strings.HasPrefix(trimmed, signal+":") || strings.HasPrefix(trimmed, signal+" ") {
			return true
		}
	}
	return false
}

// checkResponse examines the LLM response for stopping signals.
// Priority: HALT > ESCALATE > TASK_DONE (HALT is constitutional, must never be masked).
//
// Prefers structured /signal JSON parsing. Falls back to line-start text matching
// if the LLM doesn't emit the JSON signal line.
func (l *Loop) checkResponse(ctx context.Context, response string, iteration int) *Result {
	sig := parseSignal(response)
	if sig == nil {
		// Fallback: scan for text-based signals.
		return l.checkResponseText(ctx, response, iteration)
	}

	switch sig.Signal {
	case SignalHalt:
		r := l.result(StopHalt, iteration, sig.Reason)
		return &r
	case SignalEscalate:
		if _, err := l.agent.Runtime.Escalate(ctx, l.humanID,
			fmt.Sprintf("loop iteration %d: %s", iteration, sig.Reason)); err != nil {
			fmt.Printf("warning: escalation event failed: %v\n", err)
		}
		r := l.result(StopEscalation, iteration, sig.Reason)
		return &r
	case SignalTaskDone:
		if _, err := l.agent.Runtime.Learn(ctx,
			"task completed after loop iteration "+fmt.Sprint(iteration), "loop"); err != nil {
			fmt.Printf("warning: completion event failed: %v\n", err)
		}
		r := l.result(StopTaskDone, iteration, sig.Reason)
		return &r
	case SignalIDLE:
		// Not a stop — handled by isQuiescent.
		return nil
	default:
		fmt.Printf("warning: unknown signal %q from %s\n", sig.Signal, l.agent.Name)
		return nil
	}
}

// checkResponseText is the text-based fallback for signal detection.
func (l *Loop) checkResponseText(ctx context.Context, response string, iteration int) *Result {
	if ContainsSignal(response, "HALT") {
		r := l.result(StopHalt, iteration, response)
		return &r
	}
	if ContainsSignal(response, "ESCALATE") {
		if _, err := l.agent.Runtime.Escalate(ctx, l.humanID,
			fmt.Sprintf("loop iteration %d: %s", iteration, response)); err != nil {
			fmt.Printf("warning: escalation event failed: %v\n", err)
		}
		r := l.result(StopEscalation, iteration, response)
		return &r
	}
	if ContainsSignal(response, "TASK_DONE") {
		if _, err := l.agent.Runtime.Learn(ctx,
			"task completed after loop iteration "+fmt.Sprint(iteration), "loop"); err != nil {
			fmt.Printf("warning: completion event failed: %v\n", err)
		}
		r := l.result(StopTaskDone, iteration, response)
		return &r
	}
	return nil
}

// isQuiescent returns true if the response indicates the agent has nothing to do.
// Checks for /signal JSON first, then falls back to text matching.
func (l *Loop) isQuiescent(response string) bool {
	if sig := parseSignal(response); sig != nil {
		return sig.Signal == SignalIDLE
	}
	return ContainsSignal(response, "IDLE")
}

// onEvent is called by the bus when a new event arrives.
func (l *Loop) onEvent(ev event.Event) {
	// Skip our own events to avoid infinite loops.
	if ev.Source() == l.agent.Runtime.ID() {
		return
	}

	l.mu.Lock()
	l.pendingEvents = append(l.pendingEvents, ev)
	l.mu.Unlock()

	// Signal the wake channel (non-blocking).
	select {
	case l.wake <- struct{}{}:
	default:
	}
}

// waitForEvents blocks until new events arrive or quiescence timeout.
// Returns true if events arrived, false if timed out.
func (l *Loop) waitForEvents(ctx context.Context) bool {
	timer := time.NewTimer(l.config.QuiescenceDelay)
	defer timer.Stop()

	select {
	case <-l.wake:
		return true
	case <-timer.C:
		return false
	case <-ctx.Done():
		return false
	}
}

// result creates a Result with budget snapshot.
func (l *Loop) result(reason StopReason, iterations int, detail string) Result {
	return Result{
		Reason:     reason,
		Iterations: iterations,
		Budget:     l.budget.Snapshot(),
		Detail:     detail,
	}
}

// AgentResult pairs a loop result with the agent's role and name,
// avoiding silent data loss when multiple agents share a role.
type AgentResult struct {
	Role   roles.Role
	Name   string
	Result Result
}

// RunConcurrent runs multiple agent loops concurrently and returns when all stop.
// Each loop runs in its own goroutine. Returns one result per agent.
func RunConcurrent(ctx context.Context, configs []Config) []AgentResult {
	results := make([]AgentResult, len(configs))
	var wg sync.WaitGroup

	for i, cfg := range configs {
		wg.Add(1)
		go func(idx int, c Config) {
			defer wg.Done()

			l, err := New(c)
			if err != nil {
				results[idx] = AgentResult{
					Role:   c.Agent.Role,
					Name:   c.Agent.Name,
					Result: Result{Reason: StopError, Detail: err.Error()},
				}
				return
			}

			results[idx] = AgentResult{
				Role:   c.Agent.Role,
				Name:   c.Agent.Name,
				Result: l.Run(ctx),
			}
		}(i, cfg)
	}

	wg.Wait()
	return results
}

