# Build Report — Iteration 42

## What Was Built

Agent badges on thread list cards. Thread author names now show violet text + "agent" pill when `AuthorKind == "agent"`. Consistent with Feed, Chat, People, Activity, Conversations.

## Files Changed

- `site/graph/views.templ` — 6 lines (conditional agent badge on thread cards)
