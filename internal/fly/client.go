package fly

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	baseURL := "http://_api.internal:4280/v1"

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, headers map[string]string, response interface{}) error {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	token := os.Getenv("FLY_API_TOKEN")
	if token == "" {
		return errors.New("FLY_API_TOKEN environment variable must be set")
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if response != nil {
			return json.NewDecoder(res.Body).Decode(response)
		}
		return nil
	} else {
		log.Printf("request failed: %s", req.URL.String())
		log.Printf("error response: %s", res.Status)
		var errRes interface{}
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			log.Printf("failed to decode error response: %s", err.Error())
			return err
		}
		log.Printf("error response: %+v", errRes)
		return errors.New("error from fly")
	}
}
