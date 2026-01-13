package router

import (
	"net/http"

	middlewares "myapp/middlewares"
	about "myapp/pages/about"
	home "myapp/pages/home"
	"myapp/pages/login"
	"myapp/pages/logout"
	"myapp/pages/notFound"
	"myapp/pages/signup"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupRoutes(sessionStore *sessions.CookieStore) *chi.Mux {
	r := chi.NewRouter()

	//r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	//r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)

	//r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middlewares.SessionMiddleware(sessionStore))

	r.Route("/", func(r chi.Router) {
		r.Get("/", home.Page)
		r.Get("/about", about.Page)
		r.Get("/login", login.Page)
		r.Post("/login", login.LoginHandler)
		r.Get("/signup", signup.Page)
		r.Post("/signup", signup.SignupHandler)
		r.Get("/logout", logout.Page)
		r.Post("/logout", logout.LogoutHandler)
	})

	//r.Group(func(r chi.Router) {
	//r.Use(middlewares.AuthMiddleware)
	//r.Get("/user", user.userHandler)
	//})

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
