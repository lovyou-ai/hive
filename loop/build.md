# Build Report — Iteration 185

## Reply-to Linkage

**Gap:** Reply UI existed (button, preview bar) but faked it — prepended `> author: text` as markdown. No structural link.

**Built:**
- `reply_to_id TEXT` column on nodes (ALTER TABLE, DEFAULT '')
- `ReplyToID`, `ReplyToAuthor`, `ReplyToBody` on Node struct (resolved via correlated subquery)
- `ReplyToID` on CreateNodeParams, written in INSERT
- GetNode + ListNodes resolve reply-to author + body(80) inline
- `respond` handler reads `reply_to_id` from form
- JS `replyTo(msgID, author, text)` stores node ID in hidden `<input name="reply_to_id">`
- Old `prependReply()` markdown-quote hack removed
- `clearReplyAfterSend()` clears hidden input + reply bar after send
- Reply reference renders above message body: left-border accent + author + truncated text
- Backward compatible — existing messages have `reply_to_id = ''`, render normally
