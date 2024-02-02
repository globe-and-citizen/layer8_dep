package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/interceptor/internals"
	"time"

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
		"static":              js.FuncOf(getStatic),
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

		headersMap := make(map[string]string)
		js.Global().Get("Object").Call("entries", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			headersMap[args[0].Index(0).String()] = args[0].Index(1).String()
			return nil
		}))

		// set the content-type to application/json if it's undefined
		if _, ok := headersMap["Content-Type"]; !ok {
			headersMap["Content-Type"] = "application/json"
		}

		go func() {
			var res *utils.Response

			switch strings.ToLower(headersMap["Content-Type"]) {
			case "application/json":
				// setting the body to an empty string if it's undefined
				body := options.Get("body")
				if body.String() == "<undefined>" {
					body = js.ValueOf(map[string]interface{}{})
				}

				// convert the body to a map
				bodyMap := map[string]interface{}{}
				err := json.Unmarshal([]byte(body.String()), &bodyMap)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return
				}

				// encode the body to json
				bodyByte, err := json.Marshal(bodyMap)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return
				}

				// forward request to the layer8 proxy server
				res = L8Client.Do(
					url, utils.NewRequest(method, headersMap, bodyByte),
					userSymmetricKey, false, UpJWT, UUID)
			case "multipart/form-data":
				headersMap["Content-Type"] = "application/layer8.buffer+json"

				body := options.Get("body")
				if body.String() == "<undefined>" || body.String() == "null" {
					reject.Invoke(js.Global().Get("Error").New("No body provided to fetch call."))
					return
				}

				// get data from formdata
				var (
					dataLength = js.Global().Get("Array").Call("from", body.Call("keys")).Get("length").Int()
					formdata   = make(map[string]interface{}, dataLength)
				)

				js.Global().Get("Array").Call("from", body.Call("keys")).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					var (
						key       = args[0].String()
						value     = body.Call("get", key)
						valueType = value.Get("constructor").Get("name").String()
					)

					switch valueType {
					case "File":
						value.Call("arrayBuffer").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
							buff := make([]byte, value.Get("size").Int())
							js.CopyBytesToGo(buff, js.Global().Get("Uint8Array").New(args[0]))

							data := map[string]interface{}{
								"_type": "File",
								"name":  value.Get("name").String(),
								"size":  value.Get("size").Int(),
								"type":  value.Get("type").String(),
								"buff":  base64.StdEncoding.EncodeToString(buff),
							}

							// because formdata can have multiple entries for the same key
							// each key is an array of maps
							if val, ok := formdata[key]; !ok {
								formdata[key] = []map[string]interface{}{data}
							} else {
								formdata[key] = append(val.([]map[string]interface{}), data)
							}
							return nil
						}))
					case "String":
						data := map[string]interface{}{
							"_type": "String",
							"value": value.String(),
						}

						if val, ok := formdata[key]; !ok {
							formdata[key] = []map[string]interface{}{data}
						} else {
							formdata[key] = append(val.([]map[string]interface{}), data)
						}
					case "Number":
						data := map[string]interface{}{
							"_type": "Number",
							"value": value.Float(),
						}

						if val, ok := formdata[key]; !ok {
							formdata[key] = []map[string]interface{}{data}
						} else {
							formdata[key] = append(val.([]map[string]interface{}), data)
						}
					case "Boolean":
						data := map[string]interface{}{
							"_type": "Boolean",
							"value": value.Bool(),
						}

						if val, ok := formdata[key]; !ok {
							formdata[key] = []map[string]interface{}{data}
						} else {
							formdata[key] = append(val.([]map[string]interface{}), data)
						}
					default:
						reject.Invoke(js.Global().Get("Error").New(fmt.Sprintf("Unsupported type: %s", valueType)))
						return nil
					}

					return nil
				}))

				// wait for the formdata to be populated, this is a hacky way to do it, but it works for now
				// having tried using a channel, it fails with a "fatal error: all goroutines are asleep - deadlock!"
				// TODO: find a better way to do this
				time.Sleep(100 * time.Millisecond)

				// encode the body to json
				bodyByte, err := json.Marshal(formdata)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return
				}

				// forward request to the layer8 proxy server
				res = L8Client.Do(
					url, utils.NewRequest(method, headersMap, bodyByte),
					userSymmetricKey, false, UpJWT, UUID)
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

func getStatic(this js.Value, args []js.Value) interface{} {
	url := args[0].String()

	return js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]

		go func() {
			resp := L8Client.Do(
				url, utils.NewRequest("GET", make(map[string]string), nil),
				userSymmetricKey, true, UpJWT, UUID)

			// convert response body to js arraybuffer
			jsBody := js.Global().Get("Uint8Array").New(len(resp.Body))
			js.CopyBytesToJS(jsBody, resp.Body)

			// create a map of the response headers
			resHeaders := js.Global().Get("Headers").New()
			for k, v := range resp.Headers {
				resHeaders.Call("append", js.ValueOf(k), js.ValueOf(v))
			}

			blob := js.Global().Get("Blob").New([]interface{}{jsBody}, js.ValueOf(map[string]interface{}{
				"type": resHeaders.Call("get", js.ValueOf("content-type")),
			}))
			objURL := js.Global().Get("URL").Call("createObjectURL", blob)

			resolve.Invoke(objURL)
		}()

		return nil
	}))
}
