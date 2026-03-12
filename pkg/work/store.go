package work

import (
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// Task represents a work item derived from a work.task.created event.
type Task struct {
	ID          types.EventID
	Title       string
	Description string
	CreatedBy   types.ActorID
}

// TaskStore creates and queries tasks as auditable events on the shared graph.
type TaskStore struct {
	store   store.Store
	factory *event.EventFactory
	signer  event.Signer
}

// NewTaskStore creates a new TaskStore backed by the given event store.
func NewTaskStore(s store.Store, factory *event.EventFactory, signer event.Signer) *TaskStore {
	return &TaskStore{store: s, factory: factory, signer: signer}
}

// Create records a work.task.created event on the graph and returns the task.
// The caller must supply at least one cause (typically the current chain head).
func (ts *TaskStore) Create(
	source types.ActorID,
	title, description string,
	causes []types.EventID,
	convID types.ConversationID,
) (Task, error) {
	if title == "" {
		return Task{}, fmt.Errorf("title is required")
	}
	content := TaskCreatedContent{
		Title:       title,
		Description: description,
		CreatedBy:   source,
	}
	ev, err := ts.factory.Create(EventTypeTaskCreated, source, content, causes, convID, ts.store, ts.signer)
	if err != nil {
		return Task{}, fmt.Errorf("create task event: %w", err)
	}
	stored, err := ts.store.Append(ev)
	if err != nil {
		return Task{}, fmt.Errorf("append task event: %w", err)
	}
	return Task{
		ID:          stored.ID(),
		Title:       title,
		Description: description,
		CreatedBy:   source,
	}, nil
}

// List returns up to limit work.task.created events as Tasks.
func (ts *TaskStore) List(limit int) ([]Task, error) {
	if limit <= 0 {
		limit = 20
	}
	page, err := ts.store.ByType(EventTypeTaskCreated, limit, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	tasks := make([]Task, 0, len(page.Items()))
	for _, ev := range page.Items() {
		c, ok := ev.Content().(TaskCreatedContent)
		if !ok {
			continue
		}
		tasks = append(tasks, Task{
			ID:          ev.ID(),
			Title:       c.Title,
			Description: c.Description,
			CreatedBy:   c.CreatedBy,
		})
	}
	return tasks, nil
}

// Assign records a work.task.assigned event on the graph.
// source is the actor performing the assignment (may equal assignee for self-assignment).
func (ts *TaskStore) Assign(
	source types.ActorID,
	taskID types.EventID,
	assignee types.ActorID,
	causes []types.EventID,
	convID types.ConversationID,
) error {
	content := TaskAssignedContent{
		TaskID:     taskID,
		AssignedTo: assignee,
		AssignedBy: source,
	}
	ev, err := ts.factory.Create(EventTypeTaskAssigned, source, content, causes, convID, ts.store, ts.signer)
	if err != nil {
		return fmt.Errorf("create assign event: %w", err)
	}
	if _, err := ts.store.Append(ev); err != nil {
		return fmt.Errorf("append assign event: %w", err)
	}
	return nil
}

// Complete records a work.task.completed event on the graph.
// source is the actor completing the task (typically the assignee).
func (ts *TaskStore) Complete(
	source types.ActorID,
	taskID types.EventID,
	summary string,
	causes []types.EventID,
	convID types.ConversationID,
) error {
	content := TaskCompletedContent{
		TaskID:      taskID,
		CompletedBy: source,
		Summary:     summary,
	}
	ev, err := ts.factory.Create(EventTypeTaskCompleted, source, content, causes, convID, ts.store, ts.signer)
	if err != nil {
		return fmt.Errorf("create complete event: %w", err)
	}
	if _, err := ts.store.Append(ev); err != nil {
		return fmt.Errorf("append complete event: %w", err)
	}
	return nil
}

// GetByAssignee returns tasks assigned to the given actor.
// It scans work.task.assigned events and joins each to its work.task.created event.
func (ts *TaskStore) GetByAssignee(assignee types.ActorID) ([]Task, error) {
	// Fetch all assigned events; no SQL join available in-memory so filter in code.
	page, err := ts.store.ByType(EventTypeTaskAssigned, 1000, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("fetch assigned events: %w", err)
	}
	var tasks []Task
	for _, ev := range page.Items() {
		c, ok := ev.Content().(TaskAssignedContent)
		if !ok || c.AssignedTo != assignee {
			continue
		}
		created, err := ts.store.Get(c.TaskID)
		if err != nil {
			continue // task event missing — skip
		}
		cc, ok := created.Content().(TaskCreatedContent)
		if !ok {
			continue
		}
		tasks = append(tasks, Task{
			ID:          created.ID(),
			Title:       cc.Title,
			Description: cc.Description,
			CreatedBy:   cc.CreatedBy,
		})
	}
	return tasks, nil
}
