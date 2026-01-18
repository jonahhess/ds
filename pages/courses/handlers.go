package courses

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/layouts"
	"net/http"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {

	 err := layouts.
	 Base("Courses", Courses(DB)).
	 Render(r.Context(), w)
	 
	 if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllCourses(ctx context.Context, DB *sql.DB) sql.Result{
	query := "select title from courses"
	result, err := DB.Exec(query)
	if err != nil {
		return nil
	}
	fmt.Println(result)
	return result
}
