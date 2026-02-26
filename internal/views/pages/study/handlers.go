package study

import (
	"net/http"

	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(w http.ResponseWriter, r *http.Request) {
	if err := layouts.
		Base("Study", Study()).
		Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
