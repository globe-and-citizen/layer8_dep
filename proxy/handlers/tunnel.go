package handlers

import (
	"fmt"
	"io"
	"net/http"
)

// Tunnel forwards the request to the service provider's backend
func Tunnel(w http.ResponseWriter, r *http.Request) {
	// get the service provider's backend url
	host := r.Header.Get("X-Forwarded-Host")
	scheme := r.Header.Get("X-Forwarded-Proto")
	url := scheme + "://" + host + r.URL.Path

	fmt.Println("\n\n*************")
	fmt.Println(r.Method) // > GET  | > POST
	fmt.Println(r.URL)    // (http://localhost:5000/api/v1 ) > /api/v1

	// backendURL := fmt.Sprintf("%s://%s/%s", r.URL.Scheme, r.URL.Host, r.URL.Path)
	println(url)

	// create the request
	req, err := http.NewRequest(r.Method, url, r.Body)
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

	fmt.Println("\nReceived response from 8000:", url, " of code: ", res.StatusCode)

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
