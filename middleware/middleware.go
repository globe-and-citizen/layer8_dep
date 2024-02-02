package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	// "globe-and-citizen/layer8/utils" (Dep)
	"syscall/js"

	utils "github.com/globe-and-citizen/layer8-utils"
)

const VERSION = "1.0.3"

var (
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
	js.Global().Set("ServeStatic", js.FuncOf(static))
	js.Global().Set("ProcessMultipart", js.FuncOf(multipart))
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
		return
	}

	clientUUID := headers.Get("x-client-uuid").String()

	ss, err := privKey_ECDH.GetECDHSharedSecret(userPubJWKConverted)
	if err != nil {
		fmt.Println("Unable to get ECDH shared secret", err.Error())
		return
	}

	fmt.Println("shared secret: ", ss)
	// spSymmetricKey = ss

	UUIDMapOfKeys = append(UUIDMapOfKeys, map[string]*utils.JWK{clientUUID: ss})

	ss_b64, err := ss.ExportAsBase64()
	if err != nil {
		fmt.Println("Unable to export shared secret as base64", err.Error())
		return
	}

	MpJWT := headers.Get("mp_jwt").String()
	fmt.Println("MpJWT at SP BE (Middleware): ", MpJWT)

	UUIDMapOfJWTs = append(UUIDMapOfJWTs, map[string]string{clientUUID: MpJWT})

	// jsBody := request.Get("body")
	// if jsBody.String() == "<undefined>" {
	// 	println("body not defined")
	// 	return
	// }

	// object := js.Global().Get("JSON").Call("parse", jsBody)

	// data := object.Get("data").String()
	// fmt.Println("data: ", data)

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

// 	data, err := base64.URLEncoding.DecodeString(jsBody.Get("data").String())
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
// 			"data": base64.URLEncoding.EncodeToString(b),
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

	// continue to next middleware/handler if it's a request for static files
	if headers.Get("x-static").String() != "<undefined>" || headers.Get("X-Static").String() != "<undefined>" {
		next.Invoke()
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

	var body string

	req.Call("on", "data", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		body += args[0].Call("toString").String()
		return nil
	}))

	req.Call("on", "end", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// parse body and decrypt the "data" field
		var enc map[string]interface{}
		json.Unmarshal([]byte(body), &enc)

		data, err := base64.URLEncoding.DecodeString(enc["data"].(string))
		if err != nil {
			fmt.Println("error decoding request:", err.Error())
			res.Set("statusText", "Could not decode request: "+err.Error())
			res.Set("statusCode", 500)
			return nil
		}

		b, err := spSymmetricKey.SymmetricDecrypt(data)
		if err != nil {
			fmt.Println("error decrypting request:", err.Error())
			res.Set("statusText", "Could not decrypt request: "+err.Error())
			res.Set("statusCode", 500)
			return nil
		}

		// parse the decrypted data into a request object
		jreq, err := utils.FromJSONRequest(b)
		if err != nil {
			fmt.Println("error serializing json request:", err.Error())
			res.Set("statusText", "Could not decode request: "+err.Error())
			res.Set("statusCode", 500)
			return nil
		}

		switch strings.ToLower(jreq.Headers["Content-Type"]) {
		case "application/layer8.buffer+json": // this is used for multipart/form-data
			var (
				reqBody  map[string]interface{}
				formData = js.Global().Get("FormData").New()
			)

			json.Unmarshal(jreq.Body, &reqBody)

			randomBytes := make([]byte, 16)
			_, err = rand.Read(randomBytes)
			if err != nil {
				fmt.Println("error generating random bytes:", err.Error())
				res.Set("statusCode", 500)
				res.Set("statusMessage", "Could not generate random bytes: "+err.Error())
				return nil
			}
			boundary := fmt.Sprintf("----Layer8FormBoundary%s", base64.StdEncoding.EncodeToString(randomBytes))

			for k, v := range reqBody {
				// formdata can have multiple entries with the same key
				// that is why each key from the interceptor is a slice
				// of maps containing all the values for that key
				// hence the O(n^2) complexity (i.e. 2 for loops)
				for _, val := range v.([]interface{}) {
					val := val.(map[string]interface{})

					switch val["_type"].(string) {
					case "File":
						buff, err := base64.StdEncoding.DecodeString(val["buff"].(string))
						if err != nil {
							fmt.Println("error decoding file buffer:", err.Error())
							res.Set("statusCode", 500)
							res.Set("statusMessage", "Could not decode file buffer: "+err.Error())
							return nil
						}

						// converting the byte array to a uint8array so that it can be sent to the next
						// handler as a file object
						uInt8Array := js.Global().Get("Uint8Array").New(val["size"].(float64))
						js.CopyBytesToJS(uInt8Array, buff)

						file := js.Global().Get("File").New(
							[]interface{}{uInt8Array},
							val["name"].(string),
							map[string]interface{}{"type": val["type"].(string)},
						)
						formData.Call("append", k, file)
					case "String":
						formData.Call("append", k, val["value"].(string))
					case "Number":
						formData.Call("append", k, val["value"].(float64))
					case "Boolean":
						formData.Call("append", k, val["value"].(bool))
					}
				}
			}

			headers.Set("Content-Type", "multipart/form-data; boundary="+boundary)
			req.Set("body", formData)
		default:
			var reqBody map[string]interface{}
			json.Unmarshal(jreq.Body, &reqBody)

			req.Set("body", reqBody)
			headers.Set("Content-Type", "application/json")
		}

		// set the method and headers
		req.Set("method", jreq.Method)
		for k, v := range jreq.Headers {
			if strings.ToLower(k) == "content-type" {
				continue
			}
			headers.Set(k, v)
		}

		// continue to next middleware/handler
		next.Invoke()
		return nil
	}))

	// OVERWRITE THE SEND FUNCTION
	res.Set("send", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			data = args[0]
			b    []byte
			err  error
		)

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

		b, err = spSymmetricKey.SymmetricEncrypt(b)
		if err != nil {
			println("error encrypting response:", err.Error())
			res.Set("statusCode", 500)
			res.Set("statusMessage", "Could not encrypt response")
			return nil
		}

		// Send response
		res.Set("statusCode", jres.Status)
		res.Set("statusMessage", jres.StatusText)
		res.Call("set", js.ValueOf(map[string]interface{}{
			"content-type": "application/json",
			"mp_JWT":       MpJWT,
		}))
		res.Call("end", js.Global().Get("JSON").Call("stringify", js.ValueOf(map[string]interface{}{
			"data": base64.URLEncoding.EncodeToString(b),
		})))

		return nil
	}))

	return nil
}

