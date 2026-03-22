# Scout Report — Iteration 35

## Map

Conversations: primitive (31), interface (32), participant (33), live polling (34). The full stack works. But the experience between sending a message and receiving a reply is dead air — no typing indicator, no feedback, no signal that the Mind is processing. The human sends a message and stares at silence for 10-30 seconds.

Also: messages start at the top on page load (no scroll-to-bottom), and you have to click Send instead of pressing Enter.

## Gap Type

Missing UX feedback (the system works but feels broken)

## The Gap

The conversation view has no presence indicator. When a human sends a message in a conversation with agent participants, nothing signals that a response is coming. The 3-second polling interval + Mind generation time (10-30s) creates a silence gap that feels like the system is broken.

Secondary: no scroll-to-bottom on initial load, no enter-to-send.

## Why This Gap Over Others

This is lesson 29 continued — the feedback loop technically closes (polling picks up responses) but the *experience* of the feedback loop is broken. The user's mental model is: "I sent a message, is anyone there?" Without a thinking indicator, the product feels abandoned.

End-to-end testing of `cmd/reply` is blocked by ANTHROPIC_API_KEY. Conversation types add structure to something that doesn't feel right yet. Fix the feeling first.

## What "Filled" Looks Like

After the user sends a message in a conversation with agent participants, a "thinking" bubble appears (animated dots, violet styling to match agent identity). The bubble has a 60-second timeout — if no agent response arrives, it fades out. When the poll picks up a new agent message, the thinking bubble is removed and replaced by the real response.

Also: page scrolls to bottom on load. Enter key sends messages.
