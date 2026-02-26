package question

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func QuestionNewPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))
		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err := layouts.Base("Add Question", NewQuestion(courseID, lessonIndex, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func QuestionCreate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))
		lessonIndex, _ := strconv.Atoi(chi.URLParam(r,"lessonIndex"))

		var createdBy int
		err := db.QueryRow("SELECT created_by FROM courses WHERE id = ?", courseID).Scan(&createdBy)
		if err != nil || createdBy != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		questionText := r.FormValue("question_text")

		if questionText == "" {
			http.Error(w, "Question text is required", http.StatusBadRequest)
			return
		}

		var lessonID int
		err = db.QueryRow("SELECT id from lessons WHERE course_id = ? AND lesson_index = ?", courseID, lessonIndex).Scan(&lessonID)
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
		// Insert question
		result, err := db.Exec(
			"INSERT INTO questions (quiz_id, text, created_by) VALUES (?, ?, ?)",
			quizID, questionText, userID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		questionID, _ := result.LastInsertId()

		// Insert answers
		correctAnswerIndex := r.FormValue("correct_answer")
		if correctAnswerIndex == "" {
			http.Error(w, "Please select a correct answer", http.StatusBadRequest)
			return
		}

		answerCount := 0
		for i := 1; i <= 4; i++ {
			answerText := r.FormValue("answer_" + strconv.Itoa(i))
			if answerText == "" {
				continue
			}
			answerCount++

			result, err := db.Exec(
				"INSERT INTO answers (question_id, text) VALUES (?, ?)",
				questionID, answerText,
			)
			if err != nil {
				http.Error(w, "Failed to insert answer", http.StatusInternalServerError)
				return
			}

			answerID, err := result.LastInsertId()
			if err != nil {
				http.Error(w, "Failed to get answer ID", http.StatusInternalServerError)
				return
			}

			// Mark as correct if matches selected index
			if strconv.Itoa(i) == correctAnswerIndex {
				_, err := db.Exec(
					"INSERT INTO correct_answers (question_id, answer_id) VALUES (?, ?)",
					questionID, answerID,
				)
				if err != nil {
					http.Error(w, "Failed to mark correct answer", http.StatusInternalServerError)
					return
				}
			}
		}

		if answerCount < 2 {
			http.Error(w, "Please provide at least 2 answers", http.StatusBadRequest)
			return
		}

		lessonIndex, _ = strconv.Atoi(chi.URLParam(r, "lessonIndex"))
		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID)+"/lessons/"+strconv.Itoa(lessonIndex)+"/quiz", http.StatusSeeOther)
	}
}

func QuestionEditPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		questionID, _ := strconv.Atoi(chi.URLParam(r, "questionID"))
		courseID, _ := strconv.Atoi(chi.URLParam(r, "courseID"))
		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))

		var questionText string
		err := db.QueryRow("SELECT text FROM questions WHERE id = ?", questionID).Scan(&questionText)
		if err != nil {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err = layouts.Base("Edit Question", EditQuestion(courseID, lessonIndex, questionID, questionText, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func QuestionUpdate(db *sql.DB) http.HandlerFunc {
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

		questionID, _ := strconv.Atoi(chi.URLParam(r, "questionID"))
		questionText := r.FormValue("question_text")

		if questionText == "" {
			http.Error(w, "Question text is required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(
			"UPDATE questions SET text = ? WHERE id = ?",
			questionText, questionID,
		)
		if err != nil {
			http.Error(w, "Failed to update question", http.StatusInternalServerError)
			return
		}

		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))
		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID)+"/lessons/"+strconv.Itoa(lessonIndex)+"/quiz", http.StatusSeeOther)
	}
}

func QuestionDelete(db *sql.DB) http.HandlerFunc {
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

		questionID, _ := strconv.Atoi(chi.URLParam(r, "questionID"))
		_, err = db.Exec("DELETE FROM questions WHERE id = ?", questionID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		lessonIndex, _ := strconv.Atoi(chi.URLParam(r, "lessonIndex"))
		http.Redirect(w, r, "/creator/courses/"+strconv.Itoa(courseID)+"/lessons/"+strconv.Itoa(lessonIndex)+"/quiz", http.StatusSeeOther)
	}
}
