package home

import (
	"myapp/layouts"
	"myapp/types"
	"net/http"
)

func HandleomeHandler(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Home", Home(types.User)).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
