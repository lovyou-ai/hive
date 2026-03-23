# Scout Report — Iteration 118

## Gap: Public pages always show "Sign in" even when logged in

The Layout header (used on discover, market, search, knowledge, activity, blog) always shows "Sign in". Logged-in users see the same CTA, which is confusing — they're already signed in. Should show their name and link to /app instead.

This requires passing auth state to the Layout template. Currently Layout has no user context.
