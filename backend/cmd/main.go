package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manassp21/drivex/backend/database"
	"github.com/manassp21/drivex/backend/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	ctx := context.Background()
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer pool.Close()

	database.DB = pool
	log.Println("connected to drivex db")

	r := routes.NewRouter()

	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("server failed:", err)
	}
}