package new

import (
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err := layouts.Base("Create New Course", NewCourse(csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}