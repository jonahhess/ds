package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/params"
	about "github.com/jonahhess/ds/internal/views/pages/about"
	"github.com/jonahhess/ds/internal/views/pages/course"
	"github.com/jonahhess/ds/internal/views/pages/courses"
	"github.com/jonahhess/ds/internal/views/pages/creator"
	creatorCourse "github.com/jonahhess/ds/internal/views/pages/creatorCourse"
	home "github.com/jonahhess/ds/internal/views/pages/home"
	"github.com/jonahhess/ds/internal/views/pages/login"
	"github.com/jonahhess/ds/internal/views/pages/myAvailableCourses"
	"github.com/jonahhess/ds/internal/views/pages/myCourse"
	"github.com/jonahhess/ds/internal/views/pages/myCourses"
	"github.com/jonahhess/ds/internal/views/pages/myLesson"
	"github.com/jonahhess/ds/internal/views/pages/myQuiz"
	"github.com/jonahhess/ds/internal/views/pages/notFound"
	"github.com/jonahhess/ds/internal/views/pages/profile"
	"github.com/jonahhess/ds/internal/views/pages/review"
	"github.com/jonahhess/ds/internal/views/pages/signup"
	"github.com/jonahhess/ds/internal/views/pages/study"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

// MethodOverrideMiddleware handles _method form field to override HTTP method
func MethodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err == nil {
				if method := r.FormValue("_method"); method != "" {
					r.Method = method
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes(sessionStore *sessions.CookieStore, DB *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(auth.SessionMiddleware(sessionStore))
	r.Use(auth.CSRFMiddleware())
	r.Use(MethodOverrideMiddleware)
	r.Use(auth.OptionalUserMiddleware)

	r.Group( func(r chi.Router) {
		r.Get("/", home.Page)
		r.Get("/about", about.Page)
		r.Get("/login", login.Page)
		r.Post("/login", auth.LoginHandler(DB))
		r.Get("/signup", signup.Page)
		r.Post("/signup", signup.SignupHandler(DB))
		r.Post("/logout", auth.LogoutHandler)
		r.NotFound(notFound.Page)
	})
	
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuthMiddleware)
		
		// Profile routes
		r.Route("/profile", func(r chi.Router) {
			r.Get("/", profile.ViewPage(DB))
			r.Get("/edit", profile.EditPage(DB))
			r.Patch("/", profile.Update(DB))
			r.Get("/password", profile.ChangePasswordPage(DB))
			r.Post("/password", profile.ChangePasswordHandler(DB))
			r.Get("/delete", profile.DeletePage(DB))
			r.Delete("/", profile.Delete(DB))
		})

		// Creator routes
		r.Route("/creator", func(r chi.Router) {
			r.Get("/", creator.Page(DB))
			
			// Creator courses management
			r.Route("/courses", func(r chi.Router) {
				r.Get("/new", creatorCourse.NewPage())
				r.Post("/", creatorCourse.Create(DB))
				
				r.Route("/{courseID}", func(r chi.Router) {
					r.Use(params.Int("courseID"))
					r.Use(auth.RequireCreatorMiddleware(DB))
					
					r.Get("/", creatorCourse.DetailPage(DB))
					r.Get("/edit", creatorCourse.EditPage(DB))
					r.Patch("/", creatorCourse.Update(DB))
					r.Delete("/", creatorCourse.Delete(DB))
					
					// Creator lesson management
					r.Route("/lessons/{lessonIndex}", func(r chi.Router) {
						r.Use(params.Int("lessonIndex"))
						r.Get("/new", creatorCourse.LessonNewPage())
						r.Post("/", creatorCourse.LessonCreate(DB))
						r.Get("/edit", creatorCourse.LessonEditPage(DB))
						r.Patch("/", creatorCourse.LessonUpdate(DB))
						r.Delete("/", creatorCourse.LessonDelete(DB))
						
						// Creator quiz management
						r.Route("/quizzes", func(r chi.Router) {
							r.Get("/new", creatorCourse.QuizNewPage(DB))
							r.Post("/", creatorCourse.QuizCreate(DB))
							
							r.Route("/{quizID}", func(r chi.Router) {
								r.Use(params.Int("quizID"))
								r.Delete("/", creatorCourse.QuizDelete(DB))
								
								// Creator question/answer management
								r.Route("/questions", func(r chi.Router) {
									r.Get("/new", creatorCourse.QuestionNewPage(DB))
									r.Post("/", creatorCourse.QuestionCreate(DB))
									
									r.Route("/{questionID}", func(r chi.Router) {
										r.Use(params.Int("questionID"))
										r.Get("/edit", creatorCourse.QuestionEditPage(DB))
										r.Patch("/", creatorCourse.QuestionUpdate(DB))
										r.Delete("/", creatorCourse.QuestionDelete(DB))
									})
								})
							})
						})
					})
				})
			})
		})

		r.Route("/courses", func(r chi.Router) {
			r.Route("/{courseID}", func(r chi.Router) {
				r.Use(params.Int("courseID"))
				r.Post("/enroll", course.Enroll(DB))
				r.Get("/", course.Page(DB))
			})
			r.Get("/",courses.Page(DB))
		})
		r.Get("/study", study.Page)
		r.Route("/review", func(r chi.Router) {
			r.Get("/", review.Page(DB))
			r.Get("/next", review.NextCard(DB))
			r.Get("/complete", review.Complete(DB))
			r.Route("/card/{questionID}", func(r chi.Router) {
				r.Use(params.Int("questionID"))
				r.Get("/answer", review.ShowAnswer(DB))
				r.Post("/rate", review.RateCard(DB))
			})
		})

		r.Route("/users/{userID}", func(r chi.Router) {
		r.Use(auth.RequireAuthMiddleware)
			r.Use(auth.RequireMatchingUserID)
			r.Route("/courses", func(r chi.Router) {
				r.Get("/add", myAvailableCourses.Page(DB))
				r.Route("/{courseID}", func(r chi.Router) {
					r.Use(params.Int("courseID"))
					r.Get("/", myCourse.Page(DB))
					r.Post("/remove", myCourse.Remove(DB))
					r.Route("/lessons/{lessonIndex}", func(r chi.Router) {
						r.Use(params.Int("lessonIndex"))
						r.Route("/quizzes/{quizID}", func(r chi.Router) {
							r.Use(params.Int("quizID"))
							r.Get("/",myQuiz.Page(DB))
							r.Post("/", myQuiz.Submit(DB))
						})
						r.Get("/", myLesson.Page(DB))
						})	
				})
				r.Get("/", myCourses.Page(DB))
			})
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
