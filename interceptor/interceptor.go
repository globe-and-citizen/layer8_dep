package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"syscall/js"

	utilities "github.com/globe-and-citizen/layer8-utils"
)

const VERSION = "1.0.2"

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
)

func main() {
	// keep the Go thread alive
	c := make(chan struct{}, 0)
	Layer8Scheme = "http"
	Layer8Host = "localhost"
	Layer8Port = "5000"
	Layer8Version = "1.0"

	// expose the layer8 functionality the global scope
	js.Global().Set("layer8", js.ValueOf(map[string]interface{}{
		"testWASM":          js.FuncOf(testWASM),
		"genericGetRequest": js.FuncOf(genericGetRequest),
		"genericPost":       js.FuncOf(genericPost),
		"fetch":             js.FuncOf(fetch),
	}))

	fmt.Println("WASM interceptor loaded.")

	// Wait indefinitely
	<-c
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

func fetch(this js.Value, args []js.Value) interface{} {
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
				res := &utilities.Response{
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
