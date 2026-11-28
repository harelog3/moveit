package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/harelog3/moveit/internal/config"
	"github.com/harelog3/moveit/internal/db"
	"github.com/harelog3/moveit/internal/tmdb"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// WireDeps builds the application services from the loaded config.
func WireDeps(c config.Config) {
	_ = tmdb.NewService(c.TMDBApiKey)
}

func main() {
	// Load backend/.env when present (relative to the working directory) so the
	// binary runs outside `just`. Real environment variables take precedence,
	// and a missing file is fine because config.LoadEnv validates everything.
	_ = godotenv.Load(".env")

	// ctx is canceled on Ctrl+C or SIGTERM so in-flight work can shut down.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// env config
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	// db
	conn, err := pgx.Connect(ctx, env.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	_ = db.New(conn)

	WireDeps(env)
}
