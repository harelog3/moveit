# 🐘 Database migration recipes, powered by goose (https://pressly.github.io/goose/).
# Requires the goose CLI (`brew install goose`) and reads DATABASE_URL from backend/.env.

# Module recipes run from this file's directory (backend/just/), so move up to backend/.
set working-directory := '..'

# Migrations directory (relative to backend/)
migrations := "db/migrations"

# Runs goose against the database in backend/.env
[private]
goose +args:
    @set -a && . ./.env && set +a && \
        GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URL" GOOSE_MIGRATION_DIR={{ migrations }} \
        goose {{ args }}

# Creates a new SQL migration in db/migrations (sequential numbering)
new name:
    @echo "📝 Creating migration {{ name }}..."
    @mkdir -p {{ migrations }}
    @goose -dir {{ migrations }} -s create {{ name }} sql
    @echo "✅ Migration created"

# Applies every pending migration
up: (goose "up")
    @echo "✅ Database up to date"

# Applies only the next pending migration
up-one: (goose "up-by-one")
    @echo "✅ Next migration applied"

# Rolls back the last applied migration
down: (goose "down")
    @echo "↩️  Last migration rolled back"

# Rolls back the last migration and applies it again
redo: (goose "redo")
    @echo "🔁 Last migration redone"

# Rolls back every migration (⚠️ drops the schema data)
[confirm("⚠️ This drops the whole schema. Do you wish to continue?")]
reset: (goose "reset")
    @echo "💣 All migrations rolled back"

# Shows which migrations are applied and which are pending
status: (goose "status")

# Prints the current database version
version: (goose "version")

# Checks the migration files without touching the database
validate: (goose "validate")
    @echo "✅ Migrations are valid"
