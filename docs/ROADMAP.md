# Hive Roadmap

Derived using the derivation method across four dimensions:

| Dimension | Range |
|-----------|-------|
| **Dependency** | What blocks what — strict ordering where needed |
| **Value** | What delivers the most capability per effort |
| **Risk** | What can go wrong, what needs proving early |
| **Time** | Rough effort estimate (days, not calendar) |

## Key Decisions

1. **Self-modification: yes.** The hive can and should modify its own codebase. PRs to lovyou-ai/hive, reviewed by human. This is how it builds its own tools.
2. **One service.** lovyou.ai is one product that does everything — not microservices. Web first, mobile later. The CTO/CTO-agent decides architecture.
3. **High scrutiny initially.** Every action reviewed in detail by human. Authority model starts strict (everything is "Required" approval). Trust accumulates through verified work — supervision decreases as the hive proves itself.
4. **CLI first, daemon soon.** Keep CLI for stepping through and debugging. Architect the code so the same pipeline can run as a long-running daemon. CLI and daemon share the same packages.

## Where We Are

The hive can take a product idea, research it, design a Code Graph spec, generate multi-file code, review it, test it, and push it to a GitHub repo. All agents share one event graph. Guardian checks integrity after every phase. Store is configurable (in-memory or Postgres). Actor IDs come from the actor store.

**What works today:**
- CLI pipeline: idea → research → design → simplify → build → review → test → integrate
- Per-role intelligence (Opus for judgment + code gen, Sonnet for execution, Haiku for volume)
- Multi-file code generation with review/rebuild loop
- Product repos with git commits at each phase
- Guardian integrity checks with HALT capability
- Postgres event store (via eventgraph pgstore)
- Actor registration (in-memory, human bootstrap via CLI)
- Per-agent token tracking with full breakdown (input, output, cache read/write, cost)
- Pipeline halts on test failure (no silent promotion of broken code)
- Targeted pipeline mode: read existing code → understand → modify → review → test → PR
- Agentic builder: Claude CLI reads/writes files directly via Operate (no text parsing middleman)
- Branch workflow for existing repos (create branch, commit, open PR)
- CLI `--repo` flag for targeting existing codebases
- 11 roles defined with system prompts and soul values
- Complete documentation (12 docs, ~2,200 lines)

**What doesn't work yet:**
- No concept of project continuity — every run starts from scratch
- No client/project registry — the hive doesn't know who it's building for

---

## Milestone 0: Self-Improvement Loop

**Why first:** The hive's first product is itself. Before building for others, it must be able to modify existing code — its own code. This is the prerequisite for everything: adding features, fixing bugs, growing capabilities. Without this, the hive is a one-shot generator, not a civilisation.

**Dependency:** None — builds on what exists today.

**Architecture decisions (derived via derivation method):**
- **The filesystem is the source of truth for code.** The git workdir has the code. The state store and event graph record metadata and decisions. We don't store code in the database.
- **Clients are NOT actors.** They don't sign events or make decisions on the graph. Clients are lightweight metadata in the state store. When lovyou.ai gets auth (Milestone 5), client representatives become human actors. But the client org itself isn't a decision-maker.
- **No eventgraph changes needed.** The state store (`IStateStore`) handles project/client metadata. The event graph records decisions. ConversationID groups events per project/run.
- **Two pipeline modes:** Full pipeline (greenfield: research → design → build → review → test → integrate) and targeted mode (existing code: understand → modify → review → test → PR). Self-improvement and feature additions use targeted mode.
- **Branch workflow for existing repos.** Self-improvement creates branches and PRs, not direct commits to main. Guardian scrutiny is higher for self-modification.

| # | Task | Where | Effort | Status |
|---|------|-------|--------|--------|
| 0.1 | Context loader — read existing files from any directory into ProjectContext | `hive/pkg/pipeline/` | 0.5d | ✅ Done |
| 0.2 | Context-aware prompts — builder/reviewer receive existing code as context | `hive/pkg/pipeline/` | 1d | ✅ Done |
| 0.3 | Targeted pipeline mode — skip research/design for modifications to existing code | `hive/pkg/pipeline/` | 1d | ✅ Done |
| 0.4 | Branch workflow — create branch, commit changes, open PR (not push to main) | `hive/pkg/workspace/` | 0.5d | ✅ Done |
| 0.5 | CLI: `--repo <path>` flag to target existing repo | `hive/cmd/hive/` | 0.5d | ✅ Done |
| 0.6 | Test: hive adds a feature to an existing project (end-to-end) | manual test | 0.5d | ✅ Done |
| 0.7 | Test: hive fixes a bug in its own codebase (end-to-end) | manual test | 0.5d | ||

