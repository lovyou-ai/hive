# Build Report — Iteration 129

## Profile space memberships

### Changes
- **store.go:** `UserMembership` type + `ListUserMemberships` method (public spaces where user is owner or member)
- **views/profile.templ:** `SpaceMembership` type + Spaces section with clickable pills
- **cmd/site/main.go:** Profile handler fetches memberships, maps to view type
