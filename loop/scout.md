# Scout Report — Iteration 196

## Gap: Repost attribution on Following feed

**Source:** social-spec.md PostCard — "↻ reposted_by.name reposted" header. Phase 3 composition.

**Current state:** Following feed includes posts reposted by followed users, but doesn't show WHO reposted them. Posts appear as if they're regular feed items — no context for why you're seeing them.

**What's needed:**
1. Handler: when building the Following feed, track which posts were included via repost (not direct authorship)
2. Store: resolve reposter names for those posts (who among the followed users reposted this?)
3. Template: "↻ username reposted" header above the FeedCard when attribution exists

**From the spec:**
```
if post.reposted_by {
    Layout(row, gap: xs, class: "text-xs text-muted mb-1", [
      Display("↻"), Display(post.reposted_by.name + " reposted")
    ])
}
```

**Approach:** In the Following filter, build a `repostedBy map[string]string` (nodeID → reposter name). Pass it through FeedView to FeedCard. When set, render the attribution header. Use ResolveUserNames to get display names from IDs.

**Risk:** Low. Handler logic + template addition. No schema changes. No new store methods needed beyond existing ones.
