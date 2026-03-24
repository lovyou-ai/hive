# Hive Limitations — What Agents Can and Cannot Catch

Honest accounting of what still needs human eyes. Updated as capabilities change.

## What the hive catches well

| Category | How | Example |
|----------|-----|---------|
| Missing code patterns | Critic greps for consistent treatment | Allowlist miss (iter 230) |
| Build breaks | `go build ./...` before closing task | Every builder commit |
| State machine violations | Critic checks preconditions in diffs | Progress handler missing guard (iter 230) |
| Entity pipeline completeness | Observer greps for kind constants, handlers, routes | All 13 kinds verified |
| Test failures | `go test ./...` before closing task | Every builder commit |
| API correctness | Builder verifies JSON responses, auth, status codes | API client has error handling |
| Stale task detection | Scout checks board, throttles at max 3 | Scout throttle in iter 227 |

## What the hive catches sometimes

| Category | How | Limitation |
|----------|-----|-----------|
| Omission errors in distant code | Critic reviews diffs, not full files | Only catches if the omission is within ~100 lines of the change. The intend allowlist (400 lines away) was missed by the Critic. |
| Product gap identification | Scout reads state.md + CLAUDE.md | Scout creates tasks from what it knows. It can't discover unknown-unknowns. |
| Spec compliance | Observer reads specs and code | Only checks what's explicitly specified. Implicit requirements are invisible. |
| Cross-entity consistency | Observer greps patterns | Can verify structural consistency (handlers, routes) but not semantic consistency (do these entities relate correctly?) |

## What the hive CANNOT catch

### Visual and aesthetic
- **Rendered appearance.** Agents read HTML/CSS, not pixels. They cannot see if the layout looks right, if spacing is off, if colors clash, if typography feels wrong.
- **"Ember Minimalism" compliance.** The visual identity is defined by feel, not code. Agents can check for CSS class names but cannot judge whether the result looks warm, intentional, alive.
- **Animation quality.** Whether a transition is too fast, too slow, or feels janky — invisible to code analysis.
- **Mobile rendering.** Responsive CSS classes can be verified but the actual mobile experience cannot be seen.

### Intuition and judgment
- **"That isn't our vibe."** Product identity is a human judgment. When Matt says the sidebar should show layers not lenses, that's a framing insight that emerges from deeply understanding the product's soul. Agents don't have this.
- **"Work isn't just a kanban board."** Structural reframing — seeing that a category is wrong, not just incomplete — requires understanding the product at a level agents can approximate but not match.
- **User confusion.** Whether a first-time user would be confused by the navigation, overwhelmed by options, or unable to find what they need. Agents can trace paths but can't feel friction.
- **Emotional resonance.** Whether the product feels caring, alive, trustworthy. The soul says "take care of your human" — whether the product EMBODIES that is a human judgment.

### Strategic and contextual
- **Market positioning.** Is this feature actually differentiated from Linear/Discord/Twitter, or just different? Agents can compare feature lists but not evaluate competitive advantage.
- **Timing.** Whether this is the right thing to build NOW, given what's happening in the market, what users are asking for, what the revenue model requires.
- **Coherence.** Whether the 13 layers feel like one product or 13 bolted-together features. This is an architectural judgment that requires seeing the whole.

### Unknown unknowns
- **What's not in the spec.** If nobody specified it, nobody checks for it. The most important features are sometimes the ones nobody thought to write down.
- **What's not in the code.** Absence is invisible. The Scout can only create tasks for gaps it can identify from existing docs and code. A gap that exists in the product but not in any document is undetectable.

## What this means for autonomy

The hive can operate autonomously for:
- **Entity pipeline additions** — mechanical, well-defined pattern
- **Handler/template consistency fixes** — greppable, verifiable
- **Bug fixes identified by Critic** — specific, code-level, actionable
- **Feature additions from detailed specs** — when the spec says exactly what to build

The hive still needs Matt for:
- **Product direction** — what to build, not how
- **Visual review** — does this look right?
- **UX judgment** — is this confusing?
- **Architectural decisions** — is this the right abstraction?
- **"Stop and re-derive"** — when the framing is wrong, not just the code
- **Soul alignment** — does this serve the mission?

## Improving over time

Each limitation is a frontier, not a wall:
- **Visual:** Multimodal models can look at screenshots. Not wired up yet.
- **UX:** User testing data (analytics, session recordings) could inform agents. Not available yet.
- **Strategy:** Competitive analysis tools exist. Not integrated yet.
- **Unknown unknowns:** More specs, more tests, more lessons reduce the space of unknowns.

The goal is not to replace Matt's judgment. It's to handle everything that ISN'T judgment, so his judgment can focus on what matters.
