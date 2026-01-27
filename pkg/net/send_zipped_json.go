package net

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/sign"
)

func SendZippedSignedJSON(url string, data any, key string) error {
	body, err := getRequestBody(data)
	if err != nil {
		return fmt.Errorf("getRequestBody: %s", err)
	}
	var hash []byte
	if key != "" {
		signFunc := sign.GetSignFunc(key)
		data, _ := io.ReadAll(body)
		hash = signFunc(data)
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("creating request error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	if len(hash) > 0 {
		base64String := base64.StdEncoding.EncodeToString(hash)
		req.Header.Set("HashSHA256", base64String)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST: %w", err)
	}
	response.Body.Close()
	return nil
}

func SendZippedJSON(url string, data any) error {
	return SendZippedSignedJSON(url, data, "")
}

func getRequestBody(data any) (io.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshalling metric error: %s", err)
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(jsonData)
	gz.Close()
	return &buf, nil
}
