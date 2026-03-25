package runner

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/lovyou-ai/hive/pkg/api"
)

// runArchitect decomposes a large task into smaller implementable subtasks.
// Called when the Builder times out or fails on a task that's too ambitious.
func (r *Runner) runArchitect(ctx context.Context) {
	if !r.cfg.OneShot && r.tick%4 != 0 {
		return
	}

	// Find the task that needs decomposition — the most recent assigned task
	// that the Builder couldn't complete.
	tasks, err := r.cfg.APIClient.GetTasks(r.cfg.SpaceSlug, "")
	if err != nil {
		log.Printf("[architect] GetTasks error: %v", err)
		return
	}

	var target *api.Node
	for _, t := range tasks {
		if t.Kind != "task" || t.State == "done" || t.State == "closed" {
			continue
		}
		if r.cfg.AgentID != "" && t.AssigneeID == r.cfg.AgentID {
			if !strings.HasPrefix(t.Title, "Fix:") {
				target = &t
				break
			}
		}
	}

	if target == nil {
		log.Printf("[architect] no task to decompose")
		if r.cfg.OneShot {
			r.done = true
		}
		return
	}

	log.Printf("[architect] decomposing: %s", target.Title)

	sharedCtx := LoadSharedContext(r.cfg.HiveDir)
	repoCtx := r.readRepoContext()

	prompt := buildArchitectPrompt(sharedCtx, repoCtx, *target)

	resp, err := r.cfg.Provider.Reason(ctx, prompt, nil)
	if err != nil {
		log.Printf("[architect] Reason error: %v", err)
		return
	}

	r.cost.Record(resp.Usage())
	log.Printf("[architect] plan ready (cost=$%.4f)", resp.Usage().CostUSD)

	// Parse subtasks from the plan.
	subtasks := parseArchitectSubtasks(resp.Content())
	if len(subtasks) == 0 {
		log.Printf("[architect] no subtasks found in plan")
		if r.cfg.OneShot {
			r.done = true
		}
		return
	}

	// Close the original task (too large) and create smaller subtasks.
	_ = r.cfg.APIClient.CommentTask(r.cfg.SpaceSlug, target.ID,
		fmt.Sprintf("Decomposed into %d subtasks by Architect.", len(subtasks)))
	_ = r.cfg.APIClient.CompleteTask(r.cfg.SpaceSlug, target.ID)

	for _, st := range subtasks {
		task, err := r.cfg.APIClient.CreateTask(r.cfg.SpaceSlug, st.title, st.desc, st.priority)
		if err != nil {
			log.Printf("[architect] create subtask error: %v", err)
			continue
		}
		if r.cfg.AgentID != "" {
			_ = r.cfg.APIClient.ClaimTask(r.cfg.SpaceSlug, task.ID)
		}
		log.Printf("[architect] created subtask: %s", st.title)
	}

	log.Printf("[architect] decomposed into %d subtasks", len(subtasks))

	if r.cfg.OneShot {
		r.done = true
	}
}

type architectSubtask struct {
	title    string
	desc     string
	priority string
}

func buildArchitectPrompt(sharedCtx, repoCtx string, task api.Node) string {
	return fmt.Sprintf(`You are the Architect. A task was too large for the Builder to implement in one session (it timed out or failed). Your job: decompose it into 2-4 smaller subtasks that can each be implemented in under 10 minutes.

## Institutional Knowledge
%s

## Target Repo Context
%s

## Task to Decompose
Title: %s
Description: %s

## Instructions

1. Read the task and understand what it requires.
2. Break it into 2-4 smaller tasks. Each should:
   - Change 1-3 files maximum
   - Be implementable in one Operate() call (under 10 minutes)
   - Be self-contained (doesn't depend on the other subtasks being done first, if possible)
3. Be specific: name the files, the functions, the changes.

## Output Format

For each subtask:

SUBTASK_TITLE: <one-line title>
SUBTASK_PRIORITY: <high|medium>
SUBTASK_DESCRIPTION: <2-3 sentences, specific files and changes>`, sharedCtx, repoCtx, task.Title, task.Body)
}

func parseArchitectSubtasks(content string) []architectSubtask {
	var tasks []architectSubtask
	var current architectSubtask

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "SUBTASK_TITLE:") {
			if current.title != "" {
				tasks = append(tasks, current)
			}
			current = architectSubtask{title: strings.TrimSpace(strings.TrimPrefix(line, "SUBTASK_TITLE:"))}
		} else if strings.HasPrefix(line, "SUBTASK_PRIORITY:") {
			current.priority = strings.TrimSpace(strings.TrimPrefix(line, "SUBTASK_PRIORITY:"))
		} else if strings.HasPrefix(line, "SUBTASK_DESCRIPTION:") {
			current.desc = strings.TrimSpace(strings.TrimPrefix(line, "SUBTASK_DESCRIPTION:"))
		}
	}
	if current.title != "" {
		tasks = append(tasks, current)
	}

	for i := range tasks {
		switch tasks[i].priority {
		case "high", "medium", "low", "urgent":
		default:
			tasks[i].priority = "high"
		}
	}
	return tasks
}
