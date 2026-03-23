# Build Report — Iteration 131

## Global activity context — node titles on /activity and profiles

- **store.go:** `ListPublicActivity` now JOINs nodes for title (same pattern as iter 127's ListOps)
- **views/activity.templ:** `ActivityItem` gains NodeID + NodeTitle fields. `activityRow` shows node title as clickable link
- **cmd/site/main.go:** Both /activity handler and profile handler now pass NodeID + NodeTitle to ActivityItem

All activity views now show what happened, not just who did what type of thing.
