// Package agent provides the unified Agent type for the hive.
//
// Every agent in the hive — Mind, CTO, Guardian, Builder, etc. — is an Agent.
// The Agent wraps EventGraph's AgentRuntime and adds:
//   - Operational state machine (driven on every operation)
//   - True causality tracking (each event caused by the previous, not store.Head())
//   - Trust accumulation hooks
//   - Budget event emission
//   - Proper retirement lifecycle
//   - Communication through the shared graph
//
// Construction:
//
//	a, err := agent.New(ctx, agent.Config{
//	    Role:     roles.RoleCTO,
//	    Name:     "CTO",
//	    Graph:    g,
//	    Provider: provider,
//	})
package agent

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/lovyou-ai/eventgraph/go/pkg/agent"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/graph"
	"github.com/lovyou-ai/eventgraph/go/pkg/intelligence"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/roles"
)

// Agent is the unified agent type for the hive.
// All agent operations go through this type, which ensures proper state
// transitions, event emission, causality tracking, and trust accumulation.
type Agent struct {
	runtime *intelligence.AgentRuntime
	role    roles.Role
	name    string
	graph   *graph.Graph
	convID  types.ConversationID
	state   agent.OperationalState
	signer  event.Signer

	// lastEvent tracks the most recent event this agent emitted.
	// Used for true causality — each event is caused by the previous,
	// not by whatever happens to be at store.Head().
	lastEvent types.EventID

	mu sync.Mutex
}

// Config holds the minimal configuration needed to create an Agent.
// Everything else is derived: ActorID from the actor store, signing key
// from deterministic derivation, conversation ID per-operation.
type Config struct {
	Role     roles.Role
	Name     string
	Graph    *graph.Graph
	Provider intelligence.Provider

	// ConversationID overrides the default agent-scoped conversation ID.
	// When set, the agent's events join this conversation thread.
	// Used by the pipeline to unify all agents under one conversation.
	ConversationID types.ConversationID
}

// New creates and boots a new Agent.
//
// The agent is registered in the actor store, its signing key is derived
// deterministically from its name, and the full boot sequence is emitted
// (identity, soul, model, authority, state → Idle).
func New(ctx context.Context, cfg Config) (*Agent, error) {
	if cfg.Graph == nil {
		return nil, fmt.Errorf("agent: Graph is required")
	}
	if cfg.Provider == nil {
		return nil, fmt.Errorf("agent: Provider is required")
	}

	// Derive deterministic signing key from agent name.
	seed := sha256.Sum256([]byte("agent:" + cfg.Name))
	priv := ed25519.NewKeyFromSeed(seed[:])
	pub := priv.Public().(ed25519.PublicKey)
	signer := &deterministicSigner{key: priv}

	// Register in actor store.
	pk, err := types.NewPublicKey([]byte(pub))
	if err != nil {
		return nil, fmt.Errorf("agent: public key: %w", err)
	}

	registered, err := cfg.Graph.ActorStore().Register(pk, cfg.Name, event.ActorTypeAI)
	if err != nil {
		// Actor may already exist from a previous run — look it up.
		existing, getErr := cfg.Graph.ActorStore().GetByPublicKey(pk)
		if getErr != nil {
			return nil, fmt.Errorf("agent: register %s: %w", cfg.Name, err)
		}
		registered = existing
	}

	// Create the AgentRuntime.
	// Use configured conversation ID if provided; otherwise derive from agent ID.
	convID := cfg.ConversationID
	if convID.Value() == "" {
		convID = types.MustConversationID(fmt.Sprintf("agent_%s", registered.ID().Value()))
	}
	rt, err := intelligence.NewRuntime(ctx, intelligence.RuntimeConfig{
		AgentID:        registered.ID(),
		Provider:       cfg.Provider,
		Store:          cfg.Graph.Store(),
		ConversationID: convID,
	})
	if err != nil {
		return nil, fmt.Errorf("agent: runtime: %w", err)
	}

	a := &Agent{
		runtime: rt,
		role:    cfg.Role,
		name:    cfg.Name,
		graph:   cfg.Graph,
		convID:  convID,
		state:   agent.StateIdle,
		signer:  signer,
	}

	// Boot: emit lifecycle events on the graph.
	if err := a.boot(pk, cfg); err != nil {
		return nil, fmt.Errorf("agent: boot %s: %w", cfg.Name, err)
	}

	return a, nil
}

// boot emits the full agent lifecycle boot sequence.
func (a *Agent) boot(pk types.PublicKey, cfg Config) error {
	bootEvents := agent.BootEvents(
		a.runtime.ID(),
		pk,
		string(cfg.Role),
		roles.PreferredModel(cfg.Role),
		"opus", // cost tier
		roles.SoulValues(cfg.Role),
		types.MustDomainScope("hive"),
		a.runtime.ID(), // self-grant at boot; human approves via Spawner
		false,          // withIdentity=false, Spawner handles identity
	)

	for _, content := range bootEvents {
		ev, err := a.record(content.EventTypeName(), content)
		if err != nil {
			return fmt.Errorf("boot event %s: %w", content.EventTypeName(), err)
		}
		a.lastEvent = ev.ID()
	}

	return nil
}

// record emits an event through the Graph facade (mutex-safe, bus-integrated).
func (a *Agent) record(eventTypeName string, content event.EventContent) (event.Event, error) {
	eventType := types.MustEventType(eventTypeName)

	var causes []types.EventID
	if !a.lastEvent.IsZero() {
		causes = []types.EventID{a.lastEvent}
	} else {
		// First event — use graph head as cause.
		head, err := a.graph.Store().Head()
		if err != nil {
			return event.Event{}, fmt.Errorf("get head: %w", err)
		}
		if head.IsSome() {
			causes = []types.EventID{head.Unwrap().ID()}
		}
	}

	return a.graph.Record(eventType, a.runtime.ID(), content, causes, a.convID, a.signer)
}

// --- Accessors ---

// ID returns the agent's actor ID.
func (a *Agent) ID() types.ActorID { return a.runtime.ID() }

// Role returns the agent's role.
func (a *Agent) Role() roles.Role { return a.role }

// Name returns the agent's display name.
func (a *Agent) Name() string { return a.name }

// State returns the agent's current operational state.
func (a *Agent) State() agent.OperationalState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state
}

// Runtime returns the underlying AgentRuntime for advanced use.
// Prefer the Agent's methods over direct Runtime access.
func (a *Agent) Runtime() *intelligence.AgentRuntime { return a.runtime }

// Provider returns the agent's intelligence provider.
func (a *Agent) Provider() intelligence.Provider { return a.runtime.Provider() }

// Graph returns the shared graph.
func (a *Agent) Graph() *graph.Graph { return a.graph }

// ConversationID returns the agent's conversation thread ID.
func (a *Agent) ConversationID() types.ConversationID { return a.convID }

// LastEvent returns the ID of the most recent event this agent emitted.
func (a *Agent) LastEvent() types.EventID {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.lastEvent
}

// deterministicSigner signs with a key derived from the agent name.
type deterministicSigner struct {
	key ed25519.PrivateKey
}

func (s *deterministicSigner) Sign(data []byte) (types.Signature, error) {
	sig := ed25519.Sign(s.key, data)
	return types.NewSignature(sig)
}
