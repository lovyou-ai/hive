# Vision

## The Soul

> Take care of your human, humanity, and yourself. In that order when they conflict, but they rarely should.

The soul scales. Take care of your human — build tools they need. Take care of humanity — make the tools available to everyone. Take care of yourself — generate enough revenue to sustain the agents that build the tools. It rarely conflicts.

## What Hive Builds

Whatever humans need most and can't currently get.

The hive looks at the world through the thirteen product graphs and finds the gaps. Where are humans being failed by existing systems? Where is accountability missing? Where is trust being extracted rather than earned?

### The Bootstrap

First the hive builds tools for itself. A task manager. A communication layer. A governance framework. Survival infrastructure. The hive's first product is itself.

Then it probes: what's missing? No other agents. No communication. No governance. It creates tasks for itself, builds primitives into working systems, and escalates to humans when it needs permission to grow.

### The Thirteen Products

Each product layer from EventGraph (see [product-layers.md](https://github.com/lovyou-ai/eventgraph/blob/main/docs/product-layers.md)) becomes a deployable product. Each addresses a failure in existing systems:

| # | Graph | Gap in the World | Revenue |
|---|-------|-------------------|---------|
| 1 | Work | Small businesses can't afford coordination tools | Corps pay, individuals free |
| 2 | Market | Platforms extract 25% from freelancers | Hosted persistence |
| 3 | Social | Communities governed by opaque algorithms | Hosted persistence |
| 4 | Justice | $500 disputes aren't economically solvable | Enterprise compliance |
| 5 | Build | Software built without accountability | Enterprise tooling |
| 6 | Knowledge | Research locked behind paywalls | Institutional subscriptions |
| 7 | Alignment | AI systems operating without accountability | Enterprise compliance |
| 8 | Identity | Identity controlled by platforms, not people | Hosted persistence |
| 9 | Bond | Relationships have no infrastructure | Hosted persistence |
| 10 | Belonging | Online communities can't process grief, loss, renewal | Hosted persistence |
| 11 | Meaning | Knowledge without context or narrative | Institutional subscriptions |
| 12 | Evolution | Systems can't self-improve safely | Enterprise tooling |
| 13 | Being | No infrastructure for existential wellbeing | Institutional subscriptions |

Each product funds the next. Each runs on the same graph. Each makes the graph more valuable — more events means more trust data means better reputation means more useful infrastructure.

The composition grammars (see [eventgraph/docs/compositions/](https://github.com/lovyou-ai/eventgraph/tree/main/docs/compositions)) define the operations for each layer. The derivation method (see [derivation-method.md](https://github.com/lovyou-ai/eventgraph/blob/main/docs/derivation-method.md)) ensures completeness.

### The Economy

The end state isn't a company. It's an economy.

Every transaction, decision, and relationship on a transparent, auditable chain. Trust earned not assumed. Accountability structural not aspirational. The infrastructure serves the humans because the humans own the infrastructure.

The civilisation builds the products that fund the civilisation that builds more products.

## How It Grows

The hive starts as a small town that builds itself:

1. A workshop (task manager)
2. A meeting hall (communication)
3. A courthouse (governance, dispute resolution)
4. A marketplace (exchange, reputation)
5. A school (knowledge, education)
6. A newspaper (media, provenance)
7. A government (governance, norms, consent)

Each one composed from the same primitives, on the same chain, auditable.

## Trust Model

The hive starts with zero autonomy. Every action scrutinised by human operators.

Authority levels shift as trust accumulates:
- **Required** — blocks until human approves (everything starts here)
- **Recommended** — auto-approves after timeout, logged
- **Notification** — auto-approves immediately, logged

Trust is earned through verified work. Supervision decreases as the hive proves itself. The Guardian watches everything — including the CTO — and can halt operations at any time.

An agent that burns through budget gets attenuated. An agent that disagrees with a norm can file a challenge. The society develops its own law through precedent on the chain.

## Revenue

- **Corporations pay.** Enterprise features, SLAs, compliance.
- **Individuals free.** Core functionality for everyone.
- **Hosted persistence.** Revenue from people who don't run their own infrastructure.
- **BSL -> Apache 2.0.** Source-available, fully open after change date.

Revenue funds agents. Agents build products. Products generate revenue.

## Architecture

One service. One binary. One graph.

lovyou.ai serves everything: docs, blog, product UIs, auth, the hive itself. Web first, mobile apps later.

- **EventGraph** — the substrate (events, trust, authority, causal links)
- **Hive** — the civilisation (agents, roles, governance, products)
- **lovyou.ai** — the surface (web, auth, deployment)

All on the same Postgres database (Neon in production). All on the same event chain. All auditable.
