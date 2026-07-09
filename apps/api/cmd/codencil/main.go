package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TheBlackHowling/codencil/apps/api/internal/auth"
	"github.com/TheBlackHowling/codencil/apps/api/internal/db"
	"github.com/TheBlackHowling/codencil/apps/api/internal/httpapi"
	"github.com/TheBlackHowling/codencil/apps/api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	database, err := db.Open(databaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	docStore := store.New(database)
	docHandler := httpapi.NewDocumentHandler(docStore)
	reviewHandler := httpapi.NewReviewHandler(docStore)

	authCfg := auth.LoadConfig()
	authMW, err := auth.NewMiddleware(docStore, authCfg)
	if err != nil {
		log.Fatalf("auth middleware: %v", err)
	}

	r := chi.NewRouter()
	r.Use(httpapi.CORSMiddleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/health", httpapi.Health)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Handler)
		docHandler.Register(r)
		reviewHandler.Register(r)
	})

	addr := ":" + port
	log.Printf("codencil api listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
