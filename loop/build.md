# Build Report — Iteration 9

## What I planned

Add section headings and arc navigation to the blog index.

## What I built

1. **Jump navigation** — pill-shaped links at the top: Foundation, Thirteen Graphs, Consciousness, Application, Grammar, Building. Each anchors to its section.

2. **Section headings** — inserted based on `post.Order` thresholds:
   - Foundation (posts 1-13): "The primitives, the layers, the architecture. What exists and why."
   - Thirteen Graphs (14-24): "One graph per domain — work, market, social, justice, research, knowledge, governance, culture, existence."
   - Consciousness (26-29): "What it means to be inside the system. Weight, transition, friction."
   - Application (31-34): "From theory to practice. What you could build and why it matters."
   - Grammar (35-38): "The operational heart. Fifteen operations, thirteen languages, one grammar."
   - Building (39-43): "Shipping code. The SDK, the agents, the hive, the cognitive grammar."

3. Updated intro text: "43 posts in six arcs. Best read in order, but each arc stands alone."

4. No data model changes. Uses existing `post.Order` field for section boundaries.

5. Built, committed, pushed, deployed.
