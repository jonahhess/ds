package template

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
	 Base("Template", Template(DB)).
	 Render(r.Context(), w)
	 
	 if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetAllTemplate(ctx context.Context, DB *sql.DB) sql.Result{
	query := "select title from template"
	result, err := DB.Exec(query)
	if err != nil {
		return nil
	}
	fmt.Println(result)
	return result
}
