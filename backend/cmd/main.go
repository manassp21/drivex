package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/manassp21/drivex/backend/database"
	"github.com/manassp21/drivex/backend/internal/storage"
	"github.com/manassp21/drivex/backend/routes"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	ctx := context.Background()

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		log.Fatal("S3_BUCKET environment variable is not set")
	}

	store, err := storage.NewS3Storage(ctx, bucket)
	if err != nil {
		log.Fatal("failed to set up S3 storage:", err)
	}
	log.Println("connected to S3 bucket:", bucket)

	r := routes.NewRouter(store)

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

	log.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("server failed:", err)
	}
}
