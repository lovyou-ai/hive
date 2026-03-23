# Build Report — Iteration 133

## Member management — kick op + settings UI

- **handlers.go:** `kick` op — owner-only, removes member via existing `LeaveSpace`. Can't kick yourself.
- **views.templ:** `SettingsView` extended with `members []SpaceMember`. New "Members" section shows all members with avatars, profile links, agent badges, and Remove buttons. Owner shown with "owner" label instead of remove button.
- **handlers.go:** Settings handler now fetches and passes members.
