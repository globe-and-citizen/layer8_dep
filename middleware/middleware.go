package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"

	// "globe-and-citizen/layer8/utils" (Dep)
	"syscall/js"

	utils "github.com/globe-and-citizen/layer8-utils"
)

const VERSION = "1.0.3"

// Unused
// type KeyPair struct {
// 	PublicKey  crypto.PublicKey
// 	PrivateKey crypto.PrivateKey
// }

var (
	// InstanceKey    *KeyPair   // Unused
	// spSymmetricKey *utils.JWK // TODO: Array of *utils.JWK (Symmetric keys)
	// MpJWT          string // TODO: Array of strings (JWTs)
	privKey_ECDH  *utils.JWK
	pubKey_ECDH   *utils.JWK
	UUIDMapOfKeys []map[string]*utils.JWK
	UUIDMapOfJWTs []map[string]string
)

func init() {
	var err error
	// generate key pair
	privKey_ECDH, pubKey_ECDH, err = utils.GenerateKeyPair(utils.ECDH)
	if err != nil {
		panic(err)
	}
}

func main() {
	c := make(chan struct{})
	fmt.Printf("L8 WASM Middleware version %s loaded.\n\n", VERSION)
	js.Global().Set("WASMMiddleware", js.FuncOf(WASMMiddleware_v2))
	js.Global().Set("TestWASM", js.FuncOf(TestWASM))
	<-c
}

func doECDHWithClient(request, response js.Value) {
	fmt.Println("TOP: ", request)
	headers := request.Get("headers")
	fmt.Println("headers: ", headers)
	userPubJWK := headers.Get("x-ecdh-init").String()
	// fmt.Println("userPubJWK: ", userPubJWK)
	userPubJWKConverted, err := utils.B64ToJWK(userPubJWK)
	if err != nil {
		fmt.Println("Failure to decode userPubJWK", err.Error())
		// response set "statusCode", 50x
		// response set "statusMessage", "err.Error()"
		return
	}

	clientUUID := headers.Get("x-client-uuid").String()
	// For testing purposes
	fmt.Println("clientUUID: ", clientUUID)

	ss, err := privKey_ECDH.GetECDHSharedSecret(userPubJWKConverted)
	if err != nil {
		fmt.Println("Unable to get ECDH shared secret", err.Error())
		// response set "statusCode", 50x
		// response set "statusMessage", "err.Error()"
		return
	}

	fmt.Println("shared secret: ", ss)
	// spSymmetricKey = ss

	UUIDMapOfKeys = append(UUIDMapOfKeys, map[string]*utils.JWK{clientUUID: ss})

	ss_b64, err := ss.ExportAsBase64()
	if err != nil {
		fmt.Println("Unable to export shared secret as base64", err.Error())
		// response set "statusCode", 50x
		// response set "statusMessage", "err.Error()"
		return
	}

	// Get mp_JWT from headers
	// Use a uuid
	MpJWT := headers.Get("mp_jwt").String()
	fmt.Println("MpJWT at SP BE (Middleware): ", MpJWT)

	UUIDMapOfJWTs = append(UUIDMapOfJWTs, map[string]string{clientUUID: MpJWT})

	response.Set("send", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// encrypt response
		jres := utils.Response{}
		jres.Body = []byte(ss_b64)
		jres.Status = 200
		jres.StatusText = "ECDH Successfully Completed!"
		// jres.Headers = make(map[string]string)
		// jres.Headers["x-shared-secret"] = ss_b64

		if err != nil {
			println("error serializing json response:", err.Error())
			response.Set("statusCode", 500)
			response.Set("statusMessage", "Failure to encode ECDH init response")
			return nil
		}

		// send response
		response.Set("statusCode", jres.Status)
		response.Set("statusMessage", jres.StatusText)

		server_pubKeyECDH, _ := pubKey_ECDH.ExportAsBase64()

		response.Call("end", server_pubKeyECDH)
		fmt.Println("SS_Server: ", ss)
		return nil
	}))

	// Send the response back to the user.
	response.Call("setHeader", "x-shared-secret", ss_b64)
	// response.Call("setHeader", "mp_JWT", MpJWT)
	result := response.Call("hasHeader", "x-shared-secret")
	fmt.Println("result: ", result)
	response.Call("send")
	return
}

