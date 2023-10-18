package main

import (
	"syscall/js"

	utilities "github.com/globe-and-citizen/layer8-utils"
	"github.com/globe-and-citizen/layer8/interceptor/internals"
)

var (
	Layer8Scheme  string
	Layer8Host    string
	Layer8Port    string
)

func fetch(this js.Value, args []js.Value) interface{} {
	url := args[0].String()
	options := js.ValueOf(map[string]interface{}{
		"method":  "GET",
		"headers": js.ValueOf(map[string]interface{}{}),
	})
	if len(args) > 1 {
		options = args[1]
	}

	method := options.Get("method").String()
	if method == "" {
		method = "GET"
	}
	headers := options.Get("headers")
	// setting headers to an empty object if it's undefined or null
	if headers.String() == "<undefined>" || headers.String() == "null" {
		headers = js.ValueOf(map[string]interface{}{})
	}
	body := options.Get("body").String()
	// setting the body to an empty string if it's undefined
	if body == "<undefined>" {
		body = ""
	}

	promise := js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// to avoid deadlock with the main thread, we need to run this in a goroutine
		go func() {
			// add headers
			headersMap := make(map[string]string)
			js.Global().Get("Object").Call("entries", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				headersMap[args[0].Index(0).String()] = args[0].Index(1).String()
				return nil
			}))

			// forward request to the layer8 proxy server
			res := internals.NewClient(Layer8Scheme, Layer8Host, Layer8Port).
				Do(url, utilities.NewRequest(method, headersMap, []byte(body)))
			if res.Status < 300 {
				resHeaders := js.Global().Get("Headers").New()
				for k, v := range res.Headers {
					resHeaders.Call("append", k, v)
				}
				args[0].Invoke(js.Global().Get("Response").New(string(res.Body), js.ValueOf(map[string]interface{}{
					"status":     res.Status,
					"statusText": res.StatusText,
					"headers":    resHeaders,
				})))
			} else {
				args[1].Invoke(js.Global().Get("Error").New(res.StatusText))
			}
		}()
		return nil
	}), js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// for rejection, we just return the error message
		return args[0].String()
	}))

	return promise
}

func main() {
	// keep the thread alive
	close := make(chan struct{}, 0)

	// expose the fetch function to the global scope
	js.Global().Set("layer8", js.ValueOf(map[string]interface{}{
		"fetch": js.FuncOf(fetch),
	}))

	// wait indefinitely
	<-close
}
