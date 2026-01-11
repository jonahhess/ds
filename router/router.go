package router

import (
	"context"
	pages "myapp/pages/home"
	"myapp/templates/layouts"
	"myapp/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Handle(
	"/static/pages/*",
	http.StripPrefix(
		"/static/pages/",
		http.FileServer(http.Dir("./pages")),
	),
)

r.Handle(
	"/static/components/*",
	http.StripPrefix(
		"/static/components/",
		http.FileServer(http.Dir("./templates/components")),
	),
)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		layouts.Base("Home",pages.HomeHead(), pages.HomePage(types.UserType(0))).Render(context.TODO(),w)	
	})

	return r
}