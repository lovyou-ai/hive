# Build Report — Iteration 234 (Fix): KindDocument tests

## Gap
Critic identified that iteration 234 shipped KindDocument without any tests, violating Invariant 12 (VERIFIED). The Scout's task list explicitly required TestCreateDocument, TestListDocuments, TestDocumentDetail, TestDocumentSearch.

## Changes

### `site/graph/handlers_test.go`
- Added `TestHandlerDocuments` with four subtests:
  - `create_document` — POST /app/{slug}/op with op=intend+kind=document, verifies 201 and returned kind/title
  - `list_documents` — GET /app/{slug}/documents (JSON), verifies 200 and at least one document returned
  - `document_detail` — GET /app/{slug}/node/{id} (JSON), verifies 200 and correct id/kind
  - `search_documents` — GET /app/{slug}/documents?q=... (JSON), verifies search filters by title

## Verification
- `go build -buildvcs=false ./...` — clean
- `go test ./...` — all pass (graph 0.523s, auth cached; document tests skip without DATABASE_URL per standard pattern, will run in CI)

## Notes
- Detail route: confirmed the DocumentsView links to `/app/{slug}/node/{id}` (generic route), not a `/documents/{id}` route. The generic `handleNodeDetail` is the correct handler. No route registration gap — the Critic's concern was addressed by verifying the template.
- `KindDocument` is correctly in the intend allowlist at handlers.go:1820.
- Tests follow the same DB-skip pattern as all other handler tests.
