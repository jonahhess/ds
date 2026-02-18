package creatorcourse

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/params"
	"github.com/jonahhess/ds/internal/views/layouts"
)

// NewPage displays the form to create a new course
func NewPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err := layouts.Base("Create New Course", NewCourse(csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

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

// EditPage displays the form to edit a course
func EditPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseIDStr := chi.URLParam(r, "courseID")
		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		var title, description sql.NullString
		err = db.QueryRow(
			"SELECT title, description FROM courses WHERE id = ?",
			courseID,
		).Scan(&title, &description)

		if err == sql.ErrNoRows {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		course := CourseForm{
			ID:          courseID,
			Title:       title.String,
			Description: description,
		}

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err = layouts.Base("Edit Course", EditCourse(course, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

// DetailPage displays course details for a creator
func DetailPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		var title, description sql.NullString
		err := db.QueryRow(
			"SELECT title, description FROM courses WHERE id = ?",
			courseID,
		).Scan(&title, &description)

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
			"SELECT lesson_index, title, text, quizzes.id FROM lessons join quizzes on lessons.id = quizzes.lesson_id WHERE course_id = ? ORDER BY lesson_index",
			courseID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var lessons []LessonDisplay
		for rows.Next() {
			var lesson LessonDisplay
			err := rows.Scan(&lesson.Index, &lesson.Title, &lesson.Text, &lesson.QuizID)
			if err != nil {
				continue
			}
			lessons = append(lessons, lesson)
		}

		course := CourseForm{
			ID:          courseID,
			Title:       title.String,
			Description: description,
		}

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err = layouts.Base(course.Title, CourseDetail(course, lessons, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
