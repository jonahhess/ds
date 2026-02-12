package lesson

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/params"
	"github.com/jonahhess/ds/internal/types"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Start(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			http.Error(w, "invalid user id", http.StatusInternalServerError)
			return
		}
		
		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "course id not found", http.StatusInternalServerError)
			return
		}

		lessonIndex, ok := params.IntFrom(ctx, "lessonIndex")
			if !ok {
			http.Error(w, "lesson id not found", http.StatusInternalServerError)
			return
		}

		// harden against url params
        // only allow lesson_index = uc.currentLesson + 1
        var current_lesson int
        err := DB.QueryRow(`
        SELECT current_lesson 
        FROM user_courses 
        WHERE user_id = ? AND course_id = ?`, userID, courseID).
        Scan(&current_lesson)

        if current_lesson + 1 != lessonIndex {
			http.Error(w, "lesson not found", http.StatusNotFound)
            return
        }
		
		// Get the user's course data
		myData, err := GetMyLessonData(DB, userID, courseID, lessonIndex)
		if err != nil {
			http.Error(w, "course not found", http.StatusNotFound)
			return
		}
		
		  if err := layouts.Base("My Lesson", Lesson(userID, courseID, lessonIndex, *myData)).
		  Render(ctx, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
		  }
	}
}

func GetMyLessonData(DB *sql.DB, userID int, courseID int, lessonIndex int) (*types.Lesson, error) {

    var lesson_id int
    err0 := DB.QueryRow("SELECT id from lessons WHERE course_id = ? AND lesson_index = ?", courseID, lessonIndex).
    Scan(&lesson_id)

    if err0 != nil {
        return nil, err0
    }

        
    var quiz_id int
    err := DB.QueryRow(`SELECT id FROM quizzes WHERE lesson_id = ?`, lesson_id).
    Scan(&quiz_id)

    if err != nil {
        return nil, err
    }

    rows, err := DB.Query(`
        SELECT
            q.id AS question_id,
            a.id AS answer_id,
            a.text AS answer_text
        FROM questions q
        LEFT JOIN answers a ON a.question_id = q.id
        WHERE q.quiz_id = ?
        ORDER BY q.id, a.id;`,quiz_id)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var answers []types.Answer
    for rows.Next() {
        var answer types.Answer
        if err := rows.Scan(&answer.QuestionID, &answer.ID, &answer.Text); err != nil {
            return nil, err
        }
        answers = append(answers, answer)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    questionRows, err3 := DB.Query(`
        SELECT id, text FROM questions 
        WHERE quiz_id = ?`, quiz_id)
    if err3 != nil {
		return nil, err
    }
    defer questionRows.Close()
    
    var questions []types.Question
    for questionRows.Next() {
		var question types.Question
        if err := questionRows.Scan(&question.ID, &question.Text); err != nil {
			return nil, err
        }
        question.Answers = FilterByQuestionID(answers, question.ID)
        questions = append(questions, question)
    }
    if err = questionRows.Err(); err != nil {
        return nil, err
    }

    var quiz types.Quiz
    quiz.Questions = questions
    
    var data types.Lesson
    err2 := DB.QueryRow(`
    SELECT title, text FROM lessons 
    WHERE course_id = ? AND lesson_index = ?`, courseID, lessonIndex).Scan(
        &data.Title,
        &data.Text,
    )
    
    if err2 != nil {
        return nil, err2
    }
    
    data.Quiz = quiz

    return &data, nil
}

func FilterByQuestionID(answers []types.Answer, questionID int) []types.Answer {
    var result []types.Answer
    for _, a := range answers {
        if a.QuestionID == questionID {
            result = append(result, a)
        }
    }
    return result
}

// func Submit()
