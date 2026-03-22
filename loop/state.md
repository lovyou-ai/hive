# Loop State

Living document. Updated by the Reflector each iteration. Read by the Scout first.

Last updated: Iteration 9, 2026-03-22.

## Current System State

Five repos, all compiling and tested:
- **eventgraph** — foundation. Postgres stores, 201 primitives, trust, authority. Complete.
- **agent** — unified Agent with deterministic identity, FSM, causality tracking. Complete.
- **work** — task store for hive agent coordination. Complete.
- **hive** — 4 agents (Strategist, Planner, Implementer, Guardian), agentic loop, budget. Complete.
- **site** — lovyou.ai on Fly.io. **Deployed and live.** Complete:
  - Blog (43 posts, 6 arcs with section nav)
  - Reference (cognitive grammar, graph grammar, 13 layer grammars, 201 primitives, 28 agent primitives)
  - Auth (Google OAuth — test mode, can be opened whenever)
  - Unified graph product (3 tables, 10 grammar operations, 5 lenses, HTMX, full CRUD)
  - Landing page: clear product description, five lens cards, three-step flow
  - SEO: meta description + OG tags on every page
  - Sitemap: 305 URLs, robots.txt
  - All secrets configured on Fly

Deploy: `fly deploy --remote-only` from site repo.

## Lessons Learned

1. **Code is truth, not docs.**
2. **Verify infra assumptions** before building.
3. **Update state.md every iteration.**
4. **Ship what you build** — every Build iteration should deploy.
5. **Try alternatives before declaring blockers.**
6. **Name iteration clusters** and recognize when they're complete.
7. **Focus on public-facing improvements** while auth is in test mode.

## Completed Clusters

- **Orient** (1-4): catch up with reality
- **Ship** (5): deploy fix
- **Discoverability** (6-8): landing page, SEO, sitemap
- **Visitor Experience** (9): blog navigation

## Known Issues

- No analytics — no way to measure traffic or behavior.
- AUDIT.md in hive/docs/ is stale. Low priority.
- Boot events re-emitted every agent run.

## What the Scout Should Focus On Next

The site is visitor-ready: clear, discoverable, navigable. Two directions:

1. **Product access** — open the auth gate, verify the app experience works for new users. The user authorized this. But first, the loop should test the anonymous fallback and onboarding flow to ensure it doesn't confuse visitors.
2. **Hive autonomy** — make the core loop self-running. Currently requires manual Claude Code CLI invocation. This is the infrastructure that would let the hive improve itself.

Pick one. The product access path has immediate user-facing impact. The hive autonomy path has compounding long-term impact.
