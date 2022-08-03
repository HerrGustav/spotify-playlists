package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL = "https://api.spotify.com/v1"
)

type Client struct {
	httpClient HttpClient
	id         string
	secret     string
	userName   string
	token      string
}

func NewClient(id, secret, user string) Client {
	return Client{
		httpClient: http.DefaultClient,
		id:         id,
		secret:     secret,
		userName:   user,
	}
}

func (c *Client) Authorize() error {
	t, err := retrieveAuthToken(c.httpClient, c.id, c.secret)
	if err != nil {
		return fmt.Errorf("failed to authorize spotify client, %w", err)
	}

	c.token = t

	return nil
}

func (c *Client) IsAuthorized() bool {
	return c.token != ""
}

type Pagination struct {
	Limit    int64  `json:"limit"`
	Next     string `json:"next"`
	Offset   int64  `json:"offset"`
	Previous string `json:"previous"`
	Total    int64  `json:"total"`
}

type Playlist struct {
	Href  string        `json:"href"`
	Items []interface{} `json:"items"`
	Pagination
}

// GetPlaylists is targeting the endpoint described here in the spotify web api docs:
// https://developer.spotify.com/documentation/web-api/reference/#/operations/get-list-users-playlists
// xxx :  this is not covering the entire functionality yet, e.g. the Pagination or the parsing of all
// existing response fields. This is just needed for the POC for now.
// ref.: https://github.com/HerrGustav/spotify-playlists/issues/1
func (c *Client) GetPlaylists() (Playlist, error) {
	if !c.IsAuthorized() {
		return Playlist{}, newError(notAuthorized, "client is not authorized", nil)
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/users/"+c.userName+"/playlists", nil)
	if err != nil {
		return Playlist{}, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Playlist{}, err
	}

	var playlists Playlist
	err = json.NewDecoder(resp.Body).Decode(&playlists)
	if err != nil {
		return Playlist{}, err
	}

	return playlists, err
}
