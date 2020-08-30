package checker

import "fmt"

type CheckResult struct {
	domain     string
	statusCode int
	error      error
}

var (
	// ErrHostNotFound error when domain name wrong
	ErrHostNotFound = fmt.Errorf("no such host")

	// ErrConnectionResetByPeer reset by peer
	ErrConnectionResetByPeer = fmt.Errorf("connection reset by peer")

	// ErrConnectionClosed server close connection without sending bytes
	ErrConnectionClosed = fmt.Errorf("the server closed connection before returning the first response byte")

	// ErrTimeout timeout error
	ErrTimeout = fmt.Errorf("timeout")
)

func (result *CheckResult) Domain() string {
	return result.domain
}

func (result *CheckResult) StatusCode() int {
	return result.statusCode
}

func (result *CheckResult) Error() error {
	return result.error
}
