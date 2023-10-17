package utilities

import (
	"encoding/json"
	"fmt"
)

type (
	// Request stores request data
	Request struct {
		Method  string            `json:"method"`
		Headers map[string]string `json:"headers"`
		Body    []byte            `json:"body"`
	}

	// Response stores response data
	Response struct {
		Status     int               `json:"status"`
		StatusText string            `json:"status_text"`
		Headers    map[string]string `json:"headers"`
		Body       []byte            `json:"body"`
	}
)

// NewRequest creates a new request.
func NewRequest(method string, headers map[string]string, body []byte) *Request {
	return &Request{
		Method:  method,
		Headers: headers,
		Body:    body,
	}
}

// ToJSON converts the request to JSON.
func (r *Request) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// ToMap converts the request to a map.
func (r *Request) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"method":  r.Method,
		"headers": r.Headers,
		"body":    r.Body,
	}
}

// ToJSON converts the response to JSON.
func (r *Response) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// ToJSONString converts the response to a JSON string.
func (r *Response) ToJSONString() (string, error) {
	resHeaders := `{}`
	if r.Headers != nil {
		b, err := json.Marshal(r.Headers)
		if err != nil {
			return "", err
		}
		resHeaders = string(b)
	}
	res := fmt.Sprintf(
		`{"status": %d,"status_text": "%s","headers": %s,"body": "%s"}`, 
		r.Status, r.StatusText, resHeaders, string(r.Body))
	return res, nil
}

// ToMap converts the response to a map.
func (r *Response) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status":      r.Status,
		"status_text": r.StatusText,
		"headers":     r.Headers,
		"body":        r.Body,
	}
}

// FromJSONRequest converts JSON to a request.
func FromJSONRequest(data []byte) (*Request, error) {
	req := &Request{}
	err := json.Unmarshal(data, req)
	return req, err
}

// FromJSONResponse converts JSON to a response.
func FromJSONResponse(data []byte) (*Response, error) {
	res := &Response{}
	err := json.Unmarshal(data, res)
	return res, err
}
