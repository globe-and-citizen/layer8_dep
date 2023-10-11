package main

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"syscall/js"

	utilities "github.com/globe-and-citizen/layer8-utils"
)

type KeyPair struct {
	PublicKey  crypto.PublicKey
	PrivateKey crypto.PrivateKey
}

var (
	InstanceKey      *KeyPair
	ClientPublicKeys map[string]crypto.PublicKey = make(map[string]crypto.PublicKey)
)

func init() {
	// generate key pair
	pri, pub, err := utilities.GenerateKeyPair(utilities.ECDH_ALGO)
	if err != nil {
		panic(err)
	}
	InstanceKey = &KeyPair{
		PublicKey:  pub,
		PrivateKey: pri,
	}
}

func main() {
	js.Global().Set("WASMMiddleware", js.FuncOf(middleware))
	<-make(chan bool)
}

func middleware(this js.Value, args []js.Value) interface{} {
	// get the request and response objects and the next function
	req := args[0]
	res := args[1]
	next := args[2]

	// check for layer8 request
	headers := req.Get("headers")
	if headers.String() == "<undefined>" {
		next.Invoke()
		return nil
	}
	isFromLayer8 := headers.Get("x-layer8-proxy").String()
	if isFromLayer8 == "<undefined>" {
		// continue for non-layer8 requests
		next.Invoke()
		return nil
	}
	
	// get the body
	jsBody := req.Get("body")
	if jsBody.String() == "<undefined>" {
		println("body not defined")
		res.Set("statusCode", 400)
		res.Set("statusMessage", "Invalid request")
		return nil
	}
	
	data, err := base64.StdEncoding.DecodeString(jsBody.Get("data").String())
	if err != nil {
		println("error decoding request:", err.Error())
		res.Set("statusCode", 500)
		res.Set("statusMessage", "Internal server error")
		return nil
	}
	
	// hardcoding a shared secret for now
	secret, err := base64.StdEncoding.DecodeString("KfbCmY2v83ptAZLLKffx0ve2Br8hkMhCkIo5RkFaNlk=")
	if err != nil {
		return &utilities.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}

	b, err := utilities.SymmetricDecrypt(data, secret)
	if err != nil {
		println("error decrypting request:", err.Error())
		res.Set("statusCode", 400)
		res.Set("statusMessage", "Could not decrypt request")
		return nil
	}
	jreq, err := utilities.FromJSONRequest(b)
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
	
	// overwrite the send function
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
		jres := utilities.Response{}
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

		b, err = utilities.SymmetricEncrypt(b, secret)
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
			"data": base64.StdEncoding.EncodeToString(b),
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
