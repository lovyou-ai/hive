package runner

// HiveWiring declares how roles communicate via the graph.
// This is the hive's nervous system — what each role reads, writes,
// and which tools it has access to. Declarative, inspectable, visualizable.
//
// To visualize: render Roles as nodes, Channels as directed edges.
// To validate: check that every Write has a matching Read somewhere.
// To extend: add a RoleWiring entry — the pipeline wires it automatically.

// HiveWiring is the complete communication topology of the hive.
type HiveWiring struct {
	Roles []RoleWiring `json:"roles"`
}

// RoleWiring describes one role's place in the hive.
type RoleWiring struct {
	Name   string    `json:"name"`
	Model  string    `json:"model"`
	Mode   string    `json:"mode"`   // "operate" or "reason"
	Reads  []Channel `json:"reads"`  // subscriptions — what it consumes
	Writes []Channel `json:"writes"` // publications — what it produces
}

// Channel is a typed communication path on the graph.
type Channel struct {
	Kind        string `json:"kind"`         // entity kind: "task", "document", "claim", "post"
	Op          string `json:"op"`           // grammar op: "intend", "assert", "express", "discuss"
	TitlePrefix string `json:"title_prefix"` // convention: "Scout Report:", "Critique:", "Build:", "Lesson:"
	Description string `json:"description"`  // human-readable purpose
}

// DefaultWiring returns the current hive topology.
// This is the source of truth for how roles communicate.
func DefaultWiring() HiveWiring {
	return HiveWiring{
		Roles: []RoleWiring{
			{
				Name:  "pm",
				Model: "sonnet",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "task", Op: "board", Description: "reads open tasks to assess work state"},
					{Kind: "document", TitlePrefix: "backlog", Description: "reads backlog for priorities"},
					{Kind: "claim", TitlePrefix: "Lesson:", Description: "reads lessons from prior cycles"},
				},
				Writes: []Channel{
					{Kind: "task", Op: "intend", TitlePrefix: "milestone", Description: "creates milestones for Architect to decompose"},
				},
			},
			{
				Name:  "scout",
				Model: "sonnet",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "task", Op: "board", Description: "reads board for active goals"},
					{Kind: "document", Description: "searches knowledge for prior work"},
				},
				Writes: []Channel{
					{Kind: "document", Op: "intend", TitlePrefix: "Scout Report:", Description: "gap report for Architect"},
					{Kind: "post", Op: "express", TitlePrefix: "Scout Report:", Description: "gap report on feed for visibility"},
				},
			},
			{
				Name:  "architect",
				Model: "sonnet",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "task", TitlePrefix: "milestone", Description: "reads PM milestones to decompose"},
					{Kind: "document", TitlePrefix: "Scout Report:", Description: "reads scout gap reports"},
				},
				Writes: []Channel{
					{Kind: "task", Op: "intend", Description: "creates subtasks for Builder"},
				},
			},
			{
				Name:  "builder",
				Model: "opus",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "task", Op: "board", Description: "claims and works assigned tasks"},
					{Kind: "document", Description: "searches knowledge for context"},
				},
				Writes: []Channel{
					{Kind: "document", Op: "intend", TitlePrefix: "Build:", Description: "build report documenting what was done"},
				},
			},
			{
				Name:  "tester",
				Model: "sonnet",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "document", TitlePrefix: "Build:", Description: "reads build report to know what to test"},
				},
				Writes: []Channel{
					{Kind: "claim", Op: "assert", TitlePrefix: "Test:", Description: "test results as verifiable claims"},
				},
			},
			{
				Name:  "critic",
				Model: "sonnet",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "document", TitlePrefix: "Build:", Description: "reads build report + git diff"},
					{Kind: "document", Description: "searches knowledge for invariants and conventions"},
				},
				Writes: []Channel{
					{Kind: "claim", Op: "assert", TitlePrefix: "Critique:", Description: "verdict (PASS/REVISE) as verifiable claim"},
					{Kind: "task", Op: "intend", TitlePrefix: "Fix:", Description: "fix tasks on REVISE"},
				},
			},
			{
				Name:  "reflector",
				Model: "sonnet",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "document", TitlePrefix: "Scout Report:", Description: "what gap was identified"},
					{Kind: "document", TitlePrefix: "Build:", Description: "what was built"},
					{Kind: "claim", TitlePrefix: "Critique:", Description: "what the critic said"},
					{Kind: "claim", TitlePrefix: "Lesson:", Description: "prior lessons to avoid repetition"},
				},
				Writes: []Channel{
					{Kind: "document", Op: "intend", TitlePrefix: "Reflection:", Description: "COVER/BLIND/ZOOM/FORMALIZE analysis"},
					{Kind: "claim", Op: "assert", TitlePrefix: "Lesson:", Description: "formalized lessons as verifiable claims"},
				},
			},
			{
				Name:  "observer",
				Model: "sonnet",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "task", Op: "board", Description: "audits board for stale/orphaned tasks"},
					{Kind: "claim", Description: "checks unlinked claims and schema violations"},
					{Kind: "document", Description: "checks knowledge layer usage"},
				},
				Writes: []Channel{
					{Kind: "task", Op: "intend", Description: "creates fix tasks for integrity issues"},
				},
			},
			{
				Name:  "council",
				Model: "sonnet",
				Mode:  "operate",
				Reads: []Channel{
					{Kind: "document", Description: "each agent searches knowledge from their perspective"},
					{Kind: "task", Op: "board", Description: "checks current work state"},
				},
				Writes: []Channel{
					{Kind: "post", Op: "express", TitlePrefix: "Council report", Description: "deliberation summary for feed"},
				},
			},
		},
	}
}
