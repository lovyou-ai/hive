# Build Report — Iteration 103

Two notification triggers added to complete the feedback loop:

1. **complete op by agent** → notifies task author: "Hive completed your task: <title>"
2. **decompose op by agent** → notifies parent task author: "Hive broke down your task: <title>"

Both only fire when the actor is an agent and the author is someone else. 2 lines of notification logic each, using the existing `notify` helper from iter 102.
