package template

import (
	"myapp/layouts"
	"myapp/types"
	"net/http"
)

func TemplateHandler(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Template", Template(types.User{})).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
