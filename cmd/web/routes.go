package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"meisa_xyz/pkg/config"
	"meisa_xyz/pkg/handlers"
	"net/http"
)

func routes(_ *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
