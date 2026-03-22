# Critique — Iteration 35

## Verdict: APPROVED

## Trace

1. Scout identified: dead air after sending a message in agent conversations
2. Builder added thinking indicator (violet dots), scroll-to-bottom, enter-to-send
3. Director feedback mid-build: "an actor is either agent or human" — shifted from message scanning to identity-based agent detection via `HasAgentParticipant`
4. Second deploy with the architectural fix
5. Deployed and healthy

## Audit

**Correctness:**
- Agent presence resolved from users table (`kind = 'agent'`), not from scanning messages. Works for new conversations with no agent messages yet. ✓
- Thinking indicator shows only in conversations with agent participants (`data-has-agent`). ✓
- Indicator hides on poll receiving new messages OR after 60-second timeout. No stale indicators. ✓
- Enter-to-send respects Shift+Enter (doesn't submit). ✓
- Scroll-to-bottom fires once on page load via inline script. ✓

**Gaps:**
- **False positive**: The thinking indicator shows even when `cmd/reply` isn't running. This is honest — it says "an agent may respond" not "an agent is responding." But it could train users to expect responses that don't come if nobody runs `cmd/reply`.
- **No poll-triggered agent detection**: If an agent joins a conversation after page load (first message from agent via poll), the `data-has-agent` attribute isn't updated dynamically. The thinking indicator won't show until page reload. Minor — conversations don't gain new agent participants mid-session in practice.
- **60-second timeout is arbitrary**: Could be too short (Mind takes longer) or too long (user thinks system is broken). Acceptable for now.

## DUAL

The director's mid-build feedback ("actor is either agent or human") caught the same pattern as iteration 33 ("who's Hive?"): the builder reached for the data when the identity system already had the answer. The architectural fix was small (10 lines) but correct — agent presence is a property of the participants, not the messages.
