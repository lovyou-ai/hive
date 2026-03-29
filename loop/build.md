# Build: Auth: email magic link as OAuth fallback

## What Was Built

### `site/auth/auth.go`

**Schema changes:**
- `users.google_id` made nullable (was `NOT NULL`) — magic-link users have no Google ID
- New table `magic_link_tokens`: `(id, token_hash, email, expires_at, used, created_at)`
- Migration added: `ALTER TABLE users ALTER COLUMN google_id DROP NOT NULL`

**New routes registered in `Register()`:**
- `GET /auth/magic-link/request` — email entry form
- `POST /auth/magic-link/request` — validate email, generate token, log link
- `GET /auth/magic-link/verify?token=...` — verify token, create session

**New handlers:**
- `handleMagicLinkRequestForm` — styled HTML form (Ember Minimalism, matches site design)
- `handleMagicLinkRequest` — validates email (must contain `@`, len ≥ 3); calls `requestMagicLink`; logs full link (`auth: magic link generated email=... link=...`); returns confirmation HTML
- `handleMagicLinkVerify` — empty token → redirect to `/auth/error?code=invalid_token`; calls `verifyMagicLink`; creates 30-day session cookie; redirects to `/app`

**New internal methods:**
- `requestMagicLink(ctx, email) (string, error)` — generates `newID()` raw token, hashes it, inserts into `magic_link_tokens` with 15-min expiry, returns raw token
- `verifyMagicLink(ctx, rawToken) (*User, error)` — atomic `UPDATE magic_link_tokens SET used=TRUE WHERE token_hash=$1 AND expires_at>NOW() AND used=FALSE RETURNING email`; handles expired/used/invalid in one query; calls `upsertUserByEmail`
- `upsertUserByEmail(ctx, email) (*User, error)` — `INSERT ... ON CONFLICT (email) DO UPDATE SET email=EXCLUDED.email RETURNING ...`; finds existing user by email or creates new one with `kind='human'`, no `google_id`

**Email delivery:** Stub — link is logged via `log.Printf`. TODO comment marks where SMTP/SendGrid wires in.

### `site/auth/auth_test.go`

No-DB tests (use `newTestAuth()`):
- `TestMagicLinkRequestInvalidEmail` — 3 subtests: empty, no_at, at_only → all return 400 before touching DB
- `TestMagicLinkVerifyMissingToken` — empty `?token=` → redirect to `/auth/error`

DB-required tests (skip without `DATABASE_URL`):
- `TestMagicLinkHappyPath` — request + verify → user created, second verify fails (used token)
- `TestMagicLinkExpiredToken` — inserts expired token, verify rejects it
- `TestMagicLinkInvalidToken` — bogus token rejected
- `TestMagicLinkIdempotentUser` — two tokens for same email → both verify resolves to same user ID

## Build Results

```
site: go.exe build -buildvcs=false ./...   → OK
site: go.exe test ./...                    → auth OK, graph OK, handlers OK
deploy: flyctl deploy --remote-only        → deployed to lovyou-ai.fly.dev
commit: 306ffe1 iter 407: Auth: email magic link as OAuth fallback
```

## Files Changed

- `site/auth/auth.go` — `magic_link_tokens` table, `google_id` nullable, 3 new routes, 3 handlers, 3 internal methods
- `site/auth/auth_test.go` — 6 new test functions (2 no-DB, 4 DB-required)
