# 🎬 MoveIT: comandos del proyecto
# Uso: `just` lista las recetas, `just <receta>` ejecuta una.

set default-list

# Nombre y ruta del binario (relativa a backend/)
bin := "bin/moveit"

# ─────────────────────────────────────────────
# 🚀 App
# ─────────────────────────────────────────────

# Corre la API con go run, cargando backend/.env
[group('api')]
[working-directory('backend')]
run:
    @echo "🚀 Corriendo proyecto..."
    @set -a && . ./.env && set +a && go run ./cmd/api

# Compila el binario en backend/bin/
[group('api')]
[working-directory('backend')]
build:
    @echo "🔨 Compilando binario en backend/{{ bin }}..."
    @go build -o {{ bin }} ./cmd/api
    @echo "✅ Binario listo"

# Ejecuta el binario compilado, cargando backend/.env
[group('api')]
[working-directory('backend')]
exec:
    @test -x {{ bin }} || { echo "❌ No existe backend/{{ bin }}, corre 'just build' primero"; exit 1; }
    @echo "▶️  Ejecutando binario..."
    @set -a && . ./.env && set +a && ./{{ bin }}

# ─────────────────────────────────────────────
# 🐳 Infra (docker compose: Postgres, s3mock, Mailpit)
# ─────────────────────────────────────────────

# Levanta los contenedores en segundo plano
[group('infra')]
[working-directory('backend')]
infra-up:
    @echo "🐳 Levantando infraestructura..."
    @docker compose up -d
    @echo "✅ Infra arriba: Postgres :5432, s3mock :9090, Mailpit :8025"

# Detiene los contenedores sin eliminarlos
[group('infra')]
[working-directory('backend')]
infra-stop:
    @echo "⏸️  Deteniendo contenedores..."
    @docker compose stop
    @echo "✅ Contenedores detenidos"

# Elimina los contenedores y la red (conserva volúmenes e imágenes)
[group('infra')]
[working-directory('backend')]
infra-down:
    @echo "🧹 Eliminando contenedores..."
    @docker compose down
    @echo "✅ Contenedores eliminados"

# Elimina todo lo de este compose: contenedores, imágenes, volúmenes y redes
[group('infra')]
[working-directory('backend')]
infra-destroy:
    @echo "💣 Destruyendo infraestructura (contenedores, imágenes, volúmenes, redes)..."
    @docker compose down --rmi all --volumes --remove-orphans
    @echo "✅ Infra destruida"

# Destruye todo y vuelve a levantar desde cero
[group('infra')]
infra-reset: infra-destroy infra-up
    @echo "🔄 Infra reiniciada desde cero"
