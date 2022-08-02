package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseTokenURL = "https://accounts.spotify.com/api/token"
	successCode  = http.StatusOK
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func createAuthHeader(id, secret string) string {
	clientCredentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", id, secret)))
	return fmt.Sprintf("Basic %s", clientCredentials)
}

func createAuthRequest(id, secret string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, baseTokenURL, bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return &http.Request{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", createAuthHeader(id, secret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	return req, nil
}

type authBody struct {
	AccessToken string `json:"access_token"`
}

func readAuthBody(body io.ReadCloser) (authBody, error) {
	if body == nil {
		return authBody{}, fmt.Errorf("auth body was empty")
	}
	var a authBody
	err := json.NewDecoder(body).Decode(&a)
	if err != nil {
		return authBody{}, fmt.Errorf("failed to decode auth body, %w", err)
	}

	return a, nil
}

func retrieveAuthToken(httpClient HttpClient, id, secret string) (string, error) {
	if httpClient == nil {
		return "", fmt.Errorf("http client can not be nil")
	}

	req, err := createAuthRequest(id, secret)
	if err != nil {
		return "", fmt.Errorf("failed to create auth request, %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute auth request, %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != successCode {
		return "", fmt.Errorf("got unexpected status code '%d', want: '%d'", resp.StatusCode, successCode)
	}

	body, err := readAuthBody(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read auth body, %w", err)
	}

	return body.AccessToken, nil
}
