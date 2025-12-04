package handlers

import (
	"html/template"
	"myapp/auth"
	"net/http"
	"path/filepath"
)

func PartialHandler(w http.ResponseWriter, r *http.Request) {
    claims, err := auth.VerifyToken(r)
    if err != nil {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    partial := r.URL.Path[len("/load-partial/"):]
    allowed := false
    for _, p := range claims.ValidPartials {
        if p == partial {
            allowed = true
            break
        }
    }
    if !allowed {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Rotate token
    if err := auth.SetTokenCookie(w); err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    tmplPath := filepath.Join("templates", "partials", partial+".html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, "Not Found", http.StatusNotFound)
        return
    }
    tmpl.Execute(w, nil)
}
