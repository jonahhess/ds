package catalog

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/types"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
	courses, _ := GetAllCourseTitles(DB)

	  if err := layouts.
	 Base("Courses", Catalog(courses)).
	 Render(r.Context(), w); err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllCourseTitles(DB *sql.DB) ([]types.Item, error){
	rows, err := DB.Query("SELECT id, title FROM courses WHERE version > 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []types.Item
	for rows.Next() {
		var item types.Item
		if err := rows.Scan(&item.ID, &item.Text); err != nil {
			return items, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return items, err
	}

	return items, nil
}
