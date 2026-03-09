// Command hive runs the product factory.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/lovyou-ai/eventgraph/go/pkg/store"

	"github.com/lovyou-ai/hive/pkg/pipeline"
)

func main() {
	name := flag.String("name", "", "Product name (derived by CTO if not provided)")
	idea := flag.String("idea", "", "Product idea (natural language description)")
	url := flag.String("url", "", "URL to research for product idea")
	spec := flag.String("spec", "", "Path to Code Graph spec file")
	workdir := flag.String("workdir", "products", "Directory for generated products")
	flag.Parse()

	if *idea == "" && *url == "" && *spec == "" {
		fmt.Fprintln(os.Stderr, "Usage: hive [--name product-name] --idea 'description' | --url 'https://...' | --spec path/to/spec.cg")
		os.Exit(1)
	}

	ctx := context.Background()

	// Create shared event graph
	s := store.NewInMemoryStore()

	// Create and run pipeline — uses Claude CLI (Max plan, flat rate)
	// CTO, Architect, Reviewer, Guardian → Opus (high judgment)
	// Builder, Tester, Integrator, Researcher → Sonnet (execution)
	p, err := pipeline.New(ctx, pipeline.Config{
		Store:   s,
		WorkDir: *workdir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "pipeline: %v\n", err)
		os.Exit(1)
	}

	// Run the pipeline — Guardian checks run automatically after each phase
	if err := p.Run(ctx, pipeline.ProductInput{
		Name:        *name,
		URL:         *url,
		Description: *idea,
		SpecFile:    *spec,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "pipeline failed: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	count, _ := s.Count()
	fmt.Printf("\nEvents recorded: %d\n", count)
	fmt.Printf("Agents active: %d\n", len(p.Agents()))
}
