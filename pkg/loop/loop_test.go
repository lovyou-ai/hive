package loop

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/bus"
	"github.com/lovyou-ai/eventgraph/go/pkg/decision"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/intelligence"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/resources"
	"github.com/lovyou-ai/hive/pkg/roles"
)

// ════════════════════════════════════════════════════════════════════════
// Mock provider — returns canned responses in sequence
// ════════════════════════════════════════════════════════════════════════

type mockProvider struct {
	responses []string
	callCount atomic.Int32
}

func newMockProvider(responses ...string) *mockProvider {
	return &mockProvider{responses: responses}
}

func (m *mockProvider) Name() string  { return "mock" }
func (m *mockProvider) Model() string { return "mock-model" }

func (m *mockProvider) Reason(_ context.Context, _ string, _ []event.Event) (decision.Response, error) {
	idx := int(m.callCount.Add(1)) - 1
	if idx >= len(m.responses) {
		idx = len(m.responses) - 1
	}
	confidence, _ := types.NewScore(0.8)
	return decision.NewResponse(m.responses[idx], confidence, decision.TokenUsage{InputTokens: 30, OutputTokens: 20}), nil
}

var _ intelligence.Provider = (*mockProvider)(nil)

// ════════════════════════════════════════════════════════════════════════
// Test helpers
// ════════════════════════════════════════════════════════════════════════

func testAgent(t *testing.T, provider intelligence.Provider) *roles.Agent {
	t.Helper()
	return testAgentWithRole(t, provider, roles.RoleBuilder, "test-builder")
}

func humanID() types.ActorID {
	return types.MustActorID("actor_00000000000000000000000000000099")
}

// ════════════════════════════════════════════════════════════════════════
// Tests
// ════════════════════════════════════════════════════════════════════════

func TestLoopTaskDone(t *testing.T) {
	// Agent immediately says TASK_DONE.
	provider := newMockProvider("I've completed the work.\n/signal {\"signal\": \"TASK_DONE\"}")
	agent := testAgent(t, provider)

	l, err := New(Config{
		Agent:   agent,
		HumanID: humanID(),
		Budget:  resources.BudgetConfig{MaxIterations: 10},
		Task:    "build a widget",
	})
	if err != nil {
		t.Fatal(err)
	}

	result := l.Run(context.Background())
	if result.Reason != StopTaskDone {
		t.Errorf("reason = %s, want %s", result.Reason, StopTaskDone)
	}
	if result.Iterations != 1 {
		t.Errorf("iterations = %d, want 1", result.Iterations)
	}
}

func TestLoopEscalation(t *testing.T) {
	provider := newMockProvider("I need human approval.\n/signal {\"signal\": \"ESCALATE\", \"reason\": \"this requires authority\"}")
	agent := testAgent(t, provider)

	l, err := New(Config{
		Agent:   agent,
		HumanID: humanID(),
		Budget:  resources.BudgetConfig{MaxIterations: 10},
	})
	if err != nil {
		t.Fatal(err)
	}

	result := l.Run(context.Background())
	if result.Reason != StopEscalation {
		t.Errorf("reason = %s, want %s", result.Reason, StopEscalation)
	}
}

func TestLoopHalt(t *testing.T) {
	provider := newMockProvider("/signal {\"signal\": \"HALT\", \"reason\": \"integrity violation detected\"}")
	agent := testAgent(t, provider)

	l, err := New(Config{
		Agent:   agent,
		HumanID: humanID(),
		Budget:  resources.BudgetConfig{MaxIterations: 10},
	})
	if err != nil {
		t.Fatal(err)
	}

	result := l.Run(context.Background())
	if result.Reason != StopHalt {
		t.Errorf("reason = %s, want %s", result.Reason, StopHalt)
	}
}

