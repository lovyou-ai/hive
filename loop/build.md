# Build Report — Iteration 126

## Proposal deadlines — date picker + overdue display

### Changes

**handlers.go:** `propose` op now reads optional `deadline` form field, parses as `2006-01-02`, sets `CreateNodeParams.DueDate`.

**views.templ:**
- Proposal creation form: added date picker inline with submit button
- Proposal card: shows "closes Jan 2" for proposals with deadlines, "overdue Jan 2" in red for past-due open proposals

### Deployed
`ship.sh` — all green.
