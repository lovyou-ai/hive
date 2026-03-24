# Critique — Iteration 202

## Unified Ontology: PASS

**Derivation chain:**
- Gap: Social spec derived from platform replacement, not from first principles. Work spec too narrow. No structural document relating them.
- Plan: Apply generator function to the root concept. Derive unified ontology.
- Output: Structural spec with derivation order, merged mode set, unified entity types.

**Correctness:**
- Derivation order (Being → Identity → Belonging → Communication → Work → ...) is logically sound. Each level requires the previous. ✓
- 10 modes are MECE for organized activity: 4 communication + 6 activity. ✓
- Grammar coverage matrix: most ops appear in 4+ modes. No orphan ops. ✓
- "Modes emerge from content" principle is architecturally sound — node kinds determine which modes surface. ✓

**The key claim: "The architecture IS the unified ontology":**
- True for data model — Nodes table supports arbitrary kinds. ✓
- True for operations — grammar ops are kind-agnostic. ✓
- True for views — adding a new lens is template + handler work. ✓
- Partially true for sidebar — currently hardcoded, would need dynamic mode detection.

**Risks:**
- 18 entity types × 10 modes = 180 cells. Not all are meaningful (many are empty). But the non-empty cells represent a large product surface.
- "Modes emerge from content" needs careful UX — showing 10 modes to a new user is overwhelming. Progressive disclosure is critical.
- The relationship between Spaces and Organizations needs resolution. Currently Spaces are the top-level container. Should Organizations contain Spaces?

**What this spec does NOT do:**
- Doesn't invalidate existing specs — places them in context.
- Doesn't prescribe build order — that's a separate decision.
- Doesn't resolve the Space ↔ Organization relationship.

## Verdict: PASS
