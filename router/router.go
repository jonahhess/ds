package router

import (
	"context"
	"net/http"

	pages "myapp/pages/home"
	"myapp/types"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Home(types.UserType(0)).Render(context.TODO(), w)
	})

	return r
}