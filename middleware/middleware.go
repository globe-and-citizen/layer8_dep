package main

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/utils"
	"syscall/js"
)

const VERSION = "1.0.3"

type KeyPair struct {
	PublicKey  crypto.PublicKey
	PrivateKey crypto.PrivateKey
}

var (
	InstanceKey *KeyPair
	//ServerJWKPair    *utils.JWK
	ClientPublicKeys map[string]crypto.PublicKey = make(map[string]crypto.PublicKey)
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

func init() {
	// generate key pair
	pri, pub, err := utils.GenerateKeyPair(utils.ECDH)
	if err != nil {
		panic(err)
	}
	InstanceKey = &KeyPair{
		PublicKey:  pub,
		PrivateKey: pri,
	}
}

type Request struct {
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
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

func TestWASM(this js.Value, args []js.Value) interface{} {
	fmt.Println("TestWasm Ran")
	return js.ValueOf("42")
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Printf("L8 WASM Middleware version %s loaded.\n\n", VERSION)
	js.Global().Set("WASMMiddleware", js.FuncOf(WASMMiddleware))
	js.Global().Set("TestWASM", js.FuncOf(TestWASM))
	<-c
}

func WASMMiddleware(this js.Value, args []js.Value) interface{} {
	// get the request and response objects and the next function
	req := args[0]
	res := args[1]
	next := args[2]

	//fmt.Println("Get Rdy to Error out...")

	// check for layer8 request
	headers := req.Get("headers")
	if headers.String() == "<undefined>" {
		next.Invoke()
		return nil
	}

	//isFromLayer8 := headers.Get("x-layer8-proxy").String()
	// if isFromLayer8 == "<undefined>" {
	// 	// continue for non-layer8 requests
	// 	next.Invoke()
	// 	return nil
	// }

	// get the body
	jsBody := req.Get("body")
	if jsBody.String() == "<undefined>" {
		println("body not defined")
		res.Set("statusCode", 400)
		res.Set("statusMessage", "Invalid request")
		return nil
	}

	data, err := base64.URLEncoding.DecodeString(jsBody.Get("data").String())
	if err != nil {
		println("error decoding request:", err.Error())
		res.Set("statusCode", 500)
		res.Set("statusMessage", "Internal server error")
		return nil
	}

	// Hardcoding a shared secret for now
	// secret, err := base64.StdEncoding.DecodeString("KfbCmY2v83ptAZLLKffx0ve2Br8hkMhCkIo5RkFaNlk=")
	// if err != nil {
	// 	return &utils.Response{
	// 		Status:     500,
	// 		StatusText: err.Error(),
	// 	}
	// }
	// b, err := utils.Dep_SymmetricDecrypt(data, secret)

	b, err := spoofedSymmetricKey.SymmetricDecrypt(data)
	if err != nil {
		println("error decrypting request:", err.Error())
		res.Set("statusCode", 400)
		res.Set("statusMessage", "Could not decrypt request")
		return nil
	}

	jreq, err := utils.FromJSONRequest(b)
	if err != nil {
		println("error serializing json request:", err.Error())
		res.Set("statusCode", 400)
		res.Set("statusMessage", "Could not decode request")
		return nil
	}

	req.Set("method", jreq.Method)
	for k, v := range jreq.Headers {
		headers.Set(k, v)
	}
	var reqBody map[string]interface{}
	json.Unmarshal(jreq.Body, &reqBody)
	req.Set("body", reqBody)

	// OVERWRITE THE SEND FUNCTION
	res.Set("send", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// convert body to readable format
		data := args[0]
		var b []byte

		if data.Type() == js.TypeObject {
			switch data.Get("constructor").Get("name").String() {
			case "Object":
				b, err = json.Marshal(parseJSObjectToMap(data))
				if err != nil {
					println("error serializing json response:", err.Error())
					res.Set("statusCode", 500)
					res.Set("statusMessage", "Could not encode response")
					return nil
				}
			case "Array":
				b, err = json.Marshal(parseJSObjectToSlice(data))
				if err != nil {
					println("error serializing json response:", err.Error())
					res.Set("statusCode", 500)
					res.Set("statusMessage", "Could not encode response")
					return nil
				}
			default:
				b = []byte(data.String())
			}
		} else {
			b = []byte(data.String())
		}

		// encrypt response
		jres := utils.Response{}
		jres.Body = b
		jres.Status = res.Get("statusCode").Int()
		jres.StatusText = res.Get("statusMessage").String()
		jres.Headers = make(map[string]string)
		if res.Get("headers").String() == "<undefined>" {
			res.Set("headers", js.ValueOf(map[string]interface{}{}))
		}
		js.Global().Get("Object").Call("keys", res.Get("headers")).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			jres.Headers[args[0].String()] = args[1].String()
			return nil
		}))
		b, err = jres.ToJSON()
		if err != nil {
			println("error serializing json response:", err.Error())
			res.Set("statusCode", 500)
			res.Set("statusMessage", "Could not encode response")
			return nil
		}

		//b, err = utils.Dep_SymmetricEncrypt(b, secret)
		b, err := spoofedSymmetricKey.SymmetricEncrypt(b)
		//fmt.Println("b: ", b)
		if err != nil {
			println("error encrypting response:", err.Error())
			res.Set("statusCode", 500)
			res.Set("statusMessage", "Could not encrypt response")
			return nil
		}

		resHeaders := make(map[string]interface{})
		for k, v := range jres.Headers {
			resHeaders[k] = v
		}

		// send response
		res.Set("statusCode", jres.Status)
		res.Set("statusMessage", jres.StatusText)
		res.Call("set", js.ValueOf(resHeaders))
		res.Call("end", js.Global().Get("JSON").Call("stringify", js.ValueOf(map[string]interface{}{
			"data": base64.URLEncoding.EncodeToString(b),
		})))
		return nil
	}))

	// continue to next middleware/handler
	next.Invoke()
	return nil
}

