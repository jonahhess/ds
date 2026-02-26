package myLesson

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
		
		myData, err := GetMyLessonData(DB, courseID, lessonIndex)
		if err != nil {
			http.Error(w, "lesson not found", http.StatusNotFound)
			return
		}
		
		  if err := layouts.Base("My Lesson", MyLesson(courseID, lessonIndex, *myData)).
		  Render(ctx, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
		  }
	}
}

func GetMyLessonData(DB *sql.DB, courseID int, lessonIndex int) (*types.Lesson, error) {

    var lesson_id int
    err := DB.QueryRow("SELECT id from lessons WHERE course_id = ? AND lesson_index = ?", courseID, lessonIndex).
    Scan(&lesson_id)

    if err != nil {
        return nil, err
    }
        
    var quiz_id int
    err = DB.QueryRow(`SELECT id FROM quizzes WHERE lesson_id = ?`, lesson_id).
    Scan(&quiz_id)

    if err != nil {
		if err == sql.ErrNoRows {
			quiz_id = 0
		} else {
			return nil, err
		}
    }
    
    var data types.Lesson
    err = DB.QueryRow(`
    SELECT title, text FROM lessons 
    WHERE course_id = ? AND lesson_index = ?`, courseID, lessonIndex).Scan(
        &data.Title,
        &data.Text,
    )
    
    if err != nil {
        return nil, err
    }
    
    data.QuizID = quiz_id
    return &data, nil
}
