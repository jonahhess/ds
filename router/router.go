package router

import (
	"context"
	"net/http"

	"myapp/auth"
	"myapp/layouts"
	home "myapp/pages/home"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(
	appCtx context.Context,
	sessionStore *sqliteStore,
) *chi.Mux {

	r := chi.NewRouter()
	r.Use(auth.AuthMiddleware(sessionStore))

	// ---- static files ----
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

	// ---- routes ----
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.UserFromContext(r.Context())
		if !ok {
			http.Error(w,"User not found",http.StatusInternalServerError)
		}
		
		err := layouts.
			Base("Home", user, home.Home(user.UserType)).
			Render(r.Context(), w)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return r
}