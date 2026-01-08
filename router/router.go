package router

import (
	"myapp/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", handlers.RootHandler)
	r.Get("/dashboard", handlers.DashboardHandler)
	r.Get("/reports", handlers.ReportsHandler)

	return r
}