# Loop State

Living document. Updated by the Reflector each iteration. Read by the Scout first.

Last updated: Iteration 15, 2026-03-22.

## Current System State

Five repos, all compiling and tested:
- **eventgraph** — foundation. Postgres stores, 201 primitives, trust, authority. Complete. Has CI.
- **agent** — unified Agent with deterministic identity, FSM, causality tracking. Complete.
- **work** — task store for hive agent coordination. Complete.
- **hive** — 4 agents (Strategist, Planner, Implementer, Guardian), agentic loop, budget. Complete. Has CI.
- **site** — lovyou.ai on Fly.io. Production-ready. Has CI. **Public spaces now supported.**

**Product features:**
- Blog (43 posts, 6 arcs with section nav)
- Reference (cognitive grammar, graph grammar, 13 layers, 201 primitives, 28 agent primitives)
- Auth (Google OAuth — test mode, can be opened whenever)
- Unified graph product (3 tables, 10 grammar ops, 5 lenses, HTMX, full CRUD)
- **Public spaces** — spaces can be private (owner only) or public (anyone can view)
- OptionalAuth for read routes, RequireAuth for write routes
- Landing page, SEO meta tags, sitemap (305 URLs), canonical redirect
- **Warm copy** — home page and footer rewritten for collaborative tone

Deploy: `fly deploy --remote-only` from site repo.

## Completed Clusters

- **Orient** (1-4): catch up with reality, fix stale docs, accumulate knowledge
- **Ship** (5): deploy fix (`--remote-only`)
- **Discoverability** (6-8): landing page, SEO, sitemap
- **Visitor Experience** (9): blog arc navigation
- **SEO Canonicalization** (10): fly.dev → lovyou.ai redirect
- **Hive Autonomy** (11-13): prompt files, run.sh, CI on hive + site
- **Product Development** (14-15): public spaces, warm aesthetics copy

## Lessons Learned

1. Code is truth, not docs.
2. Verify infra assumptions before building.
3. Update state.md every iteration.
4. Ship what you build — every Build iteration should deploy.
5. Try alternatives before declaring blockers.
6. Name iteration clusters and recognize completion.
7. Hostname middleware must exclude /health (Fly probes via internal IP).
8. Codify implicit knowledge into executable artifacts.
9. Multi-repo replace directives require CI to mirror directory structure.
10. Templ drift check catches stale generated files.
11. Start with the simplest access model (public/private) before building roles/ACLs.
12. When the founder says "that isn't our vibe," treat it as highest-priority. Tone compounds — every visitor gets the wrong impression until it's fixed.

## Vision Notes

- Agents should acquire skills dynamically (like OpenClaw).
- Auth gate can be opened to public whenever ready.
- Users provide OAuth tokens, agents build things for them via board or personal agent.
- Social product: humans and agents build MySpace-like personal pages.
- Business use: companies use the platform to build products.
- Agents and humans are peers on the social graph.
- Site vibe should be warm/collaborative (agents+humans together), NOT corporate/business-like.
- **Aesthetics direction:** lovyou2 (`/c/src/matt/lovyou2`) is inspiration. It uses a dark theme with pink accent (#f472b6), light font weights (200-300 for headings), system fonts, breathing/pulse animations, generous whitespace, and "sacred minimalism" — intelligence meets warmth. The Scout should research color theory, artist palettes, and aesthetic principles to develop a cohesive visual identity rather than copying lovyou2 directly.

## What the Scout Should Focus On Next

Copy is now warm. Visual design (colors, typography, spacing) is the next aesthetics lever. Options:

1. **Visual aesthetics** — research color theory, palettes, and design principles. Study lovyou2's dark+pink approach as inspiration. Develop a cohesive visual identity for lovyou3 that feels warm, intelligent, and collaborative. This is a research-then-build task: study first, propose, then implement.
2. **Discover page** — list public spaces so visitors can browse without knowing URLs.
3. **Open auth gate** — switch Google OAuth from test mode to production so anyone can sign up.
4. **Space settings** — allow changing visibility after creation, rename, delete spaces.

The aesthetics direction from the user suggests studying art, palettes, and color theory — not just copying lovyou2. This is a creative research task before it becomes a build task.
