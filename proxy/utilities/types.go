package utilities

import "encoding/json"

type (
	// Request stores request data
	Request struct {
		Method  string            `json:"method"`
		Url     string            `json:"url"`
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
func NewRequest(method, url string, headers map[string]string, body []byte) *Request {
	return &Request{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    body,
	}
}

// ToJSON converts the request to JSON.
func (r *Request) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// ToJSON converts the response to JSON.
func (r *Response) ToJSON() ([]byte, error) {
	return json.Marshal(r)
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
