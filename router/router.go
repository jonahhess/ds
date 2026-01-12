package router

import (
	"net/http"
	"time"

	middlewares "myapp/middlewares"
	about "myapp/pages/about"
	home "myapp/pages/home"
	"myapp/pages/login"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

func SetupRoutes(sessionStore *sessions.CookieStore) *chi.Mux {
	r := chi.NewRouter()

	// --- core middleware ---
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middlewares.SessionMiddleware(sessionStore))

	// --- routes ---
	r.Route("/", func(r chi.Router) {
		r.Get("/", home.HomeHandler)
		r.Get("/about", about.AboutHandler)
		r.Get("/login", login.LoginHandler)
		r.Post("/login", login.LoginUserHandler)
	})

	// --- static files ---
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

	return r
}
