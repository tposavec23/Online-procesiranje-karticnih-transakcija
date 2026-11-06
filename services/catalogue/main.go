package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"sevenpay/services/catalogue/internal/article"
	"sevenpay/services/catalogue/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	// Error ako se ne može pingati baza, aplikacija se ne može pokrenuti jer je to kritična greška.
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("database ping: %v", err)
	}

	// Repozitorij, servis i handler za artikle. Ovdje se povezuju slojevi aplikacije.
	repo := article.NewPostgresRepository(db)
	svc := article.NewService(repo)
	handler := article.NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Printf("catalogue service listening on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
