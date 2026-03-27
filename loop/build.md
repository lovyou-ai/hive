# Build: Wire nav links and add handler test

## Task
Wire nav links and add handler test for `/hive`.

## Changes

### `site/graph/views.templ`
- Added `<a href="/hive">Hive</a>` to `simpleFooter` nav links (after Discover, before Market)
- Added `<p class="text-sm text-warm-faint italic">The Civilization Builds.</p>` tagline to `HiveStatusPartial` (below the h1 "The Civilization Engine")

### `site/views/layout.templ`
- Added `<a href="/hive">Hive</a>` to the footer nav links (after Discover, before Agents)

### `site/graph/hive_test.go`
- Added `TestGetHive_ContainsCivilizationBuilds`: issues `GET /hive`, asserts HTTP 200 and body contains "The Civilization Builds"

## Verification
- `templ generate` — 16 updates, no errors
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test -buildvcs=false ./...` — all pass (graph: 0.086s)

## Notes
- Header (`simpleHeader` in views.templ) already had the Hive nav link from a prior iteration — not duplicated
- Test placed in `site/graph/hive_test.go` (the existing handler test file); `site/handlers/` does not exist as a package
