package courses

import (
	"database/sql"
	"myapp/layouts"
	"net/http"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
	courses, _ := GetAllCourseTitles(DB)

	  if err := layouts.
	 Base("Courses", Courses(courses)).
	 Render(r.Context(), w); err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllCourseTitles(DB *sql.DB) ([]string, error){
	rows, err := DB.Query("SELECT title FROM courses")
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
