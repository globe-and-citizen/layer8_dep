package internals

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"globe-and-citizen/layer8/utils"
	"net/http"
	"net/url"
)

type Client struct {
	proxyURL string
}

type ClientImpl interface {
	Do(url string, req *utils.Request, sharedSecret *utils.JWK) *utils.Response
}

// NewClient creates a new client with the given proxy server url
func NewClient(scheme, host, port string) ClientImpl {
	return &Client{
		proxyURL: fmt.Sprintf("%s://%s:%s", scheme, host, port),
	}
}

// Do sends a request to through the layer8 proxy server and returns a response
func (c *Client) Do(url string, req *utils.Request, sharedSecret *utils.JWK) *utils.Response {
	// Send request
	res, err := c.transfer(sharedSecret, req, url)
	if err != nil {
		return &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}
	}
	return res
}

// transfer sends the request to the remote server through the layer8 proxy server
func (c *Client) transfer(sharedSecret *utils.JWK, req *utils.Request, url string) (*utils.Response, error) {
	// encode request body
	b, err := req.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("Could not encode request: %w", err)
	}
	// send the request
	_, res := c.do(b, sharedSecret, url)
	// decode response body
	resData, err := utils.FromJSONResponse(res)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	// Perhaps it's here that you'll rehydrate the headers from the service provider?
	resData.Headers["x-custom-header-1"] = "ONE"
	resData.Headers["x-custom-header-2"] = "TWO"
	resData.Headers["x-custom-header-3"] = "THREE"

	return resData, nil
}

// do sends the request to the remote server through the layer8 proxy server
// returns a status code and response body
func (c *Client) do(data []byte, sharedSecret *utils.JWK, backendUrl string) (int, []byte) {
	// encrypt request body if a secret is provided
	// if secret != nil {
	// 	data, err = utils.Dep_SymmetricEncrypt(data, secret)

	var err error
	data, err = sharedSecret.SymmetricEncrypt(data)
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}

	data, err = json.Marshal(map[string]interface{}{
		"data": base64.URLEncoding.EncodeToString(data),
	})

	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}

	parsedURL, err := url.Parse(backendUrl)
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
		// return 500, []byte{}
	}
	// create request
	client := &http.Client{}
	r, err := http.NewRequest("POST", c.proxyURL+parsedURL.Path, bytes.NewBuffer(data))
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}
	// add headers
	r.Header.Add("X-Forwarded-Host", parsedURL.Host)
	r.Header.Add("X-Forwarded-Proto", parsedURL.Scheme)
	r.Header.Add("Content-Type", "application/json")
	// send request
	res, err := client.Do(r)
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	bufByte := buf.Bytes()

	mapB := make(map[string]interface{})
	json.Unmarshal(bufByte, &mapB)

	fmt.Println("mapB: ", mapB)
	toDecode, ok := mapB["data"].(string)
	if !ok {
		res := &utils.Response{
			Status:     500,
			StatusText: "mapB[\"data\"].(string) not 'ok'",
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}

	decoded, err := base64.URLEncoding.DecodeString(toDecode)
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}

	fmt.Println("decoded: ", decoded)
	bufByte, err = sharedSecret.SymmetricDecrypt(decoded)
	fmt.Println("bufBytes: ", string(bufByte))
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}
		resByte, _ := res.ToJSON()
		return 500, resByte
	}

	// At this point the proxy's headers have been stripped and you have the SP's response as bufByte
	return res.StatusCode, bufByte
}