func parseJSObjectToMap(obj js.Value) map[string]interface{} {
	m := map[string]interface{}{}

	keys := js.Global().Get("Object").Call("keys", obj)
	for i := 0; i < keys.Length(); i++ {
		key := keys.Index(i).String()
		val := obj.Get(key)

		switch val.Type() {
		case js.TypeNumber:
			m[key] = val.Float()
		case js.TypeBoolean:
			m[key] = val.Bool()
		case js.TypeString:
			m[key] = val.String()
		case js.TypeObject:
			if val.Get("constructor").Get("name").String() == "Array" {
				m[key] = parseJSObjectToSlice(val)
				continue
			}
			m[key] = parseJSObjectToMap(val)
		}
	}

	return m
}

func parseJSObjectToSlice(obj js.Value) []interface{} {
	var s []interface{}

	for i := 0; i < obj.Length(); i++ {
		val := obj.Index(i)

		switch val.Type() {
		case js.TypeNumber:
			s = append(s, val.Float())
		case js.TypeBoolean:
			s = append(s, val.Bool())
		case js.TypeString:
			s = append(s, val.String())
		case js.TypeObject:
			if val.Get("constructor").Get("name").String() == "Array" {
				s = append(s, parseJSObjectToSlice(val))
				continue
			}
			s = append(s, parseJSObjectToMap(val))
		}
	}

	return s
}

func Dep_Ravi_WASMMiddleware(this js.Value, args []js.Value) interface{} {
	//request := args[0]
	//response := args[1]
	next := args[2]

	fmt.Println("closer....")

	//request.Call("on", "data", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// var uint8array []byte

	// js.Global().Get("Object").Call("values", args[0]).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	theInt := args[0].Int()
	// 	uint8array = binary.AppendUvarint(uint8array, uint64(theInt))
	// 	return nil
	// }))

	// // // CHOICE: set the body of the request to be a JSON string
	// request.Set("body", js.ValueOf(string(uint8array)))

	// // // OR: process the request within the WASM module
	// // myRequest := new(Request)
	// // err := json.Unmarshal(uint8array, &myRequest)
	// // if err != nil {
	// // 	fmt.Println("damn")
	// // }
	// // fmt.Println("Method: ", string(myRequest.Method))
	// // fmt.Println("Headers: ", myRequest.Headers)
	// // fmt.Println("Body: ", string(myRequest.Body))
	// // url := request.Get("baseUrl").String()
	// // fmt.Println(url)

	// // Set any layer8 particular custom props
	// response.Set("custom_test_prop", js.ValueOf("Example string"))

	// // Replace the standard methods with the L8 equivalents
	// // so that end users don't notice any difference.
	// nativeFunc := response.Get("send")
	// response.Set("send", nativeFunc)
	// fmt.Println("Request has transitted the WASM middleware.", request)

	// 	next.Invoke()
	// 	return nil
	// }))
	next.Invoke()
	return nil
}
