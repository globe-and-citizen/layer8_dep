package internals

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	utils "github.com/globe-and-citizen/layer8-utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClientDo(t *testing.T) {
	// from "GenerateStandardToken" in server/utils/utils.go
	genToken := func(secretKey string) (string, error) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		}
		token.Claims = claims
		tokenString, err := token.SignedString([]byte(secretKey))
		assert.NoError(t, err)
		return tokenString, nil
	}

	var (
		RequestMethod  = "GET"
		RequestURL     = "https://test.layer8.com/test"
		RequestHeaders = map[string]string{
			"Content-Type":  "application/json",
			"X-Test-Header": "test",
		}
		RequestPayload, _ = json.Marshal(map[string]interface{}{
			"test": "test",
		})

		ResponseStatusCode = 200
		ResponseHeader     = map[string]string{
			"Content-Type":  "application/json",
			"X-Test-Header": "test-response",
		}
		ResponsePayload, _ = json.Marshal(map[string]interface{}{
			"test": "test-response",
		})
	)

	// generate a key pair for the server
	var sharedkey *utils.JWK
	sPri, sPub, err := utils.GenerateKeyPair(utils.ECDH)
	assert.NoError(t, err)
	assert.NotNil(t, sPri)
	assert.NotNil(t, sPub)

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {
		case "/init-tunnel":
			token, err := genToken("mock_secret")
			assert.NoError(t, err)
			assert.NotNil(t, token)

			userpub, err := utils.B64ToJWK(r.Header.Get("x-ecdh-init"))
			assert.NoError(t, err)

			sharedkey, err = sPri.GetECDHSharedSecret(userpub)
			assert.NoError(t, err)

			data, err := json.Marshal(map[string]interface{}{
				"server_pubKeyECDH": sPub,
				"up_JWT":            token,
			})
			assert.NoError(t, err)

			w.WriteHeader(200)
			w.Write(data)
		default:
			// every request to the proxy is expected to be a POST request
			assert.Equal(t, r.Method, "POST")

			pURL, err := url.Parse(RequestURL)
			assert.NoError(t, err)

			// the "X-Forwarded-For" header must be set to the original request's IP
			assert.Equal(t, pURL.Host, r.Header.Get("X-Forwarded-Host"))
			// the "X-Forwarded-Proto" header must be set to the original request's scheme
			assert.Equal(t, pURL.Scheme, r.Header.Get("X-Forwarded-Proto"))
			// the path of the request must match the path of the original request
			assert.Equal(t, pURL.Path, r.URL.Path)

			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.NotNil(t, body)

			reqBody := make(map[string]interface{})
			err = json.Unmarshal(body, &reqBody)
			assert.NoError(t, err)

			// it is expected that the body is encrypted and encoded in base64 format
			// and set to the "data" key of the request body
			assert.NotNil(t, reqBody["data"])

			// decrypt the body
			data, err := base64.StdEncoding.DecodeString(reqBody["data"].(string))
			assert.NoError(t, err)

			decrypted, err := sharedkey.SymmetricDecrypt(data)
			assert.NoError(t, err)

			req, err := utils.FromJSONRequest(decrypted)
			assert.NoError(t, err)
			assert.NotNil(t, req)
			assert.Equal(t, req.Method, RequestMethod)
			assert.Equal(t, req.Headers, RequestHeaders)
			assert.Equal(t, req.Body, RequestPayload)

			// encrypt and return response
			res := utils.Response{
				Body:       ResponsePayload,
				Headers:    ResponseHeader,
				Status:     ResponseStatusCode,
				StatusText: http.StatusText(ResponseStatusCode),
			}
			bRes, err := res.ToJSON()
			assert.NoError(t, err)
			assert.NotNil(t, bRes)

			encRes, err := sharedkey.SymmetricEncrypt(bRes)
			assert.NoError(t, err)
			assert.NotNil(t, encRes)

			resData, err := json.Marshal(map[string]interface{}{
				"data": base64.StdEncoding.EncodeToString(encRes),
			})
			assert.NoError(t, err)

			w.WriteHeader(ResponseStatusCode)
			w.Write(resData)
		}
	}))
	defer mockServer.Close()

	// init tunnel
	pri, pub, err := utils.GenerateKeyPair(utils.ECDH)
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	b64, err := pub.ExportAsBase64()
	assert.NoError(t, err)
	assert.NotNil(t, b64)

	uuid := uuid.New().String()

	iClient := &http.Client{}
	iReq, err := http.NewRequest("GET", mockServer.URL+"/init-tunnel", nil)
	assert.NoError(t, err)

	iReq.Header.Add("x-ecdh-init", b64)
	iReq.Header.Add("x-client-uuid", uuid)

	iRes, err := iClient.Do(iReq)
	assert.NoError(t, err)
	assert.Equal(t, iRes.StatusCode, 200)

	iBody, err := io.ReadAll(iRes.Body)
	assert.NoError(t, err)

	iData := make(map[string]interface{})
	err = json.Unmarshal(iBody, &iData)
	assert.NoError(t, err)

	serverjwk, err := utils.JWKFromMap(iData)
	assert.NoError(t, err)
	assert.NotNil(t, serverjwk)

	symmkey, err := pri.GetECDHSharedSecret(serverjwk)
	assert.NoError(t, err)
	assert.NotNil(t, symmkey)

	// tunnel
	client := &Client{
		proxyURL: mockServer.URL,
	}

	res := client.Do(RequestURL, utils.NewRequest(RequestMethod, RequestHeaders, RequestPayload), symmkey)
	assert.NotNil(t, res)
	assert.Equal(t, ResponseStatusCode, res.Status)
	for k, v := range ResponseHeader {
		assert.Equal(t, v, res.Headers[k])
	}
	assert.Equal(t, ResponsePayload, res.Body)
	assert.Equal(t, http.StatusText(ResponseStatusCode), res.StatusText)
}
