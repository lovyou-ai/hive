# Documentation Audit

Derivation method applied to the hive documentation itself. Last run: 2026-03-09.

## Method

1. **Identified purpose:** The docs must fully specify the hive such that humans can understand it, agents can operate within it, developers can build it, and contributors can evaluate it.
2. **Named base operations:** Define, Justify, Specify, Relate, Constrain, Sequence.
3. **Identified seven dimensions:** Audience, Abstraction, Time, Scope, Agency, Trust, Economics.
4. **Traversed each dimension** to find gaps between what should exist and what does exist.

## Current Documents

| Doc | Purpose | Lines | Status |
|-----|---------|-------|--------|
| CLAUDE.md | Agent context (read at every prompt) | 198 | Solid |
| README.md | Entry point for humans | 37 | Adequate |
| VISION.md | Why the hive exists | 155 | Solid |
| ARCHITECTURE.md | How it's built | 136 | Solid |
| ROADMAP.md | What's done, what's next | 241 | Solid |
| AGENT-RIGHTS.md | Rights, invariants, governance | 175 | Solid |
| AGENT-TOOLS.md | MCP server, agentic loop | 206 | Solid |

## Gaps Found

### High Priority

#### 1. Event Type Catalog
**Blocking for:** Tier 1 (MCP server)
**Gap:** No catalog of event types, their schemas, or when agents should emit/query them. The eventgraph has event types scattered across `event/constants.go`, `event/agent_content.go`, `event/codegraph_content.go`, but the hive doesn't document what its own events look like.
**Need:** A doc or section listing: event type → schema → who emits → who consumes → authority required. This is the contract between agents and the graph.
**Where:** `docs/EVENT-TYPES.md` or a section in AGENT-TOOLS.md

### Medium Priority

#### 2. Trust Dynamics
**Blocking for:** Tier 2 (agentic loop — agents need to evaluate trust to make authority decisions)
**Gap:** Trust is mentioned everywhere but never specified precisely. How fast does it accumulate? What actions build trust? What's the decay rate? What triggers trust loss? What specific actions shift from Required → Recommended → Notification?
**Need:** Concrete trust mechanics — not the philosophy (well-covered) but the mechanism.
**Where:** `docs/TRUST.md` or a section in AGENT-RIGHTS.md

#### 3. Human Operator Guide
**Blocking for:** Tier 2 (web dashboard) and any new human operators
**Gap:** No doc for the human. What does day-to-day operation look like? What do approval requests look like? What should you scrutinise? How do you build trust with agents? How do you interpret the event graph?
**Need:** Practical guide for the human side of the authority model.
**Where:** `docs/OPERATOR.md`

### Low Priority (Note Now, Build Later)

#### 4. Inter-Agent Dynamics
**Blocking for:** Tier 2+ (agentic loop with multiple self-directing agents)
**Gap:** How agents relate to each other — trust between agents, delegation chains, conflict resolution between agents (not just human-agent). The 28 primitives in eventgraph cover this but hive docs don't reflect it.
**Where:** Part of agentic loop spec when built

#### 5. Product Derivation Pattern
**Blocking for:** Tier 5 (first products)
**Gap:** The pattern by which each of the 13 products will be derived from its composition grammar using the derivation method. Not the products themselves — the method.
**Where:** Note in ROADMAP.md, full spec when Tier 5 begins

#### 6. Agent Growth Model
**Blocking for:** Tier 2+ (agentic loop)
**Gap:** How agents learn within a lifetime — decision tree evolution, memory accumulation, skill development. Distinct from self-modification (which changes the codebase).
**Where:** Part of agentic loop spec when built

## Dimension Coverage Matrix

| Dimension | Range | Coverage | Gap Area |
|-----------|-------|----------|----------|
| Audience | Human ↔ Agent ↔ Dev ↔ Contributor | Good for agents/devs. **Thin for human operators and contributors.** | Operator guide, CONTRIBUTING |
| Abstraction | Vision ↔ Architecture ↔ Spec ↔ Code | Good. **Missing event type spec layer.** | Event catalog |
| Time | Now ↔ Next ↔ Eventually | Good. Tiers 5-6 thin (expected). | Product derivation pattern |
| Scope | Soul ↔ Civilisation ↔ Product ↔ Feature | Good. Per-product specs not yet needed. | Note in roadmap |
| Agency | Rights ↔ Tools ↔ Governance ↔ Lifecycle | Good. **Inter-agent dynamics thin.** | Agent-to-agent relations |
| Trust | Zero ↔ Accumulating ↔ High | Philosophy good. **Mechanics thin.** | Trust dynamics doc |
| Economics | Revenue ↔ Cost ↔ Transparency ↔ Sustainability | Good. Pricing not yet specified (expected). | — |

## Consistency Check

Checked all docs for contradictions:
- Soul statement: consistent across all 7 docs ✓
- Authority model (Required/Recommended/Notification): consistent ✓
- 8 rights: consistent between CLAUDE.md and AGENT-RIGHTS.md ✓
- 10 invariants: consistent between CLAUDE.md and AGENT-RIGHTS.md ✓
- Neutrality clause: consistent between CLAUDE.md, VISION.md, AGENT-RIGHTS.md ✓
- Revenue model: consistent across VISION.md, ROADMAP.md, CLAUDE.md ✓
- Build sequence: ROADMAP tiers align with AGENT-TOOLS phases ✓
- Intelligence assignments: CLAUDE.md roles table matches ARCHITECTURE.md intelligence table ✓

## Redundancy Check

Intentional redundancy (agent context needs self-contained summary):
- CLAUDE.md duplicates key content from VISION, ARCHITECTURE, AGENT-RIGHTS. This is correct — CLAUDE.md is read by agents at every prompt and must be self-contained.

Unintentional redundancy:
- Revenue model appears in VISION.md (section), ROADMAP.md (section), and CLAUDE.md (inline). All slightly different wording but consistent content. Could consolidate if it drifts, but acceptable for now.
- Resource transparency described in both VISION.md and ROADMAP.md gap #11. VISION has the "why", ROADMAP has the "what/where". Correct separation.

## Next Actions

1. **Now:** Add event type catalog note to ROADMAP.md Tier 1 (blocking for MCP)
2. **Now:** Add product derivation pattern note to ROADMAP.md Tier 5
3. **Tier 1:** Write EVENT-TYPES.md when building MCP server
4. **Tier 2:** Write TRUST.md with concrete mechanics
5. **Tier 2:** Write OPERATOR.md for human operators
6. **Tier 2:** Expand inter-agent dynamics in agentic loop spec
