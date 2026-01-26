package myCourses

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
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
	
	myTitles, err := GetAllMyCourseTitles(DB, userID); 
	if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	 if err := layouts.
	 Base("MyCourses", MyCourses(userID, myTitles)).
	 Render(r.Context(), w);  err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllMyCourseTitles(DB *sql.DB, userID int) ([]types.Item, error){
	rows, err := DB.Query("SELECT c.id, c.title FROM user_courses uc INNER JOIN courses c ON c.id = uc.course_id WHERE user_id = ?", userID)
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

