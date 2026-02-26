package myAvailableCourses

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
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	
	myTitles, err := GetAllNotMyCourseTitles(DB, userID); 
	if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	 if err := layouts.
	 Base("My Available Courses", MyAvailableCourses(myTitles)).
	 Render(r.Context(), w);  err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllNotMyCourseTitles(DB *sql.DB, userID int) ([]types.Item, error){
	rows, err := DB.Query(`
    SELECT c.id, c.title
    FROM courses c
    WHERE NOT EXISTS (
        SELECT 1
        FROM user_courses uc
        WHERE uc.course_id = c.id
        AND uc.user_id = ?
    ) AND c.version > 0
`, userID)
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