func static(this js.Value, args []js.Value) interface{} {
	var (
		req     = args[0]
		res     = args[1]
		dir     = args[2].String()
		fs      = args[3]
		headers = req.Get("headers")

		// returns the default EncryptedImageData
		returnEncryptedImage = func() interface{} {
			arrayBuffer := js.Global().Get("Uint8Array").New(len(EncryptedImageData))
			js.CopyBytesToJS(arrayBuffer, EncryptedImageData)

			res.Set("statusCode", 200)
			res.Set("statusMessage", "OK")
			res.Set("content-type", "image/png")
			res.Call("end", arrayBuffer)
			return nil
		}
	)

	// get the file path
	path := req.Get("url").String()
	if path == "/" {
		path = "/index.html"
	}

	path, err := url.QueryUnescape(path)
	if err != nil {
		println("error url decoding path:", err.Error())
		res.Set("statusCode", 500)
		res.Set("statusMessage", "Internal Server Error")
		res.Call("end", "500 Internal Server Error")
		return nil
	}

	path = dir + path
	exists := fs.Call("existsSync", path).Bool()
	if !exists {
		res.Set("statusCode", 404)
		res.Set("statusMessage", "Not Found")
		res.Call("end", "Cannot GET "+req.Get("url").String())
		return nil
	}

	// return the default EncryptedImageData if the request is not a layer8 request
	if headers.String() == "<undefined>" || headers.Get("x-tunnel").String() == "<undefined>" {
		return returnEncryptedImage()
	}

	clientUUID := headers.Get("x-client-uuid").String()
	if clientUUID == "<undefined>" {
		return returnEncryptedImage()
	}

	var mpJWT string
	for _, v := range UUIDMapOfJWTs {
		if v[clientUUID] != "" {
			mpJWT = v[clientUUID]
		}
	}

	var sym *utils.JWK
	for _, v := range UUIDMapOfKeys {
		if v[clientUUID] != nil {
			sym = v[clientUUID]
		}
	}
	if sym == nil {
		return returnEncryptedImage()
	}

	// read the file
	buffer := fs.Call("readFileSync", path)
	b := make([]byte, buffer.Get("length").Int())
	js.CopyBytesToGo(b, buffer)

	// create a response object
	jres := utils.Response{
		Body:       b,
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Headers: map[string]string{
			"content-type": http.DetectContentType(b),
		},
	}

	b, err = jres.ToJSON()
	if err != nil {
		println("error serializing json response:", err.Error())
		res.Set("statusCode", 500)
		res.Set("statusMessage", "Internal Server Error")
		res.Call("end", "500 Internal Server Error")
		return nil
	}

	// encrypt the file
	encrypted, err := sym.SymmetricEncrypt(b)
	if err != nil {
		println("error encrypting file:", err.Error())
		res.Set("statusCode", 500)
		res.Set("statusMessage", "Internal Server Error")
		res.Call("end", "500 Internal Server Error")
		return nil
	}

	// send the response
	res.Set("statusCode", jres.Status)
	res.Set("statusMessage", jres.StatusText)
	res.Call("set", js.ValueOf(map[string]interface{}{
		"content-type": "application/json",
		"mp_JWT":       mpJWT,
	}))
	res.Call("end", js.Global().Get("JSON").Call("stringify", js.ValueOf(map[string]interface{}{
		"data": base64.URLEncoding.EncodeToString(encrypted),
	})))
	return nil
}

