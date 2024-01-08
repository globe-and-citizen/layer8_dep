// internals_test.go
package internals

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// js.Global().Get("Promise").New(js.FuncOf(func (this js.Value, resolve_reject []js.Value) interface{} {
// 	testWASM(this, []js.Value{resolve_reject[0], resolve_reject[1]})
// }))

// TestClientDo tests the Do method of the Client type
func TestClientDo(t *testing.T) {
	// from "GenerateStandardToken" in server/utils/utils.go
	genToken := func (secretKey string) (string, error) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		}
		token.Claims = claims
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return "", fmt.Errorf("could not generate standard token: %s", err)
		}
		return tokenString, nil
	}

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/init-tunnel":
			_, pub, err := utils.GenerateKeyPair(utils.ECDH)
			assert.NoError(t, err)
			assert.NotNil(t, pub)

			token, err := genToken("mock_secret")
			assert.NoError(t, err)
			assert.NotNil(t, token)

			data, err := json.Marshal(map[string]interface{}{
				"server_pubKeyECDH": pub,
				"up_JWT": token,
			})
			assert.NoError(t, err)

			w.WriteHeader(200)
			w.Write([]byte(base64.StdEncoding.EncodeToString(data)))
		case "/":
			// TODO: mock data transfer
		}
	}))
	defer mockServer.Close()

	// how to set Layer8Scheme, Layer8Host, Layer8Port??

	os.Setenv("LAYER8_SCHEME")

	mockRequest := utils.NewRequest("POST", map[string]string{
		"Content-Type": "application/json",
	}, []byte(`{"data": "mocked_request"}`))


	// client := &Client{
	// 	proxyURL: mockServer.URL,
	// }

	// response := client.Do("/path", mockRequest, mockSharedSecret)

	// if response.Status != 200 {
	// 	t.Errorf("Expected status code 200, got %d", response.Status)
	// }

	// expectedHeader1 := "ONE"
	// if response.Headers["x-custom-header-1"] != expectedHeader1 {
	// 	t.Errorf("Expected header x-custom-header-1 to be %s, got %s", expectedHeader1, response.Headers["x-custom-header-1"])
	// }

	// expectedData := "mocked_response"
	// if string(response.Data) != expectedData {
	// 	t.Errorf("Expected data to be %s, got %s", expectedData, string(response.Data))
	// }
}
