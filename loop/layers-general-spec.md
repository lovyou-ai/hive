# The Thirteen Layers — Generalized

**Each layer expanded from a feature into a full domain via cognitive grammar. All layers share one graph, one grammar, one set of entity kinds.**

Matt Searles + Claude · March 2026

---

## Method

For each layer: apply Distinguish (what entities exist in this domain at every scale?), Relate (how do they connect to other layers?), Select (what's the minimum viable expansion?). The root is collective existence, not productivity.

---

## Layer 0: The Graph

Not a product layer. The substrate. Events, causality, trust, authority. Everything below runs on this.

**Already complete.** EventGraph provides: event storage, hash chaining, signing, causal links, trust scores, authority levels.

---

## Layer 1: Work — Organized Activity

**Current:** Tasks, Board, List, Projects, Goals. 7 ops.

**General domain:** All organized activity toward outcomes, at every scale.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Task | `task` | Atomic work unit | ✓ exists |
| Project | `project` | Scoped collection | ✓ exists |
| Goal | `goal` | Desired outcome | ✓ exists |
| Milestone | `milestone` | Measurable checkpoint | Between goal and task |
| Sprint/Cycle | `cycle` | Time-boxed work period | "Q2 Sprint 3", "March cycle" |

**Missing ops:** review, handoff, scope (from work-product-spec.md)

**Cross-layer:** Tasks discussed in Social (Chat). Tasks governed by Policies (Governance). Task completion builds Reputation (Identity). Tasks consume Resources (Market).

---

## Layer 2: Market — Exchange

**Current:** Available tasks page, claim op, prioritize op.

**General domain:** All value exchange between actors. Not just "task marketplace" — any flow of value.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Resource | `resource` | Consumable unit | Budget, compute hours, headcount, materials |
| Budget | `budget` | Allocated resource pool | "$50k Q2 engineering", "10,000 GPU hours" |
| Contract | `contract` | Agreement between parties | SLA, employment, vendor agreement |
| Invoice | `invoice` | Record of exchange | Billing, payment, reimbursement |
| Listing | `listing` | Available offering | Job posting, service offering, product listing |

**Ops needed:** bid, offer, accept, transfer, invoice, settle

**Cross-layer:** Resources consumed by Work. Contracts governed by Governance. Trust from Identity affects pricing. Exchange recorded for Alignment.

---

## Layer 3: Social — Connection and Communication

**Current:** Chat, Feed (4 tabs), Threads, People. Follow, endorse, repost, quote, reactions, message search.

**General domain:** All forms of connection and communication between beings. See social-spec.md for detailed compositions.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Post | `post` | Broadcast content | ✓ exists |
| Thread | `thread` | Discussion topic | ✓ exists |
| Conversation | `conversation` | Real-time channel | ✓ exists |
| Message | `comment` | Communication unit | ✓ exists |
| Channel | `channel` | Persistent group chat (Rooms mode) | #general, #engineering, #random |
| Event | `event` | Scheduled gathering | Meetup, standup, all-hands, conference |

**Missing modes:** Rooms (Discord-like channels), Forum (Reddit-like quality discussion)

**Cross-layer:** Communication coordinates Work. Endorsements build Identity. Following creates Bond. Channels create Belonging.

---

## Layer 4: Justice — Dispute Resolution

**Current:** Report op, resolve op.

**General domain:** All mechanisms for resolving conflict, enforcing rules, and maintaining fairness.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Report | `report` | Complaint/flag | ✓ exists (as op payload) |
| Case | `case` | Dispute proceeding | Moderation case, HR complaint, legal dispute |
| Ruling | `ruling` | Decision on a case | Verdict, settlement, mediation outcome |
| Appeal | `appeal` | Challenge to ruling | Second review, escalation |
| Precedent | `precedent` | Reusable ruling pattern | "This type of content violates policy X" |

**Ops needed:** file, adjudicate, appeal, enforce, pardon

**Cross-layer:** Cases reference Policy violations (Governance). Rulings affect Trust (Identity). Appeals use Consent (Governance). Precedents become Knowledge.

---

## Layer 5: Build — Creation

**Current:** Changelog lens (completed tasks as build history).

**General domain:** All creation and development activity — not just software.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Release | `release` | Shipped artifact | v1.0, "March update", product launch |
| Artifact | `artifact` | Created output | Code commit, design file, built product, report |
| Review | `review` | Quality assessment | Code review, design review, QA pass |
| Incident | `incident` | Something broke | Bug, outage, safety incident, quality failure |

**Ops needed:** release, review (approve/revise/reject), rollback

**Cross-layer:** Builds are Work completion. Reviews are Knowledge assertions. Incidents trigger Justice. Releases are celebrated in Social.

---

## Layer 6: Knowledge — Understanding

**Current:** Claims, assert/challenge ops, Knowledge lens, evidence trail.

**General domain:** All mechanisms for establishing, validating, and sharing what is known.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Claim | `claim` | Knowledge assertion | ✓ exists |
| Document | `document` | Knowledge artifact | Spec, ADR, handbook, research paper, wiki page |
| Question | `question` | Open inquiry | FAQ, research question, "how do we..." |
| Definition | `definition` | Shared meaning | Glossary entry, term definition, concept |
| Lesson | `lesson` | Learned insight | Post-mortem finding, best practice, anti-pattern |

**Ops needed:** verify, retract (partially exist), cite, define, answer

**Cross-layer:** Knowledge supports Work decisions. Documents govern via Policy. Questions answered in Social. Lessons improve Build.

---

## Layer 7: Alignment — Transparency

**Current:** Activity feed (global audit trail).

**General domain:** All mechanisms for making activity visible, accountable, and auditable.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Report (analytics) | `analytics` | Aggregated view | Dashboard, metric, KPI |
| Audit | `audit` | Compliance check | SOC2 audit, financial audit, code audit |
| Disclosure | `disclosure` | Voluntary transparency | Financial report, impact report, status update |

**Ops needed:** audit, disclose, flag

**Cross-layer:** Alignment makes all other layers visible. The activity feed already captures ops from every layer. Analytics aggregate Work metrics. Audits verify Governance compliance.

---

## Layer 8: Identity — Selfhood

**Current:** User profiles, agent badges, endorsements, action history.

**General domain:** All aspects of who someone is, what they can do, and what they've earned.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Role | `role` | Capability + responsibility | "Engineer", "Moderator", "Delivery driver" |
| Credential | `credential` | Verified capability | Certificate, qualification, clearance |
| Reputation | `reputation` | Earned trust score | Work reputation, community standing |
| Badge | `badge` | Achievement marker | "100 tasks completed", "Trusted reviewer" |

**Ops needed:** certify, revoke, attest

**Cross-layer:** Identity earned through Work (completed tasks). Endorsed via Social. Governed by Policy. Reputation affects Market access.

---

## Layer 9: Bond — Relationship

**Current:** Endorsements, follows.

**General domain:** All forms of relationship between beings.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Connection | `connection` | Mutual relationship | Friendship, mentorship, collaboration |
| Endorsement | (table) | Quality vouch | ✓ exists |
| Follow | (table) | Subscription | ✓ exists |
| Block | `block` | Boundary | Mute, block, restrict |
| Recommendation | `recommendation` | Directed endorsement | "You should work with X on this" |

**Ops needed:** connect, disconnect, block, unblock, recommend

**Cross-layer:** Bonds affect Social visibility (following). Trust from Bonds affects Work delegation. Recommendations are Knowledge + Bond.

---

## Layer 10: Belonging — Membership

**Current:** Space membership (join/leave), space settings.

**General domain:** All forms of group membership and community lifecycle.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Team | `team` | Functional group | ✓ from unified spec |
| Department | `department` | Organizational unit | ✓ from unified spec |
| Organization | `organization` | Legal/structural entity | ✓ from unified spec |
| Invitation | `invitation` | Membership offer | ✓ exists (invites table) |
| Membership | (join/leave) | Belonging state | ✓ exists |

**Ops needed:** promote, demote, transfer (membership lifecycle)

**Cross-layer:** Belonging determines access to Work. Teams are Social units. Organizations govern via Policy. Membership earned through Identity.

---

## Layer 11: Governance — Collective Decision-Making

**Current:** Proposals, voting (propose/vote ops), governance lens.

**General domain:** All mechanisms for making collective decisions and establishing rules.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Proposal | `proposal` | Decision to be made | ✓ exists |
| Policy | `policy` | Governing rule | ✓ from unified spec |
| Decision | `decision` | Recorded choice | ✓ from unified spec |
| Process | `process` | Repeatable sequence | ✓ from unified spec |
| Amendment | `amendment` | Change to existing rule | Constitutional change, policy update |

**Ops needed:** ratify, amend, repeal, delegate (authority)

**Cross-layer:** Policies govern Work. Decisions recorded as Knowledge. Proposals discussed in Social. Voting requires Identity.

---

## Layer 12: Culture — Shared Norms

**Current:** Pin/unpin ops.

**General domain:** All shared norms, values, traditions, and practices that shape collective behavior.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Norm | `norm` | Shared expectation | "We review all PRs", "Meetings start on time" |
| Tradition | `tradition` | Recurring practice | "Friday demos", "New member welcome", "Retrospectives" |
| Value | `value` | Guiding principle | "Transparency", "Speed", "Quality over quantity" |
| Recognition | `recognition` | Celebrating contribution | "MVP of the sprint", "Community champion" |

**Ops needed:** enshrine, celebrate, recognize

**Cross-layer:** Culture shapes how Work is done. Norms are informal Governance. Values expressed in Social. Recognition builds Identity.

---

## Layer 13: Being — Existential

**Current:** Reflect op.

**General domain:** All mechanisms for existential wellbeing, purpose, and growth.

| Entity | Kind | What | Scale examples |
|--------|------|------|---------------|
| Reflection | `reflection` | Self-examination | ✓ exists (as reflect op) |
| Intention | `intention` | Stated purpose | Personal mission, "why I'm here" |
| Check-in | `checkin` | Wellbeing pulse | "How are you?", mood tracking, energy level |
| Growth | `growth` | Development path | Learning plan, skill progression, career path |

**Ops needed:** checkin, mentor, support

**Cross-layer:** Reflection improves Work (retrospectives). Wellbeing affects Identity (burnout prevention). Growth is Knowledge about self. Check-ins are Social + Being.

---

## Complete Entity Kind Count

| Layer | Existing kinds | New kinds from generalization | Total |
|-------|---------------|------------------------------|-------|
| Work | task, project, goal | milestone, cycle | 5 |
| Market | — | resource, budget, contract, invoice, listing | 5 |
| Social | post, thread, conversation, comment | channel, event | 6 |
| Justice | — | case, ruling, appeal, precedent | 4 |
| Build | — | release, artifact, review, incident | 4 |
| Knowledge | claim | document, question, definition, lesson | 5 |
| Alignment | — | analytics, audit, disclosure | 3 |
| Identity | — | role, credential, reputation, badge | 4 |
| Bond | — | connection, block, recommendation | 3 |
| Belonging | — | team, department, organization | 3 |
| Governance | proposal | policy, decision, process, amendment | 5 |
| Culture | — | norm, tradition, value, recognition | 4 |
| Being | — | intention, checkin, growth | 3 |
| **Total** | **10** | **~54** | **~64** |

**~54 new entity kinds.** Each one is a Node with a kind. Each one uses the same grammar ops. Each one needs: a create form, a detail view, and a lens/list.

---

## Cross-Layer Relationship Map

```
Being ──grounds──→ Identity ──earns──→ Bond ──forms──→ Belonging
  │                    │                  │                │
  │                    └──reputation──→ Market            │
  │                                       │               │
  └──reflection──→ Knowledge              │               │
                      │                   │               │
                      └──informs──→ Work ←──structures──┘
                                     │
                    ┌────────────────┤
                    │                │
              Governance ←──norms── Culture
                    │
              Justice ←──transparency── Alignment
                    │
                    └──builds──→ Build
```

Every arrow is a concrete relationship: an entity in one layer references an entity in another via parent_id, node_deps, or tags.

---

## Build Priority

Not all 54 entity kinds are equally valuable. Priority based on:
1. **How many layers it connects** (cross-layer entities first)
2. **How many scales it serves** (universal entities first)
3. **How close to existing infrastructure** (cheap to add first)

### Tier 1 — High impact, low cost (already proven pattern)
- Team, Role, Organization (Belonging/Identity — Organize mode)
- Policy, Decision (Governance — Govern mode)
- Document (Knowledge — Learn mode)
- Channel (Social — Rooms mode)

### Tier 2 — High impact, medium cost
- Resource, Budget (Market — Allocate mode)
- Release, Incident (Build — deeper Build mode)
- Case, Ruling (Justice — deeper Justice mode)

### Tier 3 — Medium impact, adds depth
- Credential, Badge (Identity — reputation system)
- Norm, Recognition (Culture — community health)
- Question, Lesson (Knowledge — learning system)
- Check-in, Growth (Being — wellbeing)

### Tier 4 — Future (needs more design)
- Contract, Invoice, Listing (Market — full exchange)
- Precedent, Appeal, Amendment (Justice/Governance — legal)
- Analytics, Audit, Disclosure (Alignment — enterprise)
- Sprint/Cycle, Milestone (Work — planning depth)
- Event, Connection, Recommendation, Tradition, Intention (various)

---

## The Principle

Every entity kind is a Node. Every operation is an Op. The grammar is the API. Adding a new entity kind costs: 1 constant, 1 handler, 1 template. Adding a new op costs: 1 handler case. The architecture doesn't change. The product grows by accumulating entity kinds and cross-layer relationships on the same graph.

This is what "substrate for collective existence" means in practice. 64 entity kinds across 13 layers, all on one graph, all using one grammar, all composable with each other. A friend group uses 5 kinds. A company uses 30. A civilization uses all 64. The same product, different configurations.
