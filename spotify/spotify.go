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

// Pagination is the representation of the pagination values that are
// typically included in every response of the spotify web api. This here
// is for keeping the code a bit more dry.
type Pagination struct {
	Limit    int64  `json:"limit"`
	Next     string `json:"next"`
	Offset   int64  `json:"offset"`
	Previous string `json:"previous"`
	Total    int64  `json:"total"`
}

// Playlist is a minimal representation of the response object described here:
// https://developer.spotify.com/documentation/web-api/reference/#/operations/create-playlist
// It does not aim to be the complete representation for now.
type Playlist struct {
	Collaborative bool   `json:"collaborative"`
	Description   string `json:"description"`
	Href          string `json:"href"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Public        bool   `json:"public"`
	SnapshotID    bool   `json:"snapshot_id"`
	Type          string `json:"type"`
	URI           string `json:"uri"`
	Pagination
}

// UserPlaylists is the minimal representation of this response
// // https://developer.spotify.com/documentation/web-api/reference/#/operations/get-list-users-playlists
type UserPlaylists struct {
	Href  string        `json:"href"`
	Items []interface{} `json:"items"`
	Pagination
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

// GetUserPlaylists is targeting the endpoint described here in the spotify web api docs:
// https://developer.spotify.com/documentation/web-api/reference/#/operations/get-list-users-playlists
// xxx :  this is not covering the entire functionality yet, e.g. the Pagination or the parsing of all
// existing response fields. This is just needed for the POC for now.
// ref.: https://github.com/HerrGustav/spotify-playlists/issues/1
func (c *Client) GetUserPlaylists() (UserPlaylists, error) {
	if !c.IsAuthorized() {
		return UserPlaylists{}, newError(notAuthorized, "client is not authorized", nil)
	}

	req, err := c.createAuthorizedRequest(http.MethodGet, baseURL+"/users/"+c.userName+"/playlists", nil)
	if err != nil {
		return UserPlaylists{}, err
	}

	resp, err := c.doRequest(req, http.StatusOK)
	if err != nil {
		return UserPlaylists{}, newError(requestFailed, "failed to request api", err)
	}
	defer resp.Body.Close()

	var playlists UserPlaylists
	err = json.NewDecoder(resp.Body).Decode(&playlists)
	if err != nil {
		return UserPlaylists{}, err
	}

	return playlists, err
}

type CreatePlaylistPayload struct {
	Name          string `json:"name"`
	Public        bool   `json:"public,omitempty"`
	Collaborative bool   `json:"collaborative,omitempty"`
	Description   string `json:"description,omitempty"`
}

// CreatePlaylist creates a playlist for a given user as described here:
// https://developer.spotify.com/documentation/web-api/reference/#/operations/create-playlist
// The name of the playlist is the only mandatory input.
func (c *Client) CreatePlaylist(payload CreatePlaylistPayload) (Playlist, error) {
	if payload.Name == "" {
		return Playlist{}, newError(invalidInputs, "playlist name is a required payload field", nil)
	}

	if !c.IsAuthorized() {
		return Playlist{}, newError(notAuthorized, "client is not authorized", nil)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Playlist{}, newError(internalError, "failed to marshal request payload", err)
	}

	req, err := c.createAuthorizedRequest(http.MethodPost, baseURL+"/users/"+c.userName+"/playlists", body)
	if err != nil {
		return Playlist{}, newError(internalError, "failed to create authorized request", err)
	}

	resp, err := c.doRequest(req, http.StatusCreated)
	if err != nil {
		return Playlist{}, newError(requestFailed, "failed to request api", err)
	}
	defer resp.Body.Close()

	var p Playlist
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return Playlist{}, newError(internalError, "failed to decode api response", err)
	}

	return p, nil
}
