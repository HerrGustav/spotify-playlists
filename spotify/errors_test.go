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
			want: "spotify client error code '0'(notAuthorized) - 'mock'",
		},
		" returns underlying error in output if existing": {
			err: errSpotify{
				code: notAuthorized,
				msg:  "mock",
				err:  errors.New("some error"),
			},
			want: "spotify client error code '0'(notAuthorized) - 'mock', some error",
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
