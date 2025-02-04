package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"moul.io/http2curl"
)

// HTTPClient wraps the standard http.Client
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTP client with timeout settings
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Request is a helper function to send HTTP requests
func (hc *HTTPClient) Request(method, url string, headers map[string]string, body interface{}) (*http.Response, error) {
	var requestBody io.Reader

	// Check if a body is provided and process it accordingly
	if body != nil {
		// Check if we need to URL encode the form data (application/x-www-form-urlencoded)
		if contentType, ok := headers["Content-Type"]; ok && contentType == "application/x-www-form-urlencoded" {
			formData, err := encodeFormData(body)
			if err != nil {
				return nil, err
			}
			requestBody = strings.NewReader(formData.Encode())
		} else {
			// Default to JSON encoding
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			requestBody = bytes.NewBuffer(jsonData)
		}
	}

	// Create new HTTP request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, err
	}

	// Log the curl command for debugging purposes
	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	// Set default content type if not provided
	if _, ok := headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return hc.client.Do(req)
}

func encodeFormData(body interface{}) (url.Values, error) {
	formData := url.Values{}
	if m, ok := body.(map[string]interface{}); ok {
		for key, value := range m {
			switch v := value.(type) {
			case string:
				formData.Add(key, v)
			case int:
				formData.Add(key, fmt.Sprintf("%d", v))
			case float64:
				formData.Add(key, fmt.Sprintf("%f", v))
			}
		}
		return formData, nil
	}
	return nil, fmt.Errorf("body type is not supported for form encoding")
}
