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
├── Justfile   # 🛠️ Project commands
├── backend/   # 🐹 Go API + infrastructure docker-compose
│   └── cmd/api/   # API entrypoint
└── frontend/  # 🖥️ (pending)
```

## 🚀 Getting started

Requirements: Go, Docker and [just](https://github.com/casey/just).

```bash
# 1. Start the infrastructure (PostgreSQL, s3mock and Mailpit) 🐳
just infra-up

# 2. Set up the environment variables 🔐
cp backend/.env.example backend/.env   # add your TMDB_API_KEY

# 3. Run the API ▶️
just run
```

| Service | URL |
| --- | --- |
| 🐹 API | http://localhost:8080/v1/api |
| 🐘 PostgreSQL | `localhost:5432` (user/password/db: `moveit`) |
| 🪣 s3mock | http://localhost:9090 (bucket `moveit`) |
| 📬 Mailpit UI | http://localhost:8025 (SMTP on `:1025`) |

## 🛠️ Commands

Run `just` from the repo root to list every recipe.

| Recipe | What it does |
| --- | --- |
| `just run` | Runs the API with `go run`, loading `backend/.env` |
| `just build` | Builds the binary into `backend/bin/moveit` |
| `just exec` | Runs the compiled binary, loading `backend/.env` |
| `just get <pkg>` | Adds a Go package to the API |
| `just infra-up` | Starts the containers |
| `just infra-stop` | Stops the containers without removing them |
| `just infra-down` | Removes the containers and network (keeps data) |
| `just infra-destroy` | Removes containers, images, volumes and networks (⚠️ deletes the database data) |
| `just infra-reset` | `infra-destroy` + `infra-up` |

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
