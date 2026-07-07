// Package main is the entrypoint for the Natter API server.
// It wires up the repository, service, and HTTP layers, then starts listening on :8080.
//
// Environment variables:
//   - DATABASE_URL: PostgreSQL connection string (optional; defaults to in-memory store).
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/mrl00/natter/internal/api"
	"github.com/mrl00/natter/internal/infra/memory"
	"github.com/mrl00/natter/internal/infra/postgres"
	"github.com/mrl00/natter/internal/repository"
	"github.com/mrl00/natter/internal/service"
)

func main() {
	var repo repository.Repository

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("failed to open database: %v", err)
		}
		defer db.Close()

		pg := postgres.New(db)
		if err := pg.Ping(context.Background()); err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}
		if err := pg.Migrate(context.Background()); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}

		fmt.Println("connected to postgres")
		repo = pg
	} else {
		fmt.Println("using in-memory store")
		repo = memory.New()
	}

	svc := service.New(repo)
	h := api.NewHandlers(svc)
	mux := api.NewRouter(h)

	addr := ":8080"
	fmt.Printf("natter API listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
