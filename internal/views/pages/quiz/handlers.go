package quiz

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/params"
	"github.com/jonahhess/ds/internal/types"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return
	}

	quizID, ok := params.IntFrom(ctx, "quizID")
	if !ok {
		return
	}

    isValid := validQuiz(DB, userID, quizID)
    if !isValid {
        http.Error(w, "invalid quiz", http.StatusInternalServerError)
        return
    }

	myData, err := GetMyQuiz(DB, userID, quizID)
	if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	 if err := layouts.
	 Base("MyQuiz", MyQuiz(userID, quizID, myData)).
	 Render(r.Context(), w);  err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetMyQuiz(DB *sql.DB, userID int, quizID int) (*types.Quiz, error) {
	rows, err := DB.Query(`
        SELECT
            q.id,
            q.text,
            a.id,
            a.text
        FROM questions q
        LEFT JOIN answers a ON a.question_id = q.id
        WHERE q.quiz_id = ?
        ORDER BY q.id, a.id;
    `, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quiz types.Quiz
    quiz.Questions = make([]types.Question, 0,6)
	lastQuestionID := -1

	for rows.Next() {
		var qID int
		var qText string
		var aID sql.NullInt64
		var aText sql.NullString

		if err := rows.Scan(&qID, &qText, &aID, &aText); err != nil {
			return nil, err
		}

		if qID != lastQuestionID {
			quiz.Questions = append(quiz.Questions, types.Question{
				ID:      qID,
				Text:    qText,
				Answers: []types.Answer{},
			})
			lastQuestionID = qID
		}

		if aID.Valid {
			// Get pointer to the last inserted question
			currentQuestion := &quiz.Questions[len(quiz.Questions)-1]

			currentQuestion.Answers = append(currentQuestion.Answers, types.Answer{
				ID:         int(aID.Int64),
				QuestionID: qID,
				Text:       aText.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &quiz, nil
}

func validQuiz(DB *sql.DB, userID int, quizID int) bool {
    var current_quiz bool
    if err := DB.QueryRow(`
        SELECT 
            1 
        FROM quizzes q
        LEFT JOIN lessons l ON l.id = lesson_id
        LEFT JOIN user_courses uc ON uc.course_id = l.course_id
        WHERE uc.user_id = ? AND q.id = ? AND uc.current_lesson+1 = l.lesson_index`, userID, quizID).
        Scan(&current_quiz); err != nil {
            return false
    }

    if current_quiz == false {
        return false
    }
    return true
}

// func Submit()