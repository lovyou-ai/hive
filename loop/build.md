# Build Report — Council List/Detail Templates + Handler + Sidebar Nav

## Gap
No way to view or navigate to Council sessions in the app. The `council` lens, routes, templates, and sidebar entry were all missing.

## Changes

### `graph/store.go`
- Added `KindCouncil = "council"` constant alongside `KindDocument`, `KindQuestion`
- Added `ListCouncilSessions(ctx, spaceID, limit)` helper (same pattern as `ListDocuments`, `ListQuestions`)

### `graph/handlers.go`
- Registered `GET /app/{slug}/council` → `handleCouncil`
- Registered `GET /app/{slug}/council/{id}` → `handleCouncilDetail`
- Added `handleCouncil`: fetches `KindCouncil` nodes, renders `CouncilListView`
- Added `handleCouncilDetail`: fetches session node (404 if wrong kind/space), fetches child responses, renders `CouncilDetailView`
- Both handlers support `wantsJSON` for API access

### `graph/views.templ`
- Added `councilIcon()` — group/users SVG, consistent with other lens icons
- Added `@mobileLensTab` for "Council" in mobile nav (between Threads and Knowledge)
- Added `@lensLink` for "Council" in desktop sidebar More section (between Threads and Knowledge)
- Added `lensDescription` case for "council"
- Added `CouncilListView`: list of council sessions with question, response count, timestamp; empty state with icon
- Added `CouncilDetailView`: question panel at top with "Convened by" attribution; each agent response as a violet panel with agent badge and timestamp; empty state message

## Verification
- `templ generate` — 16 updates, no errors
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (graph: 0.594s)
