package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jonahhess/ds/internal/auth"
	about "github.com/jonahhess/ds/internal/views/pages/about"
	"github.com/jonahhess/ds/internal/views/pages/course"
	"github.com/jonahhess/ds/internal/views/pages/courses"
	home "github.com/jonahhess/ds/internal/views/pages/home"
	"github.com/jonahhess/ds/internal/views/pages/login"
	"github.com/jonahhess/ds/internal/views/pages/myCourse"
	"github.com/jonahhess/ds/internal/views/pages/myCourses"
	"github.com/jonahhess/ds/internal/views/pages/notFound"
	"github.com/jonahhess/ds/internal/views/pages/review"
	"github.com/jonahhess/ds/internal/views/pages/signup"
	"github.com/jonahhess/ds/internal/views/pages/study"

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
		r.Route("/courses", func(r chi.Router) {
			r.Post("/{courseID}/enroll", course.Enroll(DB))
			r.Get("/{courseID}", course.Page(DB))
			r.Get("/",courses.Page(DB))
		})
		r.Get("/study", study.Page)
		r.Route("/review", func(r chi.Router) {
			r.Get("/", review.Page(DB))
			r.Get("/next", review.GetNextCard(DB))
			r.Post("/submit", review.SubmitAnswer(DB))
		})

		r.Route("/users/{userID}", func(r chi.Router) {
			r.Get("/courses/{courseID}", myCourse.Page(DB))
			r.Get("/courses", myCourses.Page(DB))
		})
	})

	r.Handle("/static/*",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./cmd/myapp/static")),
		),
	)

	r.Handle("/static/pages/*",
		http.StripPrefix("/static/pages/",
			http.FileServer(http.Dir("./internal/views/pages")),
		),
	)

	r.Handle("/static/components/*",
		http.StripPrefix("/static/components/",
			http.FileServer(http.Dir("./internal/views/components")),
		),
	)

	return r
}
