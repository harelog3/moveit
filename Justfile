# 🎬 MoveIT: project commands
# Usage: `just` lists the recipes, `just <recipe>` runs one.

set default-list

# Binary name and path (relative to backend/)
bin := "bin/moveit"

# ─────────────────────────────────────────────
# 🚀 API
# ─────────────────────────────────────────────

# Runs the API with go run, loading backend/.env
[group('api')]
[working-directory('backend')]
run:
    @echo "🚀 Running project..."
    @set -a && . ./.env && set +a && go run ./cmd/api

# Builds the binary into backend/bin/
[group('api')]
[working-directory('backend')]
build:
    @echo "🔨 Building binary into backend/{{ bin }}..."
    @go build -o {{ bin }} ./cmd/api
    @echo "✅ Binary ready"

# Runs the compiled binary, loading backend/.env
[group('api')]
[working-directory('backend')]
exec:
    @test -x {{ bin }} || { echo "❌ backend/{{ bin }} not found, run 'just build' first"; exit 1; }
    @echo "▶️  Running binary..."
    @set -a && . ./.env && set +a && ./{{ bin }}

# Adds a package to the API project
[group('api')]
[working-directory('backend')]
get name:
    @echo "📦 Adding package {{ name }} to API..."
    @go get {{ name }}
    @echo "✅ Package added"

# ─────────────────────────────────────────────
# 🐳 Infra (docker compose: Postgres, s3mock, Mailpit)
# ─────────────────────────────────────────────

# Starts the containers in the background
[group('infra')]
[working-directory('backend')]
infra-up:
    @echo "🐳 Starting infrastructure..."
    @docker compose up -d
    @echo "✅ Infra up: Postgres :5432, s3mock :9090, Mailpit :8025"

# Stops the containers without removing them
[group('infra')]
[working-directory('backend')]
infra-stop:
    @echo "⏸️  Stopping containers..."
    @docker compose stop
    @echo "✅ Containers stopped"

# Removes the containers and network (keeps volumes and images)
[group('infra')]
[working-directory('backend')]
infra-down:
    @echo "🧹 Removing containers..."
    @docker compose down
    @echo "✅ Containers removed"

# Removes everything from this compose: containers, images, volumes and networks
[group('infra')]
[working-directory('backend')]
infra-destroy:
    @echo "💣 Destroying infrastructure (containers, images, volumes, networks)..."
    @docker compose down --rmi all --volumes --remove-orphans
    @echo "✅ Infra destroyed"

# Destroys everything and starts again from scratch
[group('infra')]
infra-reset: infra-destroy infra-up
    @echo "🔄 Infra reset from scratch"
