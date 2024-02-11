package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"globe-and-citizen/layer8/server/resource_server/utils"

	utilities "github.com/globe-and-citizen/layer8-utils"
)

// Tunnel forwards the request to the service provider's backend
func InitTunnel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n*************")
	fmt.Println(r.Method) // > GET  | > POST
	fmt.Println(r.URL)    // (http://localhost:5000/api/v1 ) > /api/v1
	params := r.URL.Query()
	var backend string
	if _, ok := params["backend"]; ok != true {
		res := utils.BuildErrorResponse("Failed to get User. Malformed query string.", "", utils.EmptyObj{})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	} else {
		backend = params["backend"][0]
	}

	mpJWT, err := utilities.GenerateStandardToken(os.Getenv("MP_123_SECRET_KEY"))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	reqData := utilities.ReadResponseBody(r.Body)
	b64PubJWK := string(reqData)
	fmt.Println("b64PubJWK: ", b64PubJWK)
	fmt.Println("x-ecdh-init: ", r.Header.Get("x-ecdh-init"))

	backendURL := fmt.Sprintf("http://%s", backend)
	fmt.Println("User agent is attempting to initialize this backend SP: ", backendURL)

	// create the request
	req, err := http.NewRequest(r.Method, backendURL, r.Body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add headers
	for k, v := range r.Header {
		req.Header[k] = v
	}

	req.Header["x-tunnel"] = []string{"true"}
	req.Header["mp-jwt"] = []string{mpJWT}

	// send the request
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make a buffer to hold response body
	var resBodyTemp bytes.Buffer

	// Copy the response body to buffer
	resBodyTemp.ReadFrom(res.Body)

	// Convert resBodyTemp to []byte

	resBodyTempBytes := []byte(resBodyTemp.String())

	// Make a copy of the response body to send back to client
	res.Body = io.NopCloser(bytes.NewBuffer(resBodyTemp.Bytes()))

	fmt.Println("\nReceived response from 8000:", backendURL, " of code: ", res.StatusCode)

	upJWT, err := utilities.GenerateStandardToken(os.Getenv("UP_999_SECRET_KEY"))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server_pubKeyECDH, err := utilities.B64ToJWK(string(resBodyTempBytes))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// // Make a json response of server_pubKeyECDH and up_JWT and send it back to client
	data := map[string]interface{}{
		"server_pubKeyECDH": server_pubKeyECDH,
		"up_JWT":            upJWT,
	}

	fmt.Println("data (Going back to client): ", data)

	datatoSend, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(datatoSend)

}

func Tunnel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n*************")
	fmt.Println(r.Method) // > GET  | > POST
	fmt.Println(r.URL)    // (http://localhost:5000/api/v1 ) > /api/v1
	fmt.Println("Ravi Adds Path: ", r.URL.Path)
	fmt.Println("Host:", r.Header.Get("X-Forwarded-Host"))

	// backendURL := fmt.Sprintf(os.Getenv("VITE_BACKEND")+"%s", r.URL)
	backendURL := fmt.Sprintf("http://%s", r.Header.Get("X-Sp-Backend-Url")+r.URL.Path)

	// create the request
	req, err := http.NewRequest(r.Method, backendURL, r.Body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// add headers
	for k, v := range r.Header {
		req.Header[k] = v
		fmt.Println("header pairs from client (Interceptor): ", k, v)
	}
	req.Header["x-tunnel"] = []string{"true"}

	// Get up_JWT from request header and verify it
	upJWT := r.Header.Get("up-jwt") // RAVI!
	fmt.Println("up-jwt coming from client: ", upJWT)

	_, err = utilities.VerifyStandardToken(upJWT, os.Getenv("UP_999_SECRET_KEY"))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// send the request
	res, err := http.DefaultClient.Do(req) // Source of MapB Error

	if err != nil {
		fmt.Println("Error sending request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("\nReceived response from:", backendURL, " of code: ", res.StatusCode)

	// Get mp_JWT from response header and verify it
	mpJWT := res.Header.Get("mp-jwt")
	fmt.Println("mp_jwt coming from SP: ", mpJWT)

	_, err = utilities.VerifyStandardToken(mpJWT, os.Getenv("MP_123_SECRET_KEY"))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// copy response back
	for k, v := range res.Header {
		w.Header()[k] = v
		fmt.Println("header pairs from SP: ", k, v)
	}

	w.Header()["setme"] = []string{"string"}
	w.Header().Add("ME TOO?", "DO IT!")
	fmt.Println("w.Headers: ", w.Header())
	//w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)

	fmt.Println("w.Headers 2: ", w.Header())
}

func TestError(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("This is a test error.")
	fmt.Println("Test error endpoint:", err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
