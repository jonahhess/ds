package router

import (
	"context"
	"net/http"

	"myapp/layouts"
	home "myapp/pages/home"
	"myapp/types"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupRoutes(
	appCtx context.Context,
	sessionStore *sessions.CookieStore,
) *chi.Mux {

	r := chi.NewRouter()

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
		// example session usage
		session, _ := sessionStore.Get(r, "myapp")

		userType, ok := session.Values["userType"].(types.UserType)
		if !ok {
			userType = types.UserType(0)
			session.Values["userType"] = userType
		}
		_ = session.Save(r, w)

		err := layouts.
			Base("Home", home.Home(userType)).
			Render(r.Context(), w)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return r
}