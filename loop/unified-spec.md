# Unified Ontology

**Everything is organized activity. Work and Social are facets of one thing, not separate products.**

Matt Searles + Claude · March 2026

---

## The Root

The product is not "a task tracker" or "a social network." It's a platform for **purposeful collective activity** at every scale. Activity requires:

- **Actors** (Identity, Bond) — who does it
- **Structure** (Belonging, Organize) — how they're grouped
- **Coordination** (Communication) — how they align
- **Execution** (Work) — what they do
- **Direction** (Planning, Goals) — why they do it
- **Control** (Governance, Justice, Moderation) — what constrains it
- **Memory** (Knowledge, Culture) — what's retained
- **Exchange** (Market, Allocation) — what's traded
- **Accountability** (Alignment, Transparency) — what's visible
- **Purpose** (Being, Reflection) — what it's for

Every one of these is a facet of one phenomenon, not a separate product. The 13 EventGraph layers are 13 angles on purposeful collective activity.

---

## Why Work is the Root (or Close to It)

Social exists to serve coordination. Coordination exists to serve organized activity. You don't follow someone for its own sake — you follow them because their activity is relevant to yours. You don't endorse a post abstractly — you vouch for its contribution to shared understanding. Every social action has a work context, even when implicit.

But Work alone is insufficient. Pure task execution without communication, governance, or knowledge is a factory line. The product's value is that all facets coexist on one graph — the task, the conversation about the task, the policy that governs the task, the knowledge produced by the task, and the people involved are all events on the same chain.

**Work is the gravitational center. The other facets orbit it.** Not because they're subordinate, but because purposeful collective activity is the thing that connects them all.

---

## The Derivation Order

From the generator function, applying Decompose to "purposeful collective activity":

```
Level 0: Being
  "I exist, I have purpose"

Level 1: Identity + Bond
  "I am someone, I relate to others"

Level 2: Belonging + Organize
  "I join groups, groups have structure"

Level 3: Communication (Social modes)
  "I coordinate with others in my groups"
  Modes: Chat, Rooms, Square, Forum

Level 4: Work (Activity modes)
  "I organize activity toward outcomes with others"
  Modes: Execute, Organize, Govern, Plan, Learn, Allocate

Level 5: Knowledge + Market + Build
  "Activity produces artifacts, exchanges value, builds tools"

Level 6: Alignment + Culture + Justice
  "Activity requires transparency, norms, and dispute resolution"
```

This isn't a build order — it's a dependency order. You can't do Work (L4) without Communication (L3). You can't Communicate without Belonging (L2). But in practice, you build them interwoven.

---

## Merged Mode Set

Work has 6 modes. Social has 4 modes. But they're not separate — they're all modes of organized activity:

| Mode | What it does | Source | Primary layer facets |
|------|-------------|--------|---------------------|
| **Execute** | Do the work | Work | Work, Build |
| **Chat** | Real-time coordination | Social | Social, Bond |
| **Rooms** | Persistent group communication | Social | Social, Belonging |
| **Square** | Public broadcast + engagement | Social | Social, Identity, Market |
| **Forum** | Threaded discussion + quality | Social | Social, Knowledge |
| **Organize** | Structure people + groups | Work | Belonging, Identity, Bond |
| **Govern** | Policies, approvals, compliance | Work | Governance, Justice, Moderation |
| **Plan** | Goals, roadmaps, direction | Work | Work, Alignment |
| **Learn** | Retrospectives, knowledge capture | Work | Knowledge, Being, Culture |
| **Allocate** | Budgets, resources, capacity | Work | Market, Work |

**10 modes total.** Each mode is a lens on the same graph. A space can use any combination. A solo dev uses Execute + Chat. A Fortune 500 uses all 10.

---

## The Unified Entity Set

Merging Work entities (iter 201) with Social entities:

| Entity | Kind | What | Created by | Used by modes |
|--------|------|------|-----------|--------------|
| Task | `task` | Atomic work unit | Execute | Execute, Plan, Allocate |
| Post | `post` | Broadcast content | Square | Square, Forum, Learn |
| Thread | `thread` | Discussion topic | Forum, Rooms | Forum, Rooms, Learn |
| Message | `comment` | Communication unit | Chat, Rooms, Forum | All communication modes |
| Conversation | `conversation` | Real-time channel | Chat | Chat, Execute (task discussion) |
| Project | `project` | Scoped work collection | Plan | Execute, Allocate, Learn |
| Goal | `goal` | Desired outcome | Plan | Plan, Execute, Allocate |
| Role | `role` | Capability + responsibility | Organize | Organize, Govern, Execute |
| Team | `team` | Functional group | Organize | All modes |
| Department | `department` | Organizational unit | Organize | Organize, Govern, Allocate |
| Policy | `policy` | Governing rule | Govern | Govern, Organize, Learn |
| Process | `process` | Repeatable sequence | Govern | Govern, Execute, Learn |
| Decision | `decision` | Choice with rationale | Govern, Plan | Govern, Plan, Learn |
| Resource | `resource` | Consumable | Allocate | Allocate, Execute, Plan |
| Document | `document` | Knowledge artifact | Learn | Learn, Govern, Plan |
| Claim | `claim` | Knowledge assertion | Learn | Learn, Govern |
| Proposal | `proposal` | Governance proposal | Govern | Govern |
| Organization | `organization` | Legal/structural entity | Organize | All modes |

