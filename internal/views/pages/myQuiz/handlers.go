package myQuiz

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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
		
		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			return
		}
		
		lessonIndex, ok := params.IntFrom(ctx, "lessonIndex")
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
		 return
	}

	csrfToken := auth.CSRFToken(r)
	 if err := layouts.
	 Base("MyQuiz", MyQuiz(userID, courseID, lessonIndex, quizID, myData, []int{}, csrfToken)).
	 Render(ctx, w);  err != nil {
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

func Submit(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			return
		}

		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			return
		}

		lessonIndex, ok := params.IntFrom(ctx, "lessonIndex")
		if !ok {
			return
		}

		quizID, ok := params.IntFrom(ctx, "quizID")
		if !ok {
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		// main loop
		score := 0
		total := 0
		var mistakes []int

		for key, values := range r.PostForm {
			intKey, err :=  strconv.Atoi(key)
			if err != nil {
				continue
			}
			
			total++
			
			for _, value := range values {
				var exists int
				if err := DB.QueryRow(`
				SELECT 1 
				FROM correct_answers 
				WHERE question_id = ? AND answer_id = ?`, key, value).
				Scan(&exists); err == nil {
					score += 1
				} else {
					mistakes = append(mistakes, intKey)
				}
			}
		}
		// end of main loop

		if total == 0 {
			http.Error(w, "no questions in quiz", http.StatusInternalServerError)
			return
		}

		if float32(score) / float32(total) >= 0.8 {
			
			// Create review cards for each question in the quiz
			questionIDs := make([]int, 0, len(r.PostForm))
			for key := range r.PostForm {
				questionID, err := strconv.Atoi(key)
				if err != nil {
					continue // Skip invalid keys
				}
				questionIDs = append(questionIDs, questionID)
			}
			
			if len(questionIDs) == 0 {
				http.Error(w, "no questions in quiz", http.StatusInternalServerError)
				return
			}
		
			// Use a transaction for atomicity
			tx, err := DB.Begin()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer tx.Rollback()
		
			stmt, err := tx.Prepare(`
				INSERT OR IGNORE INTO reviewcards 
				(user_id, question_id, review_at, interval, easiness, repetitions, successes, reviews)
				VALUES (?, ?, CURRENT_TIMESTAMP, 1, 2.5, 0, 0, 0)
			`)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer stmt.Close()
		
			for _, questionID := range questionIDs {
				_, err := stmt.Exec(userID, questionID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			
			if _, err := tx.Exec(`
			UPDATE user_courses 
			SET current_lesson = current_lesson + 1 
			WHERE user_id = ? AND course_id = ?`, userID, courseID); err != nil {
				http.Error(w, "cannot update current lesson", http.StatusInternalServerError)
				return
			}
			
			if err := tx.Commit(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/courses/%d", courseID), http.StatusSeeOther)
			return
		} else {		
			isValid := validQuiz(DB, userID, quizID)
			if !isValid {
				http.Error(w, "invalid quiz", http.StatusInternalServerError)
				return
			}

			myData, err := GetMyQuiz(DB, userID, quizID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
					
			csrfToken := auth.CSRFToken(r)
			if err := layouts.
				Base("MyQuiz", MyQuiz(userID, courseID, lessonIndex, quizID, myData, mistakes, csrfToken)).
				Render(ctx, w);  err != nil {
		 			http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
	}
}