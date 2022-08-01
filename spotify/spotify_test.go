package spotify

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
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
						Body: io.NopCloser(bytes.NewBufferString(`{"access_token": "test_token"}`)),
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
