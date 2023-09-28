package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"syscall/js"
)

const VERSION = "1.0.2"

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("testWASM", js.FuncOf(testWASM))
	js.Global().Set("pingProxy", js.FuncOf(pingProxy))
	js.Global().Set("genericPost", js.FuncOf(genericPost))
	// js.Global().Set("", js.FuncOf())
	fmt.Println("WASM interceptor")
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

func pingProxy(this js.Value, args []js.Value) interface{} {
	url := args[0]
	fmt.Println("HERE: ", url.String())
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			// Main function body
			res, err := http.Get(url.String())
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

// func genericGET() interface{}{
// }

func genericPost(this js.Value, args []js.Value) interface{} {
	url := args[0]
	body_text_as_str := args[1].String()
	body_as_ioReader := strings.NewReader(body_text_as_str)
	fmt.Println("Interceptor will now POST to this url: ", url.String())
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]
		go func() {
			// Main function body
			res, err := http.Post(url.String(), "text/html", body_as_ioReader)
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
