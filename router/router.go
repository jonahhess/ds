package router

import (
	"database/sql"
	"net/http"
	"time"

	"myapp/components/navbar"
	middlewares "myapp/middlewares"
	about "myapp/pages/about"
	"myapp/pages/courses"
	home "myapp/pages/home"
	"myapp/pages/login"
	"myapp/pages/notFound"
	"myapp/pages/review"
	"myapp/pages/signup"
	"myapp/pages/study"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

func SetupRoutes(sessionStore *sessions.CookieStore, DB *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middlewares.SessionMiddleware(sessionStore))

	r.Route("/", func(r chi.Router) {
		r.Get("/", home.Page)
		r.Get("/about", about.Page)
		r.Get("/login", login.Page)
		r.Post("/login", login.LoginHandler)
		r.Get("/signup", signup.Page)
		r.Post("/signup", signup.SignupHandler)
		r.Post("/logout", navbar.LogoutHandler)
	})
	
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Get("/study", study.Page)
		r.Get("/courses",courses.Page(DB))
		r.Get("/review", review.Page(DB))
	})

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

	r.NotFound(notFound.Page)

	return r
}
