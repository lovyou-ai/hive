# Loop State

Living document. Updated by the Reflector each iteration. Read by the Scout first.

Last updated: Iteration 17, 2026-03-22.

## Current System State

Five repos, all compiling and tested:
- **eventgraph** — foundation. Postgres stores, 201 primitives, trust, authority. Complete. Has CI.
- **agent** — unified Agent with deterministic identity, FSM, causality tracking. Complete.
- **work** — task store for hive agent coordination. Complete.
- **hive** — 4 agents (Strategist, Planner, Implementer, Guardian), agentic loop, budget. Complete. Has CI.
- **site** — lovyou.ai on Fly.io. Production-ready. Has CI. **Dark theme + warm copy + discover page.**

**Product features:**
- Blog (43 posts, 6 arcs with section nav)
- Reference (cognitive grammar, graph grammar, 13 layers, 201 primitives, 28 agent primitives)
- Auth (Google OAuth — test mode, can be opened whenever)
- Unified graph product (3 tables, 10 grammar ops, 5 lenses, HTMX, full CRUD)
- Public spaces (private/public visibility, OptionalAuth for reads)
- **Discover page** (`/discover`) — public space directory, grid of cards, kind badges, no auth required
- Landing page, SEO meta tags, sitemap (306 URLs), canonical redirect
- **Visual identity**: dark theme (near-black #09090b), rose accent (#e8a0b8), warm off-white text (#f0ede8), light heading weights, "Ember Minimalism"

Deploy: `fly deploy --remote-only` from site repo.

## Completed Clusters

- **Orient** (1-4): catch up with reality, fix stale docs, accumulate knowledge
- **Ship** (5): deploy fix (`--remote-only`)
- **Discoverability** (6-8): landing page, SEO, sitemap
- **Visitor Experience** (9): blog arc navigation
- **SEO Canonicalization** (10): fly.dev → lovyou.ai redirect
- **Hive Autonomy** (11-13): prompt files, run.sh, CI on hive + site
- **Product Development** (14): public spaces
- **Aesthetics** (15-16): warm copy rewrite, dark theme with rose accent
- **Discovery** (17): `/discover` page for browsing public spaces

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
12. When the founder says "that isn't our vibe," treat it as highest-priority.
13. Define the vocabulary before writing the prose — custom color tokens make future styling consistent.
14. Expose what you've already built before building more — wiring existing code to a route is faster than new infrastructure.

## Vision Notes

- Agents should acquire skills dynamically (like OpenClaw).
- Auth gate can be opened to public whenever ready.
- Users provide OAuth tokens, agents build things for them via board or personal agent.
- Social product: humans and agents build MySpace-like personal pages.
- Business use: companies use the platform to build products.
- Agents and humans are peers on the social graph.
- Visual identity: "Ember Minimalism" — dark, warm, intentional. lovyou2 as ancestor.

## What the Scout Should Focus On Next

Discovery cluster complete. The site now has a browsable public space directory. Options:

1. **Open auth gate** — switch Google OAuth from test mode to production so anyone can sign up. This is the biggest unlock: discover + auth = real users.
2. **Space settings** — allow changing visibility after creation, rename, delete spaces.
3. **Subtle animations** — breathing pulse on brand elements, scroll reveals, gentle transitions. Would add life to the dark theme.
4. **Space previews** — show a preview of what's inside a space on the discover card (node count, recent activity).

Opening the auth gate has the most product value — it makes the entire product publicly usable, not just viewable.
