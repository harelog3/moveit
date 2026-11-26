# CLAUDE.md

## Project

MoveIT is a Go rewrite of [Harel-ITESO/mov-e-backend](https://github.com/Harel-ITESO/mov-e-backend), a NestJS + Prisma Letterboxd clone (movie ratings, rating likes, follows, profiles, TMDB integration). The original repo is the **functional reference**: when porting a feature, read the original source first, e.g.

```bash
gh api repos/Harel-ITESO/mov-e-backend/contents/src/modules/ratings/ratings.service.ts -H 'Accept: application/vnd.github.raw'
gh api repos/Harel-ITESO/mov-e-backend/contents/prisma/schema.prisma -H 'Accept: application/vnd.github.raw'
```

The rewrite is intentionally smaller. Port behavior, not NestJS structure.

## Layout

Monorepo:

- `backend/`: Go module `github.com/harelog3/moveit` and `docker-compose.yaml` for local infra.
- `frontend/`: not created yet; it will be added later.

## Scope rules

- **No AWS, no Redis, no deploy, no Dockerfile for the app.** Docker is only for local infra: PostgreSQL, adobe/s3mock (S3-compatible, `:9090`) and Mailpit (SMTP `:1025`, UI `:8025`).
- Replacements for the original stack: DynamoDB email verifications become a Postgres table; Redis sessions and cache become Postgres or in-memory; Brevo becomes plain SMTP to Mailpit; Localstack S3 becomes s3mock through an S3-compatible client with a custom endpoint and path-style addressing.
- Don't add infrastructure, CI or deployment config unless asked.

## Commands

Common tasks are `just` recipes, run from the repo root (`just` lists them). The root `Justfile` only loads modules; each module file runs its recipes from its own directory:

- `backend/api.just` → `just api <recipe>`
- `backend/infra.just` → `just infra <recipe>`

```bash
just infra up                 # start Postgres, s3mock, Mailpit
just api run                  # run the API (cmd/api), loading backend/.env
just api build                # build the binary into backend/bin/moveit
just api exec                 # run the built binary, loading backend/.env
```

Checks, run from `backend/`:

```bash
go build ./... && go vet ./...
go test ./...
gofmt -l .                    # must print nothing
```

## Go conventions

- Prefer the standard library: `net/http` with `ServeMux` method and path patterns (`GET /v1/api/movies/{id}`), `log/slog`, `database/sql` or pgx. Add a dependency only when it clearly pays off.
- Put application code under `backend/internal/<domain>` (e.g. `auth`, `user`, `movie`, `rating`, `follow`, `profile`, `file`), with shared infra in packages like `internal/config`, `internal/database`, `internal/mail`, `internal/storage`, `internal/tmdb`.
- Use dependency injection through constructors. No package-level mutable state and no `init()` side effects.
- `context.Context` is the first parameter of anything that does I/O.
- Wrap errors with `fmt.Errorf("doing x: %w", err)`, map domain errors to HTTP status codes at the handler layer, and never leak internal errors to clients.
- Load config from environment variables once at startup and fail fast on missing required values. `backend/.env.example` documents them.
- Write schema changes as SQL migrations. Never mutate the schema ad hoc.
- Write table-driven tests with `testing`. Test handlers with `httptest`.
- Validate request input explicitly (the original relied on class-validator with whitelist and forbidNonWhitelisted, so reject unknown JSON fields).

## Language and documentation

- Everything written into the repo is in English: code comments, log and error messages, script and recipe output, docs (README, CLAUDE.md) and commit messages. This holds even when the conversation is in another language.
- Comment what isn't obvious from the code (why, constraints, gotchas). Skip comments that restate the code.
- Keep the README and CLAUDE.md in sync when commands, setup steps or project structure change.

## API

- Keep the `/v1/api` prefix and the original endpoint semantics unless a change is intentional and noted.
- JSON in and out, with consistent error bodies.
- Hash passwords with bcrypt. Session cookies are httpOnly.

## Git

- Conventional commits (`feat:`, `fix:`, `chore:`, `refactor:`, `test:`, `docs:`).
- Work on feature branches, with `main` as the base.
- Never commit `.env` or secrets.
