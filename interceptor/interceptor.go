package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/interceptor/internals"
	"globe-and-citizen/layer8/utils"
	"net/http"
	"syscall/js"
)

// Declare global constants
const INTERCEPTOR_VERSION = "1.0.0"

// Declare global variables
var (
	Layer8Scheme     string = "http"
	Layer8Host       string = "localhost"
	Layer8Port       string = "5001" // 5001 is Proxy & 5000 is Auth server.
	Layer8Version    string = "1.0.0"
	Counter          int    = 0
	ETunnelFlag      bool   = false
	privJWK_ecdh     *utils.JWK
	pubJWK_ecdh      *utils.JWK
	userSymmetricKey *utils.JWK
)

func main() {
	// Create channel to keep the Go thread alive
	c := make(chan struct{}, 0)

	// Expose layer8 functionality to the front end Javascript
	js.Global().Set("layer8", js.ValueOf(map[string]interface{}{
		"testWASM":         js.FuncOf(testWASM),
		"persistenceCheck": js.FuncOf(persistenceCheck),
		// "genericGetRequest": js.FuncOf(genericGetRequest),
		// "genericPost":       js.FuncOf(genericPost),
		"fetch": js.FuncOf(fetch),
	}))

	// Initialize the encrypted tunnel
	initializeECDHTunnel()

	// Developer Warnings:
	fmt.Println("WARNING: wasm_exec.js is versioned and has some breaking changes. Ensure you are using the correct version.")

	// Wait indefinitely
	<-c
}

// Utility function to test promise resolution / rejection from the console.
func testWASM(this js.Value, args []js.Value) interface{} {
	var promise_logic = func(this js.Value, resolve_reject []js.Value) interface{} {
		resolve := resolve_reject[0]
		reject := resolve_reject[1]
		if len(args) == 0 {
			reject.Invoke(js.ValueOf("Promise rejection occurs if not arguments are passed. Pass an argument."))
			return nil
		}
		go func() {
			resolve.Invoke(js.ValueOf(fmt.Sprintf("WASM Interceptor version %s successfully loaded. Argument passed: %v. To test promise rejection, call with no argument.", INTERCEPTOR_VERSION, args[0])))
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(promise_logic))
	return promise
}

func persistenceCheck(this js.Value, args []js.Value) interface{} {
	var promise_logic = func(this js.Value, resolve_reject []js.Value) interface{} {
		resolve := resolve_reject[0]
		go func() {
			Counter++
			fmt.Println("WASM Counter: ", Counter)
			resolve.Invoke(js.ValueOf(Counter))
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(promise_logic))
	return promise
}

func initializeECDHTunnel() {
	go func() {
		var err error
		privJWK_ecdh, pubJWK_ecdh, err = utils.GenerateKeyPair(utils.ECDH)
		if err != nil {
			fmt.Println(err.Error())
			ETunnelFlag = false
			return
		}

		b64PubJWK, err := pubJWK_ecdh.ExportAsBase64()
		if err != nil {
			fmt.Println(err.Error())
			ETunnelFlag = false
			return
		}

		ProxyURL := fmt.Sprintf("%s://%s:%s", Layer8Scheme, Layer8Host, Layer8Port)
		client := &http.Client{}
		req, err := http.NewRequest("GET", ProxyURL, bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Println(err.Error())
			ETunnelFlag = false
			return
		}
		req.Header.Add("x-ecdh-init", b64PubJWK)
		req.Header.Add("X-client-id", "1")

		// send request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			ETunnelFlag = false
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 401 {
			fmt.Printf("User not authorized\n")
			ETunnelFlag = false
			return
		}

		// TODO: For some reason I am unable to put (or access?) custom response headers coming from
		// either the backend OR the proxy... Therefore, I've sent along the backend's public key in the
		// response's body.
		// for k, v := range resp.Header {
		// 	fmt.Println("header pairs from SP:", k, v)
		// }
		// fmt.Println("resp.Header: ", resp.Header.Get("Content-Length"))

		Respbody := utils.ReadResponseBody(resp.Body)
		//fmt.Println("response body: ", string(Respbody))
		server_pubKeyECDH, err := utils.B64ToJWK(string(Respbody))
		if err != nil {
			fmt.Println(err.Error())
			ETunnelFlag = false
			return
		}

		userSymmetricKey, err = privJWK_ecdh.GetECDHSharedSecret(server_pubKeyECDH)
		if err != nil {
			fmt.Println(err.Error())
			ETunnelFlag = false
			return
		}

		// TODO: Send an encrypted ping / confirmation to the server using the shared secret
		// just like the 1. Syn 2. Syn/Ack 3. Ack flow in a TCP handshake
		ETunnelFlag = true
		fmt.Println("Encrypted tunnel successfully established.")
		return
	}()

	return
}

func fetch(this js.Value, args []js.Value) interface{} {
	var promise_logic = func(this js.Value, resolve_reject []js.Value) interface{} {
		resolve := resolve_reject[0]
		reject := resolve_reject[1]

		if !ETunnelFlag {
			reject.Invoke(js.Global().Get("Error").New("The Encrypted tunnel is closed. Reload page."))
			return nil
		}

		if len(args) == 0 {
			reject.Invoke(js.Global().Get("Error").New("No URL provided to fetch call."))
			return nil
		}

		url := args[0].String()
		if len(url) <= 0 {
			reject.Invoke(js.Global().Get("Error").New("Invalid URL provided to fetch call."))
			return nil
		}

		options := js.ValueOf(map[string]interface{}{
			"method":  "GET", // Set HTTP "GET" request to be the default
			"headers": js.ValueOf(map[string]interface{}{}),
		})

		if len(args) > 1 {
			options = args[1]
		}

		method := options.Get("method").String()
		if method == "" {
			method = "GET"
		}

		// Set headers to an empty object if it is 'undefined' or 'null'
		headers := options.Get("headers")
		if headers.String() == "<undefined>" || headers.String() == "null" {
			headers = js.ValueOf(map[string]interface{}{})
		}

		// setting the body to an empty string if it's undefined
		body := options.Get("body").String()
		if body == "<undefined>" {
			body = ""
		}

		headersMap := make(map[string]string)
		js.Global().Get("Object").Call("entries", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			headersMap[args[0].Index(0).String()] = args[0].Index(1).String()
			return nil
		}))

		go func() {
			// forward request to the layer8 proxy server
			res := internals.NewClient(Layer8Scheme, Layer8Host, Layer8Port).
				Do(url, utils.NewRequest(method, headersMap, []byte(body)), userSymmetricKey)

			if res.Status >= 100 && res.Status < 300 { // Handle Success & Default Rejection
				resHeaders := js.Global().Get("Headers").New()

				for k, v := range res.Headers {
					resHeaders.Call("append", js.ValueOf(k), js.ValueOf(v))
				}

				resolve.Invoke(js.Global().Get("Response").New(string(res.Body), js.ValueOf(map[string]interface{}{
					"status":     res.Status,
					"statusText": res.StatusText,
					"headers":    resHeaders,
				})))
				return
			}

			bodyMap := make(map[string]interface{})
			json.Unmarshal(res.Body, &bodyMap)
			reject.Invoke(js.Global().Get("Error").New(res.StatusText, js.ValueOf(map[string]interface{}{
				"cause": js.ValueOf(bodyMap),
			})))
			return
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(promise_logic))
	return promise
}
