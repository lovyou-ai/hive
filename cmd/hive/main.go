// Command hive runs hive agents.
//
// New mode (runner): one process per agent role, polls lovyou.ai.
//
//	go run ./cmd/hive --role builder --repo ../site --space hive
//
// Legacy mode (runtime): spawns all agents, coordinates via event graph.
//
//	go run ./cmd/hive --human Matt --idea "description"
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lovyou-ai/eventgraph/go/pkg/actor"
	"github.com/lovyou-ai/eventgraph/go/pkg/actor/pgactor"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/intelligence"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/store/pgstore"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/api"
	"github.com/lovyou-ai/hive/pkg/hive"
	"github.com/lovyou-ai/hive/pkg/runner"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Runner mode flags.
	role := flag.String("role", "", "Agent role (builder, scout, critic, monitor). Enables runner mode.")
	space := flag.String("space", "hive", "lovyou.ai space slug")
	apiBase := flag.String("api", "https://lovyou.ai", "lovyou.ai API base URL")
	budget := flag.Float64("budget", 10.0, "Daily budget in USD")
	agentID := flag.String("agent-id", "", "Agent's lovyou.ai user ID (filters task assignment)")
	oneShot := flag.Bool("one-shot", false, "Work one task then exit (for testing)")

	// Shared flags.
	repo := flag.String("repo", "", "Path to repo for Operate (default: current dir)")

	// Legacy runtime mode flags.
	human := flag.String("human", "", "Human operator name (legacy runtime mode)")
	idea := flag.String("idea", "", "Seed idea for agents (legacy runtime mode)")
	storeDSN := flag.String("store", "", "Store DSN (legacy runtime mode)")
	autoApprove := flag.Bool("yes", false, "Auto-approve authority (legacy runtime mode)")
	flag.Parse()

	if *role != "" {
		return runRunner(*role, *space, *apiBase, *repo, *budget, *agentID, *oneShot)
	}
	if *human != "" {
		return runLegacy(*human, *idea, *storeDSN, *autoApprove, *repo)
	}

	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  Runner mode:  hive --role builder --repo ../site [--space hive] [--budget 10]")
	fmt.Fprintln(os.Stderr, "  Legacy mode:  hive --human Matt --idea 'description' [--store postgres://...]")
	return fmt.Errorf("specify --role or --human")
}

// ─── Runner mode ─────────────────────────────────────────────────────

func runRunner(role, space, apiBase, repoPath string, budget float64, agentID string, oneShot bool) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Resolve API key.
	apiKey := os.Getenv("LOVYOU_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("LOVYOU_API_KEY required")
	}

	// Resolve repo path.
	if repoPath == "" {
		repoPath = "."
	}
	absRepo, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("resolve repo: %w", err)
	}

	// Find hive directory (for loading role prompts).
	hiveDir := findHiveDir()

	// Create intelligence provider.
	model := runner.ModelForRole(role)
	provider, err := intelligence.New(intelligence.Config{
		Provider:     "claude-cli",
		Model:        model,
		MaxBudgetUSD: budget,
		SystemPrompt: runner.LoadRolePrompt(hiveDir, role),
	})
	if err != nil {
		return fmt.Errorf("provider: %w", err)
	}

	// Create API client.
	client := api.New(apiBase, apiKey)

	// Load role prompt (also passed as instruction context, not just system prompt).
	rolePrompt := runner.LoadRolePrompt(hiveDir, role)

	// Create and run the runner.
	r := runner.New(runner.Config{
		Role:       role,
		AgentID:    agentID,
		SpaceSlug:  space,
		RepoPath:   absRepo,
		APIClient:  client,
		Provider:   provider,
		RolePrompt: rolePrompt,
		BudgetUSD:  budget,
		OneShot:    oneShot,
	})

	log.Printf("hive agent starting: role=%s model=%s space=%s repo=%s agent-id=%s one-shot=%v",
		role, model, space, absRepo, agentID, oneShot)
	return r.Run(ctx)
}

// findHiveDir returns the hive repo directory by walking up from cwd.
func findHiveDir() string {
	// Try cwd first.
	cwd, _ := os.Getwd()
	if _, err := os.Stat(filepath.Join(cwd, "agents")); err == nil {
		return cwd
	}
	// Try parent (if running from cmd/hive).
	parent := filepath.Dir(filepath.Dir(cwd))
	if _, err := os.Stat(filepath.Join(parent, "agents")); err == nil {
		return parent
	}
	return cwd
}

