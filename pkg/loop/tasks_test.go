package loop

import (
	"testing"
)

func TestIsMetaTaskBody(t *testing.T) {
	tests := []struct {
		title       string
		description string
		want        bool
	}{
		// Positive — each meta-task pattern should be detected.
		{"op=complete task-123", "", true},
		{"Close task build feature", "", true},
		{"Mark done: implement auth", "", true},
		{"Close the following tasks", "", true},

		// Pattern in description rather than title.
		{"Follow-up work", "op=complete the previous task", true},
		{"Wrap up", "close task xyz after review", true},
		{"Update status", "mark done all items in backlog", true},
		{"Batch", "close the following: task-1, task-2", true},

		// Case-insensitive matching.
		{"OP=COMPLETE task", "", true},
		{"CLOSE TASK now", "", true},
		{"MARK DONE", "", true},
		{"CLOSE THE FOLLOWING", "", true},

		// Negative — genuine new task descriptions.
		{"Build authentication module", "implement OAuth2 flow", false},
		{"Fix bug in task parser", "the parser drops empty lines", false},
		{"Add tests for loop", "cover edge cases in signal parsing", false},
		{"Deploy to production", "run fly deploy after build", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got := isMetaTaskBody(tt.title, tt.description)
		if got != tt.want {
			t.Errorf("isMetaTaskBody(%q, %q) = %v, want %v", tt.title, tt.description, got, tt.want)
		}
	}
}

func TestParseTaskCommandsMetaTaskNotFiltered(t *testing.T) {
	// parseTaskCommands itself does NOT filter meta-tasks — that happens in
	// execTaskCreate. Verify the command is still parsed so the guard can fire.
	response := `/task create {"title": "close task xyz", "description": "mark done"}`
	commands := parseTaskCommands(response)
	if len(commands) != 1 {
		t.Fatalf("got %d commands, want 1", len(commands))
	}
	if commands[0].Action != "create" {
		t.Errorf("action = %q, want %q", commands[0].Action, "create")
	}
}
