package navbar

import (
	"myapp/utils"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, ok := utils.SessionFromContext(r.Context())
	if !ok {
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
