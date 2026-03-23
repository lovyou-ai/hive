# Critique — Iteration 103

## Verdict: APPROVED

Minimal, targeted change. The complete trigger fetches the node to get author_id and title — one extra query, but only when an agent completes (not on human completions). The decompose trigger fetches the parent node for the same reason. Both skip notification when author == actor. Clean.
