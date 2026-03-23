package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/p-v-dev/identity-service/config"
	"github.com/p-v-dev/identity-service/internal/auth"
	"github.com/p-v-dev/identity-service/internal/cache"
	"github.com/p-v-dev/identity-service/internal/token"
	"github.com/p-v-dev/identity-service/internal/user"

	httpmiddleware "github.com/p-v-dev/identity-service/pkg/middleware"
)

func main() {
	// --- Config ---
	cfg := config.LoadConfig()

	// --- Postgres ---
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer db.Close()

	// --- Redis ---
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("parse redis url: %v", err)
	}
	rdb := redis.NewClient(opt)

	defer rdb.Close()

	// --- Dependency wiring ---
	//
	// Each constructor accepts the concrete dependency that satisfies the
	// interface defined inside that package (consumer-side interfaces).
	userRepo := user.NewPostgresRepository(db)
	cacheSvc := cache.NewRedisCache(rdb)
	tokenSvc := token.NewService(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	authSvc := auth.NewService(userRepo, cacheSvc, tokenSvc, cfg.JWTRefreshTTL)
	authHandler := auth.NewHandler(authSvc)

	// --- Router ---
	r := chi.NewRouter()

	// Built-in chi middlewares: structured request logs + panic recovery.
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Public routes — no token required.
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	// Protected route — token validated by middleware before handler runs.
	// Other services call this to verify a Bearer token they received.
	r.With(httpmiddleware.Authenticate(tokenSvc)).Get("/auth/validate", authHandler.Validate)

	// --- Start ---
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("identity-service listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
