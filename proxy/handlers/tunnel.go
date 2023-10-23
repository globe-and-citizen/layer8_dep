package handlers

import (
	"io"
	"log"
	"net/http"
)

// Tunnel forwards the request to the service provider's backend
func Tunnel(w http.ResponseWriter, r *http.Request) {
	// get the service provider's backend url
	host := r.Header.Get("X-Forwarded-Host")
	scheme := r.Header.Get("X-Forwarded-Proto")
	url := scheme + "://" + host + r.URL.Path

	log.Println("Forwarding request to:", url)

	// create the request
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		log.Println("Error creating request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// add headers
	req.Header.Add("X-Layer8-Proxy", "true")
	for k, v := range r.Header {
		req.Header[k] = v
	}
	// send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Received response from:", url, res.StatusCode)

	// copy response
	for k, v := range res.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}
