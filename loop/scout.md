# Scout Report — Iteration 189

## Gap: Message Search (Phase 1 item 6 — final Chat Foundation item)

**Source:** social-spec.md Phase 1, board milestone. Last remaining item before Phase 2.

**Current state:** Chat lens has `?q=` search that filters conversation *titles* only (line 653 of handlers.go). No way to search message *bodies* across conversations. No operator syntax.

**What's needed:**
1. `SearchMessages(spaceID, query, limit)` store method — ILIKE on message bodies, returns messages with conversation context (convo title, author name, space slug)
2. Operator parsing: `from:username` filters by author, plain text matches body
3. Chat lens search: when query contains text, search messages (not just conversation titles) and show results with context
4. Results link to the conversation containing the match

**Approach:** Follow existing pattern — every lens already has `?q=` search. Chat gets message-level search. Store does the heavy lifting with a JOIN query (nodes as messages JOIN nodes as conversations). View renders results as message cards with conversation context.

**Composition:** `Search(Navigate(messages))` — traverse the message graph, filter by query. Uses the same `ListNodesParams.Query` ILIKE pattern already in the codebase.

**Risk:** Low. New store method, handler logic, template section. No schema changes. No breaking changes.

**Estimated scope:** ~100 lines store + handler, ~50 lines template.
