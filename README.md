# 🎬 MoveIT

Reescritura en **Go** de [mov-e-backend](https://github.com/Harel-ITESO/mov-e-backend), un proyecto escolar inspirado en Letterboxd 🍿: los usuarios califican y comentan películas, dan like a reseñas, siguen a otros usuarios y arman su perfil con sus películas favoritas.

> 🧪 Versión reducida y enfocada en el código: sin deploy, sin AWS y sin dockerizar la app (por ahora).

## 🧱 Stack

| Pieza | Tecnología |
| --- | --- |
| 🐹 Backend | Go |
| 🐘 Base de datos | PostgreSQL |
| 🪣 Archivos (avatares) | [adobe/s3mock](https://github.com/adobe/S3Mock) (API compatible con S3) |
| 📬 Correo | [Mailpit](https://mailpit.axllent.org/) (SMTP local con UI) |
| 🎞️ Datos de películas | [TMDB API](https://developer.themoviedb.org/) |
| 🖥️ Frontend | Próximamente en este mismo repo |

## 🔄 Original vs. reescritura

| Original (NestJS) | MoveIT (Go) |
| --- | --- |
| NestJS + TypeScript | Go 🐹 |
| Prisma ORM | SQL + migraciones |
| Redis para sesiones y caché | PostgreSQL / memoria |
| DynamoDB para verificación de email | Tabla en PostgreSQL |
| Localstack (S3, SES, DynamoDB) | s3mock |
| Brevo (SMTP) | Mailpit |
| AWS ECS + Dockerfile + GitHub Actions | ❌ Fuera de alcance por ahora |

## 📁 Estructura

```
moveit/
├── backend/   # 🐹 API en Go + docker-compose de infraestructura
└── frontend/  # 🖥️ (pendiente)
```

## 🚀 Cómo correrlo

Requisitos: Go y Docker.

```bash
cd backend

# 1. Levanta la infraestructura (PostgreSQL, s3mock y Mailpit) 🐳
docker compose up -d

# 2. Configura las variables de entorno 🔐
cp .env.example .env   # agrega tu TMDB_API_KEY

# 3. Corre la API ▶️
go run .
```

| Servicio | URL |
| --- | --- |
| 🐹 API | http://localhost:8080/v1/api |
| 🐘 PostgreSQL | `localhost:5432` (usuario/contraseña/db: `moveit`) |
| 🪣 s3mock | http://localhost:9090 (bucket `moveit`) |
| 📬 Mailpit UI | http://localhost:8025 (SMTP en `:1025`) |

## 🗃️ Modelo de datos

- 👤 **user**: username, email, contraseña (hash), nombre, ubicación, sitio web, bio, avatar y hasta 3 películas favoritas.
- 🎞️ **movie**: se guarda solo cuando alguien interactúa con ella (id de TMDB, título, géneros, sinopsis, póster, año y duración).
- ⭐ **rating**: calificación y comentario de un usuario a una película (una por usuario y película).
- ❤️ **rating_like**: likes a las reseñas.
- 🤝 **user_follow**: quién sigue a quién.

## 🗺️ Roadmap

- [ ] 🏗️ Base del proyecto (config, servidor HTTP, conexión a DB, migraciones)
- [ ] 🔐 Autenticación (registro, login, sesiones) + ✉️ verificación de email
- [ ] 👤 Usuarios y cuenta (editar datos, películas favoritas)
- [ ] 🖼️ Archivos (subida de avatares a s3mock)
- [ ] 🎞️ Películas (búsqueda y detalle vía TMDB)
- [ ] ⭐ Calificaciones y ❤️ likes
- [ ] 🤝 Follows
- [ ] 🪪 Perfiles públicos
- [ ] 🖥️ Frontend
