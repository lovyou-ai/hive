# Build Report — Iteration 102

Notifications system. New `notifications` table (id, user_id, op_id, space_id, message, read, created_at).

**Store:** CreateNotification, ListNotifications, UnreadCount, MarkNotificationsRead.

**Triggers:** `assign` op notifies assignee. `intend` with assignee notifies them. `respond` in conversations notifies all other participants.

**UI:** Unread count badge on dashboard header. `/app/notifications` page with read/unread styling, space links, relative timestamps. Viewing the page marks all as read.

**Dashboard:** Now takes `unreadCount` param.

Shipped via `ship.sh` (foreground). 10 tables total.
