package utils

// for processing and responding to http requests

import (
	"io"
)

func ReadResponseBody(body io.ReadCloser) []byte {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return bodyBytes
}
