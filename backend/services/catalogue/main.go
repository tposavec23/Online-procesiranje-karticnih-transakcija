package main

import (
	"context"
	"log"
	"net/http"

	"sevenpay/backend/services/catalogue/internal/catalogue"
	"sevenpay/backend/services/catalogue/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
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

	// Error if the database cannot be pinged; the application cannot start because this is a critical error.
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("database ping: %v", err)
	}

	repo := catalogue.NewArticleRepository(db)
	svc := catalogue.NewArticleService(repo)
	handler := catalogue.NewHandler(svc)
	cartRepo := catalogue.NewCartRepository(db)
	cartSvc := catalogue.NewCartService(cartRepo)
	cartHandler := catalogue.NewCartHandler(cartSvc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	cartHandler.RegisterRoutes(mux)

	log.Printf("catalogue service listening on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
