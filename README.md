# goCloudStorage

A self-hosted cloud storage app written in Go with server-rendered HTML + HTMX, passkey/GitHub authentication, and SQLite-backed metadata.

## What this project does

goCloudStorage provides a lightweight personal cloud storage experience:

- Register/login with **WebAuthn passkeys** or **GitHub OAuth**
- Upload files to per-user directories on disk
- Persist file metadata and auth/session data in SQLite
- Browse files from a dashboard and download/delete owned files

The application is intentionally simple: one Go service, one SQLite database (`cloud.db`), and local filesystem storage under `./storage/users/`.

## Design choices

### 1. Monolith over microservices

Everything (routing, auth flows, file handling, template rendering, DB access) lives in one process.  
This keeps deployment and local development very easy for a personal/self-hosted project.

### 2. HTMX + server templates instead of SPA framework

UI is rendered with Go templates in `web/templates` and progressively enhanced with HTMX interactions.

- Less frontend build complexity
- Fast iteration on backend-driven UI
- Partial HTML responses for actions like auth/file updates

### 3. SQLite for portability

`modernc.org/sqlite` is used as an embedded DB, so the app runs without an external database service.  
Schema is managed with SQL migration files in `migrations/` and applied at startup.

### 4. Local disk storage with DB metadata

Actual file bytes are stored on disk, while metadata (owner, path, MIME type, size, timestamps) is in `files` table.

- Good for simplicity and local hosting
- Clear separation between binary storage and queryable metadata

### 5. Dual authentication strategy

Auth supports both:

- **Passkeys** via `github.com/go-webauthn/webauthn`
- **GitHub OAuth** for social login

Sessions are stored server-side in DB and linked to a `session_id` cookie.

## How it works

### Request flow (high level)

1. App starts, validates GitHub OAuth env vars, opens SQLite DB, runs migrations, and parses templates.
2. Routes are registered in `main.go`.
3. Protected routes are wrapped by `requireAuth`, which resolves user from `session_id`.
4. Handlers in `controller_*.go` call service/db helpers and render templates or return fragments.

### Authentication flow

- **Passkey register/login**
  - Begin endpoint generates WebAuthn options + session challenge
  - Finish endpoint validates credential response
  - On success, creates DB session + secure auth cookie

- **GitHub OAuth**
  - User is redirected to GitHub with state token
  - Callback exchanges code for token and resolves/creates user
  - Session is created and written to cookie

### File flow

- Upload endpoint accepts multipart file
- `StorageService` sanitizes filename and writes to `./storage/users/{userID}/`
- If filename already exists, dedup suffix like `(1)` is appended
- DB metadata row is inserted in `files`
- Download/delete endpoints verify ownership via user context + DB lookup

## Project structure

```text
.
├── main.go                   # bootstrapping + route registration
├── config.go                 # constants + env wiring
├── controller_*.go           # HTTP handlers by domain
├── service_auth.go           # WebAuthn/OAuth integration logic
├── service_storage.go        # file IO logic
├── db_*queries.go            # SQL operations
├── db_migrationrunner.go     # startup migration runner
├── migrations/               # SQL schema migration files
├── web/templates/            # layouts/pages/components
└── storage/users/            # uploaded user files
```

## Data model (core tables)

- `users`
- `oauth_accounts`
- `passkeys`
- `webauthn_sessions`
- `sessions`
- `files`
- `shares` (schema exists; sharing handlers are not finished)

See `migrations/001_initial.sql` for full schema.

## Setup and run

## Prerequisites

- Go **1.25+**

## Environment

Create `.env` (or export vars) with:

```env
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
```

> Note: The app currently hardcodes WebAuthn RP settings for `localhost` in `service_auth.go`.

## Run

```bash
go run .
```

Server starts on `http://localhost:8080`.

## Current status

Implemented:

- HTMX-based UI and dashboard rendering
- Passkey auth flow
- GitHub OAuth login flow
- Session-based protected routes
- File upload/download/delete

Not fully implemented yet:

- File sharing endpoints/UI flow
- User profile/list/delete handlers
- File rename/move/search

## Future features and improvements

1. Complete file sharing (create/list/revoke shares, permissions, shared-with-me view).
2. Add file operations: rename, move, copy, bulk actions, and search/filter.
3. Improve security hardening (strict cookie settings by environment, CSRF strategy for state-changing endpoints, secure token storage/rotation).
4. Add background jobs for cleanup (expired WebAuthn sessions, expired auth state tokens, stale sessions).
5. Add tests (unit tests for services/query helpers + integration tests for auth/file flows).
6. Add observability (structured logs, request tracing, basic metrics/health reporting).
7. Improve UX (drag-and-drop, file previews, richer error states).
