package creatorcourse

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/params"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func LessonNewPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}
		csrfToken := auth.CSRFTokenFromContext(ctx)
		err := layouts.Base("Create New Lesson", NewLesson(courseID, csrfToken)).Render(ctx, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func LessonCreate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		title := r.Form.Get("title")
		text := r.Form.Get("text")

		if title == "" || text == "" {
			http.Error(w, "Title and text are required", http.StatusBadRequest)
			return
		}

		var maxIndex int
		db.QueryRow("SELECT COALESCE(MAX(lesson_index), 1) FROM lessons WHERE course_id = ?", courseID).Scan(&maxIndex)

		_, err = db.Exec(
			"INSERT INTO lessons (course_id, lesson_index, title, text, created_by) VALUES (?, ?, ?, ?, ?)",
			courseID, maxIndex+1, title, text, userID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/creator/courses/%d/lessons/%d", courseID, maxIndex+1), http.StatusSeeOther)
	}
}

func LessonEditPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}
		lessonIndex, ok := params.IntFrom(ctx, "lessonIndex")
		if !ok {
			http.Error(w, "Invalid lesson index", http.StatusBadRequest)
			return
		}

		var title, text sql.NullString
		err := db.QueryRow(
			"SELECT title, text FROM lessons WHERE course_id = ? AND lesson_index = ?",
			courseID, lessonIndex,
		).Scan(&title.String, &text.String)

		if err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}

		csrfToken := auth.CSRFTokenFromContext(ctx)
		lesson := LessonForm{CourseID: courseID, LessonIndex: lessonIndex, Title: title.String, Text: text.String}
		err = layouts.Base("Edit Lesson", EditLesson(lesson, csrfToken)).Render(ctx, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func LessonUpdate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}
		lessonIndex, ok := params.IntFrom(ctx, "lessonIndex")
		if !ok {
			http.Error(w, "Invalid lesson index", http.StatusBadRequest)
			return
		}

		var createdBy int
		db.QueryRow("SELECT created_by FROM courses WHERE id = ?", courseID).Scan(&createdBy)
		if createdBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		title := r.Form.Get("title")
		text := r.Form.Get("text")
		if title == "" || text == "" {
			http.Error(w, "Title and text are required", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(
			"UPDATE lessons SET title = ?, text = ? WHERE course_id = ? AND lesson_index = ?",
			title, text, courseID, lessonIndex,
		)

		http.Redirect(w, r, fmt.Sprintf("/creator/courses/%d/lessons/%d", courseID, lessonIndex), http.StatusSeeOther)
	}
}

func LessonDelete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

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

		lessonIndex, ok := params.IntFrom(ctx, "lessonIndex")
		if !ok {
			http.Error(w, "Invalid lesson index", http.StatusBadRequest)
			return
		}
		_, err = db.Exec("DELETE FROM lessons WHERE course_id = ? AND lesson_index = ?", courseID, lessonIndex)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/creator/courses/%d/lessons/%d", courseID, lessonIndex), http.StatusSeeOther)
	}
}
