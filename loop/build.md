# Fix: update stale /hive tests after HivePage redesign

## What changed

**`site/graph/hive_test.go`**
- Removed `TestGetHive_ContainsCivilizationBuilds` — asserted old tagline "The Civilization Builds" that HivePage no longer renders; replaced with `TestGetHive_ContainsHiveFeed` which checks for the `hive-feed` element that HivePage always renders
- Removed `TestGetHive_RendersMetrics` — asserted old stat card labels ("Features shipped", "Total autonomous spend", "Avg cost") that HivePage no longer renders
- Removed `TestGetHive_RendersCurrentlyBuilding` — asserted "Currently building" section with "Idle" text and task title; this section no longer exists in HivePage
- Removed unused `"fmt"` import (was only used by the dropped metrics test)

## Why

The new `HivePage` renders: iteration counter, phase pill, last build title, phase timeline (`hive-feed` HTMX element), and recent commits. It no longer fetches posts for stat cards or tasks for a "currently building" section. `TestHiveDashboard` already covers the core HivePage content assertions (iteration number, phase, build title).

## Verification

- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test -buildvcs=false ./...` — all pass (graph tests skip without DB, as expected)
