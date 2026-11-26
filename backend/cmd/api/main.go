package main

import (
	"context"
	"log"

	"github.com/harelog3/moveit/internal/config"
	"github.com/harelog3/moveit/internal/db"
	"github.com/jackc/pgx/v5"
)

func main() {

	ctx := context.Background()

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

}
