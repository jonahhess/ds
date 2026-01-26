package template

import (
	"database/sql"
	"myapp/auth"
	"myapp/layouts"
	"myapp/types"
	"net/http"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return
	}

	myData, err := GetAllTemplate(DB, userID)
	if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	 if err := layouts.
	 Base("MyTemplate", MyTemplate(userID, myData)).
	 Render(r.Context(), w);  err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllTemplate(DB *sql.DB, userID int) ([]types.Item, error){

	rows, err := DB.Query("SELECT * FROM template")
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
	}
	if err = rows.Err(); err != nil {
		return items, err
	}
	
	return items, nil
}