// ─── Legacy runtime mode ────────────────────────────────────────────

func runLegacy(humanName, idea, dsn string, autoApprove bool, repoPath string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	var pool *pgxpool.Pool
	if dsn != "" {
		fmt.Fprintf(os.Stderr, "Postgres: %s\n", dsn)
		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
		defer pool.Close()
	}

	s, err := openStore(ctx, pool)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	defer s.Close()

	actors, err := openActorStore(ctx, pool)
	if err != nil {
		return fmt.Errorf("actor store: %w", err)
	}

	if pool != nil {
		fmt.Fprintln(os.Stderr, "WARNING: CLI key derivation is insecure for persistent Postgres stores.")
	}
	humanID, err := registerHuman(actors, humanName)
	if err != nil {
		return fmt.Errorf("register human: %w", err)
	}

	if err := bootstrapGraph(s, humanID); err != nil {
		return fmt.Errorf("bootstrap graph: %w", err)
	}

	if repoPath == "" {
		repoPath = "."
	}

	rt, err := hive.New(ctx, hive.Config{
		Store:       s,
		Actors:      actors,
		HumanID:     humanID,
		AutoApprove: autoApprove,
		RepoPath:    repoPath,
	})
	if err != nil {
		return fmt.Errorf("runtime: %w", err)
	}

	for _, def := range hive.StarterAgents(humanName) {
		if err := rt.Register(def); err != nil {
			return fmt.Errorf("register %s: %w", def.Name, err)
		}
	}

	if err := rt.Run(ctx, idea); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	count, _ := s.Count()
	fmt.Fprintf(os.Stderr, "Events recorded: %d\n", count)
	return nil
}

// ─── Store helpers ───────────────────────────────────────────────────

func openStore(ctx context.Context, pool *pgxpool.Pool) (store.Store, error) {
	if pool == nil {
		fmt.Fprintln(os.Stderr, "Store: in-memory")
		return store.NewInMemoryStore(), nil
	}
	fmt.Fprintln(os.Stderr, "Store: postgres")
	return pgstore.NewPostgresStoreFromPool(ctx, pool)
}

func openActorStore(ctx context.Context, pool *pgxpool.Pool) (actor.IActorStore, error) {
	if pool == nil {
		fmt.Fprintln(os.Stderr, "Actor store: in-memory")
		return actor.NewInMemoryActorStore(), nil
	}
	fmt.Fprintln(os.Stderr, "Actor store: postgres")
	return pgactor.NewPostgresActorStoreFromPool(ctx, pool)
}

func registerHuman(actors actor.IActorStore, displayName string) (types.ActorID, error) {
	h := sha256.Sum256([]byte("human:" + displayName))
	priv := ed25519.NewKeyFromSeed(h[:])
	pub := priv.Public().(ed25519.PublicKey)

	pk, err := types.NewPublicKey([]byte(pub))
	if err != nil {
		return types.ActorID{}, fmt.Errorf("public key: %w", err)
	}

	a, err := actors.Register(pk, displayName, event.ActorTypeHuman)
	if err != nil {
		return types.ActorID{}, err
	}
	return a.ID(), nil
}

func bootstrapGraph(s store.Store, humanID types.ActorID) error {
	head, err := s.Head()
	if err != nil {
		return fmt.Errorf("check head: %w", err)
	}
	if head.IsSome() {
		return nil
	}

	fmt.Fprintln(os.Stderr, "Bootstrapping event graph...")
	registry := event.DefaultRegistry()
	bsFactory := event.NewBootstrapFactory(registry)

	signer := &bootstrapSigner{humanID: humanID}
	bootstrap, err := bsFactory.Init(humanID, signer)
	if err != nil {
		return fmt.Errorf("create genesis event: %w", err)
	}
	if _, err := s.Append(bootstrap); err != nil {
		return fmt.Errorf("append genesis event: %w", err)
	}
	fmt.Fprintln(os.Stderr, "Event graph bootstrapped.")
	return nil
}

type bootstrapSigner struct {
	humanID types.ActorID
}

func (b *bootstrapSigner) Sign(data []byte) (types.Signature, error) {
	h := sha256.Sum256([]byte("signer:" + b.humanID.Value()))
	priv := ed25519.NewKeyFromSeed(h[:])
	sig := ed25519.Sign(priv, data)
	return types.NewSignature(sig)
}
