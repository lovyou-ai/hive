package work_test

import (
	"testing"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
	"github.com/lovyou-ai/hive/pkg/work"
)

var (
	testActor = types.MustActorID("actor_00000000000000000000000000000001")
	testConv  = types.MustConversationID("conv_00000000000000000000000000000001")
)

type testSigner struct{}

func (testSigner) Sign(data []byte) (types.Signature, error) {
	sig := make([]byte, 64)
	copy(sig, data)
	return types.MustSignature(sig), nil
}

// setupStore bootstraps an in-memory graph and returns the store and genesis causes.
func setupStore(t *testing.T) (*store.InMemoryStore, []types.EventID) {
	t.Helper()
	s := store.NewInMemoryStore()
	registry := event.DefaultRegistry()
	work.RegisterWithRegistry(registry)
	bf := event.NewBootstrapFactory(registry)
	boot, err := bf.Init(testActor, testSigner{})
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	stored, err := s.Append(boot)
	if err != nil {
		t.Fatalf("append genesis: %v", err)
	}
	return s, []types.EventID{stored.ID()}
}

// newTaskStore creates a TaskStore against the given store.
func newTaskStore(t *testing.T, s *store.InMemoryStore) *work.TaskStore {
	t.Helper()
	registry := event.DefaultRegistry()
	work.RegisterWithRegistry(registry)
	factory := event.NewEventFactory(registry)
	return work.NewTaskStore(s, factory, testSigner{})
}

func TestTaskStore_Create(t *testing.T) {
	s, causes := setupStore(t)
	ts := newTaskStore(t, s)

	task, err := ts.Create(testActor, "Fix the auth bug", "login fails on mobile", causes, testConv)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if task.Title != "Fix the auth bug" {
		t.Errorf("Title = %q; want %q", task.Title, "Fix the auth bug")
	}
	if task.Description != "login fails on mobile" {
		t.Errorf("Description = %q; want %q", task.Description, "login fails on mobile")
	}
	if task.CreatedBy != testActor {
		t.Errorf("CreatedBy = %v; want %v", task.CreatedBy, testActor)
	}
}

func TestTaskStore_Create_RequiresTitle(t *testing.T) {
	s, causes := setupStore(t)
	ts := newTaskStore(t, s)

	_, err := ts.Create(testActor, "", "no title", causes, testConv)
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestTaskStore_List_Empty(t *testing.T) {
	s, _ := setupStore(t)
	ts := newTaskStore(t, s)

	tasks, err := ts.List(20)
	if err != nil {
		t.Fatalf("List (empty): %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks; got %d", len(tasks))
	}
}

func TestTaskStore_List(t *testing.T) {
	s, causes := setupStore(t)
	ts := newTaskStore(t, s)

	_, err := ts.Create(testActor, "Task Alpha", "", causes, testConv)
	if err != nil {
		t.Fatalf("Create Alpha: %v", err)
	}
	_, err = ts.Create(testActor, "Task Beta", "some detail", causes, testConv)
	if err != nil {
		t.Fatalf("Create Beta: %v", err)
	}

	tasks, err := ts.List(20)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks; got %d", len(tasks))
	}
}

func TestTaskStore_List_RespectsLimit(t *testing.T) {
	s, causes := setupStore(t)
	ts := newTaskStore(t, s)

	for i := 0; i < 5; i++ {
		if _, err := ts.Create(testActor, "Task", "", causes, testConv); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}

	tasks, err := ts.List(3)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks (limit=3); got %d", len(tasks))
	}
}
