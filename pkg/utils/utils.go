package utils

import (
	"bytes"
	"io"
	"net/http"
)

func SendHTTPRequest(method, url, format string, body *bytes.Buffer) (io.ReadCloser, error) {
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	return res.Body, nil
}
