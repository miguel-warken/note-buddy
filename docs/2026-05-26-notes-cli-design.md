# Notes CLI V1 Design

## Context

The project is a Go application focused on fast note capture during class through a terminal interface. The first version should optimize for low friction capture, durable local persistence, and easy inspection of saved sessions before adding AI enrichment or a web UI.

This document defines only phase 1 of the broader product:

1. Interactive CLI capture
2. Local SQLite persistence
3. Basic inspection of drafts and saved sessions

AI enrichment and the web interface are explicitly out of scope for V1.

## Goals

- Capture notes in a terminal-first workflow that feels natural during class
- Persist every entered line immediately to reduce risk of data loss
- Organize notes as class sessions with required subject and optional title
- Support inspection of both saved sessions and unfinished drafts
- Keep the implementation simple enough to teach clear Go application structure

## Non-Goals

- AI enrichment
- Download/export features
- Multi-user support
- Sync across devices
- Editing or deleting individual note lines after they are written
- Complex command systems inside the interactive mode

## User Experience

### Start capture

The user runs:

```text
notes capture
```

The CLI asks for:

- subject: required
- title: optional

Once the prompts are completed, the CLI opens an interactive capture session.

### Capture behavior

- Each line entered by pressing Enter is written immediately to a draft cache in SQLite
- Normal input is treated as note content
- The capture model is append-only in V1
- `/save` finalizes the session
- `/cancel` discards the draft session
- `Ctrl+D` should behave like a save shortcut

### Draft recovery

If a draft session already exists when the user starts a new capture, the CLI should not silently replace it. Instead, it should inform the user that an unfinished session exists and offer:

- resume
- discard

V1 should allow only one active draft at a time.

### Session inspection

The CLI should support:

```text
notes list
notes show <id>
```

`notes list` shows:

- internal numeric ID for saved sessions
- subject
- optional title
- started or saved time context
- number of lines
- status (`draft` or `saved`)

If an active draft exists, it may be shown in the list with a literal draft marker instead of a numeric session ID, since draft and saved records live in separate tables in V1.

`notes show <id>` shows:

- session metadata
- reconstructed running text
- individual lines with timestamps

In V1, `notes show <id>` targets saved sessions only. Active drafts are managed through the capture recovery flow rather than direct inspection commands.

## Recommended Project Layout

```text
notes/
├─ cmd/
│  └─ notes/
│     └─ main.go
├─ internal/
│  ├─ app/
│  │  ├─ capture_service.go
│  │  ├─ list_service.go
│  │  └─ show_service.go
│  ├─ cli/
│  │  ├─ capture_command.go
│  │  ├─ list_command.go
│  │  ├─ show_command.go
│  │  └─ prompt.go
│  ├─ domain/
│  │  ├─ session.go
│  │  ├─ draft_entry.go
│  │  └─ saved_entry.go
│  └─ store/
│     └─ sqlite/
│        ├─ db.go
│        ├─ draft_store.go
│        ├─ session_store.go
│        └─ migrations/
│           └─ 001_init.sql
├─ docs/
├─ README.md
└─ go.mod
```

### Responsibility split

- `cmd/notes`: application entrypoint and dependency wiring
- `internal/cli`: terminal interaction and output formatting
- `internal/app`: use-case orchestration
- `internal/domain`: core data structures and core invariants
- `internal/store/sqlite`: persistence implementation and schema

## Architecture

V1 should follow a simple application structure:

1. The CLI layer reads user input and prints output
2. The app layer coordinates the capture/list/show flows
3. The store layer reads and writes SQLite data
4. The domain layer defines the system vocabulary

The important boundary is that the app logic should not be mixed directly into SQL statements or prompt rendering. The CLI should ask for input and print results, but the rules for opening a draft, appending entries, saving, and recovering unfinished work should live in the app layer.

Interfaces should not be introduced preemptively. Concrete SQLite-backed structs are sufficient for V1 unless a clear need for abstraction appears during testing or a later refactor.

## Data Model

V1 uses separate draft and saved tables so the lifecycle is explicit and easy to reason about.

### `draft_sessions`

- `id`
- `subject`
- `title`
- `started_at`
- `status`

Expected status values for V1:

- `active`

### `draft_entries`

- `id`
- `draft_session_id`
- `content`
- `created_at`

### `sessions`

- `id`
- `subject`
- `title`
- `started_at`
- `saved_at`

### `session_entries`

- `id`
- `session_id`
- `content`
- `created_at`
- `position`

## Data Flow

### Open capture

1. User runs `notes capture`
2. Application checks whether an active draft exists
3. If a draft exists, the user chooses to resume or discard
4. If no draft is active, a new draft session is created

### Append line

1. User types a line and presses Enter
2. The line is inserted into `draft_entries`
3. The draft remains active until save or cancel

### Save session

Saving must happen inside a single database transaction:

1. Create a row in `sessions`
2. Copy draft entries into `session_entries`
3. Preserve entry order using `position`
4. Delete `draft_entries` for the draft
5. Delete the `draft_session`
6. Commit the transaction

If any step fails, the transaction must be rolled back.

### Cancel session

1. Delete the draft entries for the active draft
2. Delete the draft session

## Command Scope

### `notes capture`

- prompts for subject and title
- starts or resumes a draft session
- appends entries line-by-line
- accepts `/save`, `/cancel`, and `Ctrl+D`

### `notes list`

- lists saved sessions
- includes a draft indicator when relevant
- shows line counts and status

### `notes show <id>`

- shows session metadata
- reconstructs the running text by joining saved entries in order
- shows individual lines and their timestamps
- only accepts saved session IDs in V1

## Error Handling

V1 should explicitly handle:

- database open/create failure
- invalid empty subject
- failure to append a draft line
- failure during save transaction
- attempt to start capture when a draft exists
- invalid session ID in `show`

Errors should be clear and terminal-friendly. Internal SQL details should not leak into the default user messages unless debugging output is intentionally added later.

## Testing Priorities

The first tests should focus on behavior, not framework cleverness:

- create a draft session
- append a draft entry on each line
- detect an existing unfinished draft
- discard a draft session
- save a draft session into final tables
- verify save rollback on failure
- list sessions with status and entry counts
- show a session with reconstructed text and per-line timestamps

Test placement should stay close to the code under test, for example:

- `capture_service_test.go`
- `list_service_test.go`
- `show_service_test.go`

## Future Phases

### Phase 2

Add AI enrichment that reads saved sessions and generates derived learning material such as:

- topic summaries
- contextual explanations
- possible corrections or clarifications
- study questions or flashcards

### Phase 3

Add a web interface that exposes:

- raw captured notes
- enriched notes
- a download action for later export

## Open Decisions Deferred Intentionally

These questions are valid, but intentionally deferred to keep V1 small:

- exact CLI framework choice
- migration tool choice
- export format
- AI provider and prompt strategy
- web stack and authentication model

## Recommendation

Begin implementation with the core capture loop and SQLite schema before adding refinements. The best first milestone is:

1. open database
2. create schema
3. start capture
4. append lines into draft entries
5. save into final tables
6. inspect through `list` and `show`

This keeps the first learning loop short while still exercising useful application structure.
