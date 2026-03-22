# Build Report — Iteration 17

## What I planned

Add a `/discover` page so visitors can browse public spaces without login. Wire `ListPublicSpaces()` (already in the store) to a new route, view, and nav link.

## What I built

Changes across 7 files (3 templates + 3 generated + 1 Go source) in the site repo.

### New: views/discover.templ
- `DiscoverSpace` struct: Slug, Name, Description, Kind, CreatedAt
- `DiscoverPage(spaces []DiscoverSpace)` — grid layout with empty state CTA
- `discoverCard(sp DiscoverSpace)` — card with kind badge, date, name, description
- `discoverKindClass(kind string)` — dark badge colors (project=indigo, community=emerald, team=amber)
- Uses `Layout()` from views package for visual consistency with blog/reference

### Modified: cmd/site/main.go
- Hoisted `graphStore` declaration (`var graphStore *graph.Store`) above the DB block
- Changed `:=` to `=` inside the DB block so `/discover` handler can access it
- Added `GET /discover` handler: calls `graphStore.ListPublicSpaces()`, maps `graph.Space` → `views.DiscoverSpace`, renders page
- Graceful degradation: if no DB, renders empty discover page
- Added `/discover` to sitemap

### Modified: views/layout.templ
- Added "Discover" link to nav bar between "App" and "Blog"

### Modified: graph/views.templ
- Added "Discover" link to `simpleHeader()` nav
- Added "Discover" link to `appLayout()` header nav

## Verification

- `templ generate` — success (7 updates)
- `go build -o /tmp/site.exe ./cmd/site/` — success
- Committed and pushed to main
- Deployed to Fly.io — both machines healthy
