package course

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/params"
	"github.com/jonahhess/ds/internal/types"
	"github.com/jonahhess/ds/internal/views/layouts"
)

// Create handles POST request to create a new course
func Create(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		title := r.FormValue("title")
		description := r.FormValue("description")

		if title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"INSERT INTO courses (title, description, created_by) VALUES (?, ?, ?)",
			title, description, userID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		courseID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/creator/courses/"+strconv.FormatInt(courseID, 10), http.StatusSeeOther)
	}
}

// Update handles PATCH request to update course details
func Update(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseIDStr := chi.URLParam(r, "courseID")
		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		description := r.FormValue("description")

		if title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(
			"UPDATE courses SET title = ?, description = ? WHERE id = ? AND created_by = ?",
			title, description, courseID, userID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID), http.StatusSeeOther)
	}
}

// Delete handles DELETE request to remove a course
func Delete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseIDStr := chi.URLParam(r, "courseID")
		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		// Check ownership
		var createdBy int
		err = db.QueryRow("SELECT created_by FROM courses WHERE id = ?", courseID).Scan(&createdBy)
		if err == sql.ErrNoRows {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if createdBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		_, err = db.Exec("DELETE FROM courses WHERE id = ?", courseID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/creator", http.StatusSeeOther)
	}
}

func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		var title, description sql.NullString
		var version int
		err := db.QueryRow(
			"SELECT title, description, version FROM courses WHERE id = ?",
			courseID,
		).Scan(&title.String, &description.String, &version)

		if err == sql.ErrNoRows {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Fetch lessons for this course
		rows, err := db.Query(
			"SELECT lesson_index, title, text, COALESCE(quizzes.id, 0) FROM lessons LEFT JOIN quizzes on lessons.id = quizzes.lesson_id WHERE course_id = ? ORDER BY lesson_index",
			courseID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var lessons []types.LessonDisplay
		for rows.Next() {
			var lesson types.LessonDisplay
			err := rows.Scan(&lesson.Index, &lesson.Title, &lesson.Text, &lesson.QuizID)
			if err != nil {
				continue
			}
			lessons = append(lessons, lesson)
		}

		course := types.CourseForm{
			ID:          courseID,
			Title:       title.String,
			Version: 	 version,
			Description: description,
		}

		csrfToken := auth.CSRFTokenFromContext(ctx)
		err = layouts.Base(course.Title, Course(course, lessons, csrfToken)).Render(ctx, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Version(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseIDStr := chi.URLParam(r, "courseID")
		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		version, err := strconv.Atoi(r.FormValue("version"))
		if err != nil {
			version = 1
		}

		_, err = db.Exec(
			"UPDATE courses SET version = ? WHERE id = ? AND created_by = ?",
			version, courseID, userID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID), http.StatusSeeOther)
	}
}