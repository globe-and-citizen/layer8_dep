package handlers

import (
	"fmt"
	"globe-and-citizen/layer8/utils"
	"io"
	"net/http"
	"os"
)

// Tunnel forwards the request to the service provider's backend
func Tunnel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n*************")
	fmt.Println(r.Method) // > GET  | > POST
	fmt.Println(r.URL)    // (http://localhost:5000/api/v1 ) > /api/v1

	fmt.Println("r.Header: ", r.Header)
	// Get the up_JWT (999) from the request header
	upJWT := r.Header.Get("up_JWT")
	fmt.Println("upJWT: ", upJWT)

	// Verify the up_JWT (999)
	_, err := utils.VerifyStandardToken(upJWT, os.Getenv("UP_999_SECRET_KEY"))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	backendURL := fmt.Sprintf("http://localhost:8000%s", r.URL)

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

	// send the request
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("\nReceived response from 8000:", backendURL, " of code: ", res.StatusCode)

	// copy response back
	for k, v := range res.Header {
		w.Header()[k] = v
		//fmt.Println("header pairs from SP: ", k, v)
	}

	w.Header()["setme"] = []string{"string"}
	w.Header().Add("ME TOO?", "DO IT!")
	fmt.Println("w.Headers: ", w.Header())
	//w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)

	fmt.Println("w.Headers 2: ", w.Header())
}
