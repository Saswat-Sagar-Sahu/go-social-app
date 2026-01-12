package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Saswat-Sagar-Sahu/Social/internal/db"
	"github.com/Saswat-Sagar-Sahu/Social/internal/env"
	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	usersN := flag.Int("users", 500, "number of users to create")
	postsN := flag.Int("posts", 1000, "number of posts to create")
	flag.Parse()

	cfgAddr := env.GetString("DB_ADDR", "postgres://user:password@localhost:5432/social?sslmode=disable")
	maxOpen := env.GetInt("DB_MAX_OPEN_CONNS", 30)
	maxIdle := env.GetInt("DB_MAX_IDLE_CONNS", 30)
	maxIdleTime := env.GetString("DB_MAX_IDLE_TIME", "15m")

	dbConn, err := db.New(cfgAddr, maxOpen, maxIdle, maxIdleTime)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer dbConn.Close()

	storeSvc := store.NewStorage(dbConn)
	ctx := context.Background()

	log.Printf("Seeding %d users...", *usersN)
	for i := 1; i <= *usersN; i++ {
		u := &store.User{
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "password123",
		}
		if err := storeSvc.Users.Create(ctx, u); err != nil {
			log.Fatalf("create user %d: %v", i, err)
		}
		if i%50 == 0 {
			log.Printf("created %d users", i)
		}
	}

	rand.Seed(time.Now().UnixNano())
	log.Printf("Seeding %d posts...", *postsN)
	for i := 1; i <= *postsN; i++ {
		uid := int64(rand.Intn(*usersN) + 1)
		p := &store.Post{
			Title:   fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("This is the content for post %d.", i),
			UserID:  sql.NullInt64{Int64: uid, Valid: true},
			Tags:    []string{"seeded", "demo"},
		}
		if err := storeSvc.Posts.Create(ctx, p); err != nil {
			log.Fatalf("create post %d: %v", i, err)
		}
		if i%100 == 0 {
			log.Printf("created %d posts", i)
		}
	}

	log.Println("Seeding complete.")
}
