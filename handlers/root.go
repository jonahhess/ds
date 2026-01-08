package handlers

import (
	"context"
	"myapp/auth"
	"myapp/templates"
	"myapp/templates/layouts"
	"net/http"

	"github.com/a-h/templ"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
    if err := auth.SetTokenCookie(w); err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // Render templ component using generated templ components
    comp := layouts.Base("My App", templates.Root())
    if err := comp.Render(templ.InitializeContext(context.TODO()), w); err != nil {
        http.Error(w, "Template execution error", http.StatusInternalServerError)
        return
    }
}