// func WASMMiddleware(this js.Value, args []js.Value) interface{} {
// 	// get the request and response objects and the next function
// 	req := args[0]
// 	res := args[1]
// 	next := args[2]

// 	//fmt.Println("Get Rdy to Error out...")

// 	// check for layer8 request
// 	headers := req.Get("headers")
// 	if headers.String() == "<undefined>" {
// 		next.Invoke()
// 		return nil
// 	}

// 	isECDHInit := headers.Get("x-ecdh-init").String()
// 	if isECDHInit != "<undefined>" {
// 		doECDHWithClient(req, res)
// 		return nil
// 	}

// 	// get the body. This depends on the express.json
// 	jsBody := req.Get("body")
// 	if jsBody.String() == "<undefined>" {
// 		println("body not defined")
// 		res.Set("statusCode", 400)
// 		res.Set("statusMessage", "Invalid request")
// 		return nil
// 	}

// 	data, err := base64.StdEncoding.DecodeString(jsBody.Get("data").String())
// 	if err != nil {
// 		println("error decoding request:", err.Error())
// 		res.Set("statusCode", 500)
// 		res.Set("statusMessage", "Internal server error")
// 		return nil
// 	}

// 	b, err := spSymmetricKey.SymmetricDecrypt(data)
// 	if err != nil {
// 		println("error decrypting request:", err.Error())
// 		res.Set("statusCode", 400)
// 		res.Set("statusMessage", "Could not decrypt request")
// 		return nil
// 	}

// 	jreq, err := utils.FromJSONRequest(b)
// 	if err != nil {
// 		println("error serializing json request:", err.Error())
// 		res.Set("statusCode", 400)
// 		res.Set("statusMessage", "Could not decode request")
// 		return nil
// 	}

// 	req.Set("method", jreq.Method)
// 	for k, v := range jreq.Headers {
// 		headers.Set(k, v)
// 	}
// 	var reqBody map[string]interface{}
// 	json.Unmarshal(jreq.Body, &reqBody)
// 	req.Set("body", reqBody)

// 	// OVERWRITE THE SEND FUNCTION
// 	res.Set("send", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		// convert body to readable format
// 		data := args[0]
// 		var b []byte

// 		if data.Type() == js.TypeObject {
// 			switch data.Get("constructor").Get("name").String() {
// 			case "Object":
// 				b, err = json.Marshal(parseJSObjectToMap(data))
// 				if err != nil {
// 					println("error serializing json response:", err.Error())
// 					res.Set("statusCode", 500)
// 					res.Set("statusMessage", "Could not encode response")
// 					return nil
// 				}
// 			case "Array":
// 				b, err = json.Marshal(parseJSObjectToSlice(data))
// 				if err != nil {
// 					println("error serializing json response:", err.Error())
// 					res.Set("statusCode", 500)
// 					res.Set("statusMessage", "Could not encode response")
// 					return nil
// 				}
// 			default:
// 				b = []byte(data.String())
// 			}
// 		} else {
// 			b = []byte(data.String())
// 		}

// 		// encrypt response
// 		jres := utils.Response{}
// 		jres.Body = b
// 		jres.Status = res.Get("statusCode").Int()
// 		jres.StatusText = res.Get("statusMessage").String()
// 		jres.Headers = make(map[string]string)
// 		if res.Get("headers").String() == "<undefined>" {
// 			res.Set("headers", js.ValueOf(map[string]interface{}{}))
// 		}
// 		js.Global().Get("Object").Call("keys", res.Get("headers")).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 			jres.Headers[args[0].String()] = args[1].String()
// 			return nil
// 		}))
// 		b, err = jres.ToJSON()
// 		if err != nil {
// 			println("error serializing json response:", err.Error())
// 			res.Set("statusCode", 500)
// 			res.Set("statusMessage", "Could not encode response")
// 			return nil
// 		}

