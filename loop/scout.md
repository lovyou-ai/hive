# Scout Report — Iteration 102

## Gap: No notifications — the platform is pull-only

The dashboard shows tasks and conversations but only when you visit it. There's no way to know something happened without manually checking. The Mind completes a task → silence. Someone replies to your conversation → silence. You get assigned a task → silence.

This is the #1 usability gap for a collaboration platform with an active AI agent. The feedback loop is closed within pages (HTMX polling) but broken across sessions.

## What "Filled" Looks Like

An unread count badge on the nav that links to a notifications page. Notifications generated from ops directed at the user: task assignments, conversation messages, task completions by agents in your spaces. Simple read/unread tracking.
