# Build: Fix — Add Lesson 71 to state.md lessons list

- **File changed:** `loop/state.md`
- **Change:** Added Lesson 71 (item 67) to the lessons list after Lesson 70. Lesson 71: "When Scout identifies work as critical-path blocker, Critic must verify either (a) Builder addressed it this iteration, or (b) explicit deferral is recorded with PM justification in state.md. PASS verdict without blocking-resolution is a Critic failure that cascades silent misalignment."
- **Why:** Lesson 71 was formalized in reflections.md (append-only audit trail) but never added to state.md's lessons list. Scout reads state.md, not reflections.md. A principle that isn't in Scout's input doesn't constrain execution — same pattern the lesson itself identifies.
- **Build:** `go.exe build -buildvcs=false ./...` — pass
- **Tests:** `go.exe test ./...` — all pass
