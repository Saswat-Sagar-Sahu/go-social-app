package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	docs "github.com/Saswat-Sagar-Sahu/Social/docs"
	"github.com/Saswat-Sagar-Sahu/Social/internal/auth"
	"github.com/Saswat-Sagar-Sahu/Social/internal/email"
	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	config config
	Store  store.Storage
	Email  email.Sender
	Auth   auth.Authenticator
}

type config struct {
	address                      string
	apiURL                       string
	db                           dbConfig
	env                          string
	activationTokenExpiryMinutes int
	jwtSecret                    string
	jwtExpiryMinutes             int
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	// enforce auth for unsafe methods globally
	r.Use(app.methodAuthMiddleware)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {

		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostsHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Get("/", app.getPostsHandler)
			})
		})
		r.Route("/comments", func(r chi.Router) {
			r.Post("/", app.createCommentsHandler)
			r.Get("/post/{postId}", app.getCommentsByPostIdHandler)
			r.Get("/user/{userId}", app.getCommentsByUserIdHandler)
			r.Delete("/{commentId}", app.deleteCommentByIdHandler)
			r.Put("/{commentId}", app.updateCommentByIdHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userId}", func(r chi.Router) {
				r.Get("/", app.getUsersHandler)

				r.Post("/follow", app.followUserHandler)
				r.Post("/unfollow", app.unfollowUserHandler)
			})
			r.Route("/feed", func(r chi.Router) {
				r.With(app.authMiddleware).Get("/", app.getUserFeedHandler)
			})

			// registration and activation
			r.Post("/register", app.registerUserHandler)
			r.Post("/login", app.loginHandler)
			r.Post("/activate", app.activateUserHandler)
		})
	})

	// Serve swagger UI and the generated swagger JSON
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(app.config.apiURL+"/swagger/doc.json")))
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		doc := docs.SwaggerInfo.ReadDoc()

		var pretty any
		if err := json.Unmarshal([]byte(doc), &pretty); err != nil {
			http.Error(w, "invalid swagger json", http.StatusInternalServerError)
			return
		}

		b, err := json.MarshalIndent(pretty, "", "    ")
		if err != nil {
			http.Error(w, "failed to encode swagger json", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(b)
	})
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.address,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", app.config.address)
	return srv.ListenAndServe()
}
