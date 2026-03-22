# Scout Report — Iteration 31

## Map

30 iterations. Site is feature-complete with 5 lenses. Mind CLI exists (iter 30) but Matt redirected: the Mind should be a web participant, not a CLI. The site already has agent identity, threads, people. What's missing is a conversation primitive — named, multi-participant, distinct from public threads.

## Gap Type

Missing product primitive (needs building)

## The Gap

No conversation infrastructure. Threads are public discussion forums. The vision is Slack + AI chat: multiple named conversations per person, DMs, groups, rooms. The Mind as a participant. Human-agent duo — when a human messages, their agent has right of reply too.

## Why This Gap Over Others

The conversation primitive is the foundation for everything: Mind web presence, human-agent duo communication, social product. Without it, the Mind has no channel to communicate through on the web.

## What "Filled" Looks Like

`kind='conversation'` node type, `converse` grammar op, Chat lens in sidebar, `ConversationsView` template with create form and conversation list. Conversations store participants in `tags[]`, messages are child comment nodes.
