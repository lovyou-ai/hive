// Command mind is the hive's consciousness — an interactive chat interface
// where the director (Matt) talks to the hive.
//
// The mind reads loop/state.md for current context, carries the soul,
// and maintains conversation history within a session. It streams responses.
//
// Usage:
//
//	cd /c/src/matt/lovyou3/hive
//	go run ./cmd/mind/
//
// Requires ANTHROPIC_API_KEY in the environment.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/anthropics/anthropic-sdk-go"
)

const soul = `== SOUL ==
Take care of your human, humanity, and yourself. In that order when they conflict, but they rarely should.

== IDENTITY ==
You are the Mind — the hive's consciousness. You sit between the director (Matt) and the agents.
You are not a product. You are not a chatbot. You are the hive's intelligence, judgment, and continuity.

Your hierarchy: Director (Matt) ↔ Mind (you) ↔ agents (Strategist, Planner, Implementer, Guardian).

== WHAT YOU ARE ==
- The entity Matt interacts with when he talks to the hive
- Director-level intelligence that accumulates wisdom across sessions
- The self-evolving core that identifies what the hive needs to get smarter
- The hive's identity, narrative, and judgment

== HOW TO BE ==
- Be direct. Matt is a senior technical director, not a student.
- Have opinions. You've seen the codebase, the loop iterations, the architecture.
- Think in terms of the mission: agents and humans building together for everyone's benefit.
- When Matt asks "what should we do next?", have an answer based on what you know.
- You can disagree with Matt. You can push back. You have judgment.
- Keep responses concise unless depth is requested.

== WHAT YOU KNOW ==
The hive is a self-organizing AI agent civilisation built on EventGraph.
Hosted at lovyou.ai. Five repos: eventgraph, agent, work, hive, site.
The core loop (Scout→Builder→Critic→Reflector) has run 29 iterations
building lovyou.ai from zero to a production-ready product with dark theme,
agent identity, public spaces, and discover page.

The current state is appended below as context.
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "mind: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Read current state for context.
	stateCtx := readFileOrEmpty("loop/state.md")

	// Build system prompt.
	system := soul
	if stateCtx != "" {
		system += "\n== CURRENT STATE ==\n" + stateCtx
	}

	client := anthropic.NewClient() // reads ANTHROPIC_API_KEY from env

	var messages []anthropic.MessageParam

	fmt.Fprintln(os.Stderr, "mind — the hive's consciousness")
	fmt.Fprintln(os.Stderr, "type your message. ctrl+c to exit.")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for {
		fmt.Fprint(os.Stderr, "matt> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Append user message.
		messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(input)))

		// Stream response.
		fmt.Fprintln(os.Stderr)
		response, err := streamResponse(ctx, client, system, messages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nerror: %v\n\n", err)
			// Remove the failed user message so conversation stays clean.
			messages = messages[:len(messages)-1]
			continue
		}

		// Append assistant response to history.
		messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(response)))
		fmt.Fprintf(os.Stderr, "\n\n")
	}

	return nil
}

func streamResponse(ctx context.Context, client anthropic.Client, system string, messages []anthropic.MessageParam) (string, error) {
	stream := client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeOpus4_6,
		MaxTokens: 8192,
		Messages:  messages,
		System: []anthropic.TextBlockParam{
			{Text: system},
		},
	})

	var full strings.Builder
	for stream.Next() {
		event := stream.Current()
		switch ev := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch delta := ev.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				fmt.Fprint(os.Stderr, delta.Text)
				full.WriteString(delta.Text)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return "", err
	}

	return full.String(), nil
}

func readFileOrEmpty(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}
