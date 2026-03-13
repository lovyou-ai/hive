package mind

import (
	"fmt"
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// Observation is a recorded CTO decision: what was proposed and why.
type Observation struct {
	Proposed string
	Why      string
	Mode     string
}

// MindStore persists CTO observations as signed events on the event graph.
type MindStore struct {
	store   store.Store
	factory *event.EventFactory
	signer  event.Signer
}

// NewMindStore creates a MindStore backed by the given event graph infrastructure.
func NewMindStore(s store.Store, factory *event.EventFactory, signer event.Signer) *MindStore {
	return &MindStore{store: s, factory: factory, signer: signer}
}

// Save records a CTO observation as a mind.observation.created event.
func (ms *MindStore) Save(source types.ActorID, obs Observation, causes []types.EventID, convID types.ConversationID) error {
	content := ObservationCreatedContent{
		Proposed: obs.Proposed,
		Why:      obs.Why,
		Mode:     obs.Mode,
	}
	_, err := ms.factory.Create(EventTypeObservationCreated, source, content, causes, convID, ms.store, ms.signer)
	if err != nil {
		return fmt.Errorf("create observation event: %w", err)
	}
	return nil
}

// Recent returns the last n observations, newest first.
func (ms *MindStore) Recent(n int) ([]Observation, error) {
	page, err := ms.store.ByType(EventTypeObservationCreated, n, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("query observations: %w", err)
	}
	obs := make([]Observation, 0, len(page.Items()))
	for _, ev := range page.Items() {
		c, ok := ev.Content().(ObservationCreatedContent)
		if !ok {
			continue
		}
		obs = append(obs, Observation{
			Proposed: c.Proposed,
			Why:      c.Why,
			Mode:     c.Mode,
		})
	}
	return obs, nil
}

// FormatPriorDecisions formats a slice of observations for injection into CTO prompts.
// Returns an empty string if obs is empty.
func FormatPriorDecisions(obs []Observation) string {
	if len(obs) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\nPRIOR DECISIONS (from previous sessions — build on these, avoid repeating them):\n")
	for i, o := range obs {
		sb.WriteString(fmt.Sprintf("  %d. [%s] %s\n     Impact: %s\n", i+1, o.Mode, o.Proposed, o.Why))
	}
	return sb.String()
}