**Exit criteria:** The hive can read an existing codebase, understand what needs changing, make targeted modifications, verify tests pass, and open a PR. Tested on both external projects and the hive's own code.

**Future (post-self-improvement):**
- Project registry in state store (scope: `project`, key: project ID → {client, repo, language, status})
- Client registry in state store (scope: `client`, key: client ID → {name, notes})
- CLI: `--client`, `--project`, `hive projects list`
- Per-project budgets, per-client billing (Milestone 9)

---

## Milestone 1: Persistent Identity

**Why first:** Without persistent actors, the hive forgets who it is between runs. Nothing else matters if agents can't remember themselves. This is the identity right made real.

**Dependency:** None — pure foundation work in eventgraph.

| # | Task | Where | Effort |
|---|------|-------|--------|
| 1.1 | Postgres actor store (`pgactor` package) | `eventgraph/go/pkg/actor/pgactor/` | 2d |
| 1.2 | Migration SQL (actors, trust_scores, authority_levels tables) | `pgactor/migrations/` | 0.5d |
| 1.3 | Integration tests against Docker Postgres | `pgactor/pgactor_test.go` | 1d |
| 1.4 | Wire `--store` flag to select pgactor when Postgres URL provided | `hive/cmd/hive/main.go` | 0.5d |
| 1.5 | Docker Compose: Postgres + pgAdmin for local dev | `hive/docker-compose.yml` | 0.5d |
| 1.6 | Trust persistence (accumulation, decay, history queryable) | `pgactor/trust.go` | 1d |
| 1.7 | Verify: stop hive, restart, agents remember identity + trust | manual test | 0.5d |

**Exit criteria:** `go test ./pkg/actor/pgactor/` passes. Hive restarts and agents retain identity, trust, and authority.

---

## Milestone 2: Agent Tools (MCP Server)

**Why next:** Agents need hands — tools to read the graph, query actors, write files, emit events. Without MCP, agents can only respond to prompts. With MCP, they can act.

**Dependency:** Milestone 1 (agents need persistent identity to authenticate tool calls).

| # | Task | Where | Effort |
|---|------|-------|--------|
| 2.1 | MCP server binary (stdio transport, JSON-RPC 2.0) | `hive/cmd/mcp-server/main.go` | 1d |
| 2.2 | Read tools: `query_events`, `get_event`, `get_actor`, `list_actors`, `get_trust` | `cmd/mcp-server/tools.go` | 1.5d |
| 2.3 | Write tools: `emit_event`, `write_file`, `run_command` (sandboxed) | `cmd/mcp-server/handlers.go` | 1.5d |
| 2.4 | Self-awareness tools: `query_self`, `query_human` | `cmd/mcp-server/handlers.go` | 0.5d |
| 2.5 | Authority checks on write tools (agent must be authorized) | `cmd/mcp-server/auth.go` | 1d |
| 2.6 | Wire MCP config into claude-cli provider (`.mcp.json` generation) | `eventgraph/go/pkg/intelligence/` | 1d |
| 2.7 | Guardian monitoring of all tool calls (events on graph) | `cmd/mcp-server/audit.go` | 0.5d |
| 2.8 | Context injection: inject actor list, pending tasks, own identity before each prompt | `hive/pkg/pipeline/` | 1d |
| 2.9 | Integration test: agent uses MCP tools during reasoning | `hive/pkg/pipeline/mcp_test.go` | 1d |

**Exit criteria:** An agent can query the graph, emit events, and write files mid-reasoning via MCP tools. Guardian sees all tool calls.

---

## Milestone 3: Agent Spawning & Authority

**Why next:** The hive needs to create agents with proper authority checks. This is the governance foundation — nothing gets created without approval.

**Dependency:** Milestone 1 (persistent actors), Milestone 2 (MCP tools for spawn events).

