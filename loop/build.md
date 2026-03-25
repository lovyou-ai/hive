# Build Report — Iteration 234: KindDocument entity kind

## Gap
Add `KindDocument` entity kind — Wiki product foundation.

## Changes

### `site/graph/store.go`
- Added `KindDocument = "document"` constant to the node kinds block (after `KindPolicy`).

### `site/graph/handlers.go`
- Updated intend allowlist to include `KindDocument` — allows creating documents via the `intend` op.
- Added `handleDocuments` handler: lists nodes with `Kind=KindDocument`, supports JSON content negotiation and search query.
- Registered route `GET /app/{slug}/documents` mapped to `handleDocuments`.

### `site/graph/views.templ`
- Added `documentsIcon()` — document SVG icon.
- Added `DocumentsView` — list view with new document form, search, document cards, empty state.
- Added `@lensLink(space.Slug, "documents", "Docs", ...)` to sidebar "More" section (after policies).
- Added `@mobileLensTab(space.Slug, "documents", "Docs", ...)` to mobile nav "More" drawer (after policies).

## Verification
- `templ generate` — 15 updates, no errors
- `go build -buildvcs=false ./...` — clean
- `go test ./...` — all pass (`graph` 0.531s, `auth` cached)

## Notes
- Search is automatic: `Search()` queries all nodes in public spaces — documents included without change.
- Detail view: existing `handleNodeDetail` + `NodeDetailView` renders body via `renderMarkdown()`.
- Pattern: replicated `KindTeam` (iter 223) exactly.
