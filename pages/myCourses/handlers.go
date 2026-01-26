package myCourses

import (
	"database/sql"
	"myapp/auth"
	"myapp/layouts"
	"net/http"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return
	}
	
	myTitles, err := GetAllMyCourseTitles(DB, userID); 
	if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	 if err := layouts.
	 Base("MyCourses", MyCourses(myTitles)).
	 Render(r.Context(), w);  err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllMyCourseTitles(DB *sql.DB, userID int) ([]string, error){
	rows, err := DB.Query("SELECT c.title FROM user_courses uc INNER JOIN courses c ON c.id = uc.course_id WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return titles, err
		}
		titles = append(titles, title)
	}
	if err = rows.Err(); err != nil {
		return titles, err
	}

	return titles, nil
}

