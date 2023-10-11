package handlers

import (
	"github.com/globe-and-citizen/layer8-utils"
	"layer8-proxy/internals/usecases"

	"github.com/valyala/fasthttp"
)

func Welcome(ctx *fasthttp.RequestCtx) {
	usecase := ctx.UserValue("usecase").(*usecases.UseCase)

	switch string(ctx.Method()) {
	case "GET":
		next := string(ctx.QueryArgs().Peek("next"))
		token := ctx.Request.Header.Cookie("token")
		user, err := usecase.GetUserByToken(string(token))
		if token == nil || err != nil || user == nil {
			if next == "" {
				ctx.Redirect("/login", fasthttp.StatusSeeOther)
			} else {
				ctx.Redirect("/login?next="+next, fasthttp.StatusSeeOther)
			}
			return
		}

		// load the welcome page
		utilities.LoadTemplate(ctx, "assets/templates/welcome.html", map[string]interface{}{
			"User":    user,
			"HasNext": next != "",
			"Next":    next,
		})
		return
	default:
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
	}
}
