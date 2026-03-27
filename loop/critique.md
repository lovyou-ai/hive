# Critique: [hive:builder] Add dirty-loop-artifacts gate before Reflector in Execute()

**Verdict:** REVISE

**Summary:** The bug is confirmed. Here is the full analysis:

---

## Analysis

### What the commit does

Adds a `loop-clean-check` phase between `critic` and `reflector` in `NewPipelineTree`. When loop artifacts are uncommitted, `loopDirtyCheck` returns an error, `Execute` emits a diagnostic, and the reflector is blocked. The right idea, correctly placed.

### Structural correctness of the gate

The `NewPipelineTree` refactor (create `pt` first, then assign `pt.phases` so the closure can capture `pt`) is sound. `Execute`'s diagnostic emission path is already correct — any error from a phase that didn't write its own diagnostic gets one written automatically.

### The bug: wrong directory

`loopDirtyCheck` runs `git status --porcelain -- loop/` in `pt.cfg.RepoPath`. But `loop/` artifacts live in the **hive repo**, which is `pt.cfg.HiveDir`. These are different fields:

- `HiveDir` = the hive repo (`/c/src/matt/lovyou3/hive`) — confirmed at `runner.go:428`: `filepath.Join(r.cfg.HiveDir, "loop", "build.md")`
- `RepoPath` = the product repo being operated on (`--repo ../site` for site work, defaults to `.` from cwd)

When `--repo ../site` is passed (the primary pipeline mode, shown in the cmd header comments), `RepoPath` points to the site repo, which has no `loop/` directory. The check returns clean, silently passing through dirty hive artifacts. The gate does nothing in the common pipeline invocation.

The test is also affected: `Config{HiveDir: hiveDir, RepoPath: repoDir}` — it puts the git repo and dirty `loop/build.md` in `repoDir` which is assigned to `RepoPath`, not `HiveDir`. The test passes because the implementation accidentally matches the test's structure, not the production use case.

**Fix required:**
```go
// loopDirtyCheck: change RepoPath → HiveDir
func (pt *PipelineTree) loopDirtyCheck(_ context.Context) error {
    if pt.cfg.HiveDir == "" {  // was: RepoPath
        return nil
    }
    cmd := exec.Command("git", "status", "--porcelain", "--", "loop/")
    cmd.Dir = pt.cfg.HiveDir   // was: RepoPath
```

**Test fix required:**
```go
pt := &PipelineTree{
    cfg: Config{HiveDir: repoDir, RepoPath: ""},  // was: HiveDir: hiveDir, RepoPath: repoDir
}
```
(And the test's `makeHiveDir` call for diagnostics still works independently.)

---

VERDICT: REVISE

**Required fix:** In `loopDirtyCheck`, replace `pt.cfg.RepoPath` with `pt.cfg.HiveDir` (both the nil-guard and `cmd.Dir`). Update the test to set `HiveDir: repoDir` instead of `RepoPath: repoDir`. The loop artifacts are in the hive repo, not the product repo — using the wrong path silently bypasses the gate in all pipeline-mode invocations with `--repo <other-repo>`.
