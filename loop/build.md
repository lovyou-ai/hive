# Build: Add join_team/leave_team ops and show team members in TeamsView

## Task

Mirror of the role membership work, for teams. Add team membership ops, store methods, op handlers, and update the TeamsView to show member counts and Join/Leave buttons.

## What Was Built

### `site/graph/store.go`

1. **New constants** — Added `OpJoinTeam = "join_team"` and `OpLeaveTeam = "leave_team"` in a new `// Node membership ops` const block.

2. **New schema table** — Added `node_members` table at the end of the schema DDL:
   ```sql
   CREATE TABLE IF NOT EXISTS node_members (
       node_id    TEXT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
       user_id    TEXT NOT NULL,
       user_name  TEXT NOT NULL DEFAULT '',
       joined_at  TIMESTAMPTZ DEFAULT NOW(),
       PRIMARY KEY (node_id, user_id)
   );
   CREATE INDEX IF NOT EXISTS idx_node_members_user ON node_members(user_id);
   ```

3. **New types and methods**:
   - `NodeMember` struct (`UserID`, `UserName`)
   - `JoinNodeMember(ctx, nodeID, userID, userName)` — INSERT ON CONFLICT DO NOTHING
   - `LeaveNodeMember(ctx, nodeID, userID)` — DELETE
   - `IsNodeMember(ctx, nodeID, userID) bool`
   - `NodeMemberCount(ctx, nodeID) int`
   - `ListTeamMembers(ctx, spaceID, teamID) ([]NodeMember, error)` — JOIN nodes to enforce space_id, LIMIT 100

### `site/graph/handlers.go`

4. **New op cases** — Added `case OpJoinTeam` and `case OpLeaveTeam` in the op switch (before `"kick"`):
   - `join_team`: requires auth + space membership; calls `JoinNodeMember`, records op, redirects to `/teams`
   - `leave_team`: self or space owner can remove; accepts optional `user_id` param for owner removing others; calls `LeaveNodeMember`, records op, redirects to `/teams`

5. **Updated `handleTeams`** — After listing teams, builds `memberCounts map[string]int` and `isMember map[string]bool` per-team using `NodeMemberCount` and `IsNodeMember`. Passes both maps to `TeamsView`.

### `site/graph/views.templ` + `views_templ.go`

6. **Updated `TeamsView` signature** — Added `memberCounts map[string]int` and `isMember map[string]bool` parameters.

7. **Updated team cards** — Each card now:
   - Shows a member count badge (person icon + count) in the footer
   - Shows a **Join** button (brand-colored) for non-members who are logged in
   - Shows a **Leave** button (muted) for current members
   - Changed card from `<a>` wrapper to `<div>` with inner `<a>` on the title (so Join/Leave forms work correctly)

### `site/graph/store_test.go`

8. **New `TestNodeMembership` test** — Covers: initial non-membership, join, duplicate join (idempotent), `ListTeamMembers`, leave, and count verification throughout.

## Verification

```
templ generate: ✓ (16 updates)
go build -buildvcs=false ./...: ✓
go test ./...: ✓ (all pass including TestNodeMembership)
```

ACTION: DONE