func TestLoopBudgetExceeded(t *testing.T) {
	// Agent never finishes — budget stops it.
	provider := newMockProvider("working on it...", "still working...", "more work...")
	agent := testAgent(t, provider)

	l, err := New(Config{
		Agent:   agent,
		HumanID: humanID(),
		Budget:  resources.BudgetConfig{MaxIterations: 2},
	})
	if err != nil {
		t.Fatal(err)
	}

	result := l.Run(context.Background())
	if result.Reason != StopBudget {
		t.Errorf("reason = %s, want %s", result.Reason, StopBudget)
	}
	if result.Iterations != 2 {
		t.Errorf("iterations = %d, want 2", result.Iterations)
	}
}

func TestLoopContextCancelled(t *testing.T) {
	provider := newMockProvider("working...")
	agent := testAgent(t, provider)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	l, err := New(Config{
		Agent:   agent,
		HumanID: humanID(),
		Budget:  resources.BudgetConfig{MaxIterations: 100},
	})
	if err != nil {
		t.Fatal(err)
	}

	result := l.Run(ctx)
	if result.Reason != StopCancelled {
		t.Errorf("reason = %s, want %s", result.Reason, StopCancelled)
	}
}

func TestLoopQuiescence(t *testing.T) {
	// Agent says IDLE twice — triggers quiescence without bus.
	provider := newMockProvider("/signal {\"signal\": \"IDLE\"}", "/signal {\"signal\": \"IDLE\"}")
	agent := testAgent(t, provider)

	l, err := New(Config{
		Agent:             agent,
		HumanID:           humanID(),
		Budget:            resources.BudgetConfig{MaxIterations: 10},
		QuiescenceDelay:   10 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}

	result := l.Run(context.Background())
	if result.Reason != StopQuiescence {
		t.Errorf("reason = %s, want %s", result.Reason, StopQuiescence)
	}
}

func TestLoopMultiIteration(t *testing.T) {
	// Agent works for 3 iterations then says done.
	provider := newMockProvider(
		"Writing the code...",
		"Running tests...",
		"All tests pass.\n/signal {\"signal\": \"TASK_DONE\"}",
	)
	agent := testAgent(t, provider)

	var iterations []int
	l, err := New(Config{
		Agent:   agent,
		HumanID: humanID(),
		Budget:  resources.BudgetConfig{MaxIterations: 10},
		Task:    "build a widget",
		OnIteration: func(i int, _ string) {
			iterations = append(iterations, i)
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	result := l.Run(context.Background())
	if result.Reason != StopTaskDone {
		t.Errorf("reason = %s, want %s", result.Reason, StopTaskDone)
	}
	if result.Iterations != 3 {
		t.Errorf("iterations = %d, want 3", result.Iterations)
	}
	if len(iterations) != 3 {
		t.Errorf("callback count = %d, want 3", len(iterations))
	}
}

func TestLoopBudgetTracking(t *testing.T) {
	provider := newMockProvider("doing work", "/signal {\"signal\": \"TASK_DONE\"}")
	agent := testAgent(t, provider)

	l, err := New(Config{
		Agent:   agent,
		HumanID: humanID(),
		Budget:  resources.BudgetConfig{MaxTokens: 10000},
	})
	if err != nil {
		t.Fatal(err)
	}

	result := l.Run(context.Background())
	if result.Budget.TokensUsed == 0 {
		t.Error("expected non-zero token usage")
	}
	if result.Budget.Iterations != 2 {
		t.Errorf("budget iterations = %d, want 2", result.Budget.Iterations)
	}
}

func TestLoopWithBus(t *testing.T) {
	// Agent says IDLE twice, enters quiescence wait, bus delivers event, agent wakes and completes.
	// Three responses: IDLE, IDLE, then TASK_DONE after waking.
	provider := newMockProvider(
		"/signal {\"signal\": \"IDLE\"}",
		"/signal {\"signal\": \"IDLE\"}",
		"Got new event!\n/signal {\"signal\": \"TASK_DONE\"}",
	)
	agent := testAgent(t, provider)

	s := store.NewInMemoryStore()
	eventBus := bus.NewEventBus(s, 16)
	defer eventBus.Close()

	// Channel-based synchronisation — signal after iteration 2 (second IDLE)
	// which is when the loop will enter waitForEvents.
	secondIterDone := make(chan struct{}, 1)

	l, err := New(Config{
		Agent:           agent,
		HumanID:         humanID(),
		Budget:          resources.BudgetConfig{MaxIterations: 10},
		Bus:             eventBus,
		QuiescenceDelay: 5 * time.Second, // long delay — test won't hit it
		OnIteration: func(i int, _ string) {
			if i == 2 {
				select {
				case secondIterDone <- struct{}{}:
				default:
				}
			}
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan Result, 1)
	go func() {
		done <- l.Run(context.Background())
	}()

	// Wait for second IDLE iteration, then give the loop a moment to enter
	// waitForEvents (the callback fires before the loop checks quiescence).
	<-secondIterDone
	time.Sleep(50 * time.Millisecond)
	otherAgent := types.MustActorID("actor_00000000000000000000000000000002")
	mockEv := createMockEvent(t, s, otherAgent)
	eventBus.Publish(mockEv)

	result := <-done
	if result.Reason != StopTaskDone {
		t.Errorf("reason = %s, want %s (detail: %s)", result.Reason, StopTaskDone, result.Detail)
	}
}

func TestRunConcurrent(t *testing.T) {
	provider1 := newMockProvider("/signal {\"signal\": \"TASK_DONE\"}")
	provider2 := newMockProvider("/signal {\"signal\": \"TASK_DONE\"}")

	agent1 := testAgentWithRole(t, provider1, roles.RoleBuilder, "builder")
	agent2 := testAgentWithRole(t, provider2, roles.RoleTester, "tester")

	results := RunConcurrent(context.Background(), []Config{
		{Agent: agent1, HumanID: humanID(), Budget: resources.BudgetConfig{MaxIterations: 5}},
		{Agent: agent2, HumanID: humanID(), Budget: resources.BudgetConfig{MaxIterations: 5}},
	})

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	for _, ar := range results {
		if ar.Result.Reason != StopTaskDone {
			t.Errorf("%s (%s): reason = %s, want %s", ar.Role, ar.Name, ar.Result.Reason, StopTaskDone)
		}
	}
}

func TestNewLoopRequiresAgent(t *testing.T) {
	_, err := New(Config{Budget: resources.BudgetConfig{MaxIterations: 1}})
	if err == nil {
		t.Fatal("expected error for nil agent")
	}
}

// ════════════════════════════════════════════════════════════════════════
// Additional helpers
// ════════════════════════════════════════════════════════════════════════

// agentCounter generates unique actor IDs for test agents.
var agentCounter uint32

func testAgentWithRole(t *testing.T, provider intelligence.Provider, role roles.Role, name string) *roles.Agent {
	t.Helper()
	n := atomic.AddUint32(&agentCounter, 1)
	agentID := types.MustActorID(fmt.Sprintf("actor_%032d", n))
	humanID := types.MustActorID("actor_00000000000000000000000000000099")

	rt, err := intelligence.NewRuntime(context.Background(), intelligence.RuntimeConfig{
		AgentID:  agentID,
		Provider: provider,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = rt.Boot("ai", "mock-model", "standard", []string{"test"}, types.MustDomainScope("test"), humanID)
	if err != nil {
		t.Fatal(err)
	}

	return &roles.Agent{
		Runtime: rt,
		Role:    role,
		Name:    name,
	}
}

// createMockEvent creates a minimal event for bus testing.
func createMockEvent(t *testing.T, _ store.Store, source types.ActorID) event.Event {
	t.Helper()
	// Create a lightweight event using the event builder pattern.
	// For bus testing, we need a valid event that passes pattern matching.
	registry := event.DefaultRegistry()
	factory := event.NewEventFactory(registry)
	memStore := store.NewInMemoryStore()

	// Bootstrap the temp store.
	bsFactory := event.NewBootstrapFactory(registry)
	signer := &testSigner{}
	bootstrap, err := bsFactory.Init(source, signer)
	if err != nil {
		t.Fatal(err)
	}
	stored, err := memStore.Append(bootstrap)
	if err != nil {
		t.Fatal(err)
	}

	convID, _ := types.NewConversationID("conv_test_0000000000000000000000001")
	ev, err := factory.Create(
		event.EventTypeAgentActed,
		source,
		event.AgentActedContent{AgentID: source, Action: "test_event", Target: "loop"},
		[]types.EventID{stored.ID()},
		convID,
		memStore,
		signer,
	)
	if err != nil {
		t.Fatal(err)
	}
	return ev
}

func TestParseSignal(t *testing.T) {
	tests := []struct {
		response string
		want     *Signal
	}{
		// Valid /signal lines.
		{`/signal {"signal": "IDLE"}`, &Signal{Signal: "IDLE"}},
		{`/signal {"signal": "TASK_DONE"}`, &Signal{Signal: "TASK_DONE"}},
		{`/signal {"signal": "HALT", "reason": "policy violation"}`, &Signal{Signal: "HALT", Reason: "policy violation"}},
		{`/signal {"signal": "ESCALATE", "reason": "needs human"}`, &Signal{Signal: "ESCALATE", Reason: "needs human"}},
		// Signal at end of multi-line response.
		{"I checked everything.\n/signal {\"signal\": \"IDLE\"}", &Signal{Signal: "IDLE"}},
		// Case insensitive — lowercased signal normalised to upper.
		{`/signal {"signal": "idle"}`, &Signal{Signal: "IDLE"}},
		// Indented /signal line.
		{"  /signal {\"signal\": \"IDLE\"}", &Signal{Signal: "IDLE"}},

		// Invalid — no /signal line.
		{"Just some text", nil},
		{"HALT: not a signal line", nil},
		{"", nil},
		// Invalid JSON.
		{"/signal not-json", nil},
		// Empty signal field.
		{`/signal {"signal": ""}`, nil},
	}
	for _, tt := range tests {
		got := parseSignal(tt.response)
		if tt.want == nil {
			if got != nil {
				t.Errorf("parseSignal(%q) = %+v, want nil", tt.response, got)
			}
		} else {
			if got == nil {
				t.Errorf("parseSignal(%q) = nil, want %+v", tt.response, tt.want)
			} else if got.Signal != tt.want.Signal || got.Reason != tt.want.Reason {
				t.Errorf("parseSignal(%q) = %+v, want %+v", tt.response, got, tt.want)
			}
		}
	}
}

func TestContainsSignalFallback(t *testing.T) {
	// Text-based fallback — used when LLM doesn't emit /signal JSON.
	tests := []struct {
		response string
		signal   string
		want     bool
	}{
		// Positive — signal at start of line.
		{"HALT: integrity violation", "HALT", true},
		{"HALT", "HALT", true},
		{"TASK_DONE", "TASK_DONE", true},
		{"ESCALATE: authority needed", "ESCALATE", true},

		// Negative — signal embedded in prose.
		{"No HALT required", "HALT", false},
		{"brought to a halt", "HALT", false},
		{"we should not ESCALATE", "ESCALATE", false},
		{"", "HALT", false},
	}
	for _, tt := range tests {
		got := ContainsSignal(tt.response, tt.signal)
		if got != tt.want {
			t.Errorf("ContainsSignal(%q, %q) = %v, want %v", tt.response, tt.signal, got, tt.want)
		}
	}
}

type testSigner struct{}

func (s *testSigner) Sign(data []byte) (types.Signature, error) {
	return types.NewSignature(make([]byte, 64))
}
