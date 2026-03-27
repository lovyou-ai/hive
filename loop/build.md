# Build: Fix loopDirtyCheck to use HiveDir instead of RepoPath

## Gap Fixed
`loopDirtyCheck` was checking `cfg.RepoPath` instead of `cfg.HiveDir`. Loop artifacts (`loop/*.md`) live in the hive repo (`HiveDir`), not the product repo (`RepoPath`). When `--repo ../site` is passed, `RepoPath` points to the site repo which has no `loop/` directory — the gate silently passed through dirty hive artifacts. The test also had the same confusion: it put the git repo in `RepoPath`, making it accidentally match the broken implementation.

## Files Changed

### `pkg/runner/pipeline_tree.go`
- `loopDirtyCheck`: changed guard `cfg.RepoPath == ""` → `cfg.HiveDir == ""`
- `loopDirtyCheck`: changed `cmd.Dir = cfg.RepoPath` → `cmd.Dir = cfg.HiveDir`
- Updated doc comment accordingly

### `pkg/runner/pipeline_tree_test.go` (`TestLoopDirtyCheckBlocksReflector`)
- Changed `Config{HiveDir: hiveDir, RepoPath: repoDir}` → `Config{HiveDir: repoDir, RepoPath: ""}` — the git repo with dirty `loop/build.md` must be in `HiveDir`
- Changed `countDiagnostics(hiveDir)` → `countDiagnostics(repoDir)` — diagnostics are written to `HiveDir`
- Removed now-unused `hiveDir := makeHiveDir(...)` variable

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass
