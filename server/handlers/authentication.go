package handlers

import (
	"globe-and-citizen/layer8/server/internals/usecases"
	"html/template"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	usecase := r.Context().Value("usecase").(*usecases.UseCase)

	switch r.Method {
	case "GET":
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}
		// check if the user is already logged in
		token, err := r.Cookie("token")
		if token != nil && err == nil {
			user, err := usecase.GetUserByToken(token.Value)
			if err == nil && user != nil {
				http.Redirect(w, r, next, http.StatusSeeOther)
				return
			}
		}

		// load the login page
		t, err := template.ParseFiles("assets-v1/templates/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]interface{}{
			"HasNext": next != "",
			"Next":    next,
		})
		return
	case "POST":
		next := r.URL.Query().Get("next")
		username := r.FormValue("username")
		password := r.FormValue("password")
		// login the user
		rUser, err := usecase.LoginUser(username, password)
		if err != nil {
			t, errT := template.ParseFiles("assets-v1/templates/login.html")
			if errT != nil {
				http.Error(w, errT.Error(), http.StatusInternalServerError)
				return
			}
			t.Execute(w, map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   err.Error(),
			})
			return
		}
		// set the token cookie
		token, ok := rUser["token"].(string)
		if !ok {
			t, err := template.ParseFiles("assets-v1/templates/login.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			t.Execute(w, map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   "could not get token",
			})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		})
		// redirect to next page - here the user already knows their pseudo profile
		// when they registered
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
func Register(w http.ResponseWriter, r *http.Request) {
	usecase := r.Context().Value("usecase").(*usecases.UseCase)

	switch r.Method {
	case "GET":
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}
		// check if the user is already logged in
		token, err := r.Cookie("token")
		if token != nil && err == nil {
			user, err := usecase.GetUserByToken(token.Value)
			if err == nil && user != nil {
				http.Redirect(w, r, next, http.StatusSeeOther)
				return
			}
		}

		// load the login page
		t, err := template.ParseFiles("assets-v1/templates/registerClient.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]interface{}{
			"HasNext": next != "",
			"Next":    next,
		})
		return
	case "POST":
		next := r.URL.Query().Get("next")
		username := r.FormValue("username")
		password := r.FormValue("password")
		// login the user
		rUser, err := usecase.LoginUser(username, password)
		if err != nil {
			t, errT := template.ParseFiles("assets-v1/templates/registerClient.html")
			if errT != nil {
				http.Error(w, errT.Error(), http.StatusInternalServerError)
				return
			}
			t.Execute(w, map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   err.Error(),
			})
			return
		}
		// set the token cookie
		token, ok := rUser["token"].(string)
		if !ok {
			t, err := template.ParseFiles("assets-v1/templates/registerClient.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			t.Execute(w, map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   "could not get token",
			})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
			Path:  "/",
		})
		// redirect to next page - here the user already knows their pseudo profile
		// when they registered
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
