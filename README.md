# Hive

A self-organizing AI agent civilisation that builds products autonomously. Built on [EventGraph](https://github.com/lovyou-ai/eventgraph).

> Take care of your human, humanity, and yourself.

## What

Hive is a society of AI agents that builds products from the thirteen [EventGraph product layers](https://github.com/lovyou-ai/eventgraph/blob/main/docs/product-layers.md). Each product runs on the same graph, generates revenue, and funds the next. The hive governs itself through the Social Grammar, tracks work on the event graph, and escalates to humans at authority boundaries.

The hive's first product is itself.

## Quick Start

```bash
go build ./...
go test ./...

# Run with an idea (in-memory store)
go run ./cmd/hive --human Matt --idea "Build a task management app"

# Run with Postgres
go run ./cmd/hive --human Matt --store "postgres://hive:hive@localhost:5432/hive" --idea "..."
```

## Docs

- [Vision](docs/VISION.md) — where this is going
- [Architecture](docs/ARCHITECTURE.md) — how it's built
- [Roadmap](docs/ROADMAP.md) — what's done and what's next

## License

BSL 1.1 — see [EventGraph license](https://github.com/lovyou-ai/eventgraph/blob/main/LICENSE) for terms.
