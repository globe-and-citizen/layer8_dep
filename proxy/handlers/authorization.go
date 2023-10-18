package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"layer8-proxy/constants"
	"layer8-proxy/internals/usecases"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/globe-and-citizen/layer8-utils"

	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

func Authorize(ctx *fasthttp.RequestCtx) {
	usecase := ctx.UserValue("usecase").(*usecases.UseCase)

	switch string(ctx.Method()) {
	case "GET":
		var (
			clientID          = string(ctx.QueryArgs().Peek("client_id"))
			scopes            = string(ctx.QueryArgs().Peek("scope"))
			redirectURI       = string(ctx.QueryArgs().Peek("redirect_uri"))
			scopeDescriptions = []string{}
			next              string
		)
		// use the default scope if none is provided
		if scopes == "" {
			scopes = constants.READ_USER_SCOPE
		}
		// add the scope descriptions
		for _, scope := range strings.Split(scopes, ",") {
			scopeDescriptions = append(scopeDescriptions, constants.ScopeDescriptions[scope])
		}
		// get the client
		client, err := usecase.GetClient(clientID)
		if err != nil {
			log.Println(err)
			// redirect to the redirect_uri with error
			ctx.Redirect("/error?opt=invalid_client", fasthttp.StatusSeeOther)
			return
		}
		// generate the next url
		uri, err := url.Parse("/authorize")
		if err != nil {
			log.Println(err)
			ctx.Redirect("/error?opt=server_error", fasthttp.StatusSeeOther)
			return
		}
		q := uri.Query()
		q.Set("client_id", clientID)
		q.Set("scope", scopes)
		uri.RawQuery = q.Encode()
		next = uri.String()

		// check that the user is logged in
		token := ctx.Request.Header.Cookie("token")
		user, err := usecase.GetUserByToken(string(token))
		// redirect to login page if not logged in
		if token == nil || err != nil || user == nil {
			ctx.Redirect("/login?next="+next, fasthttp.StatusSeeOther)
			return
		}
		// check that the redirect_uri is valid match the client's redirect_uri
		if redirectURI != "" && client.RedirectURI != redirectURI {
			ctx.Redirect("/error?opt=redirect_uri_mismatch", fasthttp.StatusSeeOther)
			return
		}
		// load the authorize page
		utilities.LoadTemplate(ctx, "assets/templates/authorize.html", map[string]interface{}{
			"ClientName": client.Name,
			"Scopes":     scopeDescriptions,
			"Next":       next,
		})
		return
	case "POST":
		var (
			clientID = string(ctx.QueryArgs().Peek("client_id"))
			scopes   = string(ctx.QueryArgs().Peek("scope"))
			returnResult, _ = strconv.ParseBool(string(ctx.QueryArgs().Peek("return_result")))
		)
		// get authorization decision
		decision := string(ctx.FormValue("decision"))
		if decision != "allow" {
			log.Println("User denied access")
			if returnResult {
				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.SetBody([]byte(`{"redr": "/error?opt=access_denied"}`))
				return
			}
			ctx.Redirect("/error?opt=access_denied", fasthttp.StatusSeeOther)
			return
		}
		// use the default scope if none is provided
		if scopes == "" {
			scopes = constants.READ_USER_SCOPE
		}
		// get the client
		client, err := usecase.GetClient(clientID)
		if err != nil {
			log.Println(err)
			if returnResult {
				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.SetBody([]byte(`{"redr": "/error?opt=invalid_client"}`))
				return
			}
			ctx.Redirect("/error?opt=invalid_client", fasthttp.StatusSeeOther)
			return
		}
		// get user
		token := ctx.Request.Header.Cookie("token")
		user, err := usecase.GetUserByToken(string(token))
		if err != nil || user == nil {
			if returnResult {
				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.SetBody([]byte(`{"redr": "/login?next=` + string(ctx.RequestURI()) + `"}`))
				return
			}
			ctx.Redirect("/login?next="+string(ctx.RequestURI()), fasthttp.StatusSeeOther)
			return
		}
		// generate authorization url
		authURL, err := usecase.GenerateAuthorizationURL(&oauth2.Config{
			ClientID:    client.ID,
			RedirectURL: client.RedirectURI,
			Scopes:      strings.Split(scopes, ","),
		}, user.ID)
		if err != nil {
			log.Println("Server error: ", err)
			if returnResult {
				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.SetBody([]byte(`{"redr": "/error?opt=server_error"}`))
				return
			}
			ctx.Redirect("error?opt=server_error", fasthttp.StatusSeeOther)
			return
		}
		// redirect to the authorization url
		if returnResult {
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBody([]byte(`{"url": "` + authURL.String() + `"}`))
			return
		}
		ctx.Redirect(authURL.String(), fasthttp.StatusSeeOther)
		return
	default:
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}
}

func OAuthToken(ctx *fasthttp.RequestCtx) {
	usecase := ctx.UserValue("usecase").(*usecases.UseCase)

	// exchange code for token
	switch string(ctx.Method()) {
	case "POST":
		var (
			code         = string(ctx.FormValue("code"))
			redirectURI  = string(ctx.FormValue("redirect_uri"))
		)

		// decode the basic auth header
		fromBasicAuth := func(t string) (string, string, error) {
			t = strings.TrimPrefix(t, "Basic ")
			b, err := base64.StdEncoding.DecodeString(t)
			if err != nil {
				return "", "", err
			}
			// first remove the "Basic " prefix
			s := strings.SplitN(string(b), ":", 2)
			if len(s) != 2 {
				return "", "", errors.New("invalid basic auth header")
			}
			return s[0], s[1], nil
		}
		clientID, clientSecret, err := fromBasicAuth(string(ctx.Request.Header.Peek("Authorization")))
		if err != nil {
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBody([]byte(`{"error": "` + err.Error() + `"}`))
			return
		}
		
		// get the client
		client, err := usecase.GetClient(clientID)
		if err != nil {
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBody([]byte(`{"error": "invalid client"}`))
			return
		}
		// check that the client secret is correct
		if client.Secret != clientSecret {
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBody([]byte(`{"error": "invalid client secret"}`))
			return
		}
		// exchange code for token
		token, err := usecase.ExchangeCodeForToken(&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
		}, code)
		if err != nil {
			res := map[string]string{"error": err.Error()}
			resJSON, _ := json.Marshal(res)
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBody(resJSON)
			return
		}
		// return token
		bToken, err := json.Marshal(token)
		if err != nil {
			res := map[string]string{"error": err.Error()}
			resJSON, _ := json.Marshal(res)
			ctx.SetContentType("application/json")
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBody(resJSON)
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(bToken)
		return
	default:
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.SetBody([]byte(`{"error": "method not allowed"}`))
		return
	}
}

func Error(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case "GET":
		var (
			opt    = string(ctx.QueryArgs().Peek("opt"))
			opts   = []string{}
			errors = map[string]string{
				"invalid_client":        "The client is invalid.",
				"access_denied":         "The user denied the request.",
				"server_error":          "An error occurred on the server.",
				"redirect_uri_mismatch": "The redirect uri does not match the client's redirect uri.",
			}
		)
		// add the error to the list of errors
		for _, v := range strings.Split(opt, ",") {
			opts = append(opts, errors[v])
		}
		// load the error page
		utilities.LoadTemplate(ctx, "assets/templates/error.html", map[string]interface{}{
			"Errors": opts,
		})
		return
	default:
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}
}
