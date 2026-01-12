package notFound

import (
	"myapp/layouts"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("NotFound", NotFound()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
