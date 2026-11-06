package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"sevenpay/services/account/internal/account"
	"sevenpay/services/account/internal/config"
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

	// Bez baze servis ne može raditi, pa se ne pokreće ako baza nije dostupna.
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("database ping: %v", err)
	}

	// Povezivanje slojeva: repozitorij dobiva bazu, servis repozitorij, a handler servis.
	repo := account.NewUserRepository(db)
	auth := account.NewAuthService(repo)
	handler := account.NewHandler(auth)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Printf("account service listening on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
