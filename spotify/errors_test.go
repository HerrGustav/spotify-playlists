package spotify

import (
	"errors"
	"testing"
)

func TestErrorMethod(t *testing.T) {
	testcases := map[string]struct {
		err  errSpotify
		want string
	}{
		"does not return underlying error in output if not existing": {
			err: errSpotify{
				code: notAuthorized,
				msg:  "mock",
				err:  nil,
			},
			want: "spotify client error, code: '0'(notAuthorized) - 'mock'",
		},
		" returns underlying error in output if existing": {
			err: errSpotify{
				code: notAuthorized,
				msg:  "mock",
				err:  errors.New("some error"),
			},
			want: "spotify client error, code: '0'(notAuthorized) - 'mock', some error",
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			got := tc.err.Error()
			if tc.want != got {
				t.Errorf("errSpotify.Error() mismatch: \n got: '%s', \n want: '%s'", tc.want, got)
			}
		})
	}
}

func TestIs(t *testing.T) {
	mockSpotifyErr := errSpotify{
		code: notAuthorized,
	}

	testcases := map[string]struct {
		err   errSpotify
		input error
		want  bool
	}{
		"returns false if a nil error is parsed in": {
			err:  mockSpotifyErr,
			want: false,
		},
		"returns false if another error type is parsed in": {
			err:   mockSpotifyErr,
			input: errors.New("mock error"),
			want:  false,
		},
		"returns false if a spotify error with different error code is parsed in": {
			err:   mockSpotifyErr,
			input: errSpotify{code: requestFailed},
			want:  false,
		},
		"returns true if a spotify error with the same error code is parsed in": {
			err:   mockSpotifyErr,
			input: errSpotify{code: notAuthorized},
			want:  true,
		},
	}

	for testName, tc := range testcases {
		t.Run(testName, func(t *testing.T) {
			got := errors.Is(tc.err, tc.input)
			if got != tc.want {
				t.Errorf("errSpotify.Is(): unexpected result, \n - got: '%v', \n - want: '%v'", got, tc.want)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	// Check that unwrap returns something
	err := errSpotify{
		code: notAuthorized,
		msg:  "test",
		err:  errMock,
	}

	want := errMock

	got := err.Unwrap()
	if got != want {
		t.Errorf("errSpotify.Is(): unexpected result, \n - got: '%v', \n - want: '%v'", got, want)
	}
}
