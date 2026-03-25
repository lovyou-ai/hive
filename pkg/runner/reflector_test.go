package runner

import (
	"strings"
	"testing"
)

func TestParseReflectorOutput(t *testing.T) {
	t.Run("bold markdown sections", func(t *testing.T) {
		input := `Here is my reflection.

**COVER:** We shipped the entity pipeline for Goal nodes. Connects to the prior work on Project.

**BLIND:** No tests were written for the new pipeline handler.

**ZOOM:** Three consecutive entity-kind iterations — pattern is converging on fixpoint.

**FORMALIZE:** Lesson 56: Entity pipelines share a single integration test template.`

		got := parseReflectorOutput(input)

		if !strings.Contains(got["COVER"], "entity pipeline for Goal") {
			t.Errorf("COVER = %q, want 'entity pipeline for Goal'", got["COVER"])
		}
		if !strings.Contains(got["BLIND"], "No tests") {
			t.Errorf("BLIND = %q, want 'No tests'", got["BLIND"])
		}
		if !strings.Contains(got["ZOOM"], "converging on fixpoint") {
			t.Errorf("ZOOM = %q, want 'converging on fixpoint'", got["ZOOM"])
		}
		if !strings.Contains(got["FORMALIZE"], "Lesson 56") {
			t.Errorf("FORMALIZE = %q, want 'Lesson 56'", got["FORMALIZE"])
		}
	})

	t.Run("plain KEY: sections", func(t *testing.T) {
		input := "COVER: Shipped the auth fix.\nBLIND: No rollback plan.\nZOOM: Auth hardening theme.\nFORMALIZE: No new lesson."

		got := parseReflectorOutput(input)

		if got["COVER"] != "Shipped the auth fix." {
			t.Errorf("COVER = %q", got["COVER"])
		}
		if got["BLIND"] != "No rollback plan." {
			t.Errorf("BLIND = %q", got["BLIND"])
		}
		if got["ZOOM"] != "Auth hardening theme." {
			t.Errorf("ZOOM = %q", got["ZOOM"])
		}
		if got["FORMALIZE"] != "No new lesson." {
			t.Errorf("FORMALIZE = %q", got["FORMALIZE"])
		}
	})

	t.Run("missing sections return empty string", func(t *testing.T) {
		input := "**COVER:** Only this section present."

		got := parseReflectorOutput(input)

		if got["COVER"] == "" {
			t.Error("COVER should be non-empty")
		}
		if got["BLIND"] != "" {
			t.Errorf("BLIND should be empty, got %q", got["BLIND"])
		}
		if got["ZOOM"] != "" {
			t.Errorf("ZOOM should be empty, got %q", got["ZOOM"])
		}
		if got["FORMALIZE"] != "" {
			t.Errorf("FORMALIZE should be empty, got %q", got["FORMALIZE"])
		}
	})

	t.Run("empty input returns empty map", func(t *testing.T) {
		got := parseReflectorOutput("")
		if len(got) != 0 {
			t.Errorf("expected empty map, got %v", got)
		}
	})

	t.Run("section content is trimmed", func(t *testing.T) {
		input := "**COVER:**   padded content   \n**BLIND:** next"
		got := parseReflectorOutput(input)
		if got["COVER"] != "padded content" {
			t.Errorf("COVER not trimmed: %q", got["COVER"])
		}
	})
}

func TestBuildReflectorPrompt(t *testing.T) {
	prompt := buildReflectorPrompt(
		"## Scout\nGap: missing Goal entity",
		"## Build\nAdded KindGoal to store.go",
		"## Critique\nVERDICT: PASS",
		"## 2026-03-25\n**COVER:** shipped Project",
		"## Invariants\n1. IDENTITY",
	)

	// All artifact content must appear in the prompt.
	if !contains(prompt, "Gap: missing Goal entity") {
		t.Error("prompt missing scout content")
	}
	if !contains(prompt, "Added KindGoal to store.go") {
		t.Error("prompt missing build content")
	}
	if !contains(prompt, "VERDICT: PASS") {
		t.Error("prompt missing critique content")
	}
	if !contains(prompt, "shipped Project") {
		t.Error("prompt missing recent reflections")
	}
	if !contains(prompt, "## Invariants") {
		t.Error("prompt missing shared context")
	}

	// Must contain the four section headings the Reflector is expected to produce.
	for _, section := range []string{"COVER", "BLIND", "ZOOM", "FORMALIZE"} {
		if !contains(prompt, section) {
			t.Errorf("prompt missing section heading: %s", section)
		}
	}

	// Must instruct on conciseness and the BLIND priority.
	if !contains(prompt, "BLIND is the most important") {
		t.Error("prompt should highlight BLIND as most important")
	}
}

func TestFormatReflectionEntry(t *testing.T) {
	entry := formatReflectionEntry(
		"2026-03-26",
		"Shipped Goal entity kind.",
		"Integration tests not written.",
		"Entity pipeline iterations converging.",
		"Lesson 56: test the pipeline once per kind.",
	)

	// Must open with a date heading.
	if !strings.HasPrefix(entry, "## 2026-03-26") {
		t.Errorf("entry should start with '## 2026-03-26', got: %q", entry[:min(30, len(entry))])
	}

	// Must contain all four labeled sections.
	for _, label := range []string{"**COVER:**", "**BLIND:**", "**ZOOM:**", "**FORMALIZE:**"} {
		if !contains(entry, label) {
			t.Errorf("entry missing label %s", label)
		}
	}

	// Must contain the supplied content.
	if !contains(entry, "Shipped Goal entity kind.") {
		t.Error("entry missing COVER content")
	}
	if !contains(entry, "Integration tests not written.") {
		t.Error("entry missing BLIND content")
	}
	if !contains(entry, "converging.") {
		t.Error("entry missing ZOOM content")
	}
	if !contains(entry, "Lesson 56") {
		t.Error("entry missing FORMALIZE content")
	}

	// Must end with a trailing newline (append-safe).
	if !strings.HasSuffix(entry, "\n") {
		t.Error("entry must end with newline for safe appending")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
