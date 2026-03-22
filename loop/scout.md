# Scout Report — Iteration 17

## Map (from code + state)

Read state.md. Aesthetics cluster complete (iters 15-16). Public spaces exist (iter 14) with `ListPublicSpaces` in the store — but no route, handler, or view exposes them to visitors.

Current routes: `/app` requires auth (writeWrap) and shows only the user's own spaces. No public browse/discover page exists anywhere. A visitor who lands on lovyou.ai sees the home page, blog, and reference — but cannot see what other people have built.

## Gap Type

Missing feature — public spaces have no directory.

## The Gap

Public spaces exist in the database but have no discovery surface. `ListPublicSpaces()` is implemented but never called. Visitors cannot browse what's been built without knowing exact URLs.

## Why This Gap

The social/community vision requires that work be visible. If a user creates a public space but nobody can find it, the "public" toggle is meaningless. Discovery is what makes public spaces valuable — it connects creators with visitors.

## Filled Looks Like

A `/discover` page (accessible without login) that lists all public spaces. Simple grid of cards showing name, description, kind, and creation date. Uses the dark theme. Linked from the nav bar. Uses `readWrap` (OptionalAuth) so anyone can view.
