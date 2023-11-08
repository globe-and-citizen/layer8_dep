package handlers

import (
	"globe-and-citizen/layer8/proxy/internals/usecases"
	"html/template"
	"net/http"
)

// import (
// 	"globe-and-citizen/layer8/proxy/entities"
// 	"globe-and-citizen/layer8/proxy/internals/usecases"
// 	"html/template"
// 	"net/http"
// )

// func Register(w http.ResponseWriter, r *http.Request) {
// 	usecase := r.Context().Value("usecase").(*usecases.UseCase)

// 	switch r.Method {
// 	case "GET":
// 		next := r.URL.Query().Get("next")
// 		if next == "" {
// 			next = "/"
// 		}
// 		// check if the user is already logged in
// 		token, err := r.Cookie("token")
// 		if token != nil && err == nil {
// 			user, err := usecase.GetUserByToken(token.Value)
// 			if err == nil && user != nil {
// 				http.Redirect(w, r, next, http.StatusSeeOther)
// 				return
// 			}
// 		}

// 		// load the register page
// 		t, err := template.ParseFiles("assets/templates/register.html")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		t.Execute(w, map[string]interface{}{
// 			"HasNext": next != "",
// 			"Next":    next,
// 		})
// 		return
// 	case "POST":
// 		next := r.URL.Query().Get("next")
// 		if next == "" {
// 			next = "/"
// 		}
// 		user := &entities.User{
// 			AbstractUser: entities.AbstractUser{
// 				Username: r.FormValue("username"),
// 				Email:    r.FormValue("email"),
// 				Fname:    r.FormValue("fname"),
// 				Lname:    r.FormValue("lname"),
// 			},
// 			// due to using the same struct for the user and the pseudonymized data,
// 			// the validation will fail if the pseudonymized data is not present
// 			// so we set some dummy data here
// 			PsedonymizedData: entities.AbstractUser{
// 				Username: "dummy",
// 				Email:    "dummy",
// 				Fname:    "dummy",
// 				Lname:    "dummy",
// 			},
// 			Password: r.FormValue("password"),
// 		}
// 		err := user.Validate()
// 		if err != nil {
// 			t, errT := template.ParseFiles("assets/templates/register.html")
// 			if errT != nil {
// 				http.Error(w, errT.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			t.Execute(w, map[string]interface{}{
// 				"HasNext": next != "",
// 				"Next":    next,
// 				"Error":   err.Error(),
// 			})
// 			return
// 		}
// 		// register the user
// 		// rUser, err := usecase.RegisterUser(user)
// 		if err != nil {
// 			t, errT := template.ParseFiles("assets/templates/register.html")
// 			if errT != nil {
// 				http.Error(w, errT.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			t.Execute(w, map[string]interface{}{
// 				"HasNext": next != "",
// 				"Next":    next,
// 				"Error":   err.Error(),
// 			})
// 			return
// 		}
// 		// set the token cookie
// 		// token, ok := rUser["token"].(string)
// 		// if !ok {
// 		// 	t, err := template.ParseFiles("assets/templates/register.html")
// 		// 	if err != nil {
// 		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		// 		return
// 		// 	}
// 		// 	t.Execute(w, map[string]interface{}{
// 		// 		"HasNext": next != "",
// 		// 		"Next":    next,
// 		// 		"Error":   "could not get token",
// 		// 	})
// 		// 	return
// 		// }
// 		// http.SetCookie(w, &http.Cookie{
// 		// 	Name:  "token",
// 		// 	Value: token,
// 		// 	Path:  "/",
// 		// })
// 		// // redirecting to home page instead of the next page so that users can see their
// 		// // pseudo profile that they'll be identified by
// 		// if next == "/" {
// 		// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// 		// 	return
// 		// }
// 		// http.Redirect(w, r, "/?next="+next, http.StatusSeeOther)
// 		// return
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// }

func Login(w http.ResponseWriter, r *http.Request) {
	usecase := r.Context().Value("usecase").(*usecases.UseCase)

	switch r.Method {
	case "GET":
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}
		// check if the user is already logged in
		// token, err := r.Cookie("token")
		// if token != nil && err == nil {
		// 	user, err := usecase.GetUserByToken(token.Value)
		// 	if err == nil && user != nil {
		// 		http.Redirect(w, r, next, http.StatusSeeOther)
		// 		return
		// 	}
		// }

		// load the login page
		t, err := template.ParseFiles("assets/templates/login.html")
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
			t, errT := template.ParseFiles("assets/templates/login.html")
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
			t, err := template.ParseFiles("assets/templates/login.html")
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
