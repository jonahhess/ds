package notFound

import (
	"net/http"

	"github.com/jonahhess/ds/internal/errors"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(w http.ResponseWriter, r *http.Request) {
	err := layouts.
		Base("NotFound", NotFound()).
		Render(r.Context(), w)

	if err != nil {
		errors.HandleInternalError(w, r, err)
		return
	}
}
