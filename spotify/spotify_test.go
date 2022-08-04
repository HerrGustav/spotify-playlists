package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	errMock = errors.New("failed")
)

func TestAuthorize(t *testing.T) {
	testcases := map[string]struct {
		client        Client
		expectedToken string
		shouldError   bool
	}{
		"failed to authorize client": {
			client: Client{
				httpClient: &mockHttpClient{
					expectedError: errMock,
				},
				id:     "test",
				secret: "test",
			},
			shouldError: true,
		},
		"successfully authorized client": {
			client: Client{
				httpClient: &mockHttpClient{
					expectedResponse: &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"access_token": "test_token"}`)),
					},
				},
				id:     "test",
				secret: "test",
			},
			expectedToken: "test_token",
			shouldError:   false,
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			err := tc.client.Authorize()
			if tc.shouldError && err == nil {
				t.Error("spotify.Client.Authorize() got no error but was expected to fail")
			} else if err != nil && !tc.shouldError {
				t.Errorf("spotify.Client.Authorize() got unexpected error, \n got: '%s'", err.Error())
			}

			if tc.client.token != tc.expectedToken {
				t.Errorf("spotify.Client.Authorize() mismatch, \n got: '%s', \n want: '%s'", tc.client.token, tc.expectedToken)
			}
		})
	}
}

func TestGetPlaylists(t *testing.T) {
	testcases := map[string]struct {
		client      Client
		want        UserPlaylists
		shouldError bool
	}{
		"not authorized": {
			shouldError: true,
		},
		"failed to request api": {
			shouldError: true,
			client: mockAuthorizedClient(Client{
				httpClient: &mockHttpClient{
					expectedError: errMock,
				},
				id:     "test",
				secret: "test",
			}),
		},
		"successfully retrieved playlists": {
			client: mockAuthorizedClient(Client{
				httpClient: &mockHttpClient{
					expectedResponse: &http.Response{
						Body: io.NopCloser(bytes.NewBuffer(marshalInterface(t, UserPlaylists{
							Href: "test",
						}))),
					},
				},
				id:     "test",
				secret: "test",
			}),
			want: UserPlaylists{
				Href: "test",
			},
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			got, err := tc.client.GetUserPlaylists()
			if tc.shouldError && err == nil {
				t.Error("spotify.Client.TestGetPlaylists() got no error but was expected to fail")
			} else if err != nil && !tc.shouldError {
				t.Errorf("spotify.Client.TestGetPlaylists() got unexpected error, \n got: '%s'", err.Error())
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("spotify.Client.GetPlaylists() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCreatePlaylist(t *testing.T) {
	mockPayload := CreatePlaylistPayload{
		Name: "mock name",
	}
	testcases := map[string]struct {
		client         Client
		payload        CreatePlaylistPayload
		expectedError  error
		expectedOutput Playlist
	}{
		"Given inputs are not valid -- should fail": {
			client: Client{},
			expectedError: errSpotify{
				code: invalidInputs,
			},
		},
		"User not authorized -- should fail": {
			client:  Client{},
			payload: mockPayload,
			expectedError: errSpotify{
				code: notAuthorized,
			},
		},
		"Failed to create playlist -- return error": {
			client: mockAuthorizedClient(Client{
				httpClient: &mockHttpClient{
					expectedError: errMock,
				},
				id:     "test",
				secret: "test",
			}),
			payload: mockPayload,
			expectedError: errSpotify{
				code: requestFailed,
			},
		},
		"Successfully created playlist -- return expected result": {
			client: mockAuthorizedClient(Client{
				httpClient: &mockHttpClient{
					expectedResponse: createMockedHttpResponse(t, Playlist{
						ID: "1234",
					}),
				},
				id:     "test",
				secret: "test",
			}),
			payload: mockPayload,
			expectedOutput: Playlist{
				ID: "1234",
			},
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			out, err := tc.client.CreatePlaylist(tc.payload)
			checkSpotifyError(t, tc.expectedError, err)

			if diff := cmp.Diff(tc.expectedOutput, out); diff != "" {
				t.Errorf("spotify.Client.CreatePlaylist() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func checkSpotifyError(t *testing.T, want, got error) {
	// want == nil && got == nil has no effect.
	// so we skip this here
	if want == nil && got != nil {
		t.Errorf("test '%s': unexpected error - '%s' ", t.Name(), got.Error())
		return
	}

	if want != nil && !errors.Is(got, want) {
		t.Errorf("test '%s': unexpected error, \n - got: '%v' , \n - want: '%v'", t.Name(), got, want)
		return
	}
}

func marshalInterface(t *testing.T, v interface{}) []byte {
	out, err := json.Marshal(v)
	if err != nil {
		t.Errorf("marshalInterface(): failed for test '%s' to marshal input of '%v'", t.Name(), v)
		return nil
	}

	return out
}

func createMockedHttpResponse(t *testing.T, v interface{}) *http.Response {
	return &http.Response{
		Body: io.NopCloser(bytes.NewBuffer(marshalInterface(t, v))),
	}
}

// mockAuthorizedClient is basically only adding a fake token
// to the client struct and is imitating the behavior of "client.Authorize()"
// and the check at "client.IsAuthorized()". This prevents that other unit tests
// need to deal with this implementation detail.
func mockAuthorizedClient(baseClient Client) Client {
	m := baseClient // copy to not return modified inputs.
	m.token = "12345"
	return m
}
