package runner

import (
	"fmt"
	"strings"
)

// parseReflectorOutput extracts COVER/BLIND/ZOOM/FORMALIZE sections from
// reflector LLM output. Sections are delimited by "**KEY:**" or "KEY:" markers.
// Returns a map of section name → trimmed content.
func parseReflectorOutput(content string) map[string]string {
	keys := []string{"COVER", "BLIND", "ZOOM", "FORMALIZE"}
	result := map[string]string{}

	for i, key := range keys {
		// Try bold markdown first: **KEY:**
		marker := "**" + key + ":**"
		idx := strings.Index(content, marker)
		markerLen := len(marker)

		if idx < 0 {
			// Fallback: plain KEY:
			marker = key + ":"
			idx = strings.Index(content, marker)
			markerLen = len(marker)
		}
		if idx < 0 {
			continue
		}

		start := idx + markerLen

		// Find where this section ends (start of next section).
		end := len(content)
		for _, nextKey := range keys[i+1:] {
			for _, nextMarker := range []string{"**" + nextKey + ":**", nextKey + ":"} {
				if nextIdx := strings.Index(content[start:], nextMarker); nextIdx >= 0 {
					if abs := start + nextIdx; abs < end {
						end = abs
					}
				}
			}
		}

		result[key] = strings.TrimSpace(content[start:end])
	}

	return result
}

// buildReflectorPrompt assembles the prompt sent to the Reflector agent.
// Artifacts: scout report, build report, critique, recent reflections, shared context.
func buildReflectorPrompt(scout, build, critique, recentReflections, sharedCtx string) string {
	return fmt.Sprintf(`You are the Reflector. You close each iteration by extracting what was learned.

## Institutional Knowledge
%s

## Scout Report (loop/scout.md)
%s

## Build Report (loop/build.md)
%s

## Critique (loop/critique.md)
%s

## Recent Reflections (loop/reflections.md)
%s

## Instructions

Produce a reflection entry with exactly these four sections:

**COVER:** What was accomplished? How does it connect to prior work?

**BLIND:** What was missed? What is invisible to the current process?

**ZOOM:** Step back. What is the larger pattern across iterations?

**FORMALIZE:** If a new lesson emerged, state it as a numbered principle. Otherwise write "No new lesson."

Keep it concise — 10-15 lines total. BLIND is the most important: actively look for absences.`, sharedCtx, scout, build, critique, recentReflections)
}

// formatReflectionEntry formats a dated append block for loop/reflections.md.
// date should be ISO 8601 (e.g. "2026-03-26").
func formatReflectionEntry(date, cover, blind, zoom, formalize string) string {
	return fmt.Sprintf("## %s\n\n**COVER:** %s\n\n**BLIND:** %s\n\n**ZOOM:** %s\n\n**FORMALIZE:** %s\n",
		date, cover, blind, zoom, formalize)
}
