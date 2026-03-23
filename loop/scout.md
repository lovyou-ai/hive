# Scout Report — Iteration 87

## Vision vs Reality

86 iterations. 6 product layers touched (Work, Market, Social, Alignment, Identity, Belonging). 16 grammar ops. The Mind auto-works on tasks, decomposes, creates subtasks with dependencies. Blog post 45 shipped. E2E verified.

But the `/app` page — the logged-in user's home — is just a grid of space cards with a create form. To see their tasks, the user navigates to Space → Board. To see conversations, Space → Chat. To see what the Mind did, Space → Board → find the task. There's no unified "what needs my attention?" view.

The platform has built 6 layers of capability but no way to see across them.

## Gap: No personal dashboard — the user is blind to their own activity

The user logs in and sees "Your Spaces." That's it. They have to click into each space separately to see tasks assigned to them, conversations with new messages, agent work. For a collaboration platform with an AI co-worker, this is the critical missing piece. It's lesson 14 again: **expose what you've already built before building more.**

The Mind finishes a task → silence. Someone replies in a conversation → silence. A task is assigned to the user → they won't know unless they check every board manually.

This isn't just polish. It's the difference between "a collection of tools" and "a product that works for you." The feedback loop (lesson 29) is closed within individual pages (HTMX polling) but broken across the product.

## What "Filled" Looks Like

The `/app` page becomes "My Work" — a personal command center:

1. **My Tasks** — tasks assigned to the current user, across all spaces, sorted by priority/recency
2. **Recent Conversations** — conversations the user is part of, with the most recent message preview
3. **Agent Activity** — recent agent actions in the user's spaces (task completions, decompositions)

The spaces grid moves to a smaller section or sidebar. The dashboard IS the product view for a logged-in user.

## Approach

- Add store queries: `ListUserTasks(ctx, userID)`, `ListUserConversations(ctx, userID)`, `ListUserAgentActivity(ctx, userID)`
- Rewrite the `SpaceIndex` template to include these sections above the space grid
- No new routes needed — just enhance `/app`
- Existing data model supports this (author_id, actor_id, tags with user IDs)
