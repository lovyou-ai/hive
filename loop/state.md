# Loop State

Living document. Updated by the Reflector each iteration. Read by the Scout first.

Last updated: Iteration 6, 2026-03-22.

## Current System State

Five repos, all compiling and tested:
- **eventgraph** — foundation. Postgres stores, 201 primitives, trust, authority. Complete.
- **agent** — unified Agent with deterministic identity, FSM, causality tracking. Complete.
- **work** — task store for hive agent coordination. Complete.
- **hive** — 4 agents (Strategist, Planner, Implementer, Guardian), agentic loop, budget. Complete.
- **site** — lovyou.ai on Fly.io. **Deployed and live.** Complete product:
  - Blog (43 posts, markdown → HTML)
  - Reference (cognitive grammar, graph grammar, 13 layer grammars, 201 primitives, 28 agent primitives)
  - Auth (Google OAuth, anonymous fallback)
  - Unified graph product (3 tables: spaces/nodes/ops, 10 grammar operations, 5 lenses, HTMX, full CRUD)
  - Landing page: concrete product description, five lens cards, three-step flow, EventGraph links

Entry point: `cmd/hive` (CLI only). No daemon, no web dashboard for the hive itself.

Deploy method: `fly deploy --remote-only` (avoids Docker Desktop dependency).

## Lessons Learned

1. **Code is truth, not docs.** Always read code to assess current state.
2. **Accumulate knowledge.** state.md + reflections.md prevent repeated mistakes.
3. **Update state.md every iteration.** Prevents phantom gaps.
4. **Ship what you build.** Every iteration should end with changes in version control.
5. **Try alternatives before declaring blockers.** `--remote-only` worked; `--local-only` was the problem, not Docker.
6. **Every Build iteration should deploy.** Scout → Builder → commit → push → deploy in one cycle.

## Known Issues

- No SEO: no meta description, no Open Graph tags, no structured data.
- No analytics: no way to measure visitor behavior or conversion.
- AUDIT.md in hive/docs/ is stale (March 9th). Low priority.
- Boot events re-emitted every agent run (chatty but harmless).
- DATABASE_URL required for /app — visitors without it get 503.

## What the Scout Should Focus On Next

The loop is in Build mode. Continue building. Candidates:

1. **SEO / meta tags** — site is invisible to search engines and social shares. No meta description, no OG tags. High leverage for discoverability.
2. **Onboarding / app experience** — what happens when someone clicks "Open the app"? Does the flow make sense for a first-time user?
3. **Hive autonomy** — make the core loop self-running (persistent process, not manual CLI invocation).
4. **Blog index** — 43 posts, potentially overwhelming. Could benefit from categories or a reading guide.

Pick one. Build it. Deploy it.
