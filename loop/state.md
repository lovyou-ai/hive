# Loop State

Living document. Updated by the Reflector each iteration. Read by the Scout first.

Last updated: Iteration 8, 2026-03-22.

## Current System State

Five repos, all compiling and tested:
- **eventgraph** — foundation. Postgres stores, 201 primitives, trust, authority. Complete.
- **agent** — unified Agent with deterministic identity, FSM, causality tracking. Complete.
- **work** — task store for hive agent coordination. Complete.
- **hive** — 4 agents (Strategist, Planner, Implementer, Guardian), agentic loop, budget. Complete.
- **site** — lovyou.ai on Fly.io. **Deployed and live.** Complete:
  - Blog (43 posts, markdown → HTML)
  - Reference (cognitive grammar, graph grammar, 13 layer grammars, 201 primitives, 28 agent primitives)
  - Auth (Google OAuth — test mode, **can be opened whenever**)
  - Unified graph product (3 tables, 10 grammar operations, 5 lenses, HTMX, full CRUD)
  - Landing page: clear product description, five lens cards, three-step flow
  - SEO: meta description + OG tags + Twitter cards on every page
  - Sitemap: 305 URLs, robots.txt points to sitemap
  - All secrets configured on Fly (DATABASE_URL, Google OAuth)

Deploy: `fly deploy --remote-only` from site repo.

## Lessons Learned

1. **Code is truth, not docs.** Always read code to assess current state.
2. **Verify infra assumptions.** State.md said DATABASE_URL might not be set — it was. Check before building.
3. **Update state.md every iteration.** Prevents phantom gaps.
4. **Ship what you build.** Every Build iteration should deploy.
5. **Try alternatives before declaring blockers.** `--remote-only` solved the Docker problem.
6. **Name iteration clusters.** Orient (1-4), Ship (5), Discoverability (6-8). Recognize when a cluster is complete.
7. **Focus on public-facing improvements** while auth is in test mode.

## Known Issues

- AUDIT.md in hive/docs/ is stale (March 9th). Low priority.
- Boot events re-emitted every agent run (chatty but harmless).
- No analytics — no way to measure traffic or behavior.
- Blog has 43 posts with no categorization or reading guide — potentially overwhelming.

## What the Scout Should Focus On Next

Discoverability cluster is COMPLETE (landing page + SEO + sitemap). The site is visible and comprehensible.

Next cluster candidates:

1. **Visitor experience** — 43 blog posts is overwhelming. A reading guide or categorized index would help new visitors find their way. The blog is the main content visitors will discover via search.
2. **Open the auth gate** — user authorized this. Would make the app accessible to all visitors, not just Matt.
3. **Hive autonomy** — make the core loop self-running (persistent process, scheduled iterations).

Pick one cluster. Build the first iteration of it.
