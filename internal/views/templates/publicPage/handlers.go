package template

import (
	"myapp/internal/views/layouts"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Template", Template()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
