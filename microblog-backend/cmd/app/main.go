package main

import (
	"context"
	"fmt"
	"microblog-backend/internal/config"
	"microblog-backend/internal/db"
	"microblog-backend/internal/handlers/signin"
	"microblog-backend/internal/handlers/signout"
	"microblog-backend/internal/handlers/signup"
	"microblog-backend/internal/handlers/test"
	"microblog-backend/internal/logger"
	"microblog-backend/internal/mw/api"
	"microblog-backend/internal/mw/auth"
	"microblog-backend/internal/repository/credentials"
	"microblog-backend/internal/repository/session"
	"microblog-backend/internal/repository/user"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Errorf("failed to parse config: %w", err)
		os.Exit(1)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	postgresDb, err := sqlx.ConnectContext(ctx, "postgres", "host=postgres user=changeme password=changeme dbname=changeme sslmode=disable")
	if err != nil {
		log.Errorf("failed to connect to postgres: %w", err)
		os.Exit(1)
	}
	db := db.New(postgresDb)

	sessionRepo := session.New(redisClient)
	userRepo := user.New(db)
	credentialsRepo := credentials.New(db)

	auth := auth.New(sessionRepo, userRepo)

	mux := chi.NewMux()
	mux.Use(logger.Inject(log))
	mux.Use(auth.AuthMW)
	mux.MethodNotAllowed(api.MethodNotAllowed())
	mux.NotFound(api.NotFound())

	mux.Post("/test", api.HandlerMW(test.New().Handle))
	mux.Post("/signin", api.HandlerMW(signin.New(credentialsRepo, sessionRepo).Handle))
	mux.Post("/signup", api.HandlerMW(signup.New(db, userRepo, credentialsRepo).Handle))
	mux.Post("/signout", api.HandlerMW(signout.New().Handle))

	log.Infof("starting server at port=%s", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux); err != nil {
		log.Errorf("failed to serve http: %w", err)
	}
}
