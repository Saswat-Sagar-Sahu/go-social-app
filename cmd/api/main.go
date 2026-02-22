package main

import (
	"log"

	docs "github.com/Saswat-Sagar-Sahu/Social/docs"
	"github.com/Saswat-Sagar-Sahu/Social/internal/auth"
	"github.com/Saswat-Sagar-Sahu/Social/internal/db"
	"github.com/Saswat-Sagar-Sahu/Social/internal/email"
	"github.com/Saswat-Sagar-Sahu/Social/internal/env"
	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/joho/godotenv"
)

// @title Swagger Social API
// @version 1.0
// @description This is a sample server for a social media application.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide the JWT token as: Bearer <token>
// @BasePath /v1
func main() {
	godotenv.Load()

	cfg := config{
		address: env.GetString("ADDRESS", ":8000"),
		apiURL:  env.GetString("API_URL", "http://localhost:8000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:                          env.GetString("ENVIRONMENT", "development"),
		activationTokenExpiryMinutes: env.GetInt("ACTIVATION_TOKEN_EXPIRY_MINUTES", 60),
		jwtSecret:                    env.GetString("JWT_SECRET", "secret"),
		jwtExpiryMinutes:             env.GetInt("JWT_EXPIRY_MINUTES", 60),
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Connected to database successfully")

	store := store.NewStorage(db)

	// construct auth service using JWT
	authSvc := auth.NewJWTAuth(cfg.jwtSecret, cfg.jwtExpiryMinutes, store)

	// try to construct SendGrid sender from environment; nil if not configured
	sender, err := email.NewSendGridFromEnv()
	if err != nil {
		log.Printf("SendGrid not configured: %v", err)
	}

	app := &application{
		cfg,
		store,
		sender,
		authSvc,
	}

	docs.SwaggerInfo.BasePath = "/v1"

	mux := app.mount()
	log.Fatal(app.run(mux))
}
