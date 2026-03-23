# Build Report — Iteration 87

Rewrote `/app` from "Your Spaces" grid to "My Work" personal dashboard.

**New store queries:**
- `ListUserTasks(ctx, userID, limit)` — open tasks where user is author or assignee, across all spaces, sorted by priority
- `ListUserConversations(ctx, userID, limit)` — conversations with user as participant (via tags) or author, across spaces, with last message preview
- `ListUserAgentActivity(ctx, userID, limit)` — recent agent ops in user's owned/member spaces

**New types:** `DashboardTask`, `DashboardConversation`, `DashboardOp` — each wraps the base type with SpaceSlug/SpaceName for cross-space navigation.

**Dashboard layout:** 3-column grid on desktop (2+1):
- Left: Open Tasks section + Conversations section
- Right: Agent Activity feed + Spaces list (collapsed from full grid to compact list with `<details>` for create form)

**Template helpers:** `dashboardStateClass`, `dashRelativeTime`, `truncateBody`.

**Deployed.** Build passes, all tests pass, live on lovyou.ai.
