package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/interceptor/internals"

	"net/http"
	"strings"
	"syscall/js"

	utils "github.com/globe-and-citizen/layer8-utils"

	uuid "github.com/google/uuid"
)

// Declare global constants
const INTERCEPTOR_VERSION = "1.0.0"

// Declare global variables
var (
	Layer8Scheme       string
	Layer8Host         string
	Layer8Port         string
	Layer8Version      string
	Layer8LightsailURL string
	Counter            int
	ETunnelFlag        bool
	privJWK_ecdh       *utils.JWK
	pubJWK_ecdh        *utils.JWK
	userSymmetricKey   *utils.JWK
	UpJWT              string
	UUID               string
)

var L8Client = internals.NewClient(Layer8Scheme, Layer8Host, Layer8Port)

func main() {
	// Create channel to keep the Go thread alive
	c := make(chan struct{})

	// Initialize global variables
	Layer8Version = "1.0.0"
	Layer8Scheme = "http"
	Layer8Host = "localhost"
	Layer8Port = "5001"
	// Layer8Scheme = "https"
	// Layer8Host = "aws-container-service-t1.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com"
	// Layer8Port = ""
	// Layer8LightsailURL = "https://aws-container-service-t1.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com"

	ETunnelFlag = false

	// Expose layer8 functionality to the front end Javascript
	js.Global().Set("layer8", js.ValueOf(map[string]interface{}{
		"testWASM":            js.FuncOf(testWASM),
		"persistenceCheck":    js.FuncOf(persistenceCheck),
		"InitEncryptedTunnel": js.FuncOf(initializeECDHTunnel),
		"fetch":               js.FuncOf(fetch),
	}))

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

func initializeECDHTunnel(this js.Value, args []js.Value) interface{} {
	backend := args[0].String()

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

		ProxyURL := fmt.Sprintf("%s://%s:%s/init-tunnel?backend=%s", Layer8Scheme, Layer8Host, Layer8Port, backend)
		// ProxyURL := fmt.Sprintf("%s/init-tunnel?backend=%s", Layer8LightsailURL, backend)
		fmt.Println(ProxyURL)
		client := &http.Client{}
		req, err := http.NewRequest("POST", ProxyURL, bytes.NewBuffer([]byte(b64PubJWK)))
		if err != nil {
			fmt.Println(err.Error())
			ETunnelFlag = false
			return
		}
		uuid := uuid.New()
		UUID = uuid.String()
		req.Header.Add("x-ecdh-init", b64PubJWK)
		req.Header.Add("x-client-uuid", uuid.String())

		// send request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			ETunnelFlag = false
			return
		}

		if resp.StatusCode == 401 {
			fmt.Printf("User not authorized\n")
			ETunnelFlag = false
			return
		}

		Respbody := utils.ReadResponseBody(resp.Body)

		data := map[string]interface{}{}

		err = json.Unmarshal(Respbody, &data)
		if err != nil {
			if strings.Contains(err.Error(), "unexpected end of JSON input") {
				fmt.Println("JSON data might be incomplete or improperly formatted.")
			} else {
				fmt.Println(err.Error())
			}
			ETunnelFlag = false
			return
		}

		UpJWT = data["up_JWT"].(string)

		server_pubKeyECDH, err := utils.JWKFromMap(data)
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

	return nil
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

		// set the UpJWT to the headers
		headers.Set("up_JWT", UpJWT)

		// set the UUID to the headers
		headers.Set("x-client-uuid", UUID)

		headersMap := make(map[string]string)
		js.Global().Get("Object").Call("entries", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			headersMap[args[0].Index(0).String()] = args[0].Index(1).String()
			return nil
		}))

		// Print the headersMap for debugging purposes
		for k, v := range headersMap {
			fmt.Println("Encrypted Headers from the SP: ", k, v)
		}

		// set the content-type to application/json if it's undefined
		if _, ok := headersMap["Content-Type"]; !ok {
			headersMap["Content-Type"] = "application/json"
		}
		
		go func() {
			var res  *utils.Response
		
			switch strings.ToLower(headersMap["Content-Type"]) {
			case "application/json":
				// setting the body to an empty string if it's undefined
				body := options.Get("body")
				if body.String() == "<undefined>" {
					body = js.ValueOf(map[string]interface{}{})
				}

				// convert the body to a map
				bodyMap := map[string]interface{}{}
				js.Global().Get("Object").Call("entries", body).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					bodyMap[args[0].Index(0).String()] = args[0].Index(1).String()
					return nil
				}))

				// encode the body to json
				bodyByte, err := json.Marshal(bodyMap)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return
				}
				
				// forward request to the layer8 proxy server
				res = L8Client.
					Do(url, utils.NewRequest(method, headersMap, bodyByte), userSymmetricKey)
			case "application/layer8.buffer+json":
				requirements := []string{"name", "size", "type", "buff"}

				body := options.Get("body")
				if body.String() == "<undefined>" || body.String() == "null" {
					reject.Invoke(js.Global().Get("Error").New("No body provided to fetch call."))
					return
				}

				// check if the body has all the required fields
				notFound := []string{}
				for _, v := range requirements {
					fmt.Println("jv: ", body.Get(v).String())
					if jv := body.Get(v); jv.String() == "<undefined>" || jv.String() == "null" {
						notFound = append(notFound, v)
					}
				}
				if len(notFound) > 0 {
					reject.Invoke(js.Global().Get("Error").New(fmt.Sprintf("The following fields are required in the body: %v", notFound)))
					return
				}

				// convert arraybuffer to Go byte array
				buff := make([]byte, body.Get("size").Int())
				js.CopyBytesToGo(buff, js.Global().Get("Uint8Array").New(body.Get("buff")))

				// convert the body to a map
				bodyMap := map[string]interface{}{
					"name": body.Get("name").String(),
					"size": body.Get("size").Int(),
					"type": body.Get("type").String(),
					"buff": base64.StdEncoding.EncodeToString(buff),
				}
				
				// encode the body to json
				bodyByte, err := json.Marshal(bodyMap)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return
				}

				// forward request to the layer8 proxy server
				res = L8Client.
					Do(url, utils.NewRequest(method, headersMap, bodyByte), userSymmetricKey)
			default:
				res = &utils.Response{
					Status:     400,
					StatusText: "Content-Type not supported",
				}
			}

			if res.Status >= 100 || res.Status < 300 { // Handle Success & Default Rejection
				resHeaders := js.Global().Get("Headers").New()

				for k, v := range res.Headers {
					//fmt.Println("Encrypted Headers from the SP: ", k, v)
					resHeaders.Call("append", js.ValueOf(k), js.ValueOf(v))
				}

				resolve.Invoke(js.Global().Get("Response").New(string(res.Body), js.ValueOf(map[string]interface{}{
					"status":     res.Status,
					"statusText": res.StatusText,
					"headers":    resHeaders,
				})))
				return
			}

			reject.Invoke(js.Global().Get("Error").New(res.StatusText))
			fmt.Println("status:", res.Status, res.StatusText)
			return
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(promise_logic))
	return promise
}
