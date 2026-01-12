package router

import (
	"context"
	home "myapp/pages/home"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupRoutes(
	appCtx context.Context,
	sessionStore *sessions.CookieStore,
) *chi.Mux {

	r := chi.NewRouter()

	r.Handle("/static/*",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	r.Handle("/static/pages/*",
		http.StripPrefix("/static/pages/",
			http.FileServer(http.Dir("./pages")),
		),
	)

	r.Handle("/static/components/*",
		http.StripPrefix("/static/components/",
			http.FileServer(http.Dir("./components")),
		),
	)

	// middleware: get cookie, store user

	r.Get("/", home.HomeHandler)

	return r
}
