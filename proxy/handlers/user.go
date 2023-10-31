package handlers

import (
	"encoding/json"
	"log"

	"globe-and-citizen/layer8/l8_oauth/internals/usecases"

	"github.com/valyala/fasthttp"
)

// UserInfo handles requests to get a user's anonymized data
// The last step of the oauth flow
func UserInfo(ctx *fasthttp.RequestCtx) {
	usecase := ctx.UserValue("usecase").(*usecases.UseCase)

	switch string(ctx.Method()) {
	case "GET":
		token := string(ctx.QueryArgs().Peek("access_token"))
		data, err := usecase.AccessResourcesWithToken(token)
		if err != nil {
			log.Println(err)
			ctx.Error(`{"error": "invalid token"}`, fasthttp.StatusUnauthorized)
			return
		}
		b, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			ctx.Error(`{"error": "server error"}`, fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(b)
		return
	default:
		ctx.Error(`{"error": "method not allowed"}`, fasthttp.StatusMethodNotAllowed)
	}
}
