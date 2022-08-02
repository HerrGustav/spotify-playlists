package spotify

import "fmt"

type errCode int64

const (
	notAuthorized errCode = iota
)

func (e errCode) String() string {
	return [...]string{
		"notAuthorized",
	}[e]
}

type errSpotify struct {
	code errCode
	msg  string
	err  error
}

func (e errSpotify) Error() string {
	m := fmt.Sprintf("spotify client error code '%d'(%s) - '%s'", e.code, e.code.String(), e.msg)

	// if the underlying error is not existing return here,
	// because an e.err.Error() which is nil would result in a panic.
	if e.err == nil {
		return m
	}

	return fmt.Sprintf("%s, %s", m, e.err.Error())
}

func newError(code errCode, msg string, err error) errSpotify {
	return errSpotify{
		code,
		msg,
		err,
	}
}
