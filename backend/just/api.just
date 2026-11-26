# 🚀 Go API recipes.

# Module recipes run from this file's directory (backend/just/), so move up to backend/.
set working-directory := '..'

# Binary name and path (relative to backend/)
bin := "bin/moveit"

# Runs the API with go run, loading backend/.env
run:
    @echo "🚀 Running project..."
    @set -a && . ./.env && set +a && go run ./cmd/api

# Builds the binary into backend/bin/
build:
    @echo "🔨 Building binary into backend/{{ bin }}..."
    @go build -o {{ bin }} ./cmd/api
    @echo "✅ Binary ready"

# Runs the compiled binary, loading backend/.env
exec:
    @test -x {{ bin }} || { echo "❌ backend/{{ bin }} not found, run 'just api build' first"; exit 1; }
    @echo "▶️  Running binary..."
    @set -a && . ./.env && set +a && ./{{ bin }}

# Adds a package to the API project
get name:
    @echo "📦 Adding package {{ name }} to API..."
    @go get {{ name }}
    @echo "✅ Package added"
