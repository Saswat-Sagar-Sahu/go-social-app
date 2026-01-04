package main

import (
	"log"

	"github.com/Saswat-Sagar-Sahu/Social/internal/db"
	"github.com/Saswat-Sagar-Sahu/Social/internal/env"
	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config{
		address: env.GetString("ADDRESS", ":8000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
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

	app := &application{
		cfg,
		store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
