# Build Report — Social Competitive Research

**Task:** Deep research: what makes each social platform great, feature by feature

**Builder:** Claude (Sonnet 4.6) · 2026-03-26

---

## What Was Built

Created `loop/social-competitive-research.md` — comprehensive competitive analysis of 7 social platforms mapped to our 15 grammar operations, with explicit identification of where our architecture is genuinely better.

---

## Artifact: `loop/social-competitive-research.md`

### Structure

1. **Methodology** — how to read the analysis (✓/~/✗/+ markers, grammar mapping)
2. **15 Grammar Operations reference table** — all 15 ops with descriptions
3. **7 Platform analyses** — Twitter, Discord, Reddit, Slack, Instagram, TikTok, Messenger
4. **Cross-platform synthesis** — table stakes, structural advantages, unified meta-advantage
5. **Feature priority map** — P0/P1/P2/P3 priorities derived from research

### Per-Platform Coverage

| Platform | Features Documented | Key Grammar Ops |
|----------|-------------------|-----------------|
| Twitter/X | 12 features | Emit, Propagate, Endorse, Annotate, Derive, Acknowledge |
| Discord | 11 features | Channel, Delegate, Consent, Sever, Subscribe |
| Reddit | 9 features | Endorse, Respond, Channel, Delegate, Retract |
| Slack | 11 features | Channel, Acknowledge, Subscribe, Delegate, Derive |
| Instagram | 9 features | Emit(ephemeral), Annotate, Channel, Derive |
| TikTok | 8 features | Propagate, Derive, Annotate, Channel |
| Messenger | 9 features | Channel, Consent, Acknowledge, Delegate |

### Key Findings

**Six structural advantages we have that no platform has:**
1. Causal provenance on all content (cryptographic Derive chains)
2. Identity-linked endorsements with Trust Score weighting
3. Governance as a first-class primitive (Propose + Vote + Delegate)
4. Agents as peers with identity and Trust Score (not second-class bots)
5. Append-only truth (Retract/Extend never delete — history preserved)
6. Cross-layer Work + Social integration on same graph (no context switching)

**P0 gaps — must match competition to be adoption-viable:**
- @mention parsing + notifications
- Unread counts on channels
- In-chat message search
- Message grouping (visual blocks for consecutive messages)
- Notification grouping by type

**Grammar operations most underdeveloped:**
- **Annotate** — hashtags, community notes, flair, alt text are all Annotate compositions
- **Consent** — structured join consent, approval-mode channels, declared purpose
- **Delegate** — moderation with audit log, role-gated channels, speaker queues

---

## Build/Test

Research artifact only — no code changes.

`go.exe build -buildvcs=false ./...` — unchanged.
`go.exe test ./...` — unchanged.

---

## Files Changed

| File | Action |
|------|--------|
| `loop/social-competitive-research.md` | Created |
| `loop/build.md` | Updated |
