package main

//Tomorrows labour: Refactor the interceptor to have better error handling.

import (
	"bytes"
	"fmt"
	"globe-and-citizen/layer8/interceptor/internals"
	"globe-and-citizen/layer8/utils"
	"io"
	"net/http"
	"syscall/js"
)

const VERSION = "1.0.2"

var counter = 0

type Request struct {
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

var (
	Layer8Scheme  string
	Layer8Host    string
	Layer8Port    string
	Layer8Version string
	privJWK_ecdh  *utils.JWK
	pubJWK_ecdh   *utils.JWK
)

var spoofedSymmetricKey *utils.JWK = &utils.JWK{
	Key_ops: []string{"encrypt", "decrypt"},
	Kty:     "EC",
	Kid:     "SpoofedKey",
	Crv:     "P-256",
	X:       "KfbCmY2v83ptAZLLKffx0ve2Br8hkMhCkIo5RkFaNlk=",
	Y:       "",
	D:       "",
}

func main() {
	// keep the Go thread alive
	c := make(chan struct{}, 0)
	Layer8Scheme = "http"
	Layer8Host = "localhost"
	Layer8Port = "5001" // 5001 is Proxy & 5000 is Auth server.
	Layer8Version = "1.0"

	// expose the layer8 functionality the global scope
	js.Global().Set("layer8", js.ValueOf(map[string]interface{}{
		"testWASMLoaded":    js.FuncOf(testWASMLoaded),
		"testWASM":          js.FuncOf(testWASM),
		"genericGetRequest": js.FuncOf(genericGetRequest),
		"genericPost":       js.FuncOf(genericPost),
		"fetch":             js.FuncOf(fetch),
	}))

	fmt.Println("WARNING: wasm_exec.js is versioned and has some breaking changes. Ensure you are using the correct version.")
	if initializeECDHTunnel() {
		fmt.Println("ECDH successfully inited")
	} else {
		fmt.Println("ECDH failed...")
	}
	// Wait indefinitely
	<-c
}

func testWASMLoaded(this js.Value, args []js.Value) interface{} {
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		//reject := args[1]
		go func() {
			// Main function body
			//fmt.Println(string(args[2]))
			fmt.Printf("WASM Interceptor version %s successfully loaded.", VERSION)
			counter++
			resolve.Invoke(js.ValueOf(counter))
			//reject.Invoke()
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}

func testWASM(this js.Value, args []js.Value) interface{} {
	fmt.Println("Fisrt argument: ", args[0])
	fmt.Println("Second argument: ", args[1])
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		//reject := args[1]
		go func() {
			// Main function body
			//fmt.Println(string(args[2]))
			resolve.Invoke(js.ValueOf(fmt.Sprintf("WASM Interceptor version %s successfully loaded.", VERSION)))
			//reject.Invoke()
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}

func initializeECDHTunnel() bool {
	fmt.Println("Top of Init")
	var err error
	privJWK_ecdh, pubJWK_ecdh, err = utils.GenerateKeyPair(utils.ECDH)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	b64PubJWK, err := pubJWK_ecdh.ExportAsBase64()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println("b64PubJWK: ", b64PubJWK)
	//URL will need to be defaulted from the page.
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:5001", bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	req.Header.Add("x-ecdh-init", b64PubJWK)
	req.Header.Add("X-client-id", "1")

	// send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		fmt.Printf("User not authorized\n")
		return false
	}

	//resHeaders := make(map[string]interface{})
	for k, v := range resp.Header {
		fmt.Println("header pairs from SP:", k, v)
	}
	fmt.Println("resp.Header: ", resp.Header.Get("Content-Length"))

	Respbody := utils.ReadResponseBody(resp.Body)

	fmt.Println("response body: ", string(Respbody))

	server_pubKeyECDH, _ := utils.B64ToJWK(string(Respbody))

	spoofedSymmetricKey, _ = privJWK_ecdh.GetECDHSharedSecret(server_pubKeyECDH)

	fmt.Println("SS_User: ", spoofedSymmetricKey)
	return true
}

func fetch(this js.Value, args []js.Value) interface{} {
	// All async functions must return only a promise, succeed or fail. To ensure that even failure returns as a promise,
	// the resolve_reject_internals must wrap the entire function logic.

	var resolve_reject_internals = func(this js.Value, resolve_reject []js.Value) interface{} {
		resolve := resolve_reject[0]
		reject := resolve_reject[1]

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
			"method":  "GET", // Set "GET" to be the default
			"headers": js.ValueOf(map[string]interface{}{}),
		})
		if len(args) > 1 {
			options = args[1]
		}

		method := options.Get("method").String()

		if method == "" { // redundant? Already set by default above
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

		headersMap := make(map[string]string)
		js.Global().Get("Object").Call("entries", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			headersMap[args[0].Index(0).String()] = args[0].Index(1).String()
			return nil
		}))

		go func() {
			// forward request to the layer8 proxy server
			res := internals.NewClient(Layer8Scheme, Layer8Host, Layer8Port).
				Do(url, utils.NewRequest(method, headersMap, []byte(body)), spoofedSymmetricKey) //This will send ALL the headers

			if res.Status >= 100 || res.Status < 300 { // Handle Success & Default Rejection
				resHeaders := js.Global().Get("Headers").New()
				for k, v := range res.Headers {
					resHeaders.Call("append", k, v)
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
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}

func Dep_fetch(this js.Value, args []js.Value) interface{} {
	url := args[0].String()
	options := js.ValueOf(map[string]interface{}{
		"method":  js.ValueOf(""),
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

	//fmt.Println("headers outside: ", headers)
	body := options.Get("body").String()
	if body == "<undefined>" {
		body = ""
	}

	fmt.Println("body: ", body)

	promise := js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			headersMap := make(map[string]string)
			// build the headersMap
			js.Global().Get("Object").Call("keys", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				// args[0] is key & args[1] is value? or the index?
				fmt.Println("args [1]  & [2]", args[0], args[1])
				headersMap[args[0].String()] = args[1].String()
				return nil
			}))

			// testRequest := Request{
			// 	Method:  "POST",
			// 	Headers: make(map[string]string),
			// 	Body:    []byte("Hello Layer8"),
			// }

			//testRequest.Headers["x-test"] = "test header"

			//data, err := json.Marshal(testRequest)

			// if err != nil {
			// 	panic("fuck this...")
			// }

			// forward request to the layer8 proxy server
			client := &http.Client{}
			r, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))

			if err != nil {
				args[1].Invoke(js.ValueOf("Problem Creating Request"))
			}

			res, err := client.Do(r)
			if err != nil {
				res := &utils.Response{
					Status:     500,
					StatusText: err.Error(),
				}
				resByte, _ := res.ToJSON()
				fmt.Println(resByte)
				args[1].Invoke(js.ValueOf("Still and error but closer"))
			}

			if res == nil || res.Body == nil {
				fmt.Println("res or res.body was nil.")
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				fmt.Println("Server returned non-OK stauts: ", res.Status)
				args[1].Invoke(js.ValueOf(fmt.Sprintf("Server returned non-OK stauts: ", res.Status)))
				return
			}

			// buf := new(bytes.Buffer)
			// buf.ReadFrom(res.Body)
			res_byteSlice, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Server side error: ", err.Error())
				args[1].Invoke(js.ValueOf(err.Error()))
				return
			}
			args[0].Invoke(js.ValueOf(string(res_byteSlice)))

			return
		}()
		return nil
	}), js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// call reject() with the error message cast as a string.
		return args[0].String()
	}))

	return promise
}

