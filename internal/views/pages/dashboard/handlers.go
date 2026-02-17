package dashboard

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get enrolled courses count
		var enrolledCourses int
		db.QueryRow("SELECT COUNT(*) FROM user_courses WHERE user_id = ?", userID).Scan(&enrolledCourses)

		// Get lessons completed count
		var lessonsCompleted int
		db.QueryRow(
			"SELECT COUNT(*) FROM review_cards WHERE user_id = ? AND is_completed = 1",
			userID,
		).Scan(&lessonsCompleted)

		// Get reviews due count
		var reviewsDue int
		db.QueryRow(
			"SELECT COUNT(*) FROM review_cards WHERE user_id = ? AND is_completed = 0",
			userID,
		).Scan(&reviewsDue)

		// Get reviews completed count
		var reviewsCompleted int
		db.QueryRow(
			"SELECT COUNT(*) FROM review_cards WHERE user_id = ? AND is_completed = 1",
			userID,
		).Scan(&reviewsCompleted)

		stats := DashboardStats{
			UserID:           userID,
			EnrolledCourses:  enrolledCourses,
			LessonsCompleted: lessonsCompleted,
			ReviewsDue:       reviewsDue,
			ReviewsCompleted: reviewsCompleted,
		}

		w.Header().Set("Content-Type", "text/html")
		err := layouts.Base("Your Dashboard", UserDashboard(stats)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func CourseProgressPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))

		// Get course title
		var title string
		err := db.QueryRow("SELECT title FROM courses WHERE id = ?", courseID).Scan(&title)
		if err != nil {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}

		// Check user enrollment
		var enrolled int
		db.QueryRow("SELECT COUNT(*) FROM user_courses WHERE user_id = ? AND course_id = ?", userID, courseID).Scan(&enrolled)
		if enrolled == 0 {
			http.Error(w, "Not enrolled in this course", http.StatusForbidden)
			return
		}

		// Get user's current lesson
		var currentLesson sql.NullInt64
		db.QueryRow(
			"SELECT current_lesson FROM user_courses WHERE user_id = ? AND course_id = ?",
			userID, courseID,
		).Scan(&currentLesson)

		// Get total lesson count
		var totalLessons int
		db.QueryRow("SELECT COUNT(*) FROM lessons WHERE course_id = ?", courseID).Scan(&totalLessons)

		// Calculate progress percentage
		progressPercent := 0
		if totalLessons > 0 && currentLesson.Valid {
			progressPercent = int((currentLesson.Int64 / int64(totalLessons)) * 100)
		}

		progress := CourseProgress{
			CourseID:        courseID,
			Title:           title,
			CurrentLesson:   int(currentLesson.Int64),
			TotalLessons:    totalLessons,
			ProgressPercent: progressPercent,
		}

		w.Header().Set("Content-Type", "text/html")
		err = layouts.Base(title, CourseProgressDetail(progress)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
