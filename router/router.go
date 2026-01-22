package router

import (
	"database/sql"
	"net/http"
	"time"

	"myapp/auth"
	about "myapp/pages/about"
	"myapp/pages/courses"
	home "myapp/pages/home"
	"myapp/pages/login"
	"myapp/pages/myCourses"
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

	r.Use(auth.SessionMiddleware(sessionStore))
	r.Use(auth.OptionalUserMiddleware)

	r.Group( func(r chi.Router) {
		r.Get("/", home.Page)
		r.Get("/about", about.Page)
		r.Get("/login", login.Page)
		r.Post("/login", auth.LoginHandler)
		r.Get("/signup", signup.Page)
		r.Post("/signup", signup.SignupHandler)
		r.Post("/logout", auth.LogoutHandler)
		r.NotFound(notFound.Page)
	})
	
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuthMiddleware)
		r.Get("/study", study.Page)
		r.Get("/courses",courses.Page(DB))
		r.Route("/review", func(r chi.Router) {
			r.Get("/", review.Page(DB))
			r.Get("/next", review.GetNextCard(DB))
			r.Post("/submit", review.SubmitAnswer(DB))
		})

		r.Route("/users/{userID}", func(r chi.Router) {
			r.Get("/courses", myCourses.Page(DB))
		})
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

	return r
}
