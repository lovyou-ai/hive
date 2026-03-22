# Scout Report — Iteration 13

## Map (from code + infra)

Read state.md. Hive Autonomy cluster complete (iterations 11-12). Explored the site repo:

- Site deploys to production (lovyou.ai on Fly.io) but has NO CI
- No replace directives — clean go.mod, simpler than hive's CI
- Uses templ for templates — `_templ.go` files committed to git but could drift from `.templ` source
- Has a Makefile: `templ generate && go build -o site ./cmd/site/`
- No tests (no `*_test.go` files)
- Go 1.25, standard deps (templ, pq, oauth2, goldmark)

## Gap Type

Missing infrastructure — production-deployed code has no automated build verification.

## The Gap

The site repo deploys to production but has no CI. Code can be pushed that doesn't compile. Templ-generated files can drift from source templates without anyone noticing.

## Why This Gap

The site is the public face — it's what visitors see. If a push breaks the build, the next `fly deploy` fails and the site is stuck on old code. CI catches this before deploy. Also: templ source/generated drift is a subtle bug that only CI would catch (regenerate and verify).

## Filled Looks Like

Push to site/main triggers: install templ → `templ generate` → `go build`. Green check on every commit.
