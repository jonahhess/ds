package handlers

import (
	"html/template"
	"myapp/auth"
	"net/http"
	"path/filepath"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
    if err := auth.SetTokenCookie(w); err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // Parse base layout + page + JS component
    tmpl, err := template.ParseFiles(
        filepath.Join("templates", "layouts", "base.html"),
        filepath.Join("templates", "root.html"),
    )
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }

    tmpl.ExecuteTemplate(w, "base.html", map[string]string{
        "Title": "My App",
    })
}