**18 entity types.** All are Nodes. All use the same grammar ops. The kind determines which modes surface them.

---

## How the Sidebar Should Work

Not "Work" and "Social" as separate sections. Instead, the sidebar presents available modes for the current space:

```
My Work (dashboard)

Modes:
  Execute (Board, List, Triage)
  Chat
  Rooms (when channels exist)
  Feed (Square)
  Forum (when threads exist)
  Knowledge
  Governance
  Build
  Transparency

Organization: (when org entities exist)
  Teams
  Roles
  Policies
  Goals
  Resources

Spaces:
  ...
```

Modes appear when they have content. A fresh space shows Execute + Chat + Feed. As the space gains structure (teams, policies, goals), more modes appear.

---

## The Product at Each Scale

| Scale | Active modes | Key entities | What it feels like |
|-------|-------------|-------------|-------------------|
| **Solo dev** | Execute, Chat (with agent) | Tasks, Conversations | Todo app + AI pair programmer |
| **Small team (5)** | Execute, Chat, Feed, Forum | Tasks, Posts, Threads, Conversations | Linear + Slack in one place |
| **Startup (20)** | + Plan, Organize | + Projects, Goals, Roles, Teams | All-in-one workspace |
| **Mid-size (200)** | + Govern, Learn | + Policies, Decisions, Documents | Complete org infrastructure |
| **Enterprise (2000+)** | + Allocate, Rooms | + Departments, Resources, Budgets | Enterprise platform |
| **Civilizational** | All 10, across Organizations | All 18 entities, inter-org | Operating system for human activity |

The product doesn't have tiers or feature gates. The modes emerge from what entities exist. When you create a Policy node, the Govern mode appears. When you create a Goal, the Plan mode surfaces. Complexity is earned, not purchased.

---

## Grammar Coverage Across All Modes

The grammar operations are mode-independent. Here's how each manifests:

| Op | Execute | Chat | Square | Forum | Organize | Govern | Plan | Learn | Rooms | Allocate |
|----|---------|------|--------|-------|----------|--------|------|-------|-------|----------|
| **Intend** | Create task | — | Create post | Create thread | Create role | Draft policy | Set goal | Start retro | — | Request budget |
| **Respond** | Comment | Reply | Reply | Nested comment | — | Comment | — | — | Reply | — |
| **Decompose** | Subtasks | — | — | — | Sub-teams | Provisions | Milestones | Root causes | — | Line items |
| **Assign** | Assign person | — | — | — | Fill role | Assign reviewer | Own goal | Action items | — | Budget owner |
| **Claim** | Self-assign | — | — | — | Volunteer | — | — | — | — | Claim resource |
| **Complete** | Mark done | — | — | Solve | — | Ratify | Mark met | Close retro | — | Close period |
| **Review** | Approve work | — | — | — | Performance | Audit | Review progress | Peer review | — | Audit spend |
| **Delegate** | Assign agent | Transfer | — | Mod | Appoint | Delegate auth | — | — | Mod | Delegate |
| **Consent** | — | Decision | Poll | — | — | Vote | — | — | Vote | Approve |
| **Endorse** | Vouch | Endorse | Endorse | Upvote | Recommend | Support | Endorse | Validate | Endorse | — |
| **Propagate** | — | Forward | Repost | Cross-post | — | — | — | — | Cross-post | — |
| **Subscribe** | Watch | Join | Follow | Join | Join team | — | — | — | Join | — |
| **Scope** | Permissions | — | — | — | Responsibilities | Jurisdiction | Success criteria | — | — | Spending limit |
| **Reflect** | — | — | — | — | — | — | — | Post-mortem | — | — |

**The grammar is the API.** Every cell in this matrix is a `POST /app/{slug}/op` with the same handler. The op name determines the semantics. The node kind determines which mode it belongs to.

---

## Convergence Analysis

**Pass 1 — Need:**
- Current product treats Work and Social as separate sidebar sections
- 10 of 18 entity types exist (task, post, thread, comment, conversation, claim, proposal + the 3 implicit: space, user, op)
- 8 entity types missing (project, goal, role, team, department, policy, process, decision, resource, document, organization — some overlap with what spaces already are)
- Modes are implicit in lenses, not explicit in sidebar

**Pass 2 — Traverse:**
- Spaces already function as proto-organizations/teams
- Membership already functions as proto-belonging
- The sidebar already groups by "layers" — one refactor away from grouping by "modes"
- Grammar ops already work on any node kind — adding new kinds is trivial
- The existing product is 60% of the way to the unified ontology

**Fixpoint at pass 2.** The architecture IS the unified ontology. The gap is in naming and UI organization, not in data model or operations.

---

## Relationship to Existing Specs

| Spec | Relationship |
|------|-------------|
| `social-spec.md` | Becomes the Communication modes (Chat, Rooms, Square, Forum). Still correct. |
| `social-product-spec.md` | Still correct for product positioning of Social features. |
| `work-product-spec.md` | Becomes the Execute mode spec. Still correct but narrower than reality. |
| `work-spec.md` | Becomes the Execute mode compositions. Still correct. |
| `work-general-spec.md` | The Work-specific expansion. Subsumed by this unified spec. |
| **This spec** | The structural document that shows how everything relates. |

Nothing is discarded. Everything is placed in context.
