# 🎬 MoveIT: project commands
# Usage: `just` lists every recipe, `just <module> <recipe>` runs one (e.g. `just api run`).
#
# Recipes live in one module per area. A module's recipes run from the
# directory that holds their file, so each module sets its own working-directory.

# 🚀 Go API: run, build, exec, get
mod api 'backend/just/api.just'

# 🐳 Local infrastructure (docker compose: Postgres, s3mock, Mailpit)
mod infra 'backend/just/infra.just'

# 🐘 Database migrations (goose)
mod db 'backend/just/db.just'

# Lists every recipe, including the ones inside modules
[private]
default:
    @just --list --list-submodules
