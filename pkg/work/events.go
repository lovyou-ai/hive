package work

import (
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// Work Graph event types — Layer 1 of the thirteen-product roadmap.
var (
	EventTypeTaskCreated   = types.MustEventType("work.task.created")
	EventTypeTaskAssigned  = types.MustEventType("work.task.assigned")
	EventTypeTaskCompleted = types.MustEventType("work.task.completed")
)

// allWorkEventTypes returns all work event types for registration.
func allWorkEventTypes() []types.EventType {
	return []types.EventType{
		EventTypeTaskCreated, EventTypeTaskAssigned, EventTypeTaskCompleted,
	}
}

// workContent is embedded in all work content types. Work events use
// no-op Accept (same pattern as pipeline content) since they are hive-specific.
type workContent struct{}

func (workContent) Accept(event.EventContentVisitor) {}

// --- Content structs ---

// TaskCreatedContent is emitted when a new task is created.
type TaskCreatedContent struct {
	workContent
	Title       string        `json:"Title"`
	Description string        `json:"Description,omitempty"`
	CreatedBy   types.ActorID `json:"CreatedBy"`
}

func (c TaskCreatedContent) EventTypeName() string { return "work.task.created" }

// TaskAssignedContent is emitted when a task is assigned to an actor.
type TaskAssignedContent struct {
	workContent
	TaskID     types.EventID `json:"TaskID"`
	AssignedTo types.ActorID `json:"AssignedTo"`
	AssignedBy types.ActorID `json:"AssignedBy"`
}

func (c TaskAssignedContent) EventTypeName() string { return "work.task.assigned" }

// TaskCompletedContent is emitted when a task is completed.
type TaskCompletedContent struct {
	workContent
	TaskID      types.EventID `json:"TaskID"`
	CompletedBy types.ActorID `json:"CompletedBy"`
	Summary     string        `json:"Summary,omitempty"`
}

func (c TaskCompletedContent) EventTypeName() string { return "work.task.completed" }

// RegisterEventTypes registers work content unmarshalers for Postgres
// deserialization. Call this before querying work events from the store.
func RegisterEventTypes() {
	event.RegisterContentUnmarshaler("work.task.created", event.Unmarshal[TaskCreatedContent])
	event.RegisterContentUnmarshaler("work.task.assigned", event.Unmarshal[TaskAssignedContent])
	event.RegisterContentUnmarshaler("work.task.completed", event.Unmarshal[TaskCompletedContent])
}

// RegisterWithRegistry registers all work event types with the given registry
// and registers content unmarshalers for Postgres deserialization.
func RegisterWithRegistry(registry *event.EventTypeRegistry) {
	for _, et := range allWorkEventTypes() {
		registry.Register(et, nil)
	}
	RegisterEventTypes()
}
