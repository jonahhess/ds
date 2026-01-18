package study

import (
	"myapp/layouts"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Study", Study()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
