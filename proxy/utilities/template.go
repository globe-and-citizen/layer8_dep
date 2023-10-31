package utilities

import (
	"html/template"

	"github.com/valyala/fasthttp"
)

func LoadTemplate(ctx *fasthttp.RequestCtx, path string, data interface{}) {
	ctx.SetContentType("text/html")
	t, err := template.ParseFiles(path)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	t.Execute(ctx, data)
}