// 		//b, err = utils.Dep_SymmetricEncrypt(b, secret)
// 		b, err := spSymmetricKey.SymmetricEncrypt(b)
// 		//fmt.Println("b: ", b)
// 		if err != nil {
// 			println("error encrypting response:", err.Error())
// 			res.Set("statusCode", 500)
// 			res.Set("statusMessage", "Could not encrypt response")
// 			return nil
// 		}

// 		resHeaders := make(map[string]interface{})
// 		for k, v := range jres.Headers {
// 			resHeaders[k] = v
// 		}

// 		// send response
// 		res.Set("statusCode", jres.Status)
// 		res.Set("statusMessage", jres.StatusText)
// 		res.Call("set", js.ValueOf(resHeaders))
// 		res.Call("end", js.Global().Get("JSON").Call("stringify", js.ValueOf(map[string]interface{}{
// 			"data": base64.StdEncoding.EncodeToString(b),
// 		})))
// 		return nil
// 	}))

// 	// continue to next middleware/handler
// 	next.Invoke()
// 	return nil
// }

// WASM Middleware Version 2 Does not depend on the Express Body Parser//

func WASMMiddleware_v2(this js.Value, args []js.Value) interface{} {
	// Get the request and response objects and the next function
	req := args[0]
	res := args[1]
	next := args[2]

	headers := req.Get("headers")

	// proceed to next middleware/handler request is not a layer8 request
	if headers.String() == "<undefined>" || headers.Get("x-tunnel").String() == "<undefined>" {
		next.Invoke()
		return nil
	}

	// Decide if this is a redirect to ECDH init.
	isECDHInit := headers.Get("x-ecdh-init").String()
	if isECDHInit != "<undefined>" {
		doECDHWithClient(req, res)
		return nil
	}

	clientUUID := headers.Get("x-client-uuid").String()
	fmt.Println("clientUUID: ", clientUUID)
	if clientUUID == "<undefined>" {
		doECDHWithClient(req, res)
		return nil
	}

	// Get the symmetric key for this client
	var spSymmetricKey *utils.JWK
	for _, v := range UUIDMapOfKeys {
		if v[clientUUID] != nil {
			spSymmetricKey = v[clientUUID]
		}
	}
	if spSymmetricKey == nil {
		doECDHWithClient(req, res)
		return nil
	}

	// Get the JWT for this client
	var MpJWT string
	for _, v := range UUIDMapOfJWTs {
		if v[clientUUID] != "" {
			MpJWT = v[clientUUID]
		}
	}
	if MpJWT == "" {
		doECDHWithClient(req, res)
		return nil
	}

	// Parse the data into an object.
	// Note, this must take place inside a callback.
	req.Call("on", "data", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var uint8array []byte
		js.Global().Get("Object").Call("values", args[0]).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			theInt := args[0].Int()
			uint8array = binary.AppendUvarint(uint8array, uint64(theInt))
			return nil
		}))

		req.Set("body", js.ValueOf(string(uint8array)))

		jsBody := req.Get("body")
		if jsBody.String() == "<undefined>" {
			println("body not defined")
			res.Set("statusCode", 400)
			res.Set("statusMessage", "Invalid request")
			return nil
		}

		object := js.Global().Get("JSON").Call("parse", jsBody)

		data, err := base64.StdEncoding.DecodeString(object.Get("data").String())
		if err != nil {
			println("error decoding request:", err.Error())
			res.Set("statusCode", 500)
			res.Set("statusMessage", "Internal server error")
			return nil
		}

		b, err := spSymmetricKey.SymmetricDecrypt(data)
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

		// Headers not working

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

			// Encrypt response
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

			b, err := spSymmetricKey.SymmetricEncrypt(b)

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
			resHeaders["mp_JWT"] = MpJWT
			// Send response
			res.Set("statusCode", jres.Status)
			res.Set("statusMessage", jres.StatusText)
			res.Call("set", js.ValueOf(resHeaders))
			res.Call("end", js.Global().Get("JSON").Call("stringify", js.ValueOf(map[string]interface{}{
				"data": base64.StdEncoding.EncodeToString(b),
			})))
			return nil
		}))

		// Continue to next middleware/handler
		next.Invoke()
		return nil
	}))
	return nil
}

// UTILS
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
