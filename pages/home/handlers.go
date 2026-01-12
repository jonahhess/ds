package home

import (
	"myapp/layouts"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Home", Home()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
