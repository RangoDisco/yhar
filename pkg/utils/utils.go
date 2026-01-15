package utils

import (
	"bytes"
	"net/http"
)

func PrepareHTTPRequest(method, url, format string, body *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	var contentType string
	switch format {
	case "json":
		contentType = "application/json"
		break
	case "xml":
		contentType = "application/xml"
		break
	default:
		contentType = "application/json"
	}
	req.Header.Set("Content-Type", contentType)

	return req, nil
}
