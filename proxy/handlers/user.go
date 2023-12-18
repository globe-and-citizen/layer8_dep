package handlers

import (
	"encoding/json"
	"globe-and-citizen/layer8/proxy/internals/usecases"
	"log"
	"net/http"
	"strings"
)

// UserInfo handles requests to get a user's anonymized data
// The last step of the oauth flow
func UserInfo(w http.ResponseWriter, r *http.Request) {
	usecase := r.Context().Value("usecase").(*usecases.UseCase)

	switch r.Method {
	case "GET":
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		data, err := usecase.AccessResourcesWithToken(token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid token"}`))
			return
		}
		b, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "server error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}
}