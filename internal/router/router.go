package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/params"
	myCourse "github.com/jonahhess/ds/internal/views/pages/courses/course"
	"github.com/jonahhess/ds/internal/views/pages/courses/myAvailableCourses"
	"github.com/jonahhess/ds/internal/views/pages/courses/myLesson"
	"github.com/jonahhess/ds/internal/views/pages/courses/myQuiz"
	courses "github.com/jonahhess/ds/internal/views/pages/courses/root"
	creatorCourse "github.com/jonahhess/ds/internal/views/pages/creator/course"
	creatorCourseEdit "github.com/jonahhess/ds/internal/views/pages/creator/course/edit"
	creatorCourseNew "github.com/jonahhess/ds/internal/views/pages/creator/course/new"
	creator "github.com/jonahhess/ds/internal/views/pages/creator/root"

	creatorLesson "github.com/jonahhess/ds/internal/views/pages/creator/lesson"
	creatorQuestion "github.com/jonahhess/ds/internal/views/pages/creator/question"
	creatorQuiz "github.com/jonahhess/ds/internal/views/pages/creator/quiz"

	"github.com/jonahhess/ds/internal/views/pages/home/about"
	"github.com/jonahhess/ds/internal/views/pages/home/catalog"
	catalogCourse "github.com/jonahhess/ds/internal/views/pages/home/catalog/course"
	"github.com/jonahhess/ds/internal/views/pages/home/login"
	"github.com/jonahhess/ds/internal/views/pages/home/notFound"
	"github.com/jonahhess/ds/internal/views/pages/home/profile"
	home "github.com/jonahhess/ds/internal/views/pages/home/root"
	"github.com/jonahhess/ds/internal/views/pages/home/signup"
	"github.com/jonahhess/ds/internal/views/pages/home/study"
	"github.com/jonahhess/ds/internal/views/pages/review"

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
		
		r.Route("/profile", func(r chi.Router) {
			r.Get("/", profile.ViewPage(DB))
			r.Get("/edit", profile.EditPage(DB))
			r.Patch("/", profile.Update(DB))
			r.Get("/password", profile.ChangePasswordPage(DB))
			r.Post("/password", profile.ChangePasswordHandler(DB))
			r.Get("/delete", profile.DeletePage(DB))
			r.Delete("/", profile.Delete(DB))
		})

		r.Route("/creator", func(r chi.Router) {
			r.Get("/", creator.Page(DB))
			
			r.Route("/courses", func(r chi.Router) {
				r.Get("/new", creatorCourseNew.Page())
				r.Post("/", creatorCourse.Create(DB))
				
				r.Route("/{courseID}", func(r chi.Router) {
					r.Use(params.Int("courseID"))
					r.Use(auth.RequireCreatorMiddleware(DB))
					
					r.Get("/", creatorCourse.Page(DB))
					r.Get("/edit", creatorCourseEdit.Page(DB))
					r.Patch("/", creatorCourse.Update(DB))
					r.Delete("/", creatorCourse.Delete(DB))
					r.Post("/version", creatorCourse.Version(DB))
					
					r.Route("/lessons/{lessonIndex}", func(r chi.Router) {
						r.Use(params.Int("lessonIndex"))
						r.Get("/new", creatorLesson.LessonNewPage())
						r.Post("/", creatorLesson.LessonCreate(DB))
						r.Get("/edit", creatorLesson.LessonEditPage(DB))
						r.Patch("/", creatorLesson.LessonUpdate(DB))
						r.Delete("/", creatorLesson.LessonDelete(DB))
						
						r.Route("/quiz", func(r chi.Router) {
							r.Post("/", creatorQuiz.QuizCreate(DB))
							r.Get("/", creatorQuiz.Page(DB))
							r.Delete("/", creatorQuiz.QuizDelete(DB))
							
							r.Route("/questions", func(r chi.Router) {
								r.Get("/new", creatorQuestion.QuestionNewPage(DB))
								r.Post("/", creatorQuestion.QuestionCreate(DB))
								
								r.Route("/{questionID}", func(r chi.Router) {
									r.Use(params.Int("questionID"))
									r.Get("/edit", creatorQuestion.QuestionEditPage(DB))
									r.Patch("/", creatorQuestion.QuestionUpdate(DB))
									r.Delete("/", creatorQuestion.QuestionDelete(DB))
								})
							})
						})
					})
				})
			})
		})

		r.Route("/catalog", func(r chi.Router) {
			r.Route("/courses/{courseID}", func(r chi.Router) {
				r.Use(params.Int("courseID"))
				r.Post("/enroll", catalogCourse.Enroll(DB))
				r.Get("/", catalogCourse.Page(DB))
			})
			r.Get("/",catalog.Page(DB))
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

		r.Route("/courses", func(r chi.Router) {
			r.Get("/add", myAvailableCourses.Page(DB))
			r.Route("/{courseID}", func(r chi.Router) {
				r.Use(params.Int("courseID"))
				r.Get("/", myCourse.Page(DB))
				r.Post("/remove", myCourse.Remove(DB))
				r.Route("/lessons/{lessonIndex}", func(r chi.Router) {
					r.Use(params.Int("lessonIndex"))
					r.Route("/quiz", func(r chi.Router) {
						r.Get("/",myQuiz.Page(DB))
						r.Post("/", myQuiz.Submit(DB))
					})
					r.Get("/", myLesson.Page(DB))
					})	
			})
			r.Get("/", courses.Page(DB))
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