func genericGetRequest(this js.Value, args []js.Value) interface{} {
	url := args[0]
	fmt.Println("HERE: ", url.String())
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			// Main function body
			res, err := http.Get(url.String()) // http://localhost:5000/success
			if err != nil {
				fmt.Println("Failure to get ", url.String())
				reject.Invoke(js.ValueOf(err.Error()))
			}

			if res == nil || res.Body == nil {
				fmt.Println("res or res.body from proxy was nil.")
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				fmt.Println("Server returned non-OK stauts: ", res.Status)
				reject.Invoke(js.ValueOf(fmt.Sprintf("Server returned non-OK stauts: ", res.Status)))
				return
			}

			res_byteSlice, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Error reading backend public key: ", err)
				reject.Invoke(js.ValueOf(err.Error()))
				return
			}

			resolve.Invoke(js.ValueOf(string(res_byteSlice)))
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}

func genericPost(this js.Value, args []js.Value) interface{} {
	url := args[0]
	//stringifiedObject := `{"client_message": "hello, server!"}`
	stringifiedObject := args[1].String()
	fmt.Println(stringifiedObject)
	jsonBody := []byte(stringifiedObject)
	bodyReader := bytes.NewReader(jsonBody)
	//	bodyReader := bytes.NewReader(byteObject)

	fmt.Println("Interceptor will now POST to this url: ", url.String())
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			// Main function body
			res, err := http.Post(url.String(), "application/json", bodyReader)
			if err != nil {
				fmt.Println("Failure to Post ", url.String())
				reject.Invoke(js.ValueOf(err.Error()))
			}

			if res == nil || res.Body == nil {
				fmt.Println("res or res.body from proxy was nil.")
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				fmt.Println("Server returned non-OK stauts: ", res.Status)
				reject.Invoke(js.ValueOf(fmt.Sprintf("Server returned non-OK stauts: ", res.Status)))
				return
			}

			res_byteSlice, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("Server side error: ", err.Error())
				reject.Invoke(js.ValueOf(err.Error()))
				return
			}

			resolve.Invoke(js.ValueOf(string(res_byteSlice)))
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}
