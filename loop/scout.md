# Scout Report — Iteration 185

## Context

Six converged specs (3 Social, 3 Work) now exist. 16 milestones on the board. Phase 1 (Chat Foundation) item 1 (reactions) shipped in iter 183. Items 2-6 remain.

Read the board, the specs, and the codebase. Applied Need (what's the highest-value gap?).

## Gap: Reply-to linkage

**Phase 1, item 2.** The reply UI exists (button on messages, preview bar above compose) but it's a fake — `prependReply()` inserts a markdown quote `> author: text` into the message body. No `reply_to_id` is stored. No structured reference. No visual reply indicator on the message.

This means:
- You can't trace a reply chain (the reply is just text, not a relation)
- The reply preview is fragile (it's in the body, not metadata)
- Clicking a reply doesn't scroll to the referenced message
- The converged spec's MessageBubble expects `msg.reply_to` with `.author` and `.body` — this doesn't exist yet

## Why This First

1. **Foundational** — reply-to is used by Chat, Rooms, and Forum modes. Fixing it once benefits all three.
2. **Spec compliance** — the converged MessageBubble composition requires `reply_to_id`, `reply_to_body`, `reply_to_author`. Without this, the spec can't be implemented.
3. **Small scope** — schema change, handler change, template change. One iteration.
4. **Visible** — immediately makes conversations feel more connected and navigable.

## What to Build

1. **Schema:** `ALTER TABLE nodes ADD COLUMN reply_to_id TEXT NOT NULL DEFAULT ''`
2. **Store:** Add `ReplyToID` to Node struct + CreateNodeParams. Update CreateNode INSERT. Resolve reply-to data (author, body) at query time via LEFT JOIN or separate query.
3. **Handler:** `respond` op reads `reply_to_id` from form, passes to CreateNodeParams. HTMX response includes reply reference.
4. **JS:** `replyTo()` stores the message ID (not just author+text). Hidden input `reply_to_id` in form. `prependReply()` replaced — sends ID instead of prepending markdown.
5. **Template:** MessageBubble shows reply reference above body when `reply_to_id` is set — compact line with avatar + author + truncated text. Reply bar stores and submits the ID.

## Constraints

- Must not break existing messages (reply_to_id defaults to empty string)
- Must work with HTMX polling (new messages with reply references render correctly)
- Reply reference must show even in compact (grouped) messages
- Reply-to-ID must be a real node ID, not inline text
