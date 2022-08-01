package spotify

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestCreateAuthHeader(t *testing.T) {
	testcases := map[string]struct {
		id       string
		secret   string
		expected string
	}{
		"Returns expected auth header": {
			id:       "id",
			secret:   "secret",
			expected: "Basic aWQ6c2VjcmV0",
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			out := createAuthHeader(tc.id, tc.secret)
			if tc.expected != out {
				t.Errorf("unexpected result, \n - got: '%s', \n - want: '%s'", out, tc.expected)
			}
		})
	}
}
func TestReadAuthBody(t *testing.T) {
	testcases := map[string]struct {
		body        io.ReadCloser
		id          string
		secret      string
		expected    authBody
		shouldError bool
	}{
		"empty auth body was passed as input": {
			id:          "id",
			secret:      "secret",
			shouldError: true,
		},
		"decoding of the auth body failed": {
			body:        io.NopCloser(bytes.NewBufferString("1234")),
			id:          "id",
			secret:      "secret",
			shouldError: true,
		},
		"Returns expected auth header": {
			body:     io.NopCloser(bytes.NewBufferString(`{"access_token": "1234"}`)),
			id:       "id",
			secret:   "secret",
			expected: authBody{AccessToken: "1234"},
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			got, err := readAuthBody(tc.body)
			if err != nil && !tc.shouldError {
				t.Errorf("TestReadAuthBody() got unexpected error '%s'", err.Error())
			} else if err == nil && tc.shouldError {
				t.Errorf("TestReadAuthBody() did not throw an error as expected")
			}

			if tc.expected.AccessToken != got.AccessToken {
				t.Errorf("unexpected result, \n - got: '%s', \n - want: '%s'", got, tc.expected)
			}
		})
	}
}

type mockHttpClient struct {
	expectedResponse *http.Response
	expectedError    error
}

func (m *mockHttpClient) Do(r *http.Request) (*http.Response, error) {
	return m.expectedResponse, m.expectedError
}

func TestRetrieveAuthToken(t *testing.T) {
	testcases := map[string]struct {
		httpClient  HttpClient
		id          string
		secret      string
		want        string
		shouldError bool
	}{
		"no http client was passed": {
			shouldError: true,
		},
		"failed to execute request": {
			shouldError: true,
			httpClient: &mockHttpClient{
				expectedError: errors.New("mock error"),
			},
		},
		"failed to read auth body": {
			shouldError: true,
			httpClient: &mockHttpClient{
				expectedResponse: &http.Response{Body: io.NopCloser(bytes.NewBufferString("1234"))},
			},
		},
		"got expected token": {
			want: "1234",
			httpClient: &mockHttpClient{
				expectedResponse: &http.Response{Body: io.NopCloser(bytes.NewBufferString(`{"access_token": "1234"}`))},
			},
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			got, err := retrieveAuthToken(tc.httpClient, tc.id, tc.secret)
			if err != nil && !tc.shouldError {
				t.Errorf("TestRetrieveAuthToken() got unexpected error '%s'", err.Error())
			} else if err == nil && tc.shouldError {
				t.Errorf("TestRetrieveAuthToken() did not throw an error as expected")
			}

			if tc.want != got {
				t.Errorf("unexpected result, \n - got: '%s', \n - want: '%s'", got, tc.want)
			}
		})
	}
}

// func TestCreateAuthRequest(t *testing.T) {
// 	testcases := map[string]struct {
// 		id     string
// 		secret string
// 		want   http.Request
// 	}{
// 		"Returns expected auth request": {
// 			id:     "id",
// 			secret: "secret",
// 			want: http.Request{
// 				Method: http.MethodPost,
// 			},
// 		},
// 	}

// 	for testName, tc := range testcases {
// 		t.Run(testName, func(t *testing.T) {
// 			got, err := createAuthRequest(tc.id, tc.secret)
// 			if err != nil {
// 				t.Errorf("unexpected error: %v", err)
// 			}
// 			if diff := cmp.Diff(tc.want, got); diff != "" {
// 				t.Errorf("CreateAuthRequest() mismatch (-want +got):\n%s", diff)
// 			}
// 		})
// 	}
// }
