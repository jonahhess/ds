package creatorcourse

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func LessonNewPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))
		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err := layouts.Base("Create New Lesson", NewLesson(courseID, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func LessonCreate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))

		var createdBy int
		err := db.QueryRow("SELECT created_by FROM courses WHERE id = ?", courseID).Scan(&createdBy)
		if err != nil {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}
		if createdBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		title := r.FormValue("title")
		text := r.FormValue("text")

		if title == "" || text == "" {
			http.Error(w, "Title and text are required", http.StatusBadRequest)
			return
		}

		var maxIndex int
		db.QueryRow("SELECT COALESCE(MAX(lesson_index), -1) FROM lessons WHERE course_id = ?", courseID).Scan(&maxIndex)

		_, err = db.Exec(
			"INSERT INTO lessons (course_id, lesson_index, title, text, created_by) VALUES (?, ?, ?, ?, ?)",
			courseID, maxIndex+1, title, text, userID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID), http.StatusSeeOther)
	}
}

func LessonEditPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))
		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))

		var title, text string
		err := db.QueryRow(
			"SELECT title, text FROM lessons WHERE course_id = ? AND lesson_index = ?",
			courseID, lessonIndex,
		).Scan(&title, &text)

		if err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		lesson := LessonForm{CourseID: courseID, LessonIndex: lessonIndex, Title: title, Text: text}
		err = layouts.Base("Edit Lesson", EditLesson(lesson, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func LessonUpdate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))
		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))

		var createdBy int
		db.QueryRow("SELECT created_by FROM courses WHERE id = ?", courseID).Scan(&createdBy)
		if createdBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		title := r.FormValue("title")
		text := r.FormValue("text")

		db.Exec(
			"UPDATE lessons SET title = ?, text = ? WHERE course_id = ? AND lesson_index = ?",
			title, text, courseID, lessonIndex,
		)

		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID), http.StatusSeeOther)
	}
}

func LessonDelete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))

		var createdBy int
		err := db.QueryRow("SELECT created_by FROM courses WHERE id = ?", courseID).Scan(&createdBy)
		if err != nil {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}
		if createdBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))
		_, err = db.Exec("DELETE FROM lessons WHERE course_id = ? AND lesson_index = ?", courseID, lessonIndex)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID), http.StatusSeeOther)
	}
}
