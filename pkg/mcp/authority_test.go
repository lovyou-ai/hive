package mcp

import (
	"strings"
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/actor"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/trust"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// authorityTestSetup creates a minimal actor store with a registered AI agent
// and bootstraps an in-memory event store. The returned EventID can be used as
// the reason argument for Suspend / Memorial calls.
func authorityTestSetup(t *testing.T) (actor.IActorStore, types.ActorID, types.EventID) {
	t.Helper()

	actors := actor.NewInMemoryActorStore()
	s := store.NewInMemoryStore()

	// Keys must be exactly 32 bytes (Ed25519 public key size).
	agentPub, err := types.NewPublicKey([]byte("chk-agent-0000000000000000000000"))
	if err != nil {
		t.Fatal(err)
	}
	agentActor, err := actors.Register(agentPub, "CheckerTestAgent", event.ActorTypeAI)
	if err != nil {
		t.Fatal(err)
	}
	agentID := agentActor.ID()

	// Register a human actor to bootstrap the graph (required for Init).
	humanPub, err := types.NewPublicKey([]byte("chk-human-0000000000000000000000"))
	if err != nil {
		t.Fatal(err)
	}
	humanActor, err := actors.Register(humanPub, "CheckerTestHuman", event.ActorTypeHuman)
	if err != nil {
		t.Fatal(err)
	}

	registry := event.DefaultRegistry()
	bsFactory := event.NewBootstrapFactory(registry)
	signer := &testSigner{} // defined in server_test.go (same package)
	bootstrap, err := bsFactory.Init(humanActor.ID(), signer)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Append(bootstrap); err != nil {
		t.Fatal(err)
	}

	head, err := s.Head()
	if err != nil {
		t.Fatal(err)
	}
	reasonID := head.Unwrap().ID()

	return actors, agentID, reasonID
}

// ════════════════════════════════════════════════════════════════════════
// AuthorityChecker.CheckWrite
// ════════════════════════════════════════════════════════════════════════

func TestAuthorityChecker_CheckWrite_ZeroMinTrust(t *testing.T) {
	// minTrust=0 bypasses all checks — any agent is allowed to write.
	actors, agentID, _ := authorityTestSetup(t)
	checker := NewAuthorityChecker(actors, trust.NewDefaultTrustModel(), agentID, 0)
	if err := checker.CheckWrite("write_code:go"); err != nil {
		t.Errorf("zero minTrust should allow all writes, got error: %v", err)
	}
}

func TestAuthorityChecker_CheckWrite_SuspendedDenied(t *testing.T) {
	// Suspended agent is denied with "suspended" in the error message.
	actors, agentID, reasonID := authorityTestSetup(t)
	if _, err := actors.Suspend(agentID, reasonID); err != nil {
		t.Fatal(err)
	}
	checker := NewAuthorityChecker(actors, trust.NewDefaultTrustModel(), agentID, 0.1)
	err := checker.CheckWrite("write_code:go")
	if err == nil {
		t.Fatal("expected error for suspended agent, got nil")
	}
	if !strings.Contains(err.Error(), "suspended") {
		t.Errorf("error should mention 'suspended', got: %v", err)
	}
}

func TestAuthorityChecker_CheckWrite_MemorialDenied(t *testing.T) {
	// Memorial agent is denied with "memorial" in the error message.
	actors, agentID, reasonID := authorityTestSetup(t)
	if _, err := actors.Memorial(agentID, reasonID); err != nil {
		t.Fatal(err)
	}
	checker := NewAuthorityChecker(actors, trust.NewDefaultTrustModel(), agentID, 0.1)
	err := checker.CheckWrite("write_code:go")
	if err == nil {
		t.Fatal("expected error for memorial agent, got nil")
	}
	if !strings.Contains(err.Error(), "memorial") {
		t.Errorf("error should mention 'memorial', got: %v", err)
	}
}

func TestAuthorityChecker_CheckWrite_BelowThresholdDenied(t *testing.T) {
	// Active agent with trust 0.0 (default) is denied when minTrust=0.5.
	actors, agentID, _ := authorityTestSetup(t)
	checker := NewAuthorityChecker(actors, trust.NewDefaultTrustModel(), agentID, 0.5)
	err := checker.CheckWrite("write_code:go")
	if err == nil {
		t.Fatal("expected error for agent below trust threshold, got nil")
	}
	if !strings.Contains(err.Error(), "authority denied") {
		t.Errorf("error should mention 'authority denied', got: %v", err)
	}
}

func TestAuthorityChecker_CheckWrite_AboveThresholdAllowed(t *testing.T) {
	// Active agent whose trust (0.5) meets minTrust (0.3) is allowed.
	// Use a custom trust model with InitialTrust=0.5 so a fresh agent passes.
	actors, agentID, _ := authorityTestSetup(t)
	trustModel := trust.NewDefaultTrustModelWithConfig(trust.DefaultConfig{
		InitialTrust: types.MustScore(0.5),
	})
	checker := NewAuthorityChecker(actors, trustModel, agentID, 0.3)
	if err := checker.CheckWrite("write_code:go"); err != nil {
		t.Errorf("agent with trust above threshold should be allowed, got error: %v", err)
	}
}

// ════════════════════════════════════════════════════════════════════════
// WrapWriteHandler
// ════════════════════════════════════════════════════════════════════════

func TestWrapWriteHandler_NilCheckerPassesThrough(t *testing.T) {
	// nil checker: handler is always called, no authority gate applied.
	called := false
	handler := func(args map[string]any) (ToolCallResult, error) {
		called = true
		return TextResult("ok"), nil
	}

	wrapped := WrapWriteHandler(nil, "write_code:go", handler)
	result, err := wrapped(map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("nil checker should call handler unconditionally")
	}
	if result.IsError {
		t.Error("nil checker should not produce an error result")
	}
}

func TestWrapWriteHandler_DeniedDoesNotCallHandler(t *testing.T) {
	// Non-nil checker that denies: ErrorResult returned, handler never called.
	actors, agentID, _ := authorityTestSetup(t)
	// Fresh agent has trust 0.0 < minTrust 0.5 → denied.
	checker := NewAuthorityChecker(actors, trust.NewDefaultTrustModel(), agentID, 0.5)

	called := false
	handler := func(args map[string]any) (ToolCallResult, error) {
		called = true
		return TextResult("should not reach"), nil
	}

	wrapped := WrapWriteHandler(checker, "write_code:go", handler)
	result, err := wrapped(map[string]any{})
	if err != nil {
		t.Fatalf("WrapWriteHandler should return ErrorResult (not error): %v", err)
	}
	if called {
		t.Error("handler should not be called when checker denies")
	}
	if !result.IsError {
		t.Error("denied request should produce IsError=true result")
	}
	if len(result.Content) == 0 || result.Content[0].Text == "" {
		t.Error("error result should contain denial message")
	}
}
