package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"syscall/js"
)

const VERSION = "1.0.2"

type Request struct {
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Printf("L8 WASM Middleware version %s loaded.\n\n", VERSION)
	js.Global().Set("testWASM", js.FuncOf(async_test_WASM))
	js.Global().Set("WASMMiddleware", js.FuncOf(WASMMiddleware))
	<-c
}

func WASMMiddleware(this js.Value, args []js.Value) interface{} {
	request := args[0]
	response := args[1]
	next := args[2]
	var dataMap map[string]interface{}

	request.Call("on", "data", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("chunck: ", args[0])
		var uint8array []byte

		js.Global().Get("Object").Call("values", args[0]).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			theInt := args[0].Int()
			uint8array = binary.AppendUvarint(uint8array, uint64(theInt))

			return nil
		}))

		myRequest := new(Request)
		err := json.Unmarshal(uint8array, &myRequest)
		if err != nil {
			fmt.Println("damn")
		}

		fmt.Println(myRequest)

		return nil
	}))

	url := request.Get("baseUrl").String()

	fmt.Println(url)

	// Set any layer8 particular custom props
	response.Set("custom_test_prop", js.ValueOf("Example string"))

	// Replace the standard methods with the L8 equivalents
	// so that end users don't notice any difference.
	nativeFunc := response.Get("send")
	response.Set("send", nativeFunc)
	fmt.Println("Request has transitted the WASM middleware.", request)

	next.Invoke()
	return nil
}

func async_test_WASM(this js.Value, args []js.Value) interface{} {
	fmt.Println("Fisrt argument: ", args[0])
	fmt.Println("Second argument: ", args[1])
	var resolve_reject_internals = func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		//reject := args[1]
		go func() {
			// Main function body
			//fmt.Println(string(args[2]))
			resolve.Invoke(js.ValueOf(fmt.Sprintf("WASM Middleware version %s successfully loaded.", VERSION)))
			//reject.Invoke()
		}()
		return nil
	}
	promiseConstructor := js.Global().Get("Promise")
	promise := promiseConstructor.New(js.FuncOf(resolve_reject_internals))
	return promise
}
