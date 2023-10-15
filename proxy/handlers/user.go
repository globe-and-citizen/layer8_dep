package handlers

import (
	"encoding/json"
	"layer8-proxy/internals/usecases"
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

// UserInfo handles requests to get a user's anonymized data
// The last step of the oauth flow
func UserInfo(ctx *fasthttp.RequestCtx) {
	usecase := ctx.UserValue("usecase").(*usecases.UseCase)

	switch string(ctx.Method()) {
	case "GET":
		token := strings.TrimPrefix(string(ctx.Request.Header.Peek("Authorization")), "Bearer ")
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
