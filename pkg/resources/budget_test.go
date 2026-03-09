package resources

import (
	"errors"
	"testing"
	"time"
)

func TestBudgetWithinLimits(t *testing.T) {
	b := NewBudget(BudgetConfig{
		MaxTokens:     1000,
		MaxCostUSD:    5.0,
		MaxIterations: 10,
	})

	b.Record(100, 0.50)

	if err := b.Check(); err != nil {
		t.Fatalf("expected within budget, got: %v", err)
	}

	snap := b.Snapshot()
	if snap.TokensUsed != 100 {
		t.Errorf("tokens = %d, want 100", snap.TokensUsed)
	}
	if snap.CostUSD != 0.50 {
		t.Errorf("cost = %.2f, want 0.50", snap.CostUSD)
	}
	if snap.Iterations != 1 {
		t.Errorf("iterations = %d, want 1", snap.Iterations)
	}
}

func TestBudgetTokensExceeded(t *testing.T) {
	b := NewBudget(BudgetConfig{MaxTokens: 100})
	b.Record(100, 0)

	err := b.Check()
	if err == nil {
		t.Fatal("expected budget exceeded")
	}
	var budgetErr *BudgetExceededError
	if !errors.As(err, &budgetErr) {
		t.Fatalf("expected BudgetExceededError, got %T", err)
	}
	if budgetErr.Resource != "tokens" {
		t.Errorf("resource = %q, want %q", budgetErr.Resource, "tokens")
	}
}

func TestBudgetCostExceeded(t *testing.T) {
	b := NewBudget(BudgetConfig{MaxCostUSD: 1.0})
	b.Record(0, 1.0)

	err := b.Check()
	if err == nil {
		t.Fatal("expected budget exceeded")
	}
	var budgetErr *BudgetExceededError
	if !errors.As(err, &budgetErr) {
		t.Fatalf("expected BudgetExceededError, got %T", err)
	}
	if budgetErr.Resource != "cost" {
		t.Errorf("resource = %q, want %q", budgetErr.Resource, "cost")
	}
}

func TestBudgetIterationsExceeded(t *testing.T) {
	b := NewBudget(BudgetConfig{MaxIterations: 3})
	b.Record(10, 0.1)
	b.Record(10, 0.1)
	b.Record(10, 0.1)

	err := b.Check()
	if err == nil {
		t.Fatal("expected budget exceeded")
	}
	var budgetErr *BudgetExceededError
	if !errors.As(err, &budgetErr) {
		t.Fatalf("expected BudgetExceededError, got %T", err)
	}
	if budgetErr.Resource != "iterations" {
		t.Errorf("resource = %q, want %q", budgetErr.Resource, "iterations")
	}
}

func TestBudgetDurationExceeded(t *testing.T) {
	// Set start time 1 second in the past — no time.Sleep needed.
	pastStart := time.Now().Add(-1 * time.Second)
	b := NewBudgetForTest(BudgetConfig{MaxDuration: 500 * time.Millisecond}, pastStart)

	err := b.Check()
	if err == nil {
		t.Fatal("expected budget exceeded")
	}
	var budgetErr *BudgetExceededError
	if !errors.As(err, &budgetErr) {
		t.Fatalf("expected BudgetExceededError, got %T", err)
	}
	if budgetErr.Resource != "duration" {
		t.Errorf("resource = %q, want %q", budgetErr.Resource, "duration")
	}
}

func TestBudgetUnlimited(t *testing.T) {
	b := NewBudget(BudgetConfig{}) // all zeros = unlimited
	b.Record(999999, 999.99)

	if err := b.Check(); err != nil {
		t.Fatalf("unlimited budget should not exceed: %v", err)
	}
}

func TestBudgetConcurrentRecords(t *testing.T) {
	b := NewBudget(BudgetConfig{MaxTokens: 100000})

	done := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func() {
			b.Record(10, 0.01)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 100; i++ {
		<-done
	}

	snap := b.Snapshot()
	if snap.TokensUsed != 1000 {
		t.Errorf("tokens = %d, want 1000", snap.TokensUsed)
	}
	if snap.Iterations != 100 {
		t.Errorf("iterations = %d, want 100", snap.Iterations)
	}
}
