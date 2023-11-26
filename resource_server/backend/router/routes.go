package routes

import (
	"net/http"

	Ctl "globe-and-citizen/layer8/resource_server/backend/controller"
)

func RegisterRoutes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set up route for API
		switch r.URL.Path {

		case "/api/v1/register-user":
			Ctl.RegisterUserHandler(w, r)

		case "/api/v1/login-precheck":
			Ctl.LoginPrecheckHandler(w, r)

		case "/api/v1/login-user":
			Ctl.LoginUserHandler(w, r)

		case "/api/v1/profile":
			Ctl.ProfileHandler(w, r)

		case "/api/v1/verify-email":
			Ctl.VerifyEmailHandler(w, r)

		case "/api/v1/change-display-name":
			Ctl.UpdateDisplayNameHandler(w, r)

		default:
			// Return a 404 Not Found error for unknown routes
			http.NotFound(w, r)
		}
	}
}