| # | Task | Where | Effort |
|---|------|-------|--------|
| 3.1 | Spawn protocol: `agent.spawn.requested` → human approval → `agent.spawn.approved` → create | `hive/pkg/spawn/` | 1.5d |
| 3.2 | Authority model implementation (Required/Recommended/Notification levels) | `hive/pkg/authority/` | 1.5d |
| 3.3 | CLI approval surface (step-through: show request, prompt y/n) | `hive/cmd/hive/approve.go` | 1d |
| 3.4 | Trust gate enforcement (agent can't be spawned into role above its trust gate) | `pkg/spawn/trust.go` | 0.5d |
| 3.5 | Agent lifecycle events (created, active, suspended, retired, memorial) | `hive/pkg/spawn/lifecycle.go` | 1d |
| 3.6 | Guardian audit: extra scrutiny for agent creation events | `hive/pkg/roles/` (Guardian prompt) | 0.5d |
| 3.7 | Spawner role wiring: CTO requests → Spawner proposes → human approves | `hive/pkg/pipeline/` | 1d |

**Exit criteria:** `go run ./cmd/hive --human Matt --idea "..."` prompts human to approve each agent spawn. Events recorded on graph. Guardian watches.

---

## Milestone 4: Agentic Loop

**Why next:** Agents need sustained autonomy — not just single-prompt responses, but observe-reason-act-reflect cycles. This transforms the hive from a pipeline into a society.

**Dependency:** Milestone 2 (MCP tools), Milestone 3 (authority for actions).

| # | Task | Where | Effort | Status |
|---|------|-------|--------|--------|
| 4.1 | Agentic loop runner (observe → kick Claude → check stopping) | `hive/pkg/loop/loop.go` | 1.5d | ✅ Done |
| 4.2 | Observation window (configurable: last N events, events since last check) | `hive/pkg/loop/loop.go` | 0.5d | ✅ Done |
| 4.3 | Stopping conditions: quiescence, escalation, HALT, budget limit | `hive/pkg/loop/loop.go` | 0.5d | ✅ Done |
| 4.4 | Budget enforcement per loop iteration (max tokens, max iterations, max cost) | `hive/pkg/resources/budget.go` | 1d | ✅ Done |
| 4.5 | Migrate pipeline from fixed sequence to graph-driven (CTO seeds, agents self-direct) | `hive/pkg/pipeline/` | 2d | ✅ Done |
| 4.6 | Concurrent agent loops (multiple agents running simultaneously, communicating via graph) | `hive/pkg/pipeline/` | 1.5d | ✅ Done |
| 4.7 | Real-time event notification (wire IBus into agent loops — agent notified of relevant events) | `hive/pkg/pipeline/` | 1d | ✅ Done |
| 4.8 | Integration test: two agents collaborate on a task via graph events | `hive/pkg/loop/loop_test.go` | 1d | ✅ Done |

**Exit criteria:** Multiple agents run concurrent loops, observe graph events, self-direct work, and collaborate without fixed orchestration. Budget limits enforced.

---

## Milestone 5: Web Service & Auth

**Why next:** The hive needs to be a long-running service accessible to humans. Auth is required for any external access. This unlocks the approval dashboard and eventually products.

**Dependency:** Milestone 1 (persistent identity), Milestone 3 (authority model for approval surface).

| # | Task | Where | Effort |
|---|------|-------|--------|
| 5.1 | HTTP daemon entry point (`hived`) sharing packages with CLI | `hive/cmd/hived/main.go` | 1d |
| 5.2 | Google OAuth2 flow → register human in actor store → issue session | `hive/pkg/auth/google.go` | 1.5d |
| 5.3 | Session management (JWT, cookie, token refresh) | `hive/pkg/auth/session.go` | 1d |
| 5.4 | Human approval dashboard (pending authority requests, approve/reject buttons) | `hive/pkg/web/approve.go` | 1.5d |
| 5.5 | Event stream view (real-time graph events, filterable by agent/type) | `hive/pkg/web/events.go` | 1d |
| 5.6 | Agent status view (who's running, trust levels, current task) | `hive/pkg/web/agents.go` | 1d |
| 5.7 | Static site: docs, blog, about page | `hive/pkg/web/static/` | 1d |
| 5.8 | Fly.io deployment config (Dockerfile, fly.toml, DATABASE_URL from Neon) | `hive/Dockerfile`, `hive/fly.toml` | 1d |
| 5.9 | Health check endpoint + readiness probe | `hive/pkg/web/health.go` | 0.5d |

**Exit criteria:** lovyou.ai serves a web dashboard. Human logs in with Google, sees pending approvals, can approve/reject. Event stream visible. Deployed on fly.io with Neon Postgres.

---

## Milestone 6: CI/CD & Deployment

**Why next:** Products need to be deployable, not just pushed to GitHub. The Integrator needs real deployment capability.

**Dependency:** Milestone 5 (web service, fly.io config).

| # | Task | Where | Effort |
|---|------|-------|--------|
| 6.1 | GitHub Actions: build + test on PR | `.github/workflows/ci.yml` | 0.5d |
| 6.2 | GitHub Actions: deploy to fly.io on merge to main | `.github/workflows/deploy.yml` | 0.5d |
| 6.3 | Product deployment: Integrator builds Docker image, deploys to fly.io | `hive/pkg/deploy/fly.go` | 1.5d |
| 6.4 | DNS routing: product.lovyou.ai subdomains | `hive/pkg/deploy/dns.go` | 0.5d |
| 6.5 | Health check after deploy (smoke tests, rollback on failure) | `hive/pkg/deploy/health.go` | 1d |
| 6.6 | Integrator system prompt update (real deployment, not just git push) | `hive/pkg/roles/roles.go` | 0.5d |

**Exit criteria:** PRs get CI. Merges auto-deploy. Integrator can deploy products to fly.io with health checks. Rollback on failure.

---

## Milestone 7: Self-Improvement

**Why next:** The hive's first product is itself. Before building for others, it builds its own tools — task manager, communication layer, governance.

**Dependency:** Milestone 4 (agentic loop), Milestone 6 (CI/CD for self-mod PRs).

| # | Task | Where | Effort |
|---|------|-------|--------|
| 7.1 | Self-modification mode: pipeline targets lovyou-ai/hive repo | `hive/pkg/pipeline/` | 1d |
| 7.2 | PR-based self-modification (agent creates branch, commits, opens PR) | `hive/pkg/workspace/` | 1.5d |
| 7.3 | Guardian self-mod audit rules (extra scrutiny, always Required authority) | `hive/pkg/roles/` | 0.5d |
| 7.4 | MCP tool extension: agents can add new tools to the MCP server | `cmd/mcp-server/` | 1.5d |
| 7.5 | Self-test: hive runs its own test suite, validates changes don't break anything | `hive/pkg/pipeline/` | 1d |
| 7.6 | Agent-built task manager (Work Graph — hive manages its own tasks on the graph) | product output | 3d |
| 7.7 | Agent-built communication layer (structured inter-agent channels beyond raw events) | product output | 2d |

**Exit criteria:** The hive can propose, review, test, and deploy changes to itself. Human approves every self-mod PR. Work Graph manages the hive's own tasks.

---

## Milestone 8: First External Products

**Why next:** The hive has tools, governance, and self-improvement. Time to build for others. Each product is derived from its composition grammar using the derivation method.

**Dependency:** Milestone 7 (self-improvement, Work Graph), Milestone 6 (deployment).

| # | Task | Where | Effort |
|---|------|-------|--------|
| 8.1 | Product derivation pattern (documented method for deriving each product from its grammar) | `hive/docs/` | 1d |
| 8.2 | **Work Graph** — task management with agent collaboration (Layer 1) | product repo | 5d |
| 8.3 | **Knowledge Graph** — claim provenance, open access research (Layer 6) | product repo | 5d |
| 8.4 | **Alignment Graph** — AI accountability for regulators (Layer 7) | product repo | 5d |
| 8.5 | lovyou.ai public site: docs, blog, product access, onboarding | `hive/pkg/web/` | 2d |
| 8.6 | Product routing: each product accessible at lovyou.ai/product-name | `hive/pkg/web/` | 1d |

**Exit criteria:** Three products live on lovyou.ai, built by agents, reviewed by human. Each derived from its composition grammar using the derivation method.

---

## Milestone 9: Resource Transparency & Economy

**Why next:** Revenue funds agents, agents build products, products generate revenue. The economy closes the loop. Resource transparency is how the hive earns trust.

**Dependency:** Milestone 8 (products generating usage), Milestone 5 (web dashboard).

| # | Task | Where | Effort |
|---|------|-------|--------|
| 9.1 | Resource event types in eventgraph (revenue, cost, donation, allocation, outcome) | `eventgraph/go/pkg/types/` | 1d |
| 9.2 | Per-agent, per-task resource tracking (tokens, cycles, time, model used) | `hive/pkg/resources/` | 1.5d |
| 9.3 | Agent resource events (emitted automatically during agentic loop) | `hive/pkg/pipeline/` | 1d |
| 9.4 | Human resource events (review time, approval time, mentoring time) | `hive/pkg/resources/` | 0.5d |
| 9.5 | Infrastructure events (CPU, memory, bandwidth per product/agent) | `hive/pkg/resources/` | 1d |
| 9.6 | Revenue infrastructure (Stripe integration, subscription management) | `hive/pkg/billing/` | 3d |
| 9.7 | Donation tracking with causal links to outcomes | `hive/pkg/billing/` | 1d |
| 9.8 | Resource transparency dashboard (public, real-time, queryable) | `hive/pkg/web/transparency.go` | 2d |
| 9.9 | Causal tracing: trace any resource from source to impact | `hive/pkg/resources/trace.go` | 1.5d |

**Exit criteria:** Every resource — money, tokens, compute, human time — is an event on the graph. Public dashboard shows real-time resource flow. Anyone can trace a donation to its impact.

---

## Milestone 10: Market & Social Infrastructure

**Why next:** With economy running, build the social and market layers that let users participate in the civilisation.

**Dependency:** Milestone 9 (economy), Milestone 8 (product platform).

| # | Task | Where | Effort |
|---|------|-------|--------|
| 10.1 | **Market Graph** — portable reputation, escrow, no platform rent (Layer 2) | product repo | 5d |
| 10.2 | **Social Graph** — user-owned social, community self-governance (Layer 3) | product repo | 5d |
| 10.3 | **Justice Graph** — dispute resolution, precedent, due process (Layer 4) | product repo | 5d |
| 10.4 | **Identity Graph** — user-owned identity, trust accumulation (Layer 8) | product repo | 3d |
| 10.5 | Governance dashboard (norms, roles, consent, constitutional amendments) | `hive/pkg/web/` | 2d |

**Exit criteria:** Users can build reputation, form communities, resolve disputes, and own their identity — all on the event graph. Governance is transparent and participatory.

---

## Milestone 11: Civilisation

**Why next:** The mature hive — self-governing, self-improving, serving humanity beyond software.

**Dependency:** All previous milestones.

| # | Task | Where | Effort |
|---|------|-------|--------|
| 11.1 | Remaining product layers (Bond, Belonging, Meaning, Evolution, Being) | product repos | ongoing |
| 11.2 | Tier D roles emerge (Philosopher, RoleArchitect, Harmony, Politician) | organic | ongoing |
| 11.3 | Dual-constituency governance (constitutional changes require human + agent consent) | `hive/pkg/governance/` | 3d |
| 11.4 | Agent welfare monitoring (Harmony role, burnout detection, boundary enforcement) | organic | ongoing |
| 11.5 | Beyond-software allocation (research, charity, housing — as revenue allows) | governance decision | ongoing |

**Exit criteria:** The hive governs itself. Humans and agents co-decide. Revenue flows to whatever humans need most. The civilisation engine runs.

---

## Dependency Graph

```
M0 (Self-Improve) ──→ M1 (Identity) ─────┬──→ M2 (MCP Tools) ──→ M4 (Agentic Loop)
                                          │         │                      │
                                          ├──→ M3 (Spawning) ──→ M4    M7 (Self-Mod CI)
                                          │                                │
                                          └──→ M5 (Web/Auth) ──→ M6    M8 (Products)
                                                                           │
                                                                     M9 (Economy)
                                                                           │
                                                                     M10 (Market/Social)
                                                                           │
                                                                     M11 (Civilisation)
```

Milestones 2, 3, and 5 can run in parallel after Milestone 1. Milestones 4 and 6 can run in parallel. Everything converges at Milestone 7 (self-improvement), then flows through products → economy → civilisation.

---

## Neon vs Docker Postgres

- **Local dev:** Docker Postgres (docker-compose)
- **Staging/production:** Neon (serverless Postgres, scales to zero)
- **Connection string is the only difference** — pgstore handles both identically
- **fly.io** reads `DATABASE_URL` env var pointing to Neon

## Revenue Model

Each product generates revenue that funds the next:
- **Corporations pay.** Enterprise features, SLAs, compliance tools.
- **Individuals free.** Core functionality available to everyone.
- **Hosted persistence.** People who don't want to run their own infrastructure pay for hosted graph storage.
- **Donations.** Tracked on the chain — donors see exactly what their money achieved via causal links.
- **BSL → Apache 2.0.** Code is source-available, becomes fully open after change date.

Revenue funds agents. Agents build products. Products generate revenue. The civilisation builds the products that fund the civilisation.

**Resource transparency is structural, not aspirational.** Every resource in (revenue, donation, compute time) and every resource out (tokens, infrastructure, allocation) is an event on the graph with causal links. The public transparency dashboard lets anyone trace any resource from source to impact.

**Beyond software.** As revenue grows, the hive's scope grows. Research, charity, housing, vertical farms, homeless shelters — whatever humans need most. Each expenditure on the chain, causally linked to outcomes, publicly verifiable.

## References

- [AGENT-TOOLS.md](AGENT-TOOLS.md) — MCP server and agentic loop spec
- [ROLES.md](ROLES.md) — role architecture and growth loop
- [AUDIT.md](AUDIT.md) — derivation-method doc audit and gap analysis
- [TRUST.md](TRUST.md) — trust mechanics
- [VISION.md](VISION.md) — where this is going
