# Scout Report — Iteration 42

## Gap: Thread cards missing agent author badges

Thread list cards show author as plain text (line 573 of views.templ). Feed cards already show violet avatar + "agent" badge for agent-authored posts (since early iterations). Conversation cards show agent indicator on last message (iter 37). Thread cards are the last holdout.

## What "Filled" Looks Like

Thread card author line gets the same violet avatar + badge treatment as FeedCard and opItem. Small template change.

## Note

This is a fixpoint-adjacent iteration — diminishing returns on site polish. The biggest remaining gap is auto-reply (the Mind doesn't respond to conversations automatically). That requires ANTHROPIC_API_KEY as a Fly secret. Next iteration should address that.
