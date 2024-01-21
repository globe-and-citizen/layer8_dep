package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/xdg-go/pbkdf2"
)

const SaltSize = 32
const SecretSize = 32

var WorkingDirectory string

// Response is used for static shape json return
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type EmptyObj struct{}

func GetPwd() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	WorkingDirectory = dir
}

func GenerateRandomSalt(saltSize int) string {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(salt[:])
}

func SaltAndHashPassword(password string, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha1.New)
	return hex.EncodeToString(dk[:])
}

func CheckPassword(password string, salt string, hash string) bool {
	return SaltAndHashPassword(password, salt) == hash
}

func BuildResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Error:   nil,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	res := Response{
		Status:  false,
		Message: message,
		Error:   splittedError,
		Data:    data,
	}

	return res
}

func HandleError(w http.ResponseWriter, status int, message string, err error) {
	w.WriteHeader(status)
	res := BuildErrorResponse(message, err.Error(), EmptyObj{})
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func GenerateUUID() string {
	newUUID := uuid.New()

	return newUUID.String()
}
func GenerateSecret(secretSize int) string {
	var randomBytes = make([]byte, secretSize)

	_, err := rand.Read(randomBytes[:])

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(randomBytes[:])
}
