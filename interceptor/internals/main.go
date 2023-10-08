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
	
	// hardcoding a shared secret for now
	secret, err := base64.StdEncoding.DecodeString("KfbCmY2v83ptAZLLKffx0ve2Br8hkMhCkIo5RkFaNlk=")
	if err != nil {
		return &utilities.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}

	// send request
	res, err := c.transfer(secret, req, url, clientID)
	if err != nil {
		return &utilities.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}
	return res
}

// transfer sends the request to the remote server through the layer8 proxy server
func (c *Client) transfer(secret []byte, req *utilities.Request, url, clientID string) (*utilities.Response, error) {
	// encode request body
	b, err := req.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("could not encode request: %v", err)
	}
	// send the request
	_, res := c.do(b, secret, url, clientID)
	// decode response body
	resData, err := utilities.FromJSONResponse(res)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}
	return resData, nil
}

// do sends the request to the remote server through the layer8 proxy server
// returns a status code and response body
func (c *Client) do(data, secret []byte, backendUrl, clientID string) (int, []byte) {
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