func multipart(this js.Value, args []js.Value) interface{} {
	var (
		options = args[0]
		fs      = args[1]

		dest = options.Get("dest").String()
	)

	single := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			req   = args[0]
			next  = args[2]
			field = args[3].String()
		)

		if dest == "" {
			dest = "tmp"
		}
		dest = strings.Trim(dest, "/")

		// if the destination directory does not exist, create it
		if !fs.Call("existsSync", dest).Bool() {
			fs.Call("mkdirSync", dest, map[string]interface{}{"recursive": true})
		}

		body := req.Get("body")
		if body.String() == "<undefined>" {
			next.Invoke()
			return nil
		}

		file := body.Call("get", field)
		if file.String() == "<undefined>" {
			next.Invoke()
			return nil
		}

		// check that file has a File constructor
		if file.Get("constructor").Get("name").String() != "File" {
			next.Invoke()
			return nil
		}

		file.Call("arrayBuffer").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			uint8Array := js.Global().Get("Uint8Array").New(args[0])

			// write the file to the destination directory
			filePath := fmt.Sprintf("%s/%s", dest, file.Get("name").String())
			fs.Call("writeFileSync", filePath, uint8Array)

			// set the file to the request body
			req.Set("file", file)

			// continue to next middleware/handler
			next.Invoke()
			return nil
		}))

		return nil
	})

	array := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var (
			req   = args[0]
			next  = args[2]
			field = args[3].String()
		)

		if dest == "" {
			dest = "tmp"
		}
		dest = strings.Trim(dest, "/")

		// if the destination directory does not exist, create it
		if !fs.Call("existsSync", dest).Bool() {
			fs.Call("mkdirSync", dest, map[string]interface{}{"recursive": true})
		}

		body := req.Get("body")
		if body.String() == "<undefined>" {
			next.Invoke()
			return nil
		}

		files := body.Call("getAll", field)
		if files.String() == "<undefined>" {
			next.Invoke()
			return nil
		}

		// write the files to the destination directory
		fileObjs := []interface{}{}
		files.Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			file := args[0]
			index := args[1].Int()

			if file.Get("constructor").Get("name").String() != "File" {
				return nil
			}

			file.Call("arrayBuffer").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				uint8Array := js.Global().Get("Uint8Array").New(args[0])

				// write the file to the destination directory
				filePath := fmt.Sprintf("%s/%s", dest, file.Get("name").String())
				fs.Call("writeFileSync", filePath, uint8Array)

				// append the file to the fileObjs slice
				fileObjs = append(fileObjs, file)

				// if all the files have been written to the destination directory
				// set the files to the request body and continue to next middleware/handler
				if index == files.Get("length").Int()-1 {
					req.Set("files", js.ValueOf(fileObjs))
					next.Invoke()
				}
				return nil
			}))

			return nil
		}))

		return nil
	})

	return map[string]interface{}{
		"single": single,
		"array":  array,
	}
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
