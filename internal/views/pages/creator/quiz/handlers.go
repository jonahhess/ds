package quiz

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

func QuizCreate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))

		var createdBy int
		err := db.QueryRow("SELECT created_by FROM courses WHERE id = ?", courseID).Scan(&createdBy)
		if err != nil || createdBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))

		var lessonID int
		err = db.QueryRow(
			"SELECT id FROM lessons WHERE course_id = ? AND lesson_index = ?",
			courseID, lessonIndex,
		).Scan(&lessonID)
		if err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}

		_, err = db.Exec(
			"INSERT INTO quizzes (lesson_id) VALUES (?)",
			lessonID,
		)
		if err != nil {
			http.Error(w, "Failed to create quiz: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID), http.StatusSeeOther)
	}
}

func QuizDelete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))

		var createdBy int
		db.QueryRow("SELECT created_by FROM courses WHERE id = ?", courseID).Scan(&createdBy)
		if createdBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		quizID, _ := strconv.Atoi(chi.URLParam(r, "quizID"))
		db.Exec("DELETE FROM quizzes WHERE id = ?", quizID)

		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))
		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID)+"/lessons/"+strconv.Itoa(lessonIndex), http.StatusSeeOther)
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

		lessonIndex, ok := params.IntFrom(ctx, "lessonIndex")
		if !ok {
			http.Error(w, "Invalid lesson index", http.StatusBadRequest)
			return
		}

		var lessonID int
		err := db.QueryRow("SELECT id from lessons WHERE course_id = ? AND lesson_index = ?", courseID, lessonIndex).Scan(&lessonID)
		if err != nil {
			http.Error(w, "Invalid lesson id", http.StatusBadRequest)
			return
		}

		var quizID int
		err = db.QueryRow("SELECT id FROM quizzes WHERE lesson_id = ?", lessonID).Scan(&quizID)
		if err != nil {
			http.Error(w, "Invalid quiz id", http.StatusBadRequest)
			return
		}

		// Fetch quiz questions for this lesson id
		rows, err := db.Query(
			"SELECT questions.id, questions.text FROM questions WHERE quiz_id = ?",
			quizID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var questions []types.Item
		for rows.Next() {
			var question types.Item
			err := rows.Scan(&question.ID, &question.Text)
			if err != nil {
				continue
			}
		questions = append(questions, question)
		}

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err = layouts.Base("Quiz Detail", Quiz(courseID, lessonIndex, questions, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
