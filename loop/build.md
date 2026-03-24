# Build Report — Iteration 196

## Repost Attribution

**Store:**
- `GetRepostAttribution(userIDs, nodeIDs) map[string]string` — for each node, returns the user ID of the most recent reposter from the given user set. Uses `DISTINCT ON (node_id)` with `ORDER BY created_at DESC`.

**Handler:**
- Following filter: after filtering posts, identifies which posts are in the feed via repost (not direct authorship)
- Calls `GetRepostAttribution` to find which followed user reposted each
- Resolves reposter IDs to display names via `ResolveUserNames`
- Passes `repostedBy map[string]string` (nodeID → display name) to FeedView

**Template:**
- `FeedView`: accepts `repostedBy map[string]string`
- `FeedCard`: accepts `repostedByName string`
- When `repostedByName != ""`, renders "↻ username reposted" header above the card (before pin indicator)
- Uses the same ↻ arrows SVG as the repost button, 10px text, warm-faint color

**Files changed:**
- `graph/store.go` — `GetRepostAttribution`
- `graph/handlers.go` — attribution logic in Following filter, FeedView call
- `graph/views.templ` — FeedView, FeedCard signatures + attribution header
