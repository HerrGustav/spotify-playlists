package spotify

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrResponseIs(t *testing.T) {
	mockResponseErr := errResponse{}
	mockResponseErr.Err.Status = 400
	mockResponseErr.Err.Message = "mock"

	testcases := map[string]struct {
		err   errResponse
		input error
		want  bool
	}{
		"returns false if a nil error is parsed in": {
			err:  mockResponseErr,
			want: false,
		},
		"returns false if another error type is parsed in": {
			err:   mockResponseErr,
			input: errors.New("mock error"),
			want:  false,
		},
		"returns true if a response error is parsed in": {
			err:   mockResponseErr,
			input: mockResponseErr,
			want:  true,
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			got := errors.Is(tc.err, tc.input)
			if got != tc.want {
				t.Errorf("errResponse.Is(): unexpected result, \n - got: '%v', \n - want: '%v'", got, tc.want)
			}
		})
	}
}

func TestParseErrResponse(t *testing.T) {
	errNotExpected := func(t *testing.T, err error) {
		if err != nil {
			t.Errorf("unexpected error, got '%s", err.Error())
		}
	}

	testcases := map[string]struct {
		resp       *http.Response
		checkError func(t *testing.T, err error)
	}{
		"returns an empty error if an empty response is passed in": {
			checkError: errNotExpected,
		},
		"returns nil if the response is not of type errorResponse": {
			resp:       createMockedHttpResponse(t, http.StatusBadRequest, Playlist{}),
			checkError: errNotExpected,
		},
		"returns nil if no fields are set": {
			resp:       createMockedHttpResponse(t, http.StatusBadRequest, errResponse{}),
			checkError: errNotExpected,
		},
		"returns the parsed error response": {
			resp: createMockedHttpResponse(t, http.StatusBadRequest, errResponse{Err: struct {
				Status  int    "json:\"status\""
				Message string "json:\"message\""
			}{Status: 400}}),
			checkError: func(t *testing.T, err error) {
				if err == nil {
					t.Error("expected error, got none")
				}
			},
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			err := parseErrResponse(tc.resp)
			tc.checkError(t, err)
		})
	}
}
