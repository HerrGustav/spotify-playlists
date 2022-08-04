package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// errorResponse is representing an error from the spotify web api
// communicated back to our client as a response. It seems to be the only format used
// to send back an error message from their backend to our client, described e.g. here:
// https://developer.spotify.com/documentation/web-api/reference/#/operations/create-playlist
type errResponse struct {
	Err struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e errResponse) Error() string {
	return fmt.Sprintf("error response from spotify web api, status: '%d', message: '%s'", e.Err.Status, e.Err.Message)
}

func (e errResponse) Is(err error) bool {
	if parsed, ok := err.(*errResponse); ok {
		// at least one fields should be present, ideally both:
		return parsed.Err.Status != 0 || parsed.Err.Message != ""
	}

	return false
}

func parseErrResponse(resp *http.Response) error {
	if resp == nil {
		return nil
	}

	var e errResponse
	err := json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		return fmt.Errorf("failed to parse error response from spotify web api, %w", err)
	}

	if e.Err.Status == 0 && e.Err.Message == "" {
		return nil
	}

	return e
}
