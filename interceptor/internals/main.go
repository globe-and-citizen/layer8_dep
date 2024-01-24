package internals

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	// "globe-and-citizen/layer8/utils" (Dep)
	"net/http"
	"net/url"

	utils "github.com/globe-and-citizen/layer8-utils"
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
	// send the request
	res := c.do(req, sharedSecret, url, req.Headers)
	// decode response body
	resData, err := utils.FromJSONResponse(res)
	if err != nil {
		return &utils.Response{
			Status:     500,
			StatusText: err.Error(),
		}, nil
	}

	// Perhaps it's here that you'll rehydrate the headers from the service provider?
	resData.Headers["x-custom-header-1"] = "ONE"
	resData.Headers["x-custom-header-2"] = "TWO"
	resData.Headers["x-custom-header-3"] = "THREE"

	return resData, nil
}

// do sends the request to the remote server through the layer8 proxy server
// returns a status code and response body
func (c *Client) do(req *utils.Request, sharedSecret *utils.JWK, backendUrl string, headers map[string]string) []byte {
	// encrypt request body if a secret is provided
	// if secret != nil {
	// 	data, err = utils.Dep_SymmetricEncrypt(data, secret)

	var err error

	data, err := req.ToJSON()
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: fmt.Sprintf("Error marshalling request: %s", err.Error()),
		}
		resByte, _ := res.ToJSON()
		return resByte
	}

	data, err = sharedSecret.SymmetricEncrypt(data)
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
			Headers:    make(map[string]string),
		}
		resByte, _ := res.ToJSON()
		return resByte
	}

	data, err = json.Marshal(map[string]interface{}{
		"data": base64.URLEncoding.EncodeToString(data),
	})

	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
			Headers:    make(map[string]string),
		}
		resByte, _ := res.ToJSON()
		return resByte
	}

	parsedURL, err := url.Parse(backendUrl)
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
			Headers:    make(map[string]string),
		}
		resByte, _ := res.ToJSON()
		return resByte
	}
	// create request
	client := &http.Client{}
	r, err := http.NewRequest("POST", c.proxyURL+parsedURL.Path, bytes.NewBuffer(data))
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
			Headers:    make(map[string]string),
		}
		resByte, _ := res.ToJSON()
		return resByte
	}
	// add headers
	r.Header.Add("X-Forwarded-Host", parsedURL.Host)
	r.Header.Add("X-Forwarded-Proto", parsedURL.Scheme)
	r.Header.Add("Content-Type", "application/json")
	// Add custom headers being sent to the client side [Important]
	for k, v := range headers {
		r.Header.Add(k, v)
	}
	// send request
	res, err := client.Do(r)
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
			Headers:    make(map[string]string),
		}
		resByte, _ := res.ToJSON()
		return resByte
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
			Headers:    make(map[string]string),
		}
		resByte, _ := res.ToJSON()
		return resByte
	}

	decoded, err := base64.URLEncoding.DecodeString(toDecode)
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
			Headers:    make(map[string]string),
		}
		resByte, _ := res.ToJSON()
		return resByte
	}

	fmt.Println("decoded: ", decoded)
	bufByte, err = sharedSecret.SymmetricDecrypt(decoded)
	fmt.Println("bufBytes: ", string(bufByte))
	if err != nil {
		res := &utils.Response{
			Status:     500,
			StatusText: err.Error(),
			Headers:    make(map[string]string),
		}
		resByte, _ := res.ToJSON()
		return resByte
	}

	// At this point the proxy's headers have been stripped and you have the SP's response as bufByte
	return bufByte
}

// DoForm sends the formdata requests (e.g. multipart/form-data, application/x-www-form-urlencoded) 
// which includes files to the remote server through the layer8 proxy server
func (c *Client) DoForm(url string, req *utils.Request, sharedSecret *utils.JWK) *utils.Response {
	fmt.Println(req.Body)

	return nil
}