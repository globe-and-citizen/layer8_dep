package handlers

import (
	"globe-and-citizen/layer8/proxy/internals/usecases"
	"html/template"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	usecase := r.Context().Value("usecase").(*usecases.UseCase)

	switch r.Method {
	case "GET":
		next := r.URL.Query().Get("next")
		token, err := r.Cookie("token")
		if err != nil {
			if next == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
			}
			return
		}
		user, err := usecase.GetUserByToken(token.Value)
		if err != nil || user == nil {
			if next == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
			}
			return
		}

		// load the welcome page
		t, err := template.ParseFiles("assets/templates/welcome.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]interface{}{
			"User":    user,
			"HasNext": next != "",
			"Next":    next,
		})
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
