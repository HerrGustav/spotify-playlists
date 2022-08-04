package spotify

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func (c *Client) createAuthorizedRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return req, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Accept", "application/json")

	return req, nil
}

func (c *Client) doRequest(req *http.Request, expectedStatus int) (*http.Response, error) {
	if req == nil {
		// this can only happen in a internal use case of this pkg, so specify the method name:
		return nil, errors.New("doRequest(): the request as input can not be nil")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// we return only the error here, since we can ignore the resp according to the docs:
		// https://pkg.go.dev/net/http#Client.Do
		return nil, err
	}

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("unexpected status code, got: '%d' - expected: '%d', %w", resp.StatusCode, expectedStatus, parseErrResponse(resp))
	}

	return resp, nil
}
