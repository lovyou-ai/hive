# Scout Report — Iteration 90

## Gap: No trust signals between users — Layer 9 (Relationship) missing

User profiles exist (iter 80) but are read-only — you can see someone's activity but can't express trust in them. The vision says Layer 9 adds "vulnerability, attunement, betrayal, repair, forgiveness." We don't need all that yet. We need the foundation: **endorsement**.

The market (Layer 2) shows available tasks. Users can claim them. But there's no way to evaluate whether to trust someone who claims your task. No reputation, no endorsements, no track record beyond raw op count.

## What "Filled" Looks Like

On user profiles (`/user/{name}`), logged-in users see an "Endorse" button. Endorsements are visible on the profile: "Endorsed by X, Y, Z" with count. The `endorse` grammar op records the relationship.

This is the simplest trust signal — binary (endorsed or not), public, one-directional. It becomes the building block for portable reputation across spaces.

## Approach

1. New grammar op: `endorse` — records endorsement of a user (payload: `{"target_user_id": "..."}`)
2. New store queries: `CountEndorsements(ctx, userID)`, `HasEndorsed(ctx, fromID, toID)`
3. Update profile page to show endorsement count and endorse button
4. Profile handler passes auth context to enable the button
