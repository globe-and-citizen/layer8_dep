package internals

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	utilities "github.com/globe-and-citizen/layer8-utils"

	"github.com/google/uuid"
)

type Client struct {
	proxyURL string
}

type ClientImpl interface {
	Do(url string, req *utilities.Request) *utilities.Response
}

// NewClient creates a new client with the given proxy server url
func NewClient(scheme, host, port string) ClientImpl {
	return &Client{
		proxyURL: fmt.Sprintf("%s://%s:%s", scheme, host, port),
	}
}

// Do sends a request to through the layer8 proxy server and returns a response
func (c *Client) Do(url string, req *utilities.Request) *utilities.Response {
	clientID := uuid.New().String()
	// exchange key
	shared, err := c.exchangeKey(url, clientID)
	if err != nil {
		return &utilities.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}

	// send request
	res, err := c.transfer(shared, req, url, clientID)
	if err != nil {
		return &utilities.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}
	return res
}

// exchangeKey exchanges public keys with the service provider's backend through the layer8 proxy server
// for symmetric encryption and returns a shared secret
func (c *Client) exchangeKey(url, clientID string) ([]byte, error) {
	// generating a key pair
	pri, pub, err := utilities.GenerateKeyPair(utilities.ECDH_ALGO)
	if err != nil {
		return nil, err
	}
	// create and send the request
	key, err := utilities.EncodePublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("could not encode public key: %v", err)
	}
	req := utilities.NewRequest("POST", map[string]string{
		"X-Layer8-CPK": key,
	}, nil)
	reqByte, _ := req.ToJSON()
	status, res := c.do(reqByte, nil, url, clientID, true)
	if status != 200 {
		resData, _ := utilities.FromJSONResponse(res)
		return nil, fmt.Errorf("could not exchange keys: %s", resData.StatusText)
	}

	// decode response
	resData, _ := utilities.FromJSONResponse(res)
	spk, ok := resData.Headers["X-Layer8-SPK"]
	if !ok {
		spk = resData.Headers["x-layer8-spk"]
	}

	serverPub, err := utilities.DecodePublicKey(spk)
	if err != nil {
		return nil, fmt.Errorf("could not decode server public key: %v", err)
	}
	return utilities.DeriveSharedSecret(pri, serverPub), nil
}

// transfer sends the request to the remote server through the layer8 proxy server
func (c *Client) transfer(secret []byte, req *utilities.Request, url, clientID string) (*utilities.Response, error) {
	// encode request body
	b, err := req.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("could not encode request: %v", err)
	}
	// send the request
	_, res := c.do(b, secret, url, clientID, false)
	// decode response body
	resData, err := utilities.FromJSONResponse(res)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}
	return resData, nil
}

// do sends the request to the remote server through the layer8 proxy server
// returns a status code and response body
func (c *Client) do(data, secret []byte, backendUrl, clientID string, isKeyExchange bool) (int, []byte) {
	var err error

	// encrypt request body if a secret is provided
	if secret != nil {
		data, err = utilities.SymmetricEncrypt(data, secret)
		if err != nil {
			res := &utilities.Response{
				Status:     500,
				StatusText: err.Error(),
			}
			resByte, _ := res.ToJSON()
			return 500, resByte
		}
	}
	data, err = json.Marshal(map[string]interface{}{
		"data": base64.StdEncoding.EncodeToString(data),
	})
	if err != nil {
		res := &utilities.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}

	parsedURL, _ := url.Parse(backendUrl)
	// create request
	client := &http.Client{}
	r, err := http.NewRequest("POST", c.proxyURL+parsedURL.Path, bytes.NewBuffer(data))
	if err != nil {
		res := &utilities.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}
	// add headers
	r.Header.Add("X-Forwarded-Host", parsedURL.Host)
	r.Header.Add("X-Forwarded-Proto", parsedURL.Scheme)
	r.Header.Add("X-Layer8-CID", clientID)
	r.Header.Add("Content-Type", "application/json")
	if isKeyExchange {
		r.Header.Add("X-Layer8-Key-Exchange", "true")
	}
	// send request
	res, err := client.Do(r)
	if err != nil {
		res := &utilities.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	// decrypt response body if a secret is provided
	bufByte := buf.Bytes()

	if secret != nil {
		// encrypted response is encoded in base64
		// so we need to decode it first
		mapB := make(map[string]interface{})
		json.Unmarshal(bufByte, &mapB)

		decoded, err := base64.StdEncoding.DecodeString(mapB["data"].(string))
		if err != nil {
			res := &utilities.Response{
				Status:     500,
				StatusText: err.Error(),
			}
			resByte, _ := res.ToJSON()
			return 500, resByte
		}
		bufByte, err = utilities.SymmetricDecrypt(decoded, secret)
		if err != nil {
			res := &utilities.Response{
				Status:     500,
				StatusText: err.Error(),
			}
			resByte, _ := res.ToJSON()
			return 500, resByte
		}
	}
	return res.StatusCode, bufByte
}
