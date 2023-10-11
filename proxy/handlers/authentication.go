package handlers

import (
	utilities "github.com/globe-and-citizen/layer8-utils"
	"layer8-proxy/entities"
	"layer8-proxy/internals/usecases"

	"github.com/valyala/fasthttp"
)

func Register(ctx *fasthttp.RequestCtx) {
	var usecase = ctx.UserValue("usecase").(*usecases.UseCase)

	switch string(ctx.Method()) {
	case "GET":
		next := string(ctx.QueryArgs().Peek("next"))
		if next == "" {
			next = "/"
		}
		// check if the user is already logged in
		token := ctx.Request.Header.Cookie("token")
		user, err := usecase.GetUserByToken(string(token))
		if token != nil && err == nil && user != nil {
			ctx.Redirect(next, fasthttp.StatusSeeOther)
			return
		}

		// load the register page
		utilities.LoadTemplate(ctx, "assets/templates/register.html", map[string]interface{}{
			"HasNext": next != "",
			"Next":    next,
		})
		return
	case "POST":
		next := string(ctx.QueryArgs().Peek("next"))
		if next == "" {
			next = "/"
		}
		user := &entities.User{
			AbstractUser: entities.AbstractUser{
				Username: string(ctx.FormValue("username")),
				Email:    string(ctx.FormValue("email")),
				Fname:    string(ctx.FormValue("fname")),
				Lname:    string(ctx.FormValue("lname")),
			},
			// due to using the same struct for the user and the pseudonymized data,
			// the validation will fail if the pseudonymized data is not present
			// so we set some dummy data here
			PsedonymizedData: entities.AbstractUser{
				Username: "dummy",
				Email:    "dummy",
				Fname:    "dummy",
				Lname:    "dummy",
			},
			Password: string(ctx.FormValue("password")),
		}
		err := user.Validate()
		if err != nil {
			utilities.LoadTemplate(ctx, "assets/templates/register.html", map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   err.Error(),
			})
			return
		}
		// register the user
		rUser, err := usecase.RegisterUser(user)
		if err != nil {
			utilities.LoadTemplate(ctx, "assets/templates/register.html", map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   err.Error(),
			})
			return
		}
		// set the token cookie
		token, ok := rUser["token"].(string)
		if !ok {
			utilities.LoadTemplate(ctx, "assets/templates/register.html", map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   "could not get token",
			})
			return
		}
		c := new(fasthttp.Cookie)
		c.SetKey("token")
		c.SetValue(token)
		c.SetPath("/")
		ctx.Response.Header.SetCookie(c)
		// redirecting to home page instead of the next page so that users can see their
		// pseudo profile that they'll be identified by
		if next == "/" {
			ctx.Redirect("/", fasthttp.StatusSeeOther)
			return
		}
		ctx.Redirect("/?next="+next, fasthttp.StatusSeeOther)
		return
	default:
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}
}

func Login(ctx *fasthttp.RequestCtx) {
	var usecase = ctx.UserValue("usecase").(*usecases.UseCase)

	switch string(ctx.Method()) {
	case "GET":
		next := string(ctx.QueryArgs().Peek("next"))
		if next == "" {
			next = "/"
		}
		// check if the user is already logged in
		token := ctx.Request.Header.Cookie("token")
		user, err := usecase.GetUserByToken(string(token))
		if token != nil && err == nil && user != nil {
			ctx.Redirect(next, fasthttp.StatusSeeOther)
			return
		}

		// load the login page
		utilities.LoadTemplate(ctx, "assets/templates/login.html", map[string]interface{}{
			"HasNext": next != "",
			"Next":    next,
		})
		return
	case "POST":
		next := string(ctx.QueryArgs().Peek("next"))
		username := string(ctx.FormValue("username"))
		password := string(ctx.FormValue("password"))
		// login the user
		rUser, err := usecase.LoginUser(username, password)
		if err != nil {
			utilities.LoadTemplate(ctx, "assets/templates/login.html", map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   err.Error(),
			})
			return
		}
		// set the token cookie
		token, ok := rUser["token"].(string)
		if !ok {
			utilities.LoadTemplate(ctx, "assets/templates/login.html", map[string]interface{}{
				"HasNext": next != "",
				"Next":    next,
				"Error":   "could not get token",
			})
			return
		}
		c := new(fasthttp.Cookie)
		c.SetKey("token")
		c.SetValue(token)
		c.SetPath("/")
		ctx.Response.Header.SetCookie(c)
		// redirect to next page - here the user already knows their pseudo profile
		// when they registered
		ctx.Redirect(next, fasthttp.StatusSeeOther)
		return
	default:
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}
}
