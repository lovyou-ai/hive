# Scout Report — Iteration 105

## Gap: Space landing page is a blind redirect — no overview or context

When you visit `/app/{slug}`, you're immediately redirected to board or feed. There's no overview — no stats, no pinned content, no sense of what this space is about before diving into a lens. For a visitor arriving from Discover or a shared link, the first impression is a raw kanban or feed with no context.

## What "Filled" Looks Like

`/app/{slug}` shows a space overview: name, description, member count, pinned items, recent activity summary, and quick links to each lens. The redirect to board/feed becomes a fallback for spaces with no description.
