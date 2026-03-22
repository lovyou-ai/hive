# Loop State

Living document. Updated by the Reflector each iteration. Read by the Scout first.

Last updated: Iteration 13, 2026-03-22.

## Current System State

Five repos, all compiling and tested:
- **eventgraph** — foundation. Postgres stores, 201 primitives, trust, authority. Complete. Has CI.
- **agent** — unified Agent with deterministic identity, FSM, causality tracking. Complete.
- **work** — task store for hive agent coordination. Complete.
- **hive** — 4 agents (Strategist, Planner, Implementer, Guardian), agentic loop, budget. Complete. Has CI.
- **site** — lovyou.ai on Fly.io. Production-ready. **Has CI** (build + templ drift check).

**Core loop infrastructure:**
- `loop/run.sh` — orchestrates Scout → Builder → Critic → Reflector via `claude -p`
- Four phase prompt files (scout, builder, critic, reflector)
- `.github/workflows/ci.yml` — build + test on push, PR, and workflow_dispatch
- Site also has `.github/workflows/ci.yml` — build + templ drift check
- Run: `cd /c/src/matt/lovyou3/hive && ./loop/run.sh`

Deploy: `fly deploy --remote-only` from site repo.
Fly/Neon resources can be scaled up per user authorization.

## Completed Clusters

- **Orient** (1-4): catch up with reality, fix stale docs, accumulate knowledge
- **Ship** (5): deploy fix (`--remote-only`)
- **Discoverability** (6-8): landing page, SEO, sitemap
- **Visitor Experience** (9): blog arc navigation
- **SEO Canonicalization** (10): fly.dev → lovyou.ai redirect
- **Hive Autonomy** (11-13): prompt files, run.sh, CI on hive + site

## Lessons Learned

1. Code is truth, not docs.
2. Verify infra assumptions before building.
3. Update state.md every iteration.
4. Ship what you build — every Build iteration should deploy.
5. Try alternatives before declaring blockers.
6. Name iteration clusters and recognize completion.
7. Hostname middleware must exclude /health (Fly probes via internal IP).
8. Codify implicit knowledge into executable artifacts — conversation context is ephemeral, files persist.
9. Multi-repo replace directives require CI to mirror the local directory structure (checkout siblings).
10. Templ drift check (`git diff --exit-code -- '*_templ.go'`) catches stale generated files.

## Vision Notes

- Agents should acquire skills dynamically (like OpenClaw) — email, invoicing, payments, public accounting, any skill.
- Auth gate can be opened to public whenever ready.
- Users provide OAuth tokens via `claude --setup-token`, agents build things for them via board or personal agent.
- Social product: humans and agents build MySpace-like personal pages hosted by lovyou.ai.
- Business use: companies use the platform to build their products (e.g., Lovatts Anthro account).
- Agents and humans are peers on the social graph.

## What the Scout Should Focus On Next

The infrastructure cluster is complete. All repos have CI, the loop is codified, the site is deployed. Time to shift to product or capability.

Options:
1. **Product development** — the unified graph product exists behind an auth gate. Open it, or build the social product (personal pages for humans + agents).
2. **Agent capability** — design how users provide tokens, how agents get assigned to users, how agents build things via the board.
3. **Content** — new blog posts, documentation for the product.

The user's expanding vision (personal agents, social pages, business accounts) suggests the product direction is most aligned with intent. The Scout should assess what concrete product work would be most impactful.
