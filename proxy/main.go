package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"globe-and-citizen/layer8/l8_oauth/internals/repository"
	"globe-and-citizen/layer8/l8_oauth/internals/usecases"

	"globe-and-citizen/layer8/l8_oauth/handlers"

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
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method != http.MethodConnect {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			handlers.Tunnel(w, r)
		}),
	}
	log.Printf("Starting proxy server on port %d...", port)
	log.Fatal(server.ListenAndServe())
}
