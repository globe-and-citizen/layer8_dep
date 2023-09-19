package main

import (
	"fmt"
	"syscall/js"
)

const VERSION = "1.0.1"

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("testWASM", js.FuncOf(testWASM))
	// js.Global().Set("", js.FuncOf())
	// js.Global().Set("", js.FuncOf())
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
