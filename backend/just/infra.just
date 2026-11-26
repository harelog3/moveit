# 🐳 Local infrastructure recipes (backend/docker-compose.yaml).

# Module recipes run from this file's directory (backend/just/), so move up to backend/.
set working-directory := '..'

# Starts the containers in the background
up:
    @echo "🐳 Starting infrastructure..."
    @docker compose up -d
    @echo "✅ Infra up: Postgres :5432, s3mock :9090, Mailpit :8025"

# Stops the containers without removing them
stop:
    @echo "⏸️  Stopping containers..."
    @docker compose stop
    @echo "✅ Containers stopped"

# Removes the containers and network (keeps volumes and images)
down:
    @echo "🧹 Removing containers..."
    @docker compose down
    @echo "✅ Containers removed"

# Removes everything from this compose: containers, images, volumes and networks
destroy:
    @echo "💣 Destroying infrastructure (containers, images, volumes, networks)..."
    @docker compose down --rmi all --volumes --remove-orphans
    @echo "✅ Infra destroyed"

# Destroys everything and starts again from scratch
reset: destroy up
    @echo "🔄 Infra reset from scratch"
