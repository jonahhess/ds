package myCourses

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/layouts"
	"myapp/utils"
	"net/http"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {

	 err := layouts.
	 Base("MyCourses", MyCourses(DB)).
	 Render(r.Context(), w)
	 
	 if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllMyCourses(ctx context.Context, DB *sql.DB) *sql.Rows{
	userID := utils.GetUserID(ctx)
	if userID < 1 {
		return nil
	}
	
	query, err := DB.Query("select title from user_courses where user_id = ?", userID)
	if err != nil {
		return nil
	}
	fmt.Println(query)
	return query
}
