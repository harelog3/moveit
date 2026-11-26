# 🎬 MoveIT

A **Go** rewrite of [mov-e-backend](https://github.com/Harel-ITESO/mov-e-backend), a school project inspired by Letterboxd 🍿: users rate and review movies, like reviews, follow other users and build a profile with their favorite movies.

> 🧪 A smaller, code-focused version: no deploy, no AWS and no dockerized app (for now).

## 🧱 Stack

| Piece | Technology |
| --- | --- |
| 🐹 Backend | Go |
| 🐘 Database | PostgreSQL |
| 🪣 Files (avatars) | [adobe/s3mock](https://github.com/adobe/S3Mock) (S3-compatible API) |
| 📬 Email | [Mailpit](https://mailpit.axllent.org/) (local SMTP with a web UI) |
| 🎞️ Movie data | [TMDB API](https://developer.themoviedb.org/) |
| 🛠️ Task runner | [just](https://github.com/casey/just) |
| 🐘 Migrations | [goose](https://pressly.github.io/goose/) |
| 🖥️ Frontend | Coming soon to this same repo |

## 🔄 Original vs. rewrite

| Original (NestJS) | MoveIT (Go) |
| --- | --- |
| NestJS + TypeScript | Go 🐹 |
| Prisma ORM | SQL + migrations |
| Redis for sessions and cache | PostgreSQL / in-memory |
| DynamoDB for email verification | PostgreSQL table |
| Localstack (S3, SES, DynamoDB) | s3mock |
| Brevo (SMTP) | Mailpit |
| AWS ECS + Dockerfile + GitHub Actions | ❌ Out of scope for now |

## 📁 Structure

```
moveit/
├── Justfile          # 🛠️ Root of the project commands (loads the modules)
├── backend/          # 🐹 Go API + infrastructure docker-compose
│   ├── just/         # just modules: api.just, db.just, infra.just
│   ├── db/migrations/  # goose SQL migrations (also the sqlc schema)
│   ├── sqlc.yaml       # sqlc config: db/queries → internal/db
│   └── cmd/api/      # API entrypoint
└── frontend/  # 🖥️ (pending)
```

## 🚀 Getting started

Requirements: Go, Docker, [just](https://github.com/casey/just) and [goose](https://pressly.github.io/goose/) (`brew install goose`).

```bash
# 1. Start the infrastructure (PostgreSQL, s3mock and Mailpit) 🐳
just infra up

# 2. Set up the environment variables 🔐
cp backend/.env.example backend/.env   # add your TMDB_API_KEY

# 3. Apply the database migrations 🐘
just db up

# 4. Run the API ▶️
just api run
```

| Service | URL |
| --- | --- |
| 🐹 API | http://localhost:8080/v1/api |
| 🐘 PostgreSQL | `localhost:5432` (user/password/db: `moveit`) |
| 🪣 s3mock | http://localhost:9090 (bucket `moveit`) |
| 📬 Mailpit UI | http://localhost:8025 (SMTP on `:1025`) |

## 🛠️ Commands

Run `just` from the repo root to list every recipe. Recipes are grouped in modules: `just <module> <recipe>`.

| Recipe | What it does |
| --- | --- |
| `just api run` | Runs the API with `go run`, loading `backend/.env` |
| `just api build` | Builds the binary into `backend/bin/moveit` |
| `just api exec` | Runs the compiled binary, loading `backend/.env` |
| `just api get <pkg>` | Adds a Go package to the API |
| `just infra up` | Starts the containers |
| `just infra stop` | Stops the containers without removing them |
| `just infra down` | Removes the containers and network (keeps data) |
| `just infra destroy` | Removes containers, images, volumes and networks (⚠️ deletes the database data) |
| `just infra reset` | `infra destroy` + `infra up` |
| `just db new <name>` | Creates a SQL migration in `backend/db/migrations` |
| `just db up` | Applies every pending migration |
| `just db up-one` | Applies only the next pending migration |
| `just db down` | Rolls back the last migration |
| `just db redo` | Rolls back the last migration and applies it again |
| `just db reset` | Rolls back every migration (⚠️ deletes the database data) |
| `just db status` | Shows applied and pending migrations |
| `just db version` | Prints the current database version |
| `just db validate` | Checks the migration files without touching the database |

## 🗃️ Data model

- 👤 **user**: username, email, password (hashed), name, location, website, bio, avatar and up to 3 favorite movies.
- 🎞️ **movie**: stored only when someone interacts with it (TMDB id, title, genres, overview, poster, year and runtime).
- ⭐ **rating**: a user's score and review for a movie (one per user and movie).
- ❤️ **rating_like**: likes on reviews.
- 🤝 **user_follow**: who follows whom.

## 🗺️ Roadmap

- [ ] 🏗️ Project foundation (config, HTTP server, DB connection, migrations)
- [ ] 🔐 Authentication (sign up, login, sessions) + ✉️ email verification
- [ ] 👤 Users and account (edit details, favorite movies)
- [ ] 🖼️ Files (avatar uploads to s3mock)
- [ ] 🎞️ Movies (search and details via TMDB)
- [ ] ⭐ Ratings and ❤️ likes
- [ ] 🤝 Follows
- [ ] 🪪 Public profiles
- [ ] 🖥️ Frontend
