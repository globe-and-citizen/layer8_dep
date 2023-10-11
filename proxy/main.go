package main

import (
	"flag"
	"fmt"
	"layer8-proxy/handlers"
	"layer8-proxy/internals/repository"
	"layer8-proxy/internals/usecases"
	"log"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	port   = flag.Int("port", 5000, "Port to listen on")
	server = flag.String("server", "auth", "Server type to run")
)

func main() {
	flag.Parse()

	switch *server {
	case "auth":
		AuthServer(*port)
	case "proxy":
		ProxyServer(*port)
	default:
		log.Fatal("Invalid server type. Valid types are: auth, proxy")
	}
}

func AuthServer(port int) {
	log.Printf("Starting auth server on port %d...", port)

	// intialize a memory repository and usecase
	repo, err := repository.CreateRepository("memory")
	if err != nil {
		log.Fatal(err)
	}
	usecase := &usecases.UseCase{Repo: repo}

	// for testing purposes, we'll create a test client for the example service provider
	_, err = usecase.AddTestClient()
	if err != nil {
		log.Fatal(err)
	}

	
	fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), func(ctx *fasthttp.RequestCtx) {
		// allow from all origins
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")

		// add the usecase to context
		ctx.SetUserValue("usecase", usecase)
		// routing
		switch path := string(ctx.Path()); {
		case path == "" || path == "/":
			handlers.Welcome(ctx)
		case path == "/login":
			handlers.Login(ctx)
		case path == "/register":
			handlers.Register(ctx)
		case path == "/authorize":
			handlers.Authorize(ctx)
		case path == "/error":
			handlers.Error(ctx)
		case path == "/api/oauth":
			handlers.OAuthToken(ctx)
		case path == "/api/user":
			handlers.UserInfo(ctx)
		case strings.HasPrefix(path, "/assets"):
			// serve static files
			fasthttp.FSHandler("./assets", 1)(ctx)
		default:
			ctx.Error("Invalid path", fasthttp.StatusNotFound)
		}
		log.Printf("%d %s\t%s", ctx.Response.StatusCode(), ctx.Method(), ctx.Path())
	})
}

func ProxyServer(port int) {
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			handlers.Tunnel(w, r)
		}),
	}
	log.Printf("Starting proxy server on port %d...", port)
	log.Fatal(server.ListenAndServe())
}
