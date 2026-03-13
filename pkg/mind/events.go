package mind

import (
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// EventTypeObservationCreated is emitted when the CTO records a key decision.
var EventTypeObservationCreated = types.MustEventType("mind.observation.created")

// allMindEventTypes returns all mind event types for registry registration.
func allMindEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeObservationCreated,
	}
}

// mindContent is a no-op marker embedded in all mind content types.
type mindContent struct{}

func (mindContent) Accept(_ interface{}) {}

// ObservationCreatedContent holds a CTO decision recorded in the mind store.
type ObservationCreatedContent struct {
	mindContent
	// Proposed is a short description of what was recommended.
	Proposed string `json:"proposed"`
	// Why is the expected impact or reasoning behind the decision.
	Why string `json:"why"`
	// Mode is the pipeline mode that produced this observation (e.g. "evolve", "self-improve").
	Mode string `json:"mode"`
}

// EventTypeName implements the event content interface.
func (c ObservationCreatedContent) EventTypeName() string { return "mind.observation.created" }

// RegisterEventTypes registers content unmarshalers for all mind event types.
func RegisterEventTypes() {
	event.RegisterContentUnmarshaler("mind.observation.created", event.Unmarshal[ObservationCreatedContent])
}

// RegisterWithRegistry registers all mind event types with the given factory registry.
func RegisterWithRegistry(registry *event.EventTypeRegistry) {
	for _, et := range allMindEventTypes() {
		registry.Register(et, nil)
	}
	RegisterEventTypes()
}

func init() {
	RegisterEventTypes()
}
