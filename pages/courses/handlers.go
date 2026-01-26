package courses

import (
	"database/sql"
	"myapp/layouts"
	"myapp/types"
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

func GetAllCourseTitles(DB *sql.DB) ([]types.Item, error){
	rows, err := DB.Query("SELECT id, title FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var titles []types.Item
	for rows.Next() {
		var item types.Item
		if err := rows.Scan(&item.ID, &item.Text); err != nil {
			return titles, err
		}
		titles = append(titles, item)
	}
	if err = rows.Err(); err != nil {
		return titles, err
	}

	return titles, nil
}
