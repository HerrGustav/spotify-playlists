package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDoRequest(t *testing.T) {
	shouldError := func(t *testing.T, err error) {
		if err == nil {
			t.Error("spotify.Client.doRequest() mismatch, expected error, got nil")
		}
	}

	expectedStatusErr := errResponse{Err: struct {
		Status  int    "json:\"status\""
		Message string "json:\"message\""
	}{
		Status:  http.StatusBadRequest,
		Message: "something failed",
	},
	}

	testcases := map[string]struct {
		client           Client
		req              *http.Request
		expectedStatus   int
		expectedResponse *http.Response
		checkError       func(t *testing.T, err error)
	}{
		"input is a pointer and can not be nil -- return an error": {
			checkError: shouldError,
		},
		"request failed -- return the error": {
			client: Client{
				httpClient: &mockHttpClient{
					expectedError: errMock,
				},
			},
			req: &http.Request{
				Method: http.MethodPost,
			},
			checkError: shouldError,
		},
		// since the actual request hit the spotify api and we got an unexpected status code,
		// we should be able to have an actual error response as described in the docs.
		// e.g. here https://developer.spotify.com/documentation/web-api/reference/#/operations/create-playlist
		"got error code -- return the actual response error inside the error": {
			client: Client{
				httpClient: &mockHttpClient{
					expectedResponse: createMockedHttpResponse(t, http.StatusBadRequest, expectedStatusErr),
				},
			},
			expectedStatus: http.StatusCreated,
			req: &http.Request{
				Method: http.MethodPost,
			},
			checkError: func(t *testing.T, err error) {
				if err == nil {
					t.Error("spotify.Client.doRequest() mismatch, expected error, got nil")
				}

				if !errors.Is(err, expectedStatusErr) {
					t.Errorf("spotify.Client.doRequest() failed to find errResponse inside the error chain, \n - expected: '%s', \n - got: '%s'", err.Error(), expectedStatusErr.Error())
				}
			},
		},
		"got success code -- return the expected response": {
			client: Client{
				httpClient: &mockHttpClient{
					expectedResponse: createMockedHttpResponse(t, http.StatusCreated, Playlist{Name: "mock"}),
				},
			},
			req: &http.Request{
				Method: http.MethodPost,
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: createMockedHttpResponse(t, http.StatusCreated, Playlist{Name: "mock"}),
			checkError: func(t *testing.T, err error) {
				if err != nil {
					t.Errorf("spotify.Client.doRequest() mismatch, unexpected error: '%s'", err.Error())
				}
			},
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			resp, err := tc.client.doRequest(tc.req, tc.expectedStatus)
			tc.checkError(t, err)
			if err != nil {
				return
			}

			var got interface{}
			err = json.NewDecoder(resp.Body).Decode(&got)
			if err != nil {
				t.Errorf("spotify.Client.doRequest() failed to unmarshal got: '%s'", err.Error())
			}

			var want interface{}
			err = json.NewDecoder(tc.expectedResponse.Body).Decode(&want)
			if err != nil {
				t.Errorf("spotify.Client.doRequest() failed to unmarshal want: '%s'", err.Error())
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("spotify.Client.doRequest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